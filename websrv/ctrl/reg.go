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
	"time"

	"github.com/lessos/lessgo/data/rdo"
	"github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/net/email"
	"github.com/lessos/lessgo/pass"
	"github.com/lessos/lessgo/types"
	"github.com/lessos/lessgo/utils"

	"../../base/login"
	"../../base/signup"
	"../../config"
	"../../idsapi"
)

type Reg struct {
	*httpsrv.Controller
}

func (c Reg) SignUpAction() {
	c.Data["continue"] = c.Params.Get("continue")
}

func (c Reg) SignUpRegAction() {

	c.AutoRender = false

	rsp := struct {
		types.TypeMeta
		Continue string `json:"continue"`
	}{
		Continue: "/ids",
	}

	defer c.RenderJson(&rsp)

	if err := signup.Validate(c.Params); err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, err.Error()}
		return
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInternalError, "Internal Server Error"}
		return
	}

	q := base.NewQuerySet().From("ids_login").Limit(1)
	q.Where.And("email", c.Params.Get("email"))
	rsu, err := dcn.Base.Query(q)
	if err == nil && len(rsu) == 1 {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "The `Email` already exists, please choose another one"}
		return
	}

	pass, _ := pass.HashDefault(c.Params.Get("passwd"))
	uid := utils.StringEncode16(c.Params.Get("uname"), 8)

	item := map[string]interface{}{
		"uid":      uid,
		"uname":    c.Params.Get("uname"),
		"email":    c.Params.Get("email"),
		"pass":     pass,
		"name":     c.Params.Get("name"),
		"status":   1,
		"roles":    "100",
		"timezone": "UTC",                    // TODO
		"created":  base.TimeNow("datetime"), // TODO
		"updated":  base.TimeNow("datetime"), // TODO
	}
	if _, err := dcn.Base.Insert("ids_login", item); err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInternalError, "Can not write to database"}
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

	if email, err := login.EmailSetValidate(c.Params.Get("email")); err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, err.Error()}
		return
	} else {
		c.Params.Set("email", email)
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInternalError, "Internal Server Error"}
		return
	}

	q := base.NewQuerySet().From("ids_login").Limit(1)
	q.Where.And("email", c.Params.Get("email"))
	rsl, err := dcn.Base.Query(q)
	if err != nil || len(rsl) != 1 {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "Email can not found"}
		return
	}

	id := utils.StringNewRand(24)
	item := map[string]interface{}{
		"id":      id,
		"status":  0,
		"email":   c.Params.Get("email"),                 // TODO
		"expired": base.TimeNowAdd("datetime", "+3600s"), // TODO
	}
	if _, err := dcn.Base.Insert("ids_resetpass", item); err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInternalError, "Can not write to database"}
		return
	}

	mr, err := email.MailerPull("def")
	if err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInternalError, "Internal Server Error"}
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
<div>%s Account Services</div>

<div>********************************************************</div>
<div>Please do not reply to this message. Mail sent to this address cannot be answered.</div>
</body>
</html>`, config.Config.ServiceName, c.Request.Host, id, base.TimeNow("datetime"), config.Config.ServiceName)

	err = mr.SendMail(c.Params.Get("email"), c.Translate("Reset your password"), body)

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

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		return
	}

	q := base.NewQuerySet().From("ids_resetpass").Limit(1)
	q.Where.And("id", c.Params.Get("id"))
	rsr, err := dcn.Base.Query(q)
	if err != nil || len(rsr) != 1 {
		return
	}

	expired := rsr[0].Field("expired").TimeParse("datetime") //, base.TimeParse(rsr[0]["expired"].(string), "datetime")
	if expired.Before(time.Now()) {
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

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInternalError, "Internal Server Error"}
		return
	}

	q := base.NewQuerySet().From("ids_resetpass").Limit(1)
	q.Where.And("id", c.Params.Get("id"))
	rsr, err := dcn.Base.Query(q)
	if err != nil || len(rsr) != 1 {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "Token not found"}
		return
	}

	expired := rsr[0].Field("expired").TimeParse("datetime") // base.TimeParse(rsr[0]["expired"].(string), "datetime")
	if expired.Before(time.Now()) {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "Token expired"}
		return
	}

	if rsr[0].Field("email").String() != c.Params.Get("email") {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "Email is not valid"}
		return
	}

	q = base.NewQuerySet().From("ids_login").Limit(1)
	q.Where.And("email", c.Params.Get("email"))
	rsl, err := dcn.Base.Query(q)
	if err != nil || len(rsl) != 1 {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "User can not found"}
		return
	}

	q = base.NewQuerySet().From("ids_profile").Limit(1)
	q.Where.And("uid", rsl[0].Field("uid").Int())
	rspf, err := dcn.Base.Query(q)
	if err != nil || len(rspf) != 1 {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "User can not found"}
		return
	}
	// if fmt.Sprintf("%v", rspf[0].Field("birthday").String()) != c.Params.Get("birthday") {
	// 	rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "Email or Birthday is not valid"}
	// 	return
	// }

	pass, _ := pass.HashDefault(c.Params.Get("passwd"))

	itemlogin := map[string]interface{}{
		"pass":    pass,
		"updated": base.TimeNow("datetime"),
	}
	ft := base.NewFilter()
	ft.And("uid", rsl[0].Field("uid").Int())
	dcn.Base.Update("ids_login", itemlogin, ft)

	//
	delfr := base.NewFilter()
	delfr.And("id", c.Params.Get("id"))
	dcn.Base.Delete("ids_resetpass", delfr)

	rsp.Kind = "UserAuth"
}
