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
	"net/http"

	"github.com/hooto/httpsrv"
	"github.com/sysinner/incore/v2/pkg/inauth"

	"github.com/hooto/iam/v2/internal/data"
	"github.com/hooto/iam/v2/pkg/iamapi"
)

type SignOutRequest struct {
	AccessToken string `json:"access_token,omitempty"`
}

type ServiceStatusResponse struct {
	Status inauth.ServiceStatus `json:"status"`
}

// SignOut invalidates the user session and clears the auth cookie.
func SignOut(ctx *httpsrv.Context) error {

	var (
		req SignOutRequest
		rsp ServiceStatusResponse
	)
	defer ctx.RenderJson(&rsp)

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
		rsp.Status = inauth.NewServiceStatus("401", "invalid access token")
		return nil
	}

	data.Data.NewDeleter(iamapi.NsUserSession(token.Claims.Jti, uint32(token.Claims.Exp))).
		Exec()

	http.SetCookie(ctx.Response().Out, &http.Cookie{
		Name:   inauth.AppHttpHeaderKey,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	rsp.Status = inauth.NewServiceStatus("200", "ok")
	return nil
}
