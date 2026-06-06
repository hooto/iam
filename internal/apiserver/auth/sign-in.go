// Copyright 2014 Eryx <evorui at gmail dot com>, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auth

import (
	"log/slog"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hooto/httpsrv"
	"github.com/lessos/lessgo/pass"
	"github.com/sysinner/incore/v2/pkg/inauth"

	"github.com/hooto/iam/v2/internal/apiserver"
	"github.com/hooto/iam/v2/internal/data"
	"github.com/hooto/iam/v2/internal/util"
	"github.com/hooto/iam/v2/pkg/iamapi"
)

// urlHost extracts the host from a URL string.
func urlHost(requrl string) string {

	u, err := url.Parse(requrl)

	if err != nil {
		return "localhost"
	}

	if i := strings.Index(u.Host, ":"); i > 0 {
		return u.Host[:i]
	}

	return u.Host
}

type SignInRequest struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	AppId         string `json:"app_id,omitempty"`
	RedirectToken string `json:"redirect_token,omitempty"`
}

type SignInResponse struct {
	Status      inauth.ServiceStatus `json:"status"`
	RedirectUri string               `json:"redirect_uri,omitempty"`
	AccessToken string               `json:"access_token,omitempty"`
}

// SignIn handles user authentication via username/password,
// creates a session, and returns an access token or auth code.
func SignIn(ctx httpsrv.Ctx) error {

	var (
		req SignInRequest
		rsp = SignInResponse{
			RedirectUri: "/iam",
		}
	)
	defer ctx.JSON(&rsp)

	if err := ctx.Request().JsonDecode(&req); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", "Invalid request format")
		return nil
	}

	var app iamapi.AppInstance
	if req.AppId != "" {

		if rs := data.Data.NewReader(iamapi.NsAppInstance(req.AppId)).Exec(); !rs.OK() {
			rsp.Status = inauth.NewServiceStatus("404", "App not found")
			return nil
		} else {
			rs.Item().JsonDecode(&app)
			if app.ID != req.AppId {
				rsp.Status = inauth.NewServiceStatus("404", "App not found")
				return nil
			}
			if app.Url == "" {
				rsp.Status = inauth.NewServiceStatus("500", "App not setup callback-url")
				return nil
			}
			rsp.RedirectUri = app.Url
		}
	}

	req.Username = strings.ToLower(req.Username)
	if err := iamapi.UsernameValid(req.Username); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", err.Error())
		return nil
	}

	if req.Password == "" {
		rsp.Status = inauth.NewServiceStatus("400", "Username or Password can not be empty")
		return nil
	}

	user := data.UserGet(req.Username)
	if user == nil {
		rsp.Status = inauth.NewServiceStatus("400", "incorrect username or password 2")
		slog.Info("service/signin-auth fail", "user", req.Username)
		return nil
	}

	if user.Type == iamapi.UserTypeGroup {
		rsp.Status = inauth.NewServiceStatus("400", "incorrect username or password 3")
		slog.Info("service/signin-auth fail", "user", req.Username)
		return nil
	}

	denyCount, denyKey, err := apiserver.UserAuthDenyCheck(req.Username, ctx.Request())
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("400", err.Error())
		return nil
	}

	if auth := user.Keys.Get(iamapi.UserKeyDefault); auth == nil ||
		!pass.Check(req.Password, auth.String()) {
		apiserver.UserAuthDenyIncr(denyCount, denyKey)
		rsp.Status = inauth.NewServiceStatus("400", "incorrect username or password")
		slog.Info("service/signin-auth fail", "user", req.Username)
		return nil
	}

	var (
		sid = uuid.NewString()

		ttl = int64(86400 * 7)

		at = inauth.NewAccessToken()

		it = &inauth.IdentityToken{
			Roles:  user.Roles,
			Groups: data.UserGroups(req.Username),
		}

		st = inauth.SessionToken{
			AccessToken:   at,
			IdentityToken: it,
		}
	)

	at.Claims.Sub = user.Name
	at.Claims.Iat = time.Now().Unix()
	at.Claims.Exp = at.Claims.Iat + ttl
	at.Claims.Jti = sid

	accessToken, err := at.SignToken(data.KeyMgr)
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("500", err.Error())
		return nil
	}

	if app.ID != "" {
		for _, perm := range app.Permissions {
			if slices.Contains(st.IdentityToken.Permissions, perm.Permission) {
				continue
			}
			if app.User != user.Name &&
				!util.Contains(st.IdentityToken.Roles, perm.Roles) {
				continue
			}
			st.IdentityToken.Permissions = append(
				st.IdentityToken.Permissions, perm.Permission)
		}
	}

	if rs := data.Data.NewWriter(
		iamapi.NsUserSession(at.Claims.Jti, uint32(at.Claims.Exp)), nil).SetJsonValue(st).
		SetTTL(ttl * 1000).Exec(); !rs.OK() {
		rsp.Status = inauth.NewServiceStatus("500", rs.ErrorMessage())
		return nil
	}

	slog.Info("auth/sign-in session", "body", st)

	if req.AppId != "" {

		// cross-host: use auth code flow
		authCode := uuid.NewString()

		if rs := data.Data.NewWriter(iamapi.NsAuthCode(authCode), nil).
			SetJsonValue(&iamapi.SignInAuthCode{
				AccessToken: accessToken,
				Username:    user.Name,
			}).SetTTL(3600e3).Exec(); !rs.OK() {
			rsp.Status = inauth.NewServiceStatus("500", rs.ErrorMessage())
			return nil
		}

		if !strings.Contains(rsp.RedirectUri, "?") {
			rsp.RedirectUri += "?"
		} else {
			rsp.RedirectUri += "&"
		}
		rsp.RedirectUri += "code=" + authCode

	} else if len(req.RedirectToken) > 20 {

		rt := iamapi.ServiceRedirectTokenDecode(req.RedirectToken)

		if len(rt.RedirectUri) > 0 {

			rsp.RedirectUri = rt.RedirectUri

			if urlHost(rsp.RedirectUri) != urlHost(ctx.Request().URL.Host) {

				// cross-host: use auth code flow
				authCode := inauth.RandHexString(16)

				if rs := data.Data.NewWriter(iamapi.NsAuthCode(authCode), nil).
					SetJsonValue(&iamapi.SignInAuthCode{
						AccessToken: accessToken,
						Username:    user.Name,
					}).SetTTL(600e3).Exec(); !rs.OK() {
					rsp.Status = inauth.NewServiceStatus("500", rs.ErrorMessage())
					return nil
				}

				if !strings.Contains(rsp.RedirectUri, "?") {
					rsp.RedirectUri += "?"
				} else {
					rsp.RedirectUri += "&"
				}

				rsp.RedirectUri += "code=" + authCode

				if len(rt.State) > 0 {
					rsp.RedirectUri += "&state=" + url.QueryEscape(rt.State)
				}

				// don't expose access_token to browser for cross-host flow
				rsp.AccessToken = ""
			}
		}
	} else {
		rsp.AccessToken = accessToken
	}

	if req.AppId == "" {
		http.SetCookie(ctx.Response().Out, &http.Cookie{
			Name:     inauth.AppHttpHeaderKey,
			Value:    rsp.AccessToken,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Unix(at.Claims.Exp, 0),
		})
	}

	rsp.Status = inauth.NewServiceStatus("200", "ok")

	slog.Info("auth/sign-in response", "body", rsp)

	return nil
}
