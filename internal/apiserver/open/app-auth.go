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

package open

import (
	"log/slog"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/hooto/httpsrv"
	"github.com/sysinner/innerstack/v2/pkg/inapi"
	"github.com/sysinner/innerstack/v2/pkg/inauth"
	"golang.org/x/mod/semver"

	"github.com/hooto/iam/v2/internal/data"
	"github.com/hooto/iam/v2/pkg/iamapi"
)

// appAuth verifies the access token from cookie and returns the
// parsed AccessKey. It writes a 401 JSON response and returns nil on failure.
func appAuth(ctx httpsrv.Ctx) *inauth.AccessKey {

	tokenStr := ""
	cookie, err := ctx.Request().Cookie(inauth.AppHttpHeaderKey)
	if err == nil && cookie.Value != "" {
		tokenStr = cookie.Value
	}

	if tokenStr == "" {
		authHeader := ctx.Request().Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenStr = parts[1]
			}
		}
	}

	if tokenStr == "" {
		ctx.JSON(struct {
			Status inauth.ServiceStatus `json:"status"`
		}{Status: inauth.NewServiceStatus("401", "Unauthorized #1")})
		return nil
	}

	token, err := inauth.ParseAccessToken(tokenStr)
	if err != nil || token.Header.Kid == "" {
		ctx.JSON(struct {
			Status inauth.ServiceStatus `json:"status"`
		}{Status: inauth.NewServiceStatus("401", "Unauthorized #2 "+err.Error())})
		return nil
	}

	ak, err := token.Verify(data.KeyMgr)
	if err != nil {
		ctx.JSON(struct {
			Status inauth.ServiceStatus `json:"status"`
		}{Status: inauth.NewServiceStatus("401", "Unauthorized #3 "+err.Error())})
		return nil
	}

	return ak
}

// AppAuth_Verify verifies app credentials (app_id + secret_key).
// Used by third-party apps during setup to validate their IAM configuration.
func AppAuth_Verify(ctx httpsrv.Ctx) error {

	ak := appAuth(ctx)
	if ak == nil {
		return nil
	}

	var (
		req struct {
			AppId string `json:"app_id"`
		}
		rsp struct {
			Status inauth.ServiceStatus `json:"status"`
			App    *iamapi.AppInstance  `json:"app,omitempty"`
		}
	)
	defer ctx.JSON(&rsp)

	if err := ctx.Request().JsonDecode(&req); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", "Bad Argument")
		return nil
	}

	if req.AppId == "" || req.AppId != ak.Id {
		rsp.Status = inauth.NewServiceStatus("400", "app_id required")
		return nil
	}

	var app iamapi.AppInstance
	if rs := data.Data.NewReader(iamapi.NsAppInstance(req.AppId)).Exec(); rs.OK() {
		rs.Item().JsonDecode(&app)
	}

	if app.ID != req.AppId {
		rsp.Status = inauth.NewServiceStatus("404", "App not found")
		return nil
	}

	app.SecretKey = ""

	rsp.App = &app
	rsp.Status = inauth.NewServiceStatus("200", "ok")

	slog.Info("open/app-auth-verify response", "body", rsp)

	return nil
}

type AppAuth_TokenExchangeRequest struct {
	Code string `json:"code"`
}

type AppAuth_TokenExchangeResponse struct {
	Status        inauth.ServiceStatus  `json:"status"`
	AccessToken   string                `json:"access_token,omitempty"`
	IdentityToken *inauth.IdentityToken `json:"identity_token,omitempty"`
}

// AppAuth_TokenExchange exchanges a one-time auth code for access_token + identity_token.
// Called by third-party app backend using app credentials.
func AppAuth_TokenExchange(ctx httpsrv.Ctx) error {

	var (
		req AppAuth_TokenExchangeRequest
		rsp AppAuth_TokenExchangeResponse
	)
	defer ctx.JSON(&rsp)

	if err := ctx.Request().JsonDecode(&req); err != nil {
		// rsp.Status = inauth.NewServiceStatus("400", "Invalid request format")
		// return nil
	}

	if req.Code == "" {
		if ctx.Params().Value("code") == "" {
			rsp.Status = inauth.NewServiceStatus("400", "code required")
			return nil
		}
		req.Code = ctx.Params().Value("code")
	}

	// look up auth code
	var entry iamapi.SignInAuthCode
	if rs := data.Data.NewReader(iamapi.NsAuthCode(req.Code)).Exec(); !rs.OK() {
		rsp.Status = inauth.NewServiceStatus("404", "Invalid or expired auth code")
		return nil
	} else {
		rs.Item().JsonDecode(&entry)
	}

	// one-time use: delete the code
	data.Data.NewDeleter(iamapi.NsAuthCode(req.Code)).Exec()

	if entry.AccessToken == "" {
		rsp.Status = inauth.NewServiceStatus("404", "Invalid auth code")
		return nil
	}

	// parse and verify the stored access token
	token, err := inauth.ParseAccessToken(entry.AccessToken)
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("401", "Invalid access token : "+err.Error())
		return nil
	}

	if _, err := token.Verify(data.KeyMgr); err != nil {
		rsp.Status = inauth.NewServiceStatus("401", "Invalid access token : "+err.Error())
		return nil
	}

	// lookup session for identity token
	var st inauth.SessionToken
	if rs := data.Data.NewReader(iamapi.NsUserSession(token.Claims.Jti, uint32(token.Claims.Exp))).Exec(); rs.OK() {
		rs.Item().JsonDecode(&st)
	}

	rsp.Status = inauth.NewServiceStatus("200", "ok")
	rsp.AccessToken = entry.AccessToken
	if st.IdentityToken != nil {
		rsp.IdentityToken = st.IdentityToken
	}

	slog.Info("open/app-auth/token-exchange response", "body", rsp)

	return nil
}

type AppAuth_UpdateRequest struct {
	AppId       string                  `json:"app_id,omitempty" toml:"app_id,omitempty"`
	Name        string                  `json:"name" toml:"name"`
	Version     string                  `json:"version,omitempty" toml:"version,omitempty"`
	Permissions []*iamapi.AppPermission `json:"permissions" toml:"permissions,omitempty"`
}

type AppAuth_UpdateResponse struct {
	Status inauth.ServiceStatus `json:"status"`
}

func AppAuth_Update(ctx httpsrv.Ctx) error {

	ak := appAuth(ctx)
	if ak == nil {
		return nil
	}

	var (
		req AppAuth_UpdateRequest
		rsp AppAuth_UpdateResponse
	)
	defer ctx.JSON(&rsp)

	if err := ctx.Request().JsonDecode(&req); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", "Invalid request format")
		return nil
	}

	if req.AppId == "" || req.AppId != ak.Id {
		rsp.Status = inauth.NewServiceStatus("400", "app_id required")
		return nil
	}

	if req.Name != "" {
		if err := iamapi.DNSLabelValid(req.Name); err != nil {
			rsp.Status = inauth.NewServiceStatus("400", "invalid name : "+err.Error())
			return nil
		}
	}

	if err := inapi.SemverValid(req.Version); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", "invalid version : "+err.Error())
		return nil
	}

	var app iamapi.AppInstance
	if rs := data.Data.NewReader(iamapi.NsAppInstance(req.AppId)).Exec(); rs.OK() {
		rs.Item().JsonDecode(&app)
	}

	if app.ID != req.AppId {
		rsp.Status = inauth.NewServiceStatus("404", "App not found")
		return nil
	}

	if req.Name != "" && req.Name != app.Name {
		app.Name = req.Name
	}

	// version can only be updated to a higher version
	if app.Version != "" && semver.Compare(req.Version, app.Version) < 0 {
		rsp.Status = inauth.NewServiceStatus("400", "version must be greater than current version")
		return nil
	}

	if app.Version == "" ||
		semver.Compare(req.Version, app.Version) < 0 {
		app.Version = req.Version
	}

	if len(req.Permissions) > 0 {
		perms := []*iamapi.AppPermission{}
		for _, v := range req.Permissions {
			v.Permission = strings.ToLower(v.Permission)
			if err := iamapi.PermissionNameValid(v.Permission); err == nil {
				if !slices.ContainsFunc(perms, func(item *iamapi.AppPermission) bool {
					return item.Permission == v.Permission
				}) {
					perms = append(perms, v)
				}
			} else {
				slog.Warn("app update", "invalid-perm-name", v.Permission,
					"err", err.Error())
			}
		}
		sort.Slice(perms, func(i, j int) bool {
			return perms[i].Permission < perms[j].Permission
		})
		app.Permissions = perms
	}

	app.Updated = time.Now().Unix()
	slog.Info("app-auth update", "app", app)

	if rs := data.Data.NewWriter(iamapi.NsAppInstance(app.ID), app).Exec(); !rs.OK() {
		rsp.Status = inauth.NewServiceStatus("500", "Failed to update app")
		return nil
	}

	rsp.Status = inauth.NewServiceStatus("200", "ok")
	slog.Info("open/app-auth/update response", "body", rsp)

	return nil
}

type AppAuth_SessionRequest struct {
	AppId       string `json:"app_id,omitempty" toml:"app_id,omitempty"`
	AccessToken string `json:"access_token,omitempty" toml:"access_token,omitempty"`
}

type AppAuth_SessionResponse struct {
	Status        inauth.ServiceStatus  `json:"status"`
	IdentityToken *inauth.IdentityToken `json:"identity_token,omitempty"`
}

func AppAuth_Session(ctx httpsrv.Ctx) error {

	ak := appAuth(ctx)
	if ak == nil {
		return nil
	}

	var (
		req AppAuth_SessionRequest
		rsp AppAuth_SessionResponse
	)
	defer ctx.JSON(&rsp)

	if err := ctx.Request().JsonDecode(&req); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", "Invalid request format")
		return nil
	}

	if req.AppId == "" || req.AppId != ak.Id {
		rsp.Status = inauth.NewServiceStatus("400", "app_id required")
		return nil
	}

	// var app iamapi.AppInstance
	// if rs := data.Data.NewReader(iamapi.NsAppInstance(req.AppId)).Exec(); rs.OK() {
	// 	rs.Item().JsonDecode(&app)
	// }

	// if app.ID != req.AppId {
	// 	rsp.Status = inauth.NewServiceStatus("404", "App not found")
	// 	return nil
	// }

	utoken, err := inauth.ParseAccessToken(req.AccessToken)
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("401", "Invalid access token : "+err.Error())
		return nil
	}

	_, err = utoken.Verify(data.KeyMgr)
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("401", "Invalid access token : "+err.Error())
		return nil
	}

	// lookup session from DB
	var (
		key = iamapi.NsUserSession(utoken.Claims.Jti, uint32(utoken.Claims.Exp))
		st  inauth.SessionToken
	)
	if rs := data.Data.NewReader(key).Exec(); !rs.OK() {
		rsp.Status = inauth.NewServiceStatus("401", "session not found")
	} else {
		if err := rs.Item().JsonDecode(&st); err != nil {
			rsp.Status = inauth.NewServiceStatus("500", "failed to decode session")
		} else {
			rsp.Status = inauth.NewServiceStatus("200", "ok")
			rsp.IdentityToken = st.IdentityToken
		}
	}

	slog.Info("open/app-auth/session response", "body", rsp)

	return nil
}
