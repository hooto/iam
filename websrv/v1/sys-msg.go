// Copyright 2019 Eryx <evorui аt gmail dοt com>, All rights reserved.
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

package v1

import (
	"time"

	"github.com/hooto/httpsrv"
	"github.com/lessos/lessgo/types"

	"github.com/hooto/hauth/go/hauth/v1"
	"github.com/hooto/hmsg/go/hmsg/v1"
	"github.com/hooto/iam/data"
	"github.com/hooto/iam/iamapi"
	"github.com/hooto/iam/iamclient"
)

type SysMsg struct {
	*httpsrv.Controller
	us iamapi.UserSession
}

func (c *SysMsg) authValid() bool {

	//
	c.us, _ = iamclient.SessionInstance(c.Session)

	if !c.us.IsLogin() {
		c.Response.Out.WriteHeader(401)
		c.RenderJson(types.NewTypeErrorMeta(iamapi.ErrCodeUnauthorized, "Unauthorized"))
		return false
	}

	return true
}

func (c SysMsg) PostAction() {

	var (
		rsp types.TypeMeta
		set hmsg.MsgItem
	)
	defer c.RenderJson(&rsp)

	if err := c.Request.JsonDecode(&set); err != nil {
		rsp.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "InvalidArgument")
		return
	}

	if err := set.Valid(); err != nil {
		rsp.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, err.Error())
		return
	}

	//
	av, err := hauth.NewAppValidatorWithHttpRequest(c.Request.Request, data.KeyMgr)
	if err != nil {
		rsp.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, err.Error())
		return
	}

	var ak hauth.AccessKey
	if rs := data.Data.NewReader(iamapi.NsAccessKey(av.User, av.Id)).Exec(); rs.OK() {
		rs.Item().JsonDecode(&ak)
	}
	if ak.Id == "" || ak.Id != av.Id {
		rsp.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "No Auth Found, AK "+av.Id)
		return
	}
	if terr := av.SignValid(c.Request.RawBody()); terr != nil {
		rsp.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "Auth Sign Failed")
		return
	}

	if set.Created < 1 {
		set.Created = uint32(time.Now().Unix())
	}

	if rs := data.Data.NewWriter(iamapi.ObjKeyMsgQueue(set.Id), nil).SetJsonValue(set).
		SetCreateOnly(true).Exec(); !rs.OK() {
		rsp.Error = types.NewErrorMeta(iamapi.ErrCodeServerError, "server/db err "+rs.ErrorMessage())
		return
	}

	rsp.Kind = "MsgItem"
}

func (c SysMsg) ListAction() {

	if !c.authValid() || c.us.UserName != "sysadmin" {
		return
	}

	var (
		rsp types.ObjectList
	)
	defer c.RenderJson(&rsp)

	var (
		offset = iamapi.ObjKeyMsgSent("")
		cutset = iamapi.ObjKeyMsgSent("zzzzzzzz")
		limit  = int64(100)
	)

	rs := data.Data.NewRanger(offset, cutset).
		SetRevert(true).SetLimit(limit).Exec()
	if !rs.OK() {
		rsp.Error = types.NewErrorMeta(iamapi.ErrCodeServerError,
			"server/db err "+rs.ErrorMessage())
		return
	}

	for _, v := range rs.Items {
		var item hmsg.MsgItem
		if err := v.JsonDecode(&item); err == nil {
			item.Id = item.SentId()
			rsp.Items = append(rsp.Items, &item)
		}
	}

	rsp.Kind = "MsgList"
}

func (c SysMsg) ItemAction() {

	if !c.authValid() || c.us.UserName != "sysadmin" {
		return
	}

	var (
		rsp iamapi.WebServiceKind
		id  = c.Params.Value("id")
	)
	defer c.RenderJson(&rsp)

	if !hmsg.MsgIdRE.MatchString(id) {
		rsp.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "invalid id")
		return
	}

	if rs := data.Data.NewReader(iamapi.ObjKeyMsgSent(id)).Exec(); !rs.OK() {
		rsp.Error = types.NewErrorMeta(iamapi.ErrCodeServerError,
			"server/db err "+rs.ErrorMessage())
		return
	} else {
		var item hmsg.MsgItem
		if err := rs.Item().JsonDecode(&item); err == nil && item.Id != "" {
			rsp.Data = &item
		}
	}

	if rsp.Data != nil {
		rsp.Kind = "MsgItem"
	} else {
		rsp.Error = types.NewErrorMeta(iamapi.ErrCodeServerError,
			"server unknow error")
	}
}
