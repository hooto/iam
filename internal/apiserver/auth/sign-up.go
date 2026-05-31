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
	"strings"
	"time"

	"github.com/hooto/httpsrv"
	"github.com/lessos/lessgo/pass"
	"github.com/sysinner/incore/v2/pkg/inauth"

	"github.com/hooto/iam/v2/internal/apiserver"
	"github.com/hooto/iam/v2/internal/config"
	"github.com/hooto/iam/v2/internal/data"
	"github.com/hooto/iam/v2/pkg/iamapi"
)

type SignUpRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	Status   inauth.ServiceStatus `json:"status"`
	Continue string               `json:"continue,omitempty"`
	Username string               `json:"username,omitempty"`
}

// SignUp registers a new user account.
func SignUp(ctx *httpsrv.Context) error {

	var (
		req SignUpRequest
		rsp = SignUpResponse{
			Continue: "/iam/auth/sign-in",
		}
	)
	defer ctx.RenderJson(&rsp)

	if err := ctx.Request().JsonDecode(&req); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", "Invalid request format")
		return nil
	}

	if !config.AllowUserSignUp {
		rsp.Status = inauth.NewServiceStatus("403", "User Registration Disabled")
		return nil
	}

	denyCount, denyKey, err := apiserver.UserAuthDenyCheck(req.Username, ctx.Request())
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("400", err.Error())
		return nil
	}

	// validate username
	req.Username = strings.ToLower(req.Username)
	if err := iamapi.UsernameValid(req.Username); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", err.Error())
		return nil
	}

	// validate email
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if err := iamapi.EmailValid(email); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", err.Error())
		return nil
	}

	// validate password
	if len(req.Password) < 8 || len(req.Password) > 30 {
		rsp.Status = inauth.NewServiceStatus("400", "Password must be between 8 and 30 characters long")
		return nil
	}

	// check if user already exists
	var existing iamapi.User
	if obj := data.Data.NewReader(iamapi.NsUser(req.Username)).Exec(); obj.OK() {
		obj.Item().JsonDecode(&existing)
	}
	if existing.Name == req.Username {
		rsp.Status = inauth.NewServiceStatus("400", "The username already exists, please choose another one")
		return nil
	}

	auth, _ := pass.HashDefault(req.Password)

	tn := time.Now().Unix()

	user := iamapi.User{
		Name:        req.Username,
		Email:       email,
		DisplayName: strings.ToUpper(req.Username[:1]) + req.Username[1:],
		Status:      1,
		Roles:       []string{iamapi.Role_User},
		Created:     tn,
		Updated:     tn,
	}
	user.Keys.Set(iamapi.UserKeyDefault, auth)

	if !data.UserSet(&user) {
		rsp.Status = inauth.NewServiceStatus("500", "Server Error")
		return nil
	}

	rsp.Status = inauth.NewServiceStatus("200", "ok")
	rsp.Username = req.Username

	apiserver.UserAuthDenyIncr(denyCount, denyKey)

	slog.Info("service/signup ok", "user", req.Username)

	return nil
}
