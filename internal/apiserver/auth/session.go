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

package auth

import (
	"log/slog"

	"github.com/hooto/httpsrv"
	"github.com/sysinner/innerstack/v2/pkg/inauth"

	"github.com/hooto/iam/v2/internal/data"
	"github.com/hooto/iam/v2/pkg/iamapi"
)

type SessionRequest struct {
	AccessToken string `json:"access_token,omitempty"`
}

type SessionResponse struct {
	Status        inauth.ServiceStatus  `json:"status"`
	AuthClaims    *inauth.AuthClaims    `json:"auth_claims,omitempty"`
	IdentityToken *inauth.IdentityToken `json:"identity_token,omitempty"`
}

// Session retrieves the current user session info including identity token.
func Session(ctx httpsrv.Ctx) error {

	var (
		req SessionRequest
		rsp = SessionResponse{}
	)
	defer ctx.JSON(&rsp)

	ctx.Request().JsonDecode(&req)

	if req.AccessToken == "" {
		// fallback to http-only cookie
		cookie, err := ctx.Request().Cookie(inauth.AppHttpHeaderKey)
		if err != nil || cookie.Value == "" {
			rsp.Status = inauth.NewServiceStatus("401", "access token not found")
			return nil
		}
		req.AccessToken = cookie.Value
	}

	token, err := inauth.ParseAccessToken(req.AccessToken)
	if err != nil {
		slog.Info("service/session-token fail", "error", err.Error(), "access_token", token)
		rsp.Status = inauth.NewServiceStatus("401", "invalid access token : "+err.Error())
		return nil
	}

	// verify signature
	if _, err := token.Verify(data.KeyMgr); err != nil {
		rsp.Status = inauth.NewServiceStatus("401", "invalid access token : "+err.Error())
		return nil
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
			rsp.AuthClaims = &st.AccessToken.Claims
			rsp.IdentityToken = st.IdentityToken
		}
	}

	return nil
}
