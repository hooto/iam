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
	"sync"
	"time"

	"github.com/hooto/iam/v2/pkg/iamapi"
	"github.com/sysinner/incore/v2/pkg/inauth"
)

type Verifier interface {
	Setup(cfg *AppAuthConfig) error
	Config() *AppAuthConfig
	Ping() error
	Update(app *iamapi.AppInstance) error
	Auth(accessToken any) (*inauth.SessionToken, error)

	Session(accessToken any) UserSession
}

var AppVerifier Verifier = &verifier{
	sessions: make(map[string]*inauth.SessionToken),
}

type verifier struct {
	mu sync.RWMutex

	cfg *AppAuthConfig

	cfgFailCnt int
	cfgLastErr error
	cfgExpAt   time.Time

	sessions map[string]*inauth.SessionToken
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
		v.cfg.Endpoint, "/v2/open/app-auth/verify",
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

func (v *verifier) Auth(at any) (*inauth.SessionToken, error) {

	if err := v.cfg.Valid(); err != nil {
		return nil, err
	}

	accessToken := ""
	switch ats := at.(type) {
	case string:
		accessToken = ats
	case *http.Request:
		cookie, err := ats.Cookie(inauth.AppHttpHeaderKey)
		if err != nil {
			return nil, err
		}
		accessToken = cookie.Value
	default:
		return nil, errors.New("invalid access token type")
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
		v.cfg.Endpoint, "/v2/open/app-auth/session",
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
		v.cfg.Endpoint, "/v2/open/app-auth/update",
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

	accessToken := ""
	switch ats := at.(type) {
	case string:
		accessToken = ats
	case *http.Request:
		cookie, err := ats.Cookie(inauth.AppHttpHeaderKey)
		if err != nil {
			sess.authError = err
			return sess
		}
		accessToken = cookie.Value
	default:
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
		v.cfg.Endpoint, "/v2/open/app-auth/session",
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
