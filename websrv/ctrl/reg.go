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

package ctrl

import (
	"fmt"
	"strings"
	"time"

	"github.com/hooto/httpsrv"
	"github.com/lessos/lessgo/net/email"
	"github.com/lessos/lessgo/pass"
	"github.com/lessos/lessgo/types"
	"github.com/lessos/lessgo/utils"
	"github.com/lessos/lessgo/utilx"
	"github.com/lynkdb/iomix/skv"

	"github.com/hooto/iam/base/login"
	"github.com/hooto/iam/base/signup"
	"github.com/hooto/iam/config"
	"github.com/hooto/iam/iamapi"
	"github.com/hooto/iam/store"
)

type Reg struct {
	*httpsrv.Controller
}

func (c Reg) SignUpAction() {
	c.Data["user_reg_disable"] = config.UserRegistrationDisabled

	if len(c.Params.Get("redirect_token")) > 20 &&
		iamapi.ServiceRedirectTokenValid(c.Params.Get("redirect_token")) {
		c.Data["redirect_token"] = c.Params.Get("redirect_token")
	}
}

func (c Reg) SignUpRegAction() {

	rsp := struct {
		types.TypeMeta
		Continue string `json:"continue"`
	}{
		Continue: "/iam",
	}

	defer c.RenderJson(&rsp)

	if config.UserRegistrationDisabled {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "The User Registration Disabled"}
		return
	}

	if err := signup.Validate(c.Params); err != nil {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, err.Error()}
		return
	}

	uname := strings.TrimSpace(strings.ToLower(c.Params.Get("uname")))
	userid := iamapi.UserId(uname)

	var user iamapi.User
	if obj := store.Data.KvProgGet(iamapi.DataUserKey(uname)); obj.OK() {
		obj.Decode(&user)
	}

	if user.Name == uname {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "The `Username` already exists, please choose another one"}
		return
	}

	auth, _ := pass.HashDefault(c.Params.Get("passwd"))

	user = iamapi.User{
		Id:          userid,
		Name:        uname,
		Created:     types.MetaTimeNow(),
		Updated:     types.MetaTimeNow(),
		Email:       strings.TrimSpace(strings.ToLower(c.Params.Get("email"))),
		DisplayName: strings.Title(uname),
		Status:      1,
		Roles:       []uint32{100},
		Groups:      []uint32{100},
	}
	user.Keys.Set(iamapi.UserKeyDefault, auth)

	if obj := store.Data.KvProgPut(iamapi.DataUserKey(user.Name), skv.NewKvEntry(user), nil); !obj.OK() {
		rsp.Error = &types.ErrorMeta{"500", obj.Bytex().String()}
		return
	}

	rsp.Kind = "User"
}

func (c Reg) RetrieveAction() {

	if len(c.Params.Get("redirect_token")) > 20 &&
		iamapi.ServiceRedirectTokenValid(c.Params.Get("redirect_token")) {
		c.Data["redirect_token"] = c.Params.Get("redirect_token")
	}
}

func (c Reg) RetrievePutAction() {

	rsp := struct {
		types.TypeMeta
		Continue string `json:"continue"`
	}{
		Continue: "/iam",
	}

	defer c.RenderJson(&rsp)

	uemail, err := login.EmailSetValidate(c.Params.Get("email"))
	if err != nil {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, err.Error()}
		return
	}

	if c.Params.Get("username") == "" {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "User Not Found"}
		return
	}
	uname := c.Params.Get("username")
	userid := iamapi.UserId(uname)

	var user iamapi.User
	if obj := store.Data.KvProgGet(iamapi.DataUserKey(uname)); obj.OK() {
		obj.Decode(&user)
	}

	if user.Id != userid || user.Email != uemail {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "User Not Found"}
		return
	}

	reset := iamapi.UserPasswordReset{
		Id:       utils.StringNewRand(24),
		UserName: user.Name,
		Email:    uemail,
		Expired:  utilx.TimeNowAdd("atom", "+3600s"),
	}

	if obj := store.Data.KvProgPut(iamapi.DataPasswordResetKey(reset.Id), skv.NewKvEntry(reset), &skv.KvProgWriteOptions{
		Expired: uint64(time.Now().Add(3600e9).UnixNano()),
	}); !obj.OK() {
		rsp.Error = &types.ErrorMeta{"500", obj.Bytex().String()}
		return
	}

	mr, err := email.MailerPull("def")
	if err != nil {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInternalError, "Internal Server Error 001"}
		return
	}

	rtstr := ""
	if len(c.Params.Get("redirect_token")) > 20 &&
		iamapi.ServiceRedirectTokenValid(c.Params.Get("redirect_token")) {
		rtstr = "&redirect_token=" + c.Params.Get("redirect_token")
	}

	// TODO tempate
	body := fmt.Sprintf(`<html>
<body>
<div>You recently requested a password reset for your %s account. To create a new password, click on the link below:</div>
<br>
<a href="http://%s/iam/reg/pass-reset?id=%s%s">Reset My Password</a>
<br>
<div>This request was made on %s.</div>
<br>
<div>Regards,</div>
<div>%s Account Service</div>

<div>********************************************************</div>
<div>Please do not reply to this message. Mail sent to this address cannot be answered.</div>
</body>
</html>`, config.Config.ServiceName, c.Request.Host, reset.Id, rtstr, utilx.TimeNow("datetime"), config.Config.ServiceName)

	err = mr.SendMail(c.Params.Get("email"), c.Translate("Reset your password"), body)

	if err != nil {
		time.Sleep(2e9)
		err = mr.SendMail(c.Params.Get("email"), c.Translate("Reset your password"), body)
	}

	if err != nil {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInternalError, err.Error()}
	} else {
		rsp.Kind = "UserAuth"
	}
}

func (c Reg) PassResetAction() {

	if c.Params.Get("id") == "" {
		return
	}

	var reset iamapi.UserPasswordReset
	if obj := store.Data.KvProgGet(iamapi.DataPasswordResetKey(c.Params.Get("id"))); obj.OK() {
		obj.Decode(&reset)
	}

	if reset.Id != c.Params.Get("id") {
		return
	}

	c.Data["pass_reset_id"] = c.Params.Get("id")

	if len(c.Params.Get("redirect_token")) > 20 &&
		iamapi.ServiceRedirectTokenValid(c.Params.Get("redirect_token")) {
		c.Data["redirect_token"] = c.Params.Get("redirect_token")
	}
}

func (c Reg) PassResetPutAction() {

	rsp := types.TypeMeta{}

	defer c.RenderJson(&rsp)

	if c.Params.Get("id") == "" {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "Token can not be null"}
		return
	}

	if err := login.PassSetValidate(iamapi.UserPasswordSet{
		NewPassword: c.Params.Get("passwd"),
	}); err != nil {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, err.Error()}
		return
	}

	var reset iamapi.UserPasswordReset
	rsobj := store.Data.KvProgGet(iamapi.DataPasswordResetKey(c.Params.Get("id")))
	if rsobj.OK() {
		rsobj.Decode(&reset)
	}

	if reset.Id != c.Params.Get("id") {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "Token not found"}
		return
	}

	if reset.Email != c.Params.Get("email") {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "Email is not valid"}
		return
	}

	var user iamapi.User
	userid := iamapi.UserId(reset.UserName)
	uobj := store.Data.KvProgGet(iamapi.DataUserKey(reset.UserName))
	if uobj.OK() {
		uobj.Decode(&user)
	}

	if user.Id != userid {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "User Not Found"}
		return
	}

	user.Updated = types.MetaTimeNow()

	auth, _ := pass.HashDefault(c.Params.Get("passwd"))
	user.Keys.Set(iamapi.UserKeyDefault, auth)

	if obj := store.Data.KvProgPut(iamapi.DataUserKey(reset.UserName), skv.NewKvEntry(user), nil); !obj.OK() {
		rsp.Error = &types.ErrorMeta{"500", obj.Bytex().String()}
		return
	}

	store.Data.KvProgDel(iamapi.DataPasswordResetKey(reset.Id), nil)

	rsp.Kind = "UserAuth"
}
