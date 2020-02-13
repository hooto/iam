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

	"github.com/hooto/iam/iamapi"
	"github.com/hooto/iam/iamauth"
	"github.com/hooto/iam/iamclient"
	"github.com/hooto/iam/store"
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
		set iamapi.MsgItem
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
	authValidator, err := iamauth.NewAppValidatorWithHttpRequest(c.Request.Request)
	if err != nil {
		rsp.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, err.Error())
		return
	}

	var ak iamapi.AccessKey
	if rs := store.Data.NewReader(iamapi.ObjKeyAccessKey(authValidator.User, authValidator.AccessKey)).Query(); rs.OK() {
		rs.Decode(&ak)
	}
	if ak.AccessKey == "" || ak.AccessKey != authValidator.AccessKey {
		rsp.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "No Auth Found, AK "+authValidator.AccessKey)
		return
	}
	if terr := authValidator.SignValid(c.Request.RawBody, ak.AuthKey()); terr != nil {
		rsp.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "Auth Sign Failed")
		return
	}

	if set.Created < 1 {
		set.Created = uint32(time.Now().Unix())
	}

	if rs := store.Data.NewWriter(iamapi.ObjKeyMsgQueue(set.Id), set).
		ModeCreateSet(true).Commit(); !rs.OK() {
		rsp.Error = types.NewErrorMeta(iamapi.ErrCodeServerError, "server/db err "+rs.Message)
		return
	}

	rsp.Kind = "MsgItem"
}

func (c SysMsg) ListAction() {

	if !c.authValid() || c.us.UserName != "sysadmin" {
		return
	}

	var (
		rsp iamapi.MsgList
	)
	defer c.RenderJson(&rsp)

	var (
		offset = iamapi.ObjKeyMsgSent("zzzzzzzz")
		cutset = iamapi.ObjKeyMsgSent("")
		limit  = int64(100)
	)

	rs := store.Data.NewReader(nil).KeyRangeSet(offset, cutset).
		ModeRevRangeSet(true).LimitNumSet(limit).Query()
	if !rs.OK() {
		rsp.Error = types.NewErrorMeta(iamapi.ErrCodeServerError,
			"server/db err "+rs.Message)
		return
	}

	for _, v := range rs.Items {
		var item iamapi.MsgItem
		if err := v.DataValue().Decode(&item, nil); err == nil {
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
		id  = c.Params.Get("id")
	)
	defer c.RenderJson(&rsp)

	if !iamapi.MsgIdReg.MatchString(id) {
		rsp.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "invalid id")
		return
	}

	if rs := store.Data.NewReader(iamapi.ObjKeyMsgSent(id)).Query(); !rs.OK() {
		rsp.Error = types.NewErrorMeta(iamapi.ErrCodeServerError,
			"server/db err "+rs.Message)
		return
	} else {
		var item iamapi.MsgItem
		if err := rs.DataValue().Decode(&item, nil); err == nil && item.Id != "" {
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
