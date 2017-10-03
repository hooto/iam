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
	"time"

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

func (c AccountMgr) ReBalanceAction() {

	var rsp struct {
		Changed int `json:"changed"`
		Total   int `json:"total"`
	}
	defer c.RenderJson(&rsp)

	if !iamclient.SessionAccessAllowed(c.Session, "user.admin", config.Config.InstanceID) {
		return
	}

	users := []string{}

	if rs := store.PoScan("user", []byte{}, []byte{}, 10000); rs.OK() {
		rss := rs.KvList()
		for _, v := range rss {
			var user iamapi.User
			if v.Decode(&user) == nil {
				users = append(users, user.Name)
			}
		}
	}

	for _, uname := range users {

		ubs := iamapi.UserIdBytes(uname)

		k := skv.NewProgKey(iamapi.AccActiveUser, ubs, "")
		if rs := store.Data.ProgRevScan(k, k, 1000); rs.OK() {

			var (
				balance float64 = 0
				prepay  float64 = 0
			)

			rss := rs.KvList()
			for _, v := range rss {

				var aa iamapi.AccountActive
				if err := v.Decode(&aa); err == nil {
					balance += (iamapi.AccountFloat64Round(aa.Amount) - iamapi.AccountFloat64Round(aa.Payout) - iamapi.AccountFloat64Round(aa.Prepay))
					prepay += iamapi.AccountFloat64Round(aa.Prepay)
				}
			}

			balance = iamapi.AccountFloat64Round(balance)
			prepay = iamapi.AccountFloat64Round(prepay)

			var au iamapi.AccountUser
			if rs := store.Data.ProgGet(skv.NewProgKey(iamapi.AccUser, ubs)); rs.OK() {
				rs.Decode(&au)
			}

			if au.User == "" {
				au.User = uname
			}

			if au.Balance != balance || au.Prepay != prepay {

				au.Balance = balance
				au.Prepay = prepay
				au.Updated = uint64(types.MetaTimeNow())

				store.Data.ProgPut(
					skv.NewProgKey(iamapi.AccUser, ubs),
					skv.NewProgValue(au),
					nil,
				)
				rsp.Changed++
			}

			rsp.Total++
		}
	}
}

func (c AccountMgr) RechargeListAction() {

	ls := types.ObjectList{}
	defer c.RenderJson(&ls)

	if !iamclient.SessionAccessAllowed(c.Session, "user.admin", config.Config.InstanceID) {
		ls.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
		return
	}

	k := skv.NewProgKey(iamapi.AccRechargeMgr, "")
	if rs := store.Data.ProgRevScan(k, k, 100); rs.OK() {
		rss := rs.KvList()
		for _, v := range rss {

			var set iamapi.AccountRecharge
			if err := v.Decode(&set); err == nil {
				ls.Items = append(ls.Items, set)
			}
		}
	}

	ls.Kind = "AccountRechargeList"
}

func (c AccountMgr) RechargeEntryAction() {

	var set struct {
		types.TypeMeta
		iamapi.AccountRecharge
	}
	defer c.RenderJson(&set)

	id := c.Params.Get("id")
	if len(id) < 16 {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeNotFound, "Object Not Found")
		return
	}

	if !iamclient.SessionAccessAllowed(c.Session, "user.admin", config.Config.InstanceID) {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
		return
	}

	set_id := iamapi.HexStringToBytes(id)

	if rs := store.Data.ProgGet(skv.NewProgKey(iamapi.AccRechargeMgr, set_id)); rs.OK() {
		rs.Decode(&set.AccountRecharge)
	}

	if set.AccountRecharge.Id == "" || set.AccountRecharge.Id != id {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeNotFound, "Object Not Found")
		return
	}

	set.Kind = "AccountRecharge"
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

	set.Amount = iamapi.AccountFloat64Round(set.Amount)

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

	var (
		userbs   = iamapi.UserIdBytes(set.User)
		acc_user iamapi.AccountUser
	)
	if rs := store.Data.ProgGet(skv.NewProgKey(iamapi.AccUser, userbs)); rs.OK() {
		rs.Decode(&acc_user)
	} else if !rs.NotFound() {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInternalError, "Server Error")
		return
	}
	if acc_user.User == "" {
		acc_user.User = set.User
	}

	var (
		tn     = time.Now()
		set_id = append(iamapi.Uint32ToBytes(uint32(tn.Unix())), idhash.Rand(8)...) // 4 + 8
	)

	set.Id = iamapi.BytesToHexString(set_id)
	set.Created = uint64(types.MetaTimeSet(tn.UTC()))
	set.Updated = set.Created
	set.UserOpr = c.us.UserName
	set.Priority = 8

	sets := []skv.ProgKeyValue{}

	sets = append(sets, skv.ProgKeyValue{
		Key: skv.NewProgKey(iamapi.AccRechargeMgr, set_id),
		Val: skv.NewProgValue(set.AccountRecharge),
	})

	sets = append(sets, skv.ProgKeyValue{
		Key: skv.NewProgKey(iamapi.AccRechargeUser, userbs, set_id),
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
		ExpProductMax:    set.ExpProductMax,
	}
	sets = append(sets, skv.ProgKeyValue{
		Key: skv.NewProgKey(iamapi.AccActiveUser, userbs, set_id),
		Val: skv.NewProgValue(set_active),
	})

	for _, v := range sets {
		if rs := store.Data.ProgPut(v.Key, v.Val, nil); !rs.OK() {
			set.Error = types.NewErrorMeta(iamapi.ErrCodeInternalError, "IO Error")
			return
		}
	}

	acc_user.Balance += set.Amount
	acc_user.Updated = set.Updated
	if rs := store.Data.ProgPut(
		skv.NewProgKey(iamapi.AccUser, userbs),
		skv.NewProgValue(acc_user),
		nil,
	); !rs.OK() {
		set.Error = types.NewErrorMeta("500", rs.Bytex().String())
		return
	}

	set.Kind = "AccountRecharge"
}

func (c AccountMgr) RechargeSetAction() {

	var set struct {
		types.TypeMeta
		iamapi.AccountRecharge
	}
	defer c.RenderJson(&set)

	if err := c.Request.JsonDecode(&set.AccountRecharge); err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Bad Request")
		return
	}

	if !iamapi.AccountCurrencyTypeValid(set.Type) {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Invalid Type Value")
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

	var (
		set_id   = iamapi.HexStringToBytes(set.Id)
		recharge iamapi.AccountRecharge
	)

	if rs := store.Data.ProgGet(skv.NewProgKey(iamapi.AccRechargeMgr, set_id)); rs.OK() {
		rs.Decode(&recharge)
	}

	if recharge.Id == "" || recharge.Id != set.Id {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeNotFound, "Object Not Found 1")
		return
	}

	recharge.Updated = uint64(types.MetaTimeNow())
	recharge.Comment = set.Comment
	recharge.Type = set.Type
	recharge.ExpProductLimits = set.ExpProductLimits
	recharge.ExpProductMax = set.ExpProductMax

	userbs := iamapi.UserIdBytes(recharge.User)

	var active iamapi.AccountActive
	if rs := store.Data.ProgGet(skv.NewProgKey(iamapi.AccActiveUser, userbs, set_id)); rs.OK() {
		rs.Decode(&active)
	}

	if active.Id != recharge.Id {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeNotFound, "Object Not Found 2")
		return
	}

	active.Updated = recharge.Updated
	active.Options = recharge.Options
	active.ExpProductLimits = recharge.ExpProductLimits
	active.ExpProductMax = recharge.ExpProductMax

	sets := []skv.ProgKeyValue{
		{
			Key: skv.NewProgKey(iamapi.AccRechargeMgr, set_id),
			Val: skv.NewProgValue(recharge),
		},
		{
			Key: skv.NewProgKey(iamapi.AccRechargeUser, userbs, set_id),
			Val: skv.NewProgValue(recharge),
		},
		{
			Key: skv.NewProgKey(iamapi.AccActiveUser, userbs, set_id),
			Val: skv.NewProgValue(active),
		},
	}

	for _, v := range sets {
		if rs := store.Data.ProgPut(v.Key, v.Val, nil); !rs.OK() {
			set.Error = types.NewErrorMeta(iamapi.ErrCodeInternalError, "IO Error")
			return
		}
	}

	set.Kind = "AccountRecharge"
}
