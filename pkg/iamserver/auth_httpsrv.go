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
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/hooto/httpsrv"
	"github.com/sysinner/innerstack/v2/pkg/inauth"
)

type UserAuth struct {
	*httpsrv.Controller
}

// SessionAction returns the current app configuration.
func (c UserAuth) SessionAction() {
	var (
		req UserAuthSessionRequest
		rsp UserAuthSessionResponse
	)
	defer c.RenderJson(&rsp)

	if err := c.Request.JsonDecode(&req); err != nil {
		// rsp.Status = inauth.NewServiceStatus("400", "Invalid request format")
		// return
	}

	if err := AppVerifier.Ping(); err != nil {
		rsp.Status = inauth.NewServiceStatus("500", err.Error())
		slog.Warn(err.Error())
		return
	}

	cfg := AppVerifier.Config()

	rsp.AppId = cfg.AppId
	rsp.AuthBaseURL = UrlJoinPath(cfg.BaseURL, "/")
	rsp.AuthSignInURL = UrlJoinPath(cfg.BaseURL,
		"/auth/sign-in") + "?app_id=" + cfg.AppId

	token, err := AppVerifier.Auth(c.Request.Request)
	if err == nil {
		rsp.AuthClaims = &token.AccessToken.Claims
		rsp.IdentityToken = token.IdentityToken
	} else {
		slog.Warn("user-auth/session: Auth failed", "error", err)

		// Fallback: parse the access token locally to get basic user info
		// even if the IAM session lookup fails
		if cookie, cerr := c.Request.Cookie(inauth.AppHttpHeaderKey); cerr == nil {
			if at, perr := inauth.ParseAccessToken(cookie.Value); perr == nil {
				rsp.AuthClaims = &at.Claims
			}
		}
	}

	// Let the host application inject its own business-logic fields
	// (e.g. ui_mgr_allow) into the response.
	if sessionResponseHook != nil {
		rsp.Extras = sessionResponseHook(AppVerifier.Session(c.Request.Request))
	}

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

// SignInAction redirects the user to the IAM sign-in page.
func (c UserAuth) SignInAction() {
	c.AutoRender = false

	if err := AppVerifier.Config().Valid(); err != nil {
		http.Error(c.Response.Out, "App not configured : "+err.Error(), http.StatusBadRequest)
		return
	}

	cfg := AppVerifier.Config()

	signInUrl := UrlJoinPath(cfg.BaseURL, "/auth/sign-in") + "?app_id=" + cfg.AppId

	cu := c.Request.URL.Query().Get("current_url")

	if cu == "" {
		if ck, err := c.Request.Cookie(inauth.AppHttpHeaderKey + "-current-url"); err == nil {
			cu = ck.Value
		}
	}

	if cu == "" {
		if ref := c.Request.Referer(); ref != "" {
			if u, err := url.Parse(ref); err == nil && isSameSite(u, c.Request.Request) {
				cu = ref
			}
		}
	}

	if cu != "" {
		slog.Info("callback-url : " + cu)
		http.SetCookie(c.Response.Out, &http.Cookie{
			Name:     inauth.AppHttpHeaderKey + "-current-url",
			Value:    cu,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(1 * time.Hour),
		})
	}

	http.Redirect(c.Response.Out, c.Request.Request, signInUrl, http.StatusFound)
}

// CallbackAction handles the IAM redirect with auth code.
func (c UserAuth) CallbackAction() {
	c.AutoRender = false
	code := c.Request.URL.Query().Get("code")
	if code == "" {
		http.Error(c.Response.Out, "missing code parameter", http.StatusBadRequest)
		return
	}

	if err := AppVerifier.Config().Valid(); err != nil {
		http.Error(c.Response.Out, "App not configured : "+err.Error(), http.StatusBadRequest)
		return
	}

	rs, err := ExchangeAuthCode(AppVerifier.Config(), code)
	if err != nil {
		slog.Error("user-auth/callback: code exchange failed", "error", err)
		http.Error(c.Response.Out, "code exchange failed", http.StatusInternalServerError)
		return
	}

	at, err := inauth.ParseAccessToken(rs.AccessToken)
	if err != nil {
		slog.Error("user-auth/callback: parse access-token failed", "error", err)
		http.Error(c.Response.Out, "parse access-token failed", http.StatusInternalServerError)
		return
	}

	http.SetCookie(c.Response.Out, &http.Cookie{
		Name:     inauth.AppHttpHeaderKey,
		Value:    rs.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(at.Claims.Exp, 0),
	})

	if urlCookie, err := c.Request.Cookie(
		inauth.AppHttpHeaderKey + "-current-url"); err == nil && urlCookie.Value != "" {
		http.Redirect(c.Response.Out, c.Request.Request, urlCookie.Value, http.StatusFound)
		deleteCookie(c.Response.Out, inauth.AppHttpHeaderKey+"-current-url")
		return
	}

	http.Redirect(c.Response.Out, c.Request.Request, "/", http.StatusFound)
}

// SignOutAction clears the session cookie and notifies IAM.
//
// For browser navigation requests, it redirects back to the Referer so the
// user lands on the page they signed out from. Otherwise it returns a JSON
// status payload (e.g. for XHR/fetch callers).
func (c UserAuth) SignOutAction() {

	deleteCookie(c.Response.Out, inauth.AppHttpHeaderKey)

	if isBrowserNavigation(c.Request.Request) {
		c.AutoRender = false
		target := "/"
		if ref := c.Request.Referer(); ref != "" {
			if u, err := url.Parse(ref); err == nil && isSameSite(u, c.Request.Request) {
				target = ref
			}
		}
		http.Redirect(c.Response.Out, c.Request.Request, target, http.StatusFound)
		return
	}

	var rsp struct {
		Status inauth.ServiceStatus `json:"status"`
	}
	defer c.RenderJson(&rsp)

	rsp.Status = inauth.NewServiceStatus("200", "ok")
}

func deleteCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
}
