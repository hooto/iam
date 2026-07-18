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

	"github.com/sysinner/innerstack/v2/pkg/inauth"
)

type UserAuthSessionRequest struct {
	CurrentUrl string `json:"current_url"`
}

type UserAuthSessionResponse struct {
	Status inauth.ServiceStatus `json:"status"`

	AppId       string `json:"app_id,omitempty"`
	AuthBaseURL string `json:"auth_base_url,omitempty"`

	AuthSignInURL string `json:"auth_sign_in_url,omitempty"`

	AuthClaims    *inauth.AuthClaims    `json:"auth_claims,omitempty"`
	IdentityToken *inauth.IdentityToken `json:"identity_token,omitempty"`

	// Extras carries application-specific fields injected by the
	// registered SessionResponseHook. It is namespaced under the
	// "extras" key so the host app can extend /user-auth/session
	// without colliding with IAM's own fields or requiring a custom
	// marshaler.
	Extras map[string]any `json:"extras,omitempty"`
}

// SessionResponseHook lets the host application extend the
// /user-auth/session JSON response with its own business-logic fields
// (e.g. ui_mgr_allow). It receives the resolved user session (whose
// Allow/Profile methods are usable; AuthClaims/IdentityToken are nil
// when the user is not authenticated or IAM is unreachable). The
// returned map is flattened into the response as top-level keys.
// Return nil to add nothing.
type SessionResponseHook func(sess UserSession) map[string]any

var sessionResponseHook SessionResponseHook

// SetSessionResponseHook registers the application-specific session
// response injector used by the session route.
func SetSessionResponseHook(fn SessionResponseHook) {
	sessionResponseHook = fn
}

// ExchangeAuthCode calls IAM to exchange an auth code for access_token.
type AuthCodeResult struct {
	AccessToken   string
	IdentityToken *inauth.IdentityToken
}

func ExchangeAuthCode(aac *AppAuthConfig, code string) (*AuthCodeResult, error) {

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

	if err := iamPost(
		aac.BaseURL, "/v2/open/app-auth/token-exchange",
		at,
		map[string]string{
			"code": code,
		},
		&rsp,
	); err != nil {
		slog.Error("ExchangeAuthCode: request failed", "error", err)
		return nil, err
	}
	slog.Info("ExchangeAuthCode: response", "body", rsp)
	return &AuthCodeResult{
		AccessToken:   rsp.AccessToken,
		IdentityToken: rsp.IdentityToken,
	}, nil
}
