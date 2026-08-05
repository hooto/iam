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

	"github.com/hooto/iam/v2/internal/util"
	"github.com/sysinner/innerstack/v2/pkg/inauth"
)

// UserSession defines the interface for managing user authentication
// and authorization within an IAM-protected application.
type UserSession interface {
	// CheckServer validates whether the IAM server connection is properly
	// configured and reachable. Returns a non-nil error if the server
	// is misconfigured or unavailable.
	CheckServer() error

	// RequireAuth verifies that the current user is authenticated.
	// If the user is not logged in, it returns the redirect URL for the
	// sign-in page and a non-nil error. Callers should redirect the user
	// to the returned URL to initiate authentication.
	RequireAuth() (redirectURL string, err error)

	// Allow checks whether the specified user is authorized to perform
	// the given operations identified by permission strings.
	// Returns true if all permissions are granted, false otherwise.
	Allow(username string, permissions ...string) bool

	// Profile retrieves the public profile information of the currently
	// authenticated user. Returns an error if no user is authenticated
	// or the profile cannot be fetched.
	Profile() (*UserProfile, error)
}

type UserProfile struct {
	Username string `json:"username,omitempty" toml:"username,omitempty"`
	PhotoURL string `json:"photo_url,omitempty" toml:"photo_url,omitempty"`
}

type userSession struct {
	authError error

	// accessKey marks a session resolved from a CLI access-key token via
	// introspection. Such sessions must NOT take the "own user, allow all"
	// shortcut in Allow, so the key's scopes stay enforced.
	accessKey bool

	cfg *AppAuthConfig

	AuthClaims    *inauth.AuthClaims    `json:"auth_claims,omitempty"`
	IdentityToken *inauth.IdentityToken `json:"identity_token,omitempty"`
}

func (it *userSession) CheckServer() error {
	if it.cfg == nil {
		return errors.New("IAM not configured")
	}
	return it.cfg.Valid()
}

func (it *userSession) RequireAuth() (string, error) {

	if it.AuthClaims == nil || it.IdentityToken == nil {
		if it.authError != nil {
			return "", it.authError
		}
		return "", errors.New("user not authenticated")
	}

	u := UrlJoinPath(it.cfg.BaseURL,
		"/auth/sign-in") + "?app_id=" + it.cfg.AppId

	return u, nil
}

func (it *userSession) Allow(username string, permissions ...string) bool {
	if it.AuthClaims == nil || it.IdentityToken == nil {
		return false
	}

	// browser sessions: the resource owner gets full access to their own data.
	// access-key sessions skip this shortcut so the key's scopes are enforced.
	if !it.accessKey && it.AuthClaims.Sub == username {
		return true
	}

	return util.Contains(it.IdentityToken.Permissions, permissions)
}

func (it *userSession) Profile() (*UserProfile, error) {
	if it.AuthClaims == nil || it.IdentityToken == nil {
		return nil, errors.New("user not authenticated")
	}

	photoURL := UrlJoinPath(it.cfg.BaseURL,
		"/auth/photo/"+it.AuthClaims.Sub)

	return &UserProfile{
		Username: it.AuthClaims.Sub,
		PhotoURL: photoURL,
	}, nil
}
