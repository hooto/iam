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
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/hooto/iam/v2/pkg/iamapi"
	"github.com/sysinner/innerstack/v2/pkg/inauth"
)

type Verifier interface {
	Setup(cfg *AppAuthConfig) error
	Config() *AppAuthConfig
	Ping() error
	Update(app *iamapi.AppInstance) error
	Auth(accessToken any) (*inauth.SessionToken, error)

	Session(accessToken any) UserSession

	// Resolve is the single entry point that serves both browser sessions
	// (cookie/sub token, delegated to Session) and CLI access-key tokens
	// (Authorization: Bearer, resolved via /open/app-auth/introspect). The
	// app's request handler can call Resolve for every incoming request
	// without knowing which credential type the caller used.
	Resolve(accessToken any) UserSession
}

var AppVerifier Verifier = &verifier{
	sessions:   make(map[string]*inauth.SessionToken),
	akSessions: make(map[string]*akCacheEntry),
}

// accessKeyCacheTTL bounds how long an introspected access-key identity is
// trusted at the app without re-hitting IAM. It also bounds revocation
// latency: a key disabled at IAM stays usable at the app for at most this
// duration. Aligned with the 60s lifetime of an access-key-signed token.
const accessKeyCacheTTL = 60 * time.Second

// akCacheEntry caches a resolved access-key session, keyed by access-key Id.
type akCacheEntry struct {
	sess  *userSession
	expAt time.Time
}

type verifier struct {
	mu sync.RWMutex

	cfg *AppAuthConfig

	cfgFailCnt int
	cfgLastErr error
	cfgExpAt   time.Time

	sessions   map[string]*inauth.SessionToken
	akSessions map[string]*akCacheEntry
}

func (v *verifier) Setup(cfg *AppAuthConfig) error {
	v.mu.Lock()
	defer v.mu.Unlock()
	if cfg == nil {
		return errors.New("config is nil")
	}
	v.cfg = cfg
	return nil
}

func (v *verifier) Config() *AppAuthConfig {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.cfg
}

func (v *verifier) Ping() error {
	v.mu.RLock()
	defer v.mu.RUnlock()
	if err := v.cfg.Valid(); err != nil {
		return err
	}

	if time.Now().Before(v.cfgExpAt) {
		return v.cfgLastErr
	}

	if v.cfgLastErr = v.ping(); v.cfgLastErr == nil {
		v.cfgFailCnt = 0
		v.cfgExpAt = time.Now().Add(300 * time.Second)
	} else {
		v.cfgFailCnt++
		sec := 1 << (v.cfgFailCnt - 1)
		if sec > 60 || v.cfgFailCnt > 6 {
			sec = 60
		}
		v.cfgExpAt = time.Now().Add(time.Duration(sec) * time.Second)
	}
	return v.cfgLastErr
}

func (v *verifier) ping() error {
	ac, err := v.cfg.NewAppCredential()
	if err != nil {
		return err
	}
	at := ac.AuthToken()

	var rsp struct {
		Status inauth.ServiceStatus `json:"status"`
	}
	if err := iamPost(
		v.cfg.BaseURL, "/v2/open/app-auth/verify",
		at,
		map[string]string{
			"app_id": v.cfg.AppId,
		},
		&rsp,
	); err != nil {
		return err
	}

	if rsp.Status.Code != "200" {
		return errors.New("Verification failed: " + rsp.Status.Message)
	}

	return nil
}

// extractAccessToken pulls the raw access-token string from a string value or
// an *http.Request. When preferBearer is set, an "Authorization: Bearer"
// header wins over the session cookie (Resolve, for CLI access keys); otherwise
// only the cookie is consulted (Auth/Session, browser flow).
func (v *verifier) extractAccessToken(at any, preferBearer bool) (string, error) {
	switch ats := at.(type) {
	case string:
		return ats, nil
	case *http.Request:
		if preferBearer {
			if h := ats.Header.Get("Authorization"); h != "" {
				if parts := strings.SplitN(h, " ", 2); len(parts) == 2 && parts[0] == "Bearer" {
					return parts[1], nil
				}
			}
		}
		cookie, err := ats.Cookie(inauth.AppHttpHeaderKey)
		if err != nil {
			return "", err
		}
		return cookie.Value, nil
	default:
		return "", errors.New("invalid access token type")
	}
}

func (v *verifier) Auth(at any) (*inauth.SessionToken, error) {

	if err := v.cfg.Valid(); err != nil {
		return nil, err
	}

	accessToken, err := v.extractAccessToken(at, false)
	if err != nil {
		return nil, err
	}

	token, err := inauth.ParseAccessToken(accessToken)
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, errors.New("invalid access token")
	}

	v.mu.RLock()
	session, ok := v.sessions[token.Claims.Jti]
	v.mu.RUnlock()
	if ok {
		return session, nil
	}

	ac, err := v.cfg.NewAppCredential()
	if err != nil {
		return nil, err
	}
	aat := ac.AuthToken()

	var rsp struct {
		Status        inauth.ServiceStatus  `json:"status"`
		IdentityToken *inauth.IdentityToken `json:"identity_token,omitempty"`
	}

	err = iamPost(
		v.cfg.BaseURL, "/v2/open/app-auth/session",
		aat,
		map[string]string{
			"app_id":       v.cfg.AppId,
			"access_token": accessToken,
		},
		&rsp,
	)
	if err != nil {
		return nil, err
	}

	session = &inauth.SessionToken{
		AccessToken:   token,
		IdentityToken: rsp.IdentityToken,
	}
	v.mu.Lock()
	v.sessions[token.Claims.Jti] = session
	v.mu.Unlock()

	return session, nil
}

func (v *verifier) Update(app *iamapi.AppInstance) error {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.cfg == nil {
		return errors.New("config is not set")
	}
	if app == nil {
		return errors.New("app instance is nil")
	}
	if app.ID != v.cfg.AppId {
		return errors.New("app instance ID does not match config AppId")
	}

	ac, err := v.cfg.NewAppCredential()
	if err != nil {
		return err
	}
	aat := ac.AuthToken()

	var rsp struct {
		Status inauth.ServiceStatus `json:"status"`
	}

	if err = iamPost(
		v.cfg.BaseURL, "/v2/open/app-auth/update",
		aat,
		map[string]any{
			"app_id":      v.cfg.AppId,
			"version":     app.Version,
			"permissions": app.Permissions,
		},
		&rsp,
	); err != nil {
		return err
	}

	return nil
}

func (v *verifier) Session(at any) UserSession {

	sess := &userSession{
		cfg: v.cfg,
	}

	if err := v.cfg.Valid(); err != nil {
		return sess
	}

	accessToken, err := v.extractAccessToken(at, false)
	if err != nil {
		sess.authError = err
		return sess
	}

	token, err := inauth.ParseAccessToken(accessToken)
	if err != nil {
		sess.authError = err
		return sess
	}
	if token == nil {
		sess.authError = errors.New("invalid access token")
		return sess
	}

	v.mu.RLock()
	session, ok := v.sessions[token.Claims.Jti]
	v.mu.RUnlock()
	if ok {
		sess.AuthClaims = &session.AccessToken.Claims
		sess.IdentityToken = session.IdentityToken
		return sess
	}

	ac, err := v.cfg.NewAppCredential()
	if err != nil {
		sess.authError = err
		return sess
	}
	aat := ac.AuthToken()

	var rsp struct {
		Status        inauth.ServiceStatus  `json:"status"`
		IdentityToken *inauth.IdentityToken `json:"identity_token,omitempty"`
	}

	err = iamPost(
		v.cfg.BaseURL, "/v2/open/app-auth/session",
		aat,
		map[string]string{
			"app_id":       v.cfg.AppId,
			"access_token": accessToken,
		},
		&rsp,
	)
	if err != nil {
		sess.authError = err
		return sess
	}

	session = &inauth.SessionToken{
		AccessToken:   token,
		IdentityToken: rsp.IdentityToken,
	}
	v.mu.Lock()
	v.sessions[token.Claims.Jti] = session
	v.mu.Unlock()

	sess.AuthClaims = &session.AccessToken.Claims
	sess.IdentityToken = session.IdentityToken
	return sess
}

// Resolve authenticates an incoming request using either a browser session
// (cookie) or a CLI access-key token (Authorization: Bearer). Browser session
// tokens carry a Sub claim and are delegated to Session. Access-key tokens
// (no Sub) are verified at IAM via /open/app-auth/introspect; the resolved
// identity is cached by access-key Id (Kid) for accessKeyCacheTTL so repeated
// requests within the window do not re-hit IAM.
func (v *verifier) Resolve(at any) UserSession {

	sess := &userSession{cfg: v.cfg}

	if err := v.cfg.Valid(); err != nil {
		return sess
	}

	accessToken, err := v.extractAccessToken(at, true)
	if err != nil {
		sess.authError = err
		return sess
	}

	token, err := inauth.ParseAccessToken(accessToken)
	if err != nil {
		sess.authError = err
		return sess
	}
	if token == nil {
		sess.authError = errors.New("invalid access token")
		return sess
	}

	// browser session token -> reuse the existing session resolution.
	// Pass the extracted token string so a Sub token carried in the
	// Authorization header is resolved just like a cookie-based one.
	if token.Claims.Sub != "" {
		return v.Session(accessToken)
	}

	// access-key token: reject locally if already expired
	if token.IsExpired() {
		sess.authError = errors.New("access token expired")
		return sess
	}

	// cache lookup by access-key Id (Kid); a fresh hit skips the IAM call
	kid := token.Header.Kid
	if kid != "" {
		v.mu.RLock()
		entry, ok := v.akSessions[kid]
		v.mu.RUnlock()
		if ok && time.Now().Before(entry.expAt) {
			return entry.sess
		}
	}

	ac, err := v.cfg.NewAppCredential()
	if err != nil {
		sess.authError = err
		return sess
	}

	var rsp struct {
		Status        inauth.ServiceStatus  `json:"status"`
		User          string                `json:"user,omitempty"`
		IdentityToken *inauth.IdentityToken `json:"identity_token,omitempty"`
	}

	if err = iamPost(
		v.cfg.BaseURL, "/v2/open/app-auth/introspect",
		ac.AuthToken(),
		map[string]string{"access_token": accessToken},
		&rsp,
	); err != nil {
		sess.authError = err
		return sess
	}

	if rsp.Status.Code != "200" {
		if rsp.Status.Message != "" {
			sess.authError = errors.New(rsp.Status.Message)
		} else {
			sess.authError = errors.New("server error")
		}
		return sess
	}

	resolved := &userSession{
		cfg:           v.cfg,
		accessKey:     true,
		AuthClaims:    &inauth.AuthClaims{Sub: rsp.User, Exp: token.Claims.Exp},
		IdentityToken: rsp.IdentityToken,
	}

	if kid != "" {
		v.mu.Lock()
		v.akSessions[kid] = &akCacheEntry{
			sess:  resolved,
			expAt: time.Now().Add(accessKeyCacheTTL),
		}
		v.mu.Unlock()
	}

	return resolved
}
