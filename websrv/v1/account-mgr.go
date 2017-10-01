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

package v1

import (

	// "github.com/hooto/hlog4g/hlog"
	"github.com/hooto/httpsrv"
	"github.com/hooto/iam/config"
	"github.com/hooto/iam/iamapi"
	"github.com/hooto/iam/iamclient"
	"github.com/hooto/iam/store"
	"github.com/lessos/lessgo/crypto/idhash"
	"github.com/lessos/lessgo/types"
	"github.com/lynkdb/iomix/skv"
)

type AccountMgr struct {
	*httpsrv.Controller
	us iamapi.UserSession
}

func (c *AccountMgr) Init() int {

	//
	c.us, _ = iamclient.SessionInstance(c.Session)

	if !c.us.IsLogin() {
		c.Response.Out.WriteHeader(401)
		c.RenderJson(types.NewTypeErrorMeta(iamapi.ErrCodeUnauthorized, "Unauthorized"))
		return 1
	}

	return 0
}

func (c AccountMgr) RechargeNewAction() {

	var set struct {
		types.TypeMeta
		iamapi.AccountRecharge
	}
	defer c.RenderJson(&set)

	if err := c.Request.JsonDecode(&set.AccountRecharge); err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Bad Request")
		return
	}

	if set.Amount < 1 {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Invalid Amount Value")
		return
	}

	if !iamapi.AccountCurrencyTypeValid(set.Type) {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Invalid Type Value")
		return
	}

	if !iamapi.UserNameRe2.MatchString(set.User) {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Invalid Username")
		return
	}

	if len(set.ExpProductLimits) == 0 {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Invalid Product Limit")
		return
	}
	for _, plimit := range set.ExpProductLimits {
		if err := plimit.Valid(); err != nil {
			set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Invalid Product Limit")
			return
		}
	}

	if !iamclient.SessionAccessAllowed(c.Session, "user.admin", config.Config.InstanceID) {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
		return
	}

	userid := iamapi.UserId(set.User)

	var login iamapi.User
	if obj := store.PoGet("user", userid); obj.OK() {
		obj.Decode(&login)
	}
	if login.Id == "" || login.Id != userid {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "User Not Found")
		return
	}

	set.Id = idhash.RandHexString(16)
	set.Created = uint64(types.MetaTimeNow())
	set.Updated = set.Created
	set.UserOpr = c.us.UserName
	set.Priority = 8

	sets := []skv.ProgKeyValue{}

	sets = append(sets, skv.ProgKeyValue{
		Key: skv.NewProgKey("iam", "acc_recharge_ls", set.Id),
		Val: skv.NewProgValue(set.AccountRecharge),
	})

	sets = append(sets, skv.ProgKeyValue{
		Key: skv.NewProgKey("iam", "acc_recharge", userid, set.Id),
		Val: skv.NewProgValue(set.AccountRecharge),
	})

	set_active := iamapi.AccountActive{
		Id:               set.Id,
		Type:             set.Type,
		User:             set.User,
		Amount:           set.Amount,
		Priority:         set.Priority,
		Options:          set.Options,
		Created:          set.Created,
		Updated:          set.Updated,
		ExpProductLimits: set.ExpProductLimits,
	}
	sets = append(sets, skv.ProgKeyValue{
		Key: skv.NewProgKey("iam", "acc_active", userid, set_active.Id),
		Val: skv.NewProgValue(set_active),
	})

	for _, v := range sets {
		if rs := store.Data.ProgPut(v.Key, v.Val, nil); !rs.OK() {
			set.Error = types.NewErrorMeta(iamapi.ErrCodeInternalError, "IO Error")
			return
		}
	}

	login.EcoinAmount += set.Amount
	login.Updated = types.MetaTime(set.Updated)
	if rs := store.PoPut("user", userid, login, &skv.PathWriteOptions{
		Force: true,
	}); !rs.OK() {
		set.Error = types.NewErrorMeta("500", rs.Bytex().String())
		return
	}

	set.Kind = "AccountRecharge"
}
