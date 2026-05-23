// Copyright 2014 Eryx <evorui аt gmаil dοt cοm>, All rights reserved.
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

package apiserver

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hooto/httpsrv"
	"github.com/lessos/lessgo/net/email"
	"github.com/lessos/lessgo/pass"
	"github.com/sysinner/incore/v2/pkg/inauth"

	"github.com/hooto/iam/v2/internal/config"
	"github.com/hooto/iam/v2/internal/data"
	"github.com/hooto/iam/v2/pkg/iamapi"
)

type Service struct {
	*httpsrv.Controller
}

type UserSignInRequest struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	RedirectToken string `json:"redirect_token,omitempty"` // 可选的重定向令牌，包含重定向 URI 和 state
}

type UserSignInResponse struct {
	Status      inauth.ServiceStatus `json:"status"`
	RedirectUri string               `json:"redirect_uri,omitempty"`
	AccessToken string               `json:"access_token,omitempty"`
}

func (c Service) UserSignInAction() {

	var (
		req UserSignInRequest
		rsp = UserSignInResponse{
			RedirectUri: "/iam",
		}
	)
	defer c.RenderJson(&rsp)

	if err := c.Request.JsonDecode(&req); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", "Invalid request format")
		return
	}

	req.Username = strings.ToLower(req.Username)
	if err := iamapi.UsernameValid(req.Username); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", err.Error())
		return
	}

	if req.Password == "" {
		rsp.Status = inauth.NewServiceStatus("400", "Username or Password can not be empty")
		return
	}

	user := data.UserGet(req.Username)
	if user == nil {
		rsp.Status = inauth.NewServiceStatus("400", "incorrect username or password 2")
		slog.Info("service/signin-auth fail", "user", req.Username)
		return
	}

	if user.Type == iamapi.UserTypeGroup {
		rsp.Status = inauth.NewServiceStatus("400", "incorrect username or password 3")
		slog.Info("service/signin-auth fail", "user", req.Username)
		return
	}

	denyCount, denyKey, err := userAuthDenyCheck(req.Username, c.Request)
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("400", err.Error())
		return
	}

	if auth := user.Keys.Get(iamapi.UserKeyDefault); auth == nil ||
		!pass.Check(req.Password, auth.String()) {
		userAuthDenyIncr(denyCount, denyKey)
		rsp.Status = inauth.NewServiceStatus("400", "incorrect username or password")
		slog.Info("service/signin-auth fail", "user", req.Username)
		return
	}

	var (
		ttl = int64(864000)

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
	at.Claims.Exp = at.Claims.Iat + 864000
	at.Claims.Jti = uuid.NewString()

	accessToken, err := at.SignToken(data.KeyMgr)
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("500", err.Error())
		return
	}

	if rs := data.Data.NewWriter(
		iamapi.NsUserSession(at.Claims.Jti, uint32(at.Claims.Exp)), nil).SetJsonValue(st).
		SetTTL(ttl * 1000).Exec(); !rs.OK() {
		rsp.Status = inauth.NewServiceStatus("500", rs.ErrorMessage())
		return
	}

	slog.Info("user-session", "body", st)

	rsp.AccessToken = accessToken

	if len(req.RedirectToken) > 20 {

		rt := iamapi.ServiceRedirectTokenDecode(req.RedirectToken)

		if len(rt.RedirectUri) > 0 {

			rsp.RedirectUri = rt.RedirectUri

			if urlHost(rsp.RedirectUri) != urlHost(c.Request.URL.Host) {

				if strings.Index(rsp.RedirectUri, "?") == -1 {
					rsp.RedirectUri += "?"
				} else {
					rsp.RedirectUri += "&"
				}

				rsp.RedirectUri += inauth.AppHttpHeaderKey + "=" + rsp.AccessToken +
					"&expires_in=" + strconv.Itoa(int(ttl))

				if len(rt.State) > 0 {
					rsp.RedirectUri += "&state=" + url.QueryEscape(rt.State)
				}
			}
		}
	}

	http.SetCookie(c.Response.Out, &http.Cookie{
		Name:     inauth.AppHttpHeaderKey,
		Value:    rsp.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(at.Claims.Exp, 0), // time.Now().Add(time.Duration(ttl) * time.Second),
	})

	rsp.Status = inauth.NewServiceStatus("200", "ok")

	slog.Info("service/signin-auth ok", "user", user.Name)
}

type UserAuthRequest struct {
	AppID        string `json:"app_id" toml:"app_id"`
	AccessToken  string `json:"access_token,omitempty" toml:"access_token,omitempty"`
	RedirectURI  string `json:"redirect_uri" toml:"redirect_uri"`
	State        string `json:"state" toml:"state"`                 // 防 CSRF 随机串
	ResponseType string `json:"response_type" toml:"response_type"` // 固定为 "code"
}

type ServiceStatusResponse struct {
	Status inauth.ServiceStatus `json:"status"`
}

func (c Service) UserSignOutAction() {

	var (
		req UserAuthRequest
		rsp ServiceStatusResponse
	)
	defer c.RenderJson(&rsp)

	c.Request.JsonDecode(&req)

	if req.AccessToken == "" {
		// fallback to http-only cookie
		cookie, err := c.Request.Cookie(inauth.AppHttpHeaderKey)
		if err != nil || cookie.Value == "" {
			rsp.Status = inauth.NewServiceStatus("401", "access token not found")
			return
		}
		req.AccessToken = cookie.Value
	}

	token, err := inauth.ParseAccessToken(req.AccessToken)
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("401", "invalid access token")
		return
	}

	data.Data.NewDeleter(iamapi.NsUserSession(token.Claims.Jti, uint32(token.Claims.Exp))).
		Exec()

	http.SetCookie(c.Response.Out, &http.Cookie{
		Name:   inauth.AppHttpHeaderKey,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	rsp.Status = inauth.NewServiceStatus("200", "ok")
}

type UserSessionResponse struct {
	Status        inauth.ServiceStatus  `json:"status"`
	AccessToken   string                `json:"access_token,omitempty"`
	IdentityToken *inauth.IdentityToken `json:"identity_token,omitempty"`
}

func (c Service) UserSessionAction() {

	var (
		req UserAuthRequest
		rsp = UserSessionResponse{}
	)
	defer c.RenderJson(&rsp)

	c.Request.JsonDecode(&req)

	if req.AccessToken == "" {
		// fallback to http-only cookie
		cookie, err := c.Request.Cookie(inauth.AppHttpHeaderKey)
		if err != nil || cookie.Value == "" {
			rsp.Status = inauth.NewServiceStatus("401", "access token not found")
			return
		}
		req.AccessToken = cookie.Value
	}

	token, err := inauth.ParseAccessToken(req.AccessToken)
	if err != nil {
		slog.Info("service/session-token fail", "error", err.Error(), "access_token", token)
		rsp.Status = inauth.NewServiceStatus("401", "invalid access token : "+err.Error())
		return
	}

	// verify signature
	if _, err := token.Verify(data.KeyMgr); err != nil {
		rsp.Status = inauth.NewServiceStatus("401", "invalid access token : "+err.Error())
		return
	}

	// lookup session from DB
	var (
		key = iamapi.NsUserSession(token.Claims.Jti, uint32(token.Claims.Exp))
		st  inauth.SessionToken
	)
	if rs := data.Data.NewReader(key).Exec(); !rs.OK() {
		rsp.Status = inauth.NewServiceStatus("401", "session not found")
	} else {
		if err := rs.Item().JsonDecode(&st); err != nil {
			rsp.Status = inauth.NewServiceStatus("500", "failed to decode session")
		} else {
			rsp.Status = inauth.NewServiceStatus("200", "ok")
			rsp.AccessToken = req.AccessToken
			rsp.IdentityToken = st.IdentityToken
		}
	}
}

type UserSignUpRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSignUpResponse struct {
	Status   inauth.ServiceStatus `json:"status"`
	Continue string               `json:"continue,omitempty"`
	Username string               `json:"username,omitempty"`
}

func (c Service) SignUpAction() {

	var (
		req UserSignUpRequest
		rsp = UserSignUpResponse{
			Continue: "/iam/service/sign-in",
		}
	)
	defer c.RenderJson(&rsp)

	if err := c.Request.JsonDecode(&req); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", "Invalid request format")
		return
	}

	if !config.AllowUserSignUp {
		rsp.Status = inauth.NewServiceStatus("403", "User Registration Disabled")
		return
	}

	denyCount, denyKey, err := userAuthDenyCheck(req.Username, c.Request)
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("400", err.Error())
		return
	}

	// validate username
	req.Username = strings.ToLower(req.Username)
	if err := iamapi.UsernameValid(req.Username); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", err.Error())
		return
	}

	// validate email
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if err := iamapi.EmailValid(email); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", err.Error())
		return
	}

	// validate password
	if len(req.Password) < 8 || len(req.Password) > 30 {
		rsp.Status = inauth.NewServiceStatus("400", "Password must be between 8 and 30 characters long")
		return
	}

	// check if user already exists
	var existing iamapi.User
	if obj := data.Data.NewReader(iamapi.NsUser(req.Username)).Exec(); obj.OK() {
		obj.Item().JsonDecode(&existing)
	}
	if existing.Name == req.Username {
		rsp.Status = inauth.NewServiceStatus("400", "The username already exists, please choose another one")
		return
	}

	auth, _ := pass.HashDefault(req.Password)

	tn := time.Now().Unix()

	user := iamapi.User{
		Name:        req.Username,
		Email:       email,
		DisplayName: strings.ToUpper(req.Username[:1]) + req.Username[1:],
		Status:      1,
		Roles:       []string{iamapi.Role_User},
		Created:     tn,
		Updated:     tn,
	}
	user.Keys.Set(iamapi.UserKeyDefault, auth)

	if !data.UserSet(&user) {
		rsp.Status = inauth.NewServiceStatus("500", "Server Error")
		return
	}

	rsp.Status = inauth.NewServiceStatus("200", "ok")
	rsp.Username = req.Username

	userAuthDenyIncr(denyCount, denyKey)

	slog.Info("service/signup ok", "user", req.Username)
}

/**
func (c Service) SessionTokenAction() {
	var (
		req inauth.AuthRequest
		rsp inauth.AuthTokenResponse
	)
	defer c.RenderJson(&rsp)

	if err := c.Request.JsonDecode(&req); err != nil {
		rsp.Status = inauth.NewServiceStatus("401", err.Error())
		return
	}

	at, err := inauth.ParseAccessToken(req.AccessToken)
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("401", err.Error())
		return
	}

	ak, err := at.Verify(data.KeyMgr)
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("401", err.Error())
		return
	}

	user := data.UserGet(at.Claims.Sub)
	if user == nil {
		rsp.Status = inauth.NewServiceStatus("401", "incorrect username")
		return
	}

	if user.Type == iamapi.UserTypeGroup {
		rsp.Status = inauth.NewServiceStatus("401", "incorrect username")
		return
	}

	rsp.AccessToken = req.AccessToken
	rsp.IdentityToken = &inauth.IdentityToken{
		Roles:  user.Roles,
		Groups: data.UserGroups(at.Claims.Sub),
		Scopes: ak.Scopes,
	}
}

func (c Service) AuthAction() {

	var (
		req inauth.AuthLoginRequest
		rsp inauth.AuthLoginResponse
	)
	defer c.RenderJson(&rsp)

	if err := c.Request.JsonDecode(&req); err != nil {
		rsp.Error = err.Error()
		return
	}

	token, err := inauth.ParseAccessToken(req.LoginToken)
	if err != nil {
		rsp.Error = err.Error()
		return
	}

	ak, err := token.Verify(data.KeyMgr)
	if err != nil {
		rsp.Error = err.Error()
		return
	}

	if ak.User == "" {
		rsp.Error = "access-key not found"
		return
	}

	user := data.UserGet(ak.User)
	if user == nil {
		rsp.Error = "incorrect username or password"
		return
	}

	if user.Type == iamapi.UserTypeGroup {
		rsp.Error = "incorrect username or password"
		return
	}

	iat := time.Now().Unix()

	header := inauth.TokenHeader{
		Kid: ak.Id,
	}

	claims := inauth.AuthClaims{
		Jti: inauth.RandHexString(16),
		Iat: iat,
		Exp: iat + 86400,
		Sub: ak.User,
	}

	at, err := inauth.Sign(header, claims, []byte(ak.Secret))
	if err != nil {
		rsp.Error = err.Error()
		return
	}
	rsp.AccessToken = at

	rsp.IdentityToken = inauth.IdentityToken{
		Roles:  user.Roles,
		Groups: data.UserGroups(claims.Sub),
		Scopes: ak.Scopes,
	}
}
*/

type ResetPasswordTicketRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// ResetPasswordTicketAction handles forgot-password requests.
// It validates username + email, generates a reset token, and sends a reset email.
func (c Service) ResetPasswordTicketAction() {

	var (
		req ResetPasswordTicketRequest
		rsp ServiceStatusResponse
	)
	defer c.RenderJson(&rsp)

	if err := c.Request.JsonDecode(&req); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", "Invalid request format")
		return
	}

	req.Username = strings.ToLower(req.Username)
	if err := iamapi.UsernameValid(req.Username); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", err.Error())
		return
	}

	emailAddr := strings.ToLower(strings.TrimSpace(req.Email))
	if err := iamapi.EmailValid(emailAddr); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", err.Error())
		return
	}

	denyCount, denyKey, err := userAuthDenyCheck(req.Username, c.Request)
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("400", err.Error())
		return
	}
	userAuthDenyIncr(denyCount, denyKey)

	user := data.UserGet(req.Username)
	if user == nil || user.Email != emailAddr {
		// always return success to prevent user enumeration
		rsp.Status = inauth.NewServiceStatus("200", "ok")
		return
	}

	// store reset token with 1-hour TTL
	reset := iamapi.UserResetPassword{
		Id:       uuid.NewString(),
		Username: user.Name,
		Email:    emailAddr,
		Expired:  time.Now().Add(3600 * time.Second).Unix(),
	}

	if rs := data.Data.NewWriter(iamapi.NsUserResetPassword(reset.Id), nil).
		SetJsonValue(reset).SetTTL(3600e3).Exec(); !rs.OK() {
		rsp.Status = inauth.NewServiceStatus("500", "Internal server error")
		return
	}

	// send reset email
	mr, err := email.MailerPull("def")
	if err != nil {
		slog.Error("reset-password: mailer not available", "error", err)
		rsp.Status = inauth.NewServiceStatus("200", "ok")
		return
	}

	body := fmt.Sprintf(`<html>
<body>
<div>You recently requested a password reset for your %s account.</div>
<br>
<div>Your verification code is: <b>%s</b></div>
<br>
<div>This code will expire in 1 hour.</div>
<br>
<div>If you did not make this request, please ignore this email.</div>
<br>
<div>Regards,</div>
<div>%s Account Service</div>
</body>
</html>`, config.Config.ServiceName, reset.Id, config.Config.ServiceName)

	if err := mr.SendMail(emailAddr, "Reset your password", body); err != nil {
		slog.Error("reset-password: send mail failed", "error", err)
	}

	rsp.Status = inauth.NewServiceStatus("200", "ok")
}

type ResetPasswordConfirmRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

// ResetPasswordConfirmAction sets a new password using a valid reset token.
func (c Service) ResetPasswordConfirmAction() {

	var (
		req ResetPasswordConfirmRequest
		rsp ServiceStatusResponse
	)
	defer c.RenderJson(&rsp)

	if err := c.Request.JsonDecode(&req); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", "Invalid request format")
		return
	}

	if req.Token == "" {
		rsp.Status = inauth.NewServiceStatus("400", "Token is required")
		return
	}

	if len(req.Password) < 8 || len(req.Password) > 30 {
		rsp.Status = inauth.NewServiceStatus("400", "Password must be between 8 and 30 characters long")
		return
	}

	// lookup reset token
	var reset iamapi.UserResetPassword
	if rs := data.Data.NewReader(iamapi.NsUserResetPassword(req.Token)).Exec(); !rs.OK() {
		rsp.Status = inauth.NewServiceStatus("400", "Invalid or expired token")
		return
	} else {
		if err := rs.Item().JsonDecode(&reset); err != nil {
			rsp.Status = inauth.NewServiceStatus("400", "Invalid or expired token")
			return
		}
	}

	if reset.Id != req.Token {
		rsp.Status = inauth.NewServiceStatus("400", "Invalid or expired token")
		return
	}

	denyCount, denyKey, err := userAuthDenyCheck(reset.Username, c.Request)
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("400", err.Error())
		return
	}
	userAuthDenyIncr(denyCount, denyKey)

	// lookup user
	var user iamapi.User
	if rs := data.Data.NewReader(iamapi.NsUser(reset.Username)).Exec(); rs.OK() {
		rs.Item().JsonDecode(&user)
	}

	if user.Name != reset.Username {
		rsp.Status = inauth.NewServiceStatus("400", "User not found")
		return
	}

	// update password
	user.Updated = time.Now().Unix()
	auth, _ := pass.HashDefault(req.Password)
	user.Keys.Set(iamapi.UserKeyDefault, auth)

	if rs := data.Data.NewWriter(iamapi.NsUser(reset.Username), nil).
		SetJsonValue(user).SetIncr(0, "user").Exec(); !rs.OK() {
		rsp.Status = inauth.NewServiceStatus("500", "Internal server error")
		return
	}

	// delete the used reset token
	data.Data.NewDeleter(iamapi.NsUserResetPassword(reset.Id)).Exec()

	rsp.Status = inauth.NewServiceStatus("200", "ok")

	slog.Info("service/reset-password ok", "user", reset.Username)
}

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
