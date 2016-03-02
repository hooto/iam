// Copyright 2015 lessOS.com, All rights reserved.
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

package ctrl

import (
	"fmt"
	"strings"
	"time"

	"github.com/lessos/bigtree/btapi"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/net/email"
	"github.com/lessos/lessgo/pass"
	"github.com/lessos/lessgo/types"
	"github.com/lessos/lessgo/utils"
	"github.com/lessos/lessgo/utilx"

	"github.com/lessos/lessids/base/login"
	"github.com/lessos/lessids/base/signup"
	"github.com/lessos/lessids/config"
	"github.com/lessos/lessids/idsapi"
	"github.com/lessos/lessids/store"
)

type Reg struct {
	*httpsrv.Controller
}

func (c Reg) SignUpAction() {
	c.Data["continue"] = c.Params.Get("continue")
	c.Data["user_reg_disable"] = config.UserRegistrationDisabled
}

func (c Reg) SignUpRegAction() {

	rsp := struct {
		types.TypeMeta
		Continue string `json:"continue"`
	}{
		Continue: "/ids",
	}

	defer c.RenderJson(&rsp)

	if config.UserRegistrationDisabled {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "The User Registration Disabled"}
		return
	}

	if err := signup.Validate(c.Params); err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, err.Error()}
		return
	}

	uname := strings.TrimSpace(strings.ToLower(c.Params.Get("uname")))

	var user idsapi.User
	if obj := store.BtAgent.ObjectGet("/global/ids/user/" + utils.StringEncode16(uname, 8)); obj.Error == nil {
		obj.JsonDecode(&user)
	}

	if user.Meta.Name == uname {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "The `Username` already exists, please choose another one"}
		return
	}

	auth, _ := pass.HashDefault(c.Params.Get("passwd"))

	user = idsapi.User{
		Meta: types.ObjectMeta{
			ID:      utils.StringEncode16(uname, 8),
			Name:    uname,
			Created: utilx.TimeNow("atom"),
			Updated: utilx.TimeNow("atom"),
		},
		Email:    strings.TrimSpace(strings.ToLower(c.Params.Get("email"))),
		Auth:     auth,
		Name:     strings.Title(uname),
		Status:   1,
		Roles:    []uint32{100},
		Groups:   []uint32{100},
		Timezone: "UTC",
	}

	if obj := store.BtAgent.ObjectSet("/global/ids/user/"+user.Meta.ID, user, nil); obj.Error != nil {
		rsp.Error = &types.ErrorMeta{"500", obj.Error.Message}
		return
	}

	rsp.Kind = "User"
}

func (c Reg) ForgotPassAction() {
}

func (c Reg) ForgotPassPutAction() {

	rsp := struct {
		types.TypeMeta
		Continue string `json:"continue"`
	}{
		Continue: "/ids",
	}

	defer c.RenderJson(&rsp)

	uemail, err := login.EmailSetValidate(c.Params.Get("email"))
	if err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, err.Error()}
		return
	}

	if c.Params.Get("username") == "" {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "User Not Found"}
		return
	}
	userid := utils.StringEncode16(c.Params.Get("username"), 8)

	var user idsapi.User
	if obj := store.BtAgent.ObjectGet("/global/ids/user/" + userid); obj.Error == nil {
		obj.JsonDecode(&user)
	}

	if user.Meta.ID != userid || user.Email != uemail {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "User Not Found"}
		return
	}

	reset := idsapi.UserPasswordReset{
		ID:      utils.StringNewRand(24),
		UserID:  userid,
		Email:   uemail,
		Expired: utilx.TimeNowAdd("atom", "+3600s"),
	}

	if obj := store.BtAgent.ObjectSet("/global/ids/pwd-reset/"+reset.ID, reset, &btapi.ObjectWriteOptions{
		Ttl: 3600000,
	}); obj.Error != nil {
		rsp.Error = &types.ErrorMeta{"500", obj.Error.Message}
		return
	}

	mr, err := email.MailerPull("def")
	if err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInternalError, "Internal Server Error 001"}
		return
	}

	// TODO tempate
	body := fmt.Sprintf(`<html>
<body>
<div>You recently requested a password reset for your %s account. To create a new password, click on the link below:</div>
<br>
<a href="http://%s/ids/reg/pass-reset?id=%s">Reset My Password</a>
<br>
<div>This request was made on %s.</div>
<br>
<div>Regards,</div>
<div>%s Account Service</div>

<div>********************************************************</div>
<div>Please do not reply to this message. Mail sent to this address cannot be answered.</div>
</body>
</html>`, config.Config.ServiceName, c.Request.Host, reset.ID, utilx.TimeNow("datetime"), config.Config.ServiceName)

	err = mr.SendMail(c.Params.Get("email"), c.Translate("Reset your password"), body)

	if err != nil {
		time.Sleep(2e9)
		err = mr.SendMail(c.Params.Get("email"), c.Translate("Reset your password"), body)
	}

	if err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInternalError, err.Error()}
	} else {
		rsp.Kind = "UserAuth"
	}
}

func (c Reg) PassResetAction() {

	if c.Params.Get("id") == "" {
		return
	}

	var reset idsapi.UserPasswordReset
	if obj := store.BtAgent.ObjectGet("/global/ids/pwd-reset/" + c.Params.Get("id")); obj.Error == nil {
		obj.JsonDecode(&reset)
	}

	if reset.ID != c.Params.Get("id") {
		return
	}

	c.Data["pass_reset_id"] = c.Params.Get("id")
}

func (c Reg) PassResetPutAction() {

	rsp := types.TypeMeta{}

	defer c.RenderJson(&rsp)

	if c.Params.Get("id") == "" {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "Token can not be null"}
		return
	}

	if err := login.PassSetValidate(idsapi.UserPasswordSet{
		NewPassword: c.Params.Get("passwd"),
	}); err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, err.Error()}
		return
	}

	var reset idsapi.UserPasswordReset
	rsobj := store.BtAgent.ObjectGet("/global/ids/pwd-reset/" + c.Params.Get("id"))
	if rsobj.Error == nil {
		rsobj.JsonDecode(&reset)
	}

	if reset.ID != c.Params.Get("id") {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "Token not found"}
		return
	}

	if reset.Email != c.Params.Get("email") {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "Email is not valid"}
		return
	}

	var user idsapi.User
	uobj := store.BtAgent.ObjectGet("/global/ids/user/" + reset.UserID)
	if uobj.Error == nil {
		uobj.JsonDecode(&user)
	}

	if user.Meta.ID != reset.UserID {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "User Not Found"}
		return
	}

	user.Email = reset.Email
	user.Auth, _ = pass.HashDefault(c.Params.Get("passwd"))
	user.Meta.Updated = utilx.TimeNow("atom")

	if obj := store.BtAgent.ObjectSet("/global/ids/user/"+user.Meta.ID, user, &btapi.ObjectWriteOptions{
		PrevVersion: uobj.Meta.Version,
	}); obj.Error != nil {
		rsp.Error = &types.ErrorMeta{"500", obj.Error.Message}
		return
	}

	store.BtAgent.ObjectDel("/global/ids/pwd-reset/"+reset.ID, &btapi.ObjectWriteOptions{
		PrevVersion: rsobj.Meta.Version,
	})

	rsp.Kind = "UserAuth"
}
