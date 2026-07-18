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

package authfiber

import (
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/sysinner/innerstack/v2/pkg/inauth"

	"github.com/hooto/iam/v2/pkg/iamserver"
)

// RegisterAuthRoutes registers the IAM user-auth routes on a fiber router. The
// caller mounts the router at the controller prefix ("/user-auth"), reproducing
// the routes the httpsrv UserAuth controller exposed:
//
//	{prefix}/session, {prefix}/sign-in, {prefix}/callback, {prefix}/sign-out
//
// Routes accept any HTTP method (httpsrv dispatch is method-agnostic).
func RegisterAuthRoutes(router fiber.Router) {
	router.All("/session", userAuthSession)
	router.All("/sign-in", userAuthSignIn)
	router.All("/callback", userAuthCallback)
	router.All("/sign-out", userAuthSignOut)
}

var sessionResponseHook iamserver.SessionResponseHook

// SetSessionResponseHook registers the application-specific session
// response injector used by the session route.
func SetSessionResponseHook(fn iamserver.SessionResponseHook) {
	sessionResponseHook = fn
}

// accessTokenCookie is the IAM access-token cookie name.
func accessTokenCookie() string { return inauth.AppHttpHeaderKey }

// currentURLCookie is the "where to return after sign-in" cookie name.
func currentURLCookie() string { return inauth.AppHttpHeaderKey + "-current-url" }

func writeAuthJSON(c fiber.Ctx, data any) error {
	c.Set("Access-Control-Allow-Origin", "*")
	return c.JSON(data)
}

func userAuthSession(c fiber.Ctx) error {
	var (
		req iamserver.UserAuthSessionRequest
		rsp iamserver.UserAuthSessionResponse
	)

	_ = c.Bind().JSON(&req) // decode errors are tolerated, matching prior behavior

	if err := iamserver.AppVerifier.Ping(); err != nil {
		rsp.Status = inauth.NewServiceStatus("500", err.Error())
		slog.Warn(err.Error())
		return writeAuthJSON(c, &rsp)
	}

	cfg := iamserver.AppVerifier.Config()

	rsp.AppId = cfg.AppId
	rsp.AuthBaseURL = iamserver.UrlJoinPath(cfg.BaseURL, "/")
	rsp.AuthSignInURL = iamserver.UrlJoinPath(cfg.BaseURL, "/auth/sign-in") + "?app_id=" + cfg.AppId

	token, err := iamserver.AppVerifier.Auth(c.Cookies(accessTokenCookie()))
	if err == nil {
		rsp.AuthClaims = &token.AccessToken.Claims
		rsp.IdentityToken = token.IdentityToken
	} else {
		slog.Warn("user-auth/session: Auth failed", "error", err)
		// Fallback: parse the access token locally for basic user info even when
		// the IAM session lookup fails.
		if cv := c.Cookies(accessTokenCookie()); cv != "" {
			if at, perr := inauth.ParseAccessToken(cv); perr == nil {
				rsp.AuthClaims = &at.Claims
			}
		}
	}

	if sessionResponseHook != nil {
		rsp.Extras = sessionResponseHook(iamserver.AppVerifier.Session(c.Cookies(accessTokenCookie())))
	}

	if req.CurrentUrl != "" {
		c.Cookie(&fiber.Cookie{
			Name:     currentURLCookie(),
			Value:    req.CurrentUrl,
			Path:     "/",
			HTTPOnly: true,
			Expires:  time.Now().Add(1 * time.Hour),
		})
	}

	rsp.Status = inauth.NewServiceStatus("200", "ok")
	return writeAuthJSON(c, &rsp)
}

func userAuthSignIn(c fiber.Ctx) error {
	if err := iamserver.AppVerifier.Config().Valid(); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("App not configured : " + err.Error())
	}

	cfg := iamserver.AppVerifier.Config()
	signInURL := iamserver.UrlJoinPath(cfg.BaseURL, "/auth/sign-in") + "?app_id=" + cfg.AppId

	cu := c.Query("current_url")
	if cu == "" {
		cu = c.Cookies(currentURLCookie())
	}
	if cu == "" {
		if ref := c.Get("Referer"); ref != "" {
			if u, err := url.Parse(ref); err == nil && sameSiteURL(u, c.Host()) {
				cu = ref
			}
		}
	}

	if cu != "" {
		slog.Info("callback-url : " + cu)
		c.Cookie(&fiber.Cookie{
			Name:     currentURLCookie(),
			Value:    cu,
			Path:     "/",
			HTTPOnly: true,
			Expires:  time.Now().Add(1 * time.Hour),
		})
	}

	return c.Redirect().Status(fiber.StatusFound).To(signInURL)
}

func userAuthCallback(c fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).SendString("missing code parameter")
	}

	if err := iamserver.AppVerifier.Config().Valid(); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("App not configured : " + err.Error())
	}

	rs, err := iamserver.ExchangeAuthCode(iamserver.AppVerifier.Config(), code)
	if err != nil {
		slog.Error("user-auth/callback: code exchange failed", "error", err)
		return c.Status(fiber.StatusInternalServerError).SendString("code exchange failed")
	}

	at, err := inauth.ParseAccessToken(rs.AccessToken)
	if err != nil {
		slog.Error("user-auth/callback: parse access-token failed", "error", err)
		return c.Status(fiber.StatusInternalServerError).SendString("parse access-token failed")
	}

	c.Cookie(&fiber.Cookie{
		Name:     accessTokenCookie(),
		Value:    rs.AccessToken,
		Path:     "/",
		HTTPOnly: true,
		Expires:  time.Unix(at.Claims.Exp, 0),
	})

	if urlCookie := c.Cookies(currentURLCookie()); urlCookie != "" {
		deleteFiberCookie(c, currentURLCookie())
		return c.Redirect().Status(fiber.StatusFound).To(urlCookie)
	}

	return c.Redirect().Status(fiber.StatusFound).To("/")
}

func userAuthSignOut(c fiber.Ctx) error {
	deleteFiberCookie(c, accessTokenCookie())

	if browserNavigation(c) {
		target := "/"
		if ref := c.Get("Referer"); ref != "" {
			if u, err := url.Parse(ref); err == nil && sameSiteURL(u, c.Host()) {
				target = ref
			}
		}
		return c.Redirect().Status(fiber.StatusFound).To(target)
	}

	var rsp struct {
		Status inauth.ServiceStatus `json:"status"`
	}
	rsp.Status = inauth.NewServiceStatus("200", "ok")
	return writeAuthJSON(c, &rsp)
}

// deleteFiberCookie expires a cookie immediately (mirrors deleteCookie for
// http.ResponseWriter).
func deleteFiberCookie(c fiber.Ctx, name string) {
	c.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		HTTPOnly: true,
		MaxAge:   -1,
	})
}

// sameSiteURL reports whether u belongs to the same host. Relative URLs (empty
// host) are treated as same-origin. fiber equivalent of isSameSite.
func sameSiteURL(u *url.URL, host string) bool {
	if u == nil {
		return false
	}
	if u.Host == "" {
		return true
	}
	return strings.EqualFold(u.Host, host)
}

// browserNavigation reports whether the request looks like a full-page browser
// navigation rather than an AJAX/XHR call. fiber equivalent of
// isBrowserNavigation.
func browserNavigation(c fiber.Ctx) bool {
	if strings.EqualFold(c.Get("X-Requested-With"), "XMLHttpRequest") {
		return false
	}
	if c.Get("Sec-Fetch-Mode") == "navigate" {
		return true
	}
	if c.Get("Upgrade-Insecure-Requests") == "1" &&
		strings.Contains(strings.ToLower(c.Get("Accept")), "text/html") {
		return true
	}
	return false
}
