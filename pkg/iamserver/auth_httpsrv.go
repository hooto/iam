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

	"github.com/hooto/httpsrv/v2"
	"github.com/sysinner/innerstack/v2/pkg/inauth"
)

func RegisterRoutes(r httpsrv.Router) {
	r.All("/user-auth/session", UserAuthSession)
	r.All("/user-auth/sign-in", UserAuthSignIn)
	r.All("/user-auth/callback", UserAuthCallback)
	r.All("/user-auth/sign-out", UserAuthSignOut)
}

// UserAuthSession returns the current app configuration.
func UserAuthSession(ctx httpsrv.Ctx) error {
	var (
		req UserAuthSessionRequest
		rsp UserAuthSessionResponse
	)
	defer ctx.JSON(&rsp)

	_ = ctx.Bind(&req)

	if err := AppVerifier.Ping(); err != nil {
		rsp.Status = inauth.NewServiceStatus("500", err.Error())
		slog.Warn(err.Error())
		return nil
	}

	cfg := AppVerifier.Config()

	rsp.AppId = cfg.AppId
	rsp.AuthBaseURL = UrlJoinPath(cfg.BaseURL, "/")
	rsp.AuthSignInURL = UrlJoinPath(cfg.BaseURL,
		"/auth/sign-in") + "?app_id=" + cfg.AppId

	token, err := AppVerifier.Auth(ctx.Request())
	if err == nil {
		rsp.AuthClaims = &token.AccessToken.Claims
		rsp.IdentityToken = token.IdentityToken
	} else {
		slog.Warn("user-auth/session: Auth failed", "error", err)

		// Fallback: parse the access token locally to get basic user info
		// even if the IAM session lookup fails
		if cookie, cerr := ctx.Request().Cookie(inauth.AppHttpHeaderKey); cerr == nil {
			if at, perr := inauth.ParseAccessToken(cookie.Value); perr == nil {
				rsp.AuthClaims = &at.Claims
			}
		}
	}

	// Let the host application inject its own business-logic fields
	// (e.g. ui_mgr_allow) into the response.
	if sessionResponseHook != nil {
		rsp.Extras = sessionResponseHook(AppVerifier.Session(ctx.Request()))
	}

	if req.CurrentUrl != "" {
		http.SetCookie(ctx.Response(), &http.Cookie{
			Name:     inauth.AppHttpHeaderKey + "-current-url",
			Value:    req.CurrentUrl,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(1 * time.Hour),
		})
	}

	rsp.Status = inauth.NewServiceStatus("200", "ok")
	return nil
}

// UserAuthSignIn redirects the user to the IAM sign-in page.
func UserAuthSignIn(ctx httpsrv.Ctx) error {

	if err := AppVerifier.Config().Valid(); err != nil {
		http.Error(ctx.Response(), "App not configured : "+err.Error(), http.StatusBadRequest)
		return nil
	}

	cfg := AppVerifier.Config()

	signInUrl := UrlJoinPath(cfg.BaseURL, "/auth/sign-in") + "?app_id=" + cfg.AppId

	cu := ctx.Request().URL.Query().Get("current_url")

	if cu == "" {
		if ck, err := ctx.Request().Cookie(inauth.AppHttpHeaderKey + "-current-url"); err == nil {
			cu = ck.Value
		}
	}

	if cu == "" {
		if ref := ctx.Request().Referer(); ref != "" {
			if u, err := url.Parse(ref); err == nil && isSameSite(u, ctx.Request()) {
				cu = ref
			}
		}
	}

	if cu != "" {
		slog.Info("callback-url : " + cu)
		http.SetCookie(ctx.Response(), &http.Cookie{
			Name:     inauth.AppHttpHeaderKey + "-current-url",
			Value:    cu,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(1 * time.Hour),
		})
	}

	return ctx.Redirect(http.StatusFound, signInUrl)
}

// UserAuthCallback handles the IAM redirect with auth code.
func UserAuthCallback(ctx httpsrv.Ctx) error {

	code := ctx.Request().URL.Query().Get("code")
	if code == "" {
		http.Error(ctx.Response(), "missing code parameter", http.StatusBadRequest)
		return nil
	}

	if err := AppVerifier.Config().Valid(); err != nil {
		http.Error(ctx.Response(), "App not configured : "+err.Error(), http.StatusBadRequest)
		return nil
	}

	rs, err := ExchangeAuthCode(AppVerifier.Config(), code)
	if err != nil {
		slog.Error("user-auth/callback: code exchange failed", "error", err)
		http.Error(ctx.Response(), "code exchange failed", http.StatusInternalServerError)
		return nil
	}

	at, err := inauth.ParseAccessToken(rs.AccessToken)
	if err != nil {
		slog.Error("user-auth/callback: parse access-token failed", "error", err)
		http.Error(ctx.Response(), "parse access-token failed", http.StatusInternalServerError)
		return nil
	}

	http.SetCookie(ctx.Response(), &http.Cookie{
		Name:     inauth.AppHttpHeaderKey,
		Value:    rs.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(at.Claims.Exp, 0),
	})

	if urlCookie, err := ctx.Request().Cookie(
		inauth.AppHttpHeaderKey + "-current-url"); err == nil && urlCookie.Value != "" {
		deleteCookie(ctx.Response(), inauth.AppHttpHeaderKey+"-current-url")
		return ctx.Redirect(http.StatusFound, urlCookie.Value)
	}

	return ctx.Redirect(http.StatusFound, "/")
}

// UserAuthSignOut clears the session cookie and notifies IAM.
//
// For browser navigation requests, it redirects back to the Referer so the
// user lands on the page they signed out from. Otherwise it returns a JSON
// status payload (e.g. for XHR/fetch callers).
func UserAuthSignOut(ctx httpsrv.Ctx) error {

	deleteCookie(ctx.Response(), inauth.AppHttpHeaderKey)

	if isBrowserNavigation(ctx.Request()) {
		target := "/"
		if ref := ctx.Request().Referer(); ref != "" {
			if u, err := url.Parse(ref); err == nil && isSameSite(u, ctx.Request()) {
				target = ref
			}
		}
		return ctx.Redirect(http.StatusFound, target)
	}

	var rsp struct {
		Status inauth.ServiceStatus `json:"status"`
	}
	defer ctx.JSON(&rsp)

	rsp.Status = inauth.NewServiceStatus("200", "ok")
	return nil
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
