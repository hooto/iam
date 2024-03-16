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

	"github.com/hooto/iam/base/login"
	"github.com/hooto/iam/base/signup"
	"github.com/hooto/iam/config"
	"github.com/hooto/iam/data"
	"github.com/hooto/iam/iamapi"
)

type Reg struct {
	*httpsrv.Controller
}

func (c Reg) SignUpAction() {
	c.Data["user_reg_disable"] = config.UserRegistrationDisabled

	if len(c.Params.Value("redirect_token")) > 20 &&
		iamapi.ServiceRedirectTokenValid(c.Params.Value("redirect_token")) {
		c.Data["redirect_token"] = c.Params.Value("redirect_token")
	}

	c.Data["sys_version_hash"] = config.VersionHash
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

	uname := iamapi.UserNameFilter(c.Params.Value("uname")) // strings.TrimSpace(strings.ToLower(c.Params.Value("uname")))
	// userid := iamapi.UserId(uname)

	var user iamapi.User
	if obj := data.Data.NewReader(iamapi.ObjKeyUser(uname)).Query(); obj.OK() {
		obj.Decode(&user)
	}

	if user.Name == uname {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "The `Username` already exists, please choose another one"}
		return
	}

	auth, _ := pass.HashDefault(c.Params.Value("passwd"))

	user = iamapi.User{
		// Id:          userid,
		Name:        uname,
		Created:     types.MetaTimeNow(),
		Updated:     types.MetaTimeNow(),
		Email:       strings.TrimSpace(strings.ToLower(c.Params.Value("email"))),
		DisplayName: strings.Title(uname),
		Status:      1,
		Roles:       []uint32{100},
	}
	user.Keys.Set(iamapi.UserKeyDefault, auth)

	if !data.UserSet(&user) {
		rsp.Error = &types.ErrorMeta{"500", "Server Error"}
		return
	}

	rsp.Kind = "User"
}

func (c Reg) RetrieveAction() {

	if len(c.Params.Value("redirect_token")) > 20 &&
		iamapi.ServiceRedirectTokenValid(c.Params.Value("redirect_token")) {
		c.Data["redirect_token"] = c.Params.Value("redirect_token")
	}
	c.Data["sys_version_hash"] = config.VersionHash
}

func (c Reg) RetrievePutAction() {

	rsp := struct {
		types.TypeMeta
		Continue string `json:"continue"`
	}{
		Continue: "/iam",
	}

	defer c.RenderJson(&rsp)

	uemail, err := login.EmailSetValidate(c.Params.Value("email"))
	if err != nil {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, err.Error()}
		return
	}

	if c.Params.Value("username") == "" {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "User Not Found"}
		return
	}
	uname := c.Params.Value("username")
	// userid := iamapi.UserId(uname)

	var user iamapi.User
	if obj := data.Data.NewReader(iamapi.ObjKeyUser(uname)).Query(); obj.OK() {
		obj.Decode(&user)
	}

	if user.Name != uname || user.Email != uemail {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "User Not Found"}
		return
	}

	reset := iamapi.UserPasswordReset{
		Id:       utils.StringNewRand(24),
		UserName: user.Name,
		Email:    uemail,
		Expired:  utilx.TimeNowAdd("atom", "+3600s"),
	}

	if obj := data.Data.NewWriter(iamapi.ObjKeyPasswordReset(reset.Id), reset).
		ExpireSet(3600000).Commit(); !obj.OK() {
		rsp.Error = &types.ErrorMeta{"500", obj.Message}
		return
	}

	mr, err := email.MailerPull("def")
	if err != nil {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInternalError, "Internal Server Error 001"}
		return
	}

	rtstr := ""
	if len(c.Params.Value("redirect_token")) > 20 &&
		iamapi.ServiceRedirectTokenValid(c.Params.Value("redirect_token")) {
		rtstr = "&redirect_token=" + c.Params.Value("redirect_token")
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

	err = mr.SendMail(c.Params.Value("email"), c.Translate("Reset your password"), body)

	if err != nil {
		time.Sleep(2e9)
		err = mr.SendMail(c.Params.Value("email"), c.Translate("Reset your password"), body)
	}

	if err != nil {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInternalError, err.Error()}
	} else {
		rsp.Kind = "UserAuth"
	}
}

func (c Reg) PassResetAction() {

	if c.Params.Value("id") == "" {
		return
	}

	var reset iamapi.UserPasswordReset
	if obj := data.Data.NewReader(iamapi.ObjKeyPasswordReset(c.Params.Value("id"))).Query(); obj.OK() {
		obj.Decode(&reset)
	}

	if reset.Id != c.Params.Value("id") {
		return
	}

	c.Data["pass_reset_id"] = c.Params.Value("id")

	if len(c.Params.Value("redirect_token")) > 20 &&
		iamapi.ServiceRedirectTokenValid(c.Params.Value("redirect_token")) {
		c.Data["redirect_token"] = c.Params.Value("redirect_token")
	}

	c.Data["sys_version_hash"] = config.VersionHash
}

func (c Reg) PassResetPutAction() {

	rsp := types.TypeMeta{}

	defer c.RenderJson(&rsp)

	if c.Params.Value("id") == "" {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "Token can not be null"}
		return
	}

	if err := login.PassSetValidate(iamapi.UserPasswordSet{
		NewPassword: c.Params.Value("passwd"),
	}); err != nil {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, err.Error()}
		return
	}

	var reset iamapi.UserPasswordReset
	rsobj := data.Data.NewReader(iamapi.ObjKeyPasswordReset(c.Params.Value("id"))).Query()
	if rsobj.OK() {
		rsobj.Decode(&reset)
	}

	if reset.Id != c.Params.Value("id") {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "Token not found"}
		return
	}

	if reset.Email != c.Params.Value("email") {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "Email is not valid"}
		return
	}

	var user iamapi.User
	// userid := iamapi.UserId(reset.UserName)
	uobj := data.Data.NewReader(iamapi.ObjKeyUser(reset.UserName)).Query()
	if uobj.OK() {
		uobj.Decode(&user)
	}

	if user.Name != reset.UserName {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "User Not Found"}
		return
	}

	user.Updated = types.MetaTimeNow()

	auth, _ := pass.HashDefault(c.Params.Value("passwd"))
	user.Keys.Set(iamapi.UserKeyDefault, auth)

	if obj := data.Data.NewWriter(iamapi.ObjKeyUser(reset.UserName), user).
		IncrNamespaceSet("user").Commit(); !obj.OK() {
		rsp.Error = &types.ErrorMeta{"500", obj.Message}
		return
	}

	data.Data.NewWriter(iamapi.ObjKeyPasswordReset(reset.Id), nil).ModeDeleteSet(true).Commit()

	rsp.Kind = "UserAuth"
}
