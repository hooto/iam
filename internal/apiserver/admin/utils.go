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

package admin

import (
	"strings"

	"github.com/hooto/httpsrv"
	"github.com/hooto/iam/v2/internal/data"
	"github.com/hooto/iam/v2/internal/util"
	"github.com/hooto/iam/v2/pkg/iamapi"
	"github.com/sysinner/innerstack/v2/pkg/inauth"
)

// userMgrPasswdHidden is the masked sentinel returned to the frontend for the
// password field. When the admin submits it back unchanged, the password is
// left untouched on update.
const userMgrPasswdHidden = "************"

type StatusResponse struct {
	Status inauth.ServiceStatus `json:"status"`
}

func NewStatusResponse(code, message string) *StatusResponse {
	return &StatusResponse{
		Status: inauth.NewServiceStatus(code, message),
	}
}

// authCtx extracts and validates the access token from cookie or Authorization
// header. Returns the authenticated user or renders 401.
func authCtx(ctx httpsrv.Ctx) *iamapi.User {

	tokenStr := ""

	cookie, err := ctx.Request().Cookie(inauth.AppHttpHeaderKey)
	if err == nil && cookie.Value != "" {
		tokenStr = cookie.Value
	}

	if tokenStr == "" {
		tokenStr = ctx.Request().Header.Get("Authorization")
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		if tokenStr == "" {
			ctx.JSON(NewStatusResponse("401", "Unauthorized"))
			return nil
		}
	}

	token, err := inauth.ParseAccessToken(tokenStr)
	if err != nil || token.Claims.Sub == "" {
		ctx.JSON(NewStatusResponse("401", "Unauthorized"))
		return nil
	}

	if _, err := token.Verify(data.KeyMgr); err != nil {
		ctx.JSON(NewStatusResponse("401", "Unauthorized"))
		return nil
	}

	user := data.UserGet(token.Claims.Sub)
	if user == nil {
		ctx.JSON(NewStatusResponse("401", "Unauthorized"))
		return nil
	}

	return user
}

// authAdmin is authCtx with the additional requirement that the user holds the
// sysadmin role. Renders 403 and returns nil otherwise.
func authAdmin(ctx httpsrv.Ctx) *iamapi.User {

	u := authCtx(ctx)
	if u == nil {
		return nil
	}

	if !util.Contains(u.Roles, []string{iamapi.Role_Sysadmin}) {
		ctx.JSON(NewStatusResponse("403", "Forbidden"))
		return nil
	}

	return u
}
