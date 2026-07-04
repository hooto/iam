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

package user

import (
	"strings"
	"time"

	"github.com/hooto/httpsrv"
	"github.com/lessos/lessgo/pass"
	"github.com/sysinner/innerstack/v2/pkg/inauth"

	"github.com/hooto/iam/v2/internal/data"
	"github.com/hooto/iam/v2/pkg/iamapi"
)

type UserPassSetRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type UserPassSetResponse struct {
	Status inauth.ServiceStatus `json:"status"`
}

// PassSet changes the current user password.
func PassSet(ctx httpsrv.Ctx) error {

	u := authCtx(ctx)
	if u == nil {
		return nil
	}

	var (
		req UserPassSetRequest
		rsp UserPassSetResponse
	)
	defer ctx.JSON(&rsp)

	if err := ctx.Request().JsonDecode(&req); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", "Bad Request")
		return nil
	}

	if len(req.NewPassword) < 8 || len(req.NewPassword) > 30 {
		rsp.Status = inauth.NewServiceStatus("400", "Password must be between 8 and 30 characters long")
		return nil
	}

	if auth := u.Keys.Get(iamapi.UserKeyDefault); auth == nil ||
		!pass.Check(req.CurrentPassword, auth.String()) {
		rsp.Status = inauth.NewServiceStatus("400", "Current Password can not match")
		return nil
	}

	u.Updated = time.Now().Unix()
	authKey, _ := pass.HashDefault(req.NewPassword)
	u.Keys.Set(iamapi.UserKeyDefault, authKey)

	data.Data.NewWriter(iamapi.NsUser(u.Name), nil).SetJsonValue(u).Exec()
	data.UserSet(u)

	rsp.Status = inauth.NewServiceStatus("200", "ok")
	return nil
}

type UserEmailSetRequest struct {
	Auth  string `json:"auth"`
	Email string `json:"email"`
}

type UserEmailSetResponse struct {
	Status inauth.ServiceStatus `json:"status"`
}

// EmailSet changes the current user email.
func EmailSet(ctx httpsrv.Ctx) error {

	u := authCtx(ctx)
	if u == nil {
		return nil
	}

	var (
		req UserEmailSetRequest
		rsp UserEmailSetResponse
	)
	defer ctx.JSON(&rsp)

	if err := ctx.Request().JsonDecode(&req); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", "Bad Request")
		return nil
	}

	// validate email
	email := strings.TrimSpace(strings.ToLower(req.Email))
	if err := iamapi.EmailValid(email); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", err.Error())
		return nil
	}

	if auth := u.Keys.Get(iamapi.UserKeyDefault); auth == nil ||
		!pass.Check(req.Auth, auth.String()) {
		rsp.Status = inauth.NewServiceStatus("400", "Password can not match")
		return nil
	}

	u.Email = email
	u.Updated = time.Now().Unix()
	data.Data.NewWriter(iamapi.NsUser(u.Name), nil).SetJsonValue(u).Exec()
	data.UserSet(u)

	rsp.Status = inauth.NewServiceStatus("200", "ok")
	return nil
}
