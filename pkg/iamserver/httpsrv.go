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

package iamserver

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/hooto/httpsrv"
	"github.com/sysinner/incore/v2/pkg/inauth"
)

type UserAuth struct {
	*httpsrv.Controller
}

type UserAuthInfoRequest struct {
	CurrentUrl string `json:"current_url"`
}

type UserAuthInfoResponse struct {
	Status inauth.ServiceStatus `json:"status"`

	AppId        string             `json:"app_id,omitempty"`
	AuthEndpoint string             `json:"auth_endpoint,omitempty"`
	AuthClaims   *inauth.AuthClaims `json:"auth_claims,omitempty"`
	// IdentityToken *inauth.IdentityToken `json:"identity_token,omitempty"`
}

// InfoAction returns the current app configuration.
func (c UserAuth) InfoAction() {
	var (
		req UserAuthInfoRequest
		rsp UserAuthInfoResponse
	)
	defer c.RenderJson(&rsp)

	if err := c.Request.JsonDecode(&req); err != nil {
		// rsp.Status = inauth.NewServiceStatus("400", "Invalid request format")
		// return
	}

	if err := AppConfig.Verify(); err != nil {
		rsp.Status = inauth.NewServiceStatus("401", err.Error())
		slog.Info("iam config not initialized")
		return
	}

	rsp.AppId = AppConfig.AppId
	rsp.AuthEndpoint = urlJoinPath(AppConfig.Endpoint,
		"/auth/sign-in") + "?app_id=" + AppConfig.AppId

	// fallback to http-only cookie
	cookie, err := c.Request.Cookie(inauth.AppHttpHeaderKey)
	if err != nil || cookie.Value == "" {
		rsp.Status = inauth.NewServiceStatus("401", "access token not found")
		return
	}

	token, err := inauth.ParseAccessToken(cookie.Value)
	if err != nil {
		slog.Info("service/session-token fail", "error", err.Error(), "access_token", token)
		rsp.Status = inauth.NewServiceStatus("401", "invalid access token : "+err.Error())
		return
	}

	rsp.AuthClaims = &token.Claims

	if req.CurrentUrl != "" {
		http.SetCookie(c.Response.Out, &http.Cookie{
			Name:     inauth.AppHttpHeaderKey + "-current-url",
			Value:    req.CurrentUrl,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(1 * time.Hour),
		})
	}

	rsp.Status = inauth.NewServiceStatus("200", "ok")
}

// ExchangeAuthCode calls IAM to exchange an auth code for access_token.
type AuthCodeResult struct {
	AccessToken   string
	IdentityToken *inauth.IdentityToken
}

// CallbackAction handles the IAM redirect with auth code.
func (c UserAuth) CallbackAction() {
	c.AutoRender = false
	code := c.Request.URL.Query().Get("code")
	if code == "" {
		http.Error(c.Response.Out, "missing code parameter", http.StatusBadRequest)
		return
	}

	if err := AppConfig.Verify(); err != nil {
		http.Error(c.Response.Out, "App not configured : "+err.Error(), http.StatusBadRequest)
		return
	}

	rs, err := exchangeAuthCode(AppConfig, code)
	if err != nil {
		slog.Error("user-auth/callback: code exchange failed", "error", err)
		http.Error(c.Response.Out, "code exchange failed", http.StatusInternalServerError)
		return
	}

	http.SetCookie(c.Response.Out, &http.Cookie{
		Name:     inauth.AppHttpHeaderKey,
		Value:    rs.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(10 * 24 * time.Hour),
	})

	if urlCookie, err := c.Request.Cookie(
		inauth.AppHttpHeaderKey + "-current-url"); err == nil && urlCookie.Value != "" {
		http.Redirect(c.Response.Out, c.Request.Request, urlCookie.Value, http.StatusFound)
		return
	}

	http.Redirect(c.Response.Out, c.Request.Request, "/", http.StatusFound)
}

func exchangeAuthCode(aac *AppAuthConfig, code string) (*AuthCodeResult, error) {

	var rsp struct {
		Status        inauth.ServiceStatus  `json:"status"`
		AccessToken   string                `json:"access_token,omitempty"`
		IdentityToken *inauth.IdentityToken `json:"identity_token,omitempty"`
	}

	ac, err := aac.NewAppCredential()
	if err != nil {
		return nil, err
	}
	at := ac.AuthToken()

	// slog.Info("exchangeAuthCode: request", "endpoint",
	// 	urlJoinPath(AppConfig.Endpoint, "/v2/app-auth/auth-token-exchange?code="+code))

	if err := iamPost(
		AppConfig.Endpoint, "/v2/open/app-auth/token-exchange",
		at,
		map[string]string{
			"code": code,
		},
		&rsp,
	); err != nil {
		slog.Error("exchangeAuthCode: request failed", "error", err)
		return nil, err
	}
	js, _ := json.Marshal(rsp)
	slog.Info("exchangeAuthCode: response", "body", string(js))
	return &AuthCodeResult{
		AccessToken:   rsp.AccessToken,
		IdentityToken: rsp.IdentityToken,
	}, nil
}

// SignOutAction clears the session cookie and notifies IAM.
func (c UserAuth) SignOutAction() {
	var rsp struct {
		Status inauth.ServiceStatus `json:"status"`
	}
	defer c.RenderJson(&rsp)

	http.SetCookie(c.Response.Out, &http.Cookie{
		Name:     inauth.AppHttpHeaderKey,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	rsp.Status = inauth.NewServiceStatus("200", "ok")
}
