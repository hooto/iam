// Copyright 2014 Eryx <evorui аt gmаil dοt cοm>, All rights reserved.
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

package apiserver

import (
	"strings"

	"github.com/hooto/httpsrv"
	"github.com/sysinner/incore/v2/pkg/inauth"

	"github.com/hooto/iam/v2/internal/data"
	"github.com/hooto/iam/v2/pkg/iamapi"
)

func NewModule() *httpsrv.Module {

	mod := httpsrv.NewModule()

	mod.RegisterController(
		new(Service),
		new(AppAuth),
		new(Sys),
		new(AccessKey),
		new(User),
		new(AdminAccessKey),
	)

	return mod
}

type UserAuthResponse struct {
	Status inauth.ServiceStatus `json:"status"`
}

func newUserAuthResponse(code, message string) *UserAuthResponse {
	return &UserAuthResponse{
		Status: inauth.NewServiceStatus(code, message),
	}
}

// userAuth extracts and validates the access token from cookie or Authorization
// header, returning the authenticated username or empty string on failure.
func userAuth(c *httpsrv.Controller) *iamapi.User {

	tokenStr := ""

	cookie, err := c.Request.Cookie(inauth.AppHttpHeaderKey)
	if err == nil && cookie.Value != "" {
		tokenStr = cookie.Value
	}

	if tokenStr == "" {
		tokenStr = c.Request.Header.Get("Authorization")
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		if tokenStr == "" {
			c.RenderJson(newUserAuthResponse("401", "Unauthorized"))
			return nil
		}
	}

	token, err := inauth.ParseAccessToken(tokenStr)
	if err != nil || token.Claims.Sub == "" {
		c.RenderJson(newUserAuthResponse("401", "Unauthorized"))
		return nil
	}

	if _, err := token.Verify(data.KeyMgr); err != nil {
		c.RenderJson(newUserAuthResponse("401", "Unauthorized"))
		return nil
	}

	user := data.UserGet(token.Claims.Sub)
	if user == nil {
		c.RenderJson(newUserAuthResponse("401", "Unauthorized"))
		return nil
	}

	return user
}
