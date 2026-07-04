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
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hooto/httpsrv"
	"github.com/lessos/lessgo/net/email"
	"github.com/lessos/lessgo/pass"
	"github.com/sysinner/innerstack/v2/pkg/inauth"

	"github.com/hooto/iam/v2/internal/apiserver"
	"github.com/hooto/iam/v2/internal/config"
	"github.com/hooto/iam/v2/internal/data"
	"github.com/hooto/iam/v2/pkg/iamapi"
)

type Password_ResetTicketRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Password_ResetTicket handles forgot-password requests.
// It validates username + email, generates a reset token, and sends a reset email.
func Password_ResetTicket(ctx httpsrv.Ctx) error {

	var (
		req Password_ResetTicketRequest
		rsp ServiceStatusResponse
	)
	defer ctx.JSON(&rsp)

	if err := ctx.Request().JsonDecode(&req); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", "Invalid request format")
		return nil
	}

	req.Username = strings.ToLower(req.Username)
	if err := iamapi.UsernameValid(req.Username); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", err.Error())
		return nil
	}

	emailAddr := strings.ToLower(strings.TrimSpace(req.Email))
	if err := iamapi.EmailValid(emailAddr); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", err.Error())
		return nil
	}

	denyCount, denyKey, err := apiserver.UserAuthDenyCheck(req.Username, ctx.Request())
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("400", err.Error())
		return nil
	}
	apiserver.UserAuthDenyIncr(denyCount, denyKey)

	user := data.UserGet(req.Username)
	if user == nil || user.Email != emailAddr {
		// always return success to prevent user enumeration
		rsp.Status = inauth.NewServiceStatus("200", "ok")
		return nil
	}

	// store reset token with 1-hour TTL
	reset := iamapi.UserResetPassword{
		Id:       uuid.NewString(),
		Username: user.Name,
		Email:    emailAddr,
		Expired:  time.Now().Add(3600 * time.Second).Unix(),
	}

	if rs := data.Data.NewWriter(iamapi.NsUserResetPassword(reset.Id), nil).
		SetJsonValue(reset).SetTTL(3600e3).Exec(); !rs.OK() {
		rsp.Status = inauth.NewServiceStatus("500", "Internal server error")
		return nil
	}

	// send reset email
	mr, err := email.MailerPull("def")
	if err != nil {
		slog.Error("reset-password: mailer not available", "error", err)
		rsp.Status = inauth.NewServiceStatus("200", "ok")
		return nil
	}

	body := fmt.Sprintf(`<html>
<body>
<div>You recently requested a password reset for your %s account.</div>
<br>
<div>Your verification code is: <b>%s</b></div>
<br>
<div>This code will expire in 1 hour.</div>
<br>
<div>If you did not make this request, please ignore this email.</div>
<br>
<div>Regards,</div>
<div>%s Account Service</div>
</body>
</html>`, config.Config.ServiceName, reset.Id, config.Config.ServiceName)

	if err := mr.SendMail(emailAddr, "Reset your password", body); err != nil {
		slog.Error("reset-password: send mail failed", "error", err)
	}

	rsp.Status = inauth.NewServiceStatus("200", "ok")

	return nil
}

type Password_ResetConfirmRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

// Password_ResetConfirm sets a new password using a valid reset token.
func Password_ResetConfirm(ctx httpsrv.Ctx) error {

	var (
		req Password_ResetConfirmRequest
		rsp ServiceStatusResponse
	)
	defer ctx.JSON(&rsp)

	if err := ctx.Request().JsonDecode(&req); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", "Invalid request format")
		return nil
	}

	if req.Token == "" {
		rsp.Status = inauth.NewServiceStatus("400", "Token is required")
		return nil
	}

	if len(req.Password) < 8 || len(req.Password) > 30 {
		rsp.Status = inauth.NewServiceStatus("400", "Password must be between 8 and 30 characters long")
		return nil
	}

	// lookup reset token
	var reset iamapi.UserResetPassword
	if rs := data.Data.NewReader(iamapi.NsUserResetPassword(req.Token)).Exec(); !rs.OK() {
		rsp.Status = inauth.NewServiceStatus("400", "Invalid or expired token")
		return nil
	} else {
		if err := rs.Item().JsonDecode(&reset); err != nil {
			rsp.Status = inauth.NewServiceStatus("400", "Invalid or expired token")
			return nil
		}
	}

	if reset.Id != req.Token {
		rsp.Status = inauth.NewServiceStatus("400", "Invalid or expired token")
		return nil
	}

	denyCount, denyKey, err := apiserver.UserAuthDenyCheck(reset.Username, ctx.Request())
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("400", err.Error())
		return nil
	}
	apiserver.UserAuthDenyIncr(denyCount, denyKey)

	// lookup user
	var user iamapi.User
	if rs := data.Data.NewReader(iamapi.NsUser(reset.Username)).Exec(); rs.OK() {
		rs.Item().JsonDecode(&user)
	}

	if user.Name != reset.Username {
		rsp.Status = inauth.NewServiceStatus("400", "User not found")
		return nil
	}

	// update password
	user.Updated = time.Now().Unix()
	auth, _ := pass.HashDefault(req.Password)
	user.Keys.Set(iamapi.UserKeyDefault, auth)

	if rs := data.Data.NewWriter(iamapi.NsUser(reset.Username), nil).
		SetJsonValue(user).SetIncr(0, "user").Exec(); !rs.OK() {
		rsp.Status = inauth.NewServiceStatus("500", "Internal server error")
		return nil
	}

	// delete the used reset token
	data.Data.NewDeleter(iamapi.NsUserResetPassword(reset.Id)).Exec()

	rsp.Status = inauth.NewServiceStatus("200", "ok")

	slog.Info("service/reset-password ok", "user", reset.Username)

	return nil
}
