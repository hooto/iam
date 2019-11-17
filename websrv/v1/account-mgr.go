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
	"github.com/lynkdb/iomix/sko"
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

	if rs := store.Data.NewReader(nil).
		KeyRangeSet(iamapi.ObjKeyUser(""), iamapi.ObjKeyUser("")).LimitNumSet(10000).Query(); rs.OK() {
		for _, v := range rs.Items {
			var user iamapi.User
			if v.Decode(&user) == nil {
				users = append(users, user.Name)
			}
		}
	}

	for _, uname := range users {

		k1 := iamapi.ObjKeyAccFundUser(uname, "zzzzzzzz")
		k2 := iamapi.ObjKeyAccFundUser(uname, "")
		if rs := store.Data.NewReader(nil).KeyRangeSet(k1, k2).
			ModeRevRangeSet(true).LimitNumSet(1000).Query(); rs.OK() {

			var (
				balance float64 = 0
				prepay  float64 = 0
			)

			for _, v := range rs.Items {

				var aa iamapi.AccountFund
				if err := v.Decode(&aa); err == nil {
					balance += (aa.Amount - aa.Payout - aa.Prepay)
					prepay += aa.Prepay
				}
			}

			balance = iamapi.AccountFloat64Round(balance, 4)
			prepay = iamapi.AccountFloat64Round(prepay, 4)

			var au iamapi.AccountUser
			if rs := store.Data.NewReader(iamapi.ObjKeyAccUser(uname)).Query(); rs.OK() {
				rs.Decode(&au)
			}

			if au.User == "" {
				au.User = uname
			}

			if au.Balance != balance || au.Prepay != prepay {

				au.Balance = balance
				au.Prepay = prepay
				au.Updated = uint64(types.MetaTimeNow())

				store.Data.NewWriter(iamapi.ObjKeyAccUser(uname), au).Commit()
				rsp.Changed++
			}

			rsp.Total++
		}
	}
}

func (c AccountMgr) FundListAction() {

	ls := types.ObjectList{}
	defer c.RenderJson(&ls)

	if !iamclient.SessionAccessAllowed(c.Session, "user.admin", config.Config.InstanceID) {
		ls.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
		return
	}

	k1 := iamapi.ObjKeyAccFundMgr("zzzzzzzz")
	k2 := iamapi.ObjKeyAccFundMgr("")
	if rs := store.Data.NewReader(nil).KeyRangeSet(k1, k2).
		ModeRevRangeSet(true).LimitNumSet(1000).Query(); rs.OK() {
		for _, v := range rs.Items {

			var set iamapi.AccountFund
			if err := v.Decode(&set); err == nil {
				ls.Items = append(ls.Items, set)
			}
		}
	}

	ls.Kind = "AccountFundList"
}

func (c AccountMgr) FundEntryAction() {

	var set struct {
		types.TypeMeta
		iamapi.AccountFund
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

	if rs := store.Data.NewReader(iamapi.ObjKeyAccFundMgr(id)).Query(); rs.OK() {
		rs.Decode(&set.AccountFund)
	}

	if set.AccountFund.Id == "" || set.AccountFund.Id != id {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeNotFound, "Object Not Found")
		return
	}

	set.Kind = "AccountFund"
}

func (c AccountMgr) FundNewAction() {

	var set struct {
		types.TypeMeta
		iamapi.AccountFund
	}
	defer c.RenderJson(&set)

	if err := c.Request.JsonDecode(&set.AccountFund); err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Bad Request")
		return
	}

	set.Amount = iamapi.AccountFloat64Round(set.Amount, 4)

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

	var acc_user iamapi.AccountUser
	if rs := store.Data.NewReader(iamapi.ObjKeyAccUser(set.User)).Query(); rs.OK() {
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
	set.Operator = c.us.UserName
	set.Priority = 8
	set.Amount = iamapi.AccountFloat64Round(set.Amount, 4)

	acc_user.Balance = iamapi.AccountFloat64Round(acc_user.Balance+set.Amount, 4)
	acc_user.Updated = set.Updated

	sets := []sko.ClientObjectItem{
		{
			Key:   iamapi.ObjKeyAccFundMgr(set.Id),
			Value: set.AccountFund,
		},
		{
			Key:   iamapi.ObjKeyAccFundUser(set.User, set.Id),
			Value: set.AccountFund,
		},
		{
			Key:   iamapi.ObjKeyAccUser(set.User),
			Value: acc_user,
		},
	}

	for _, v := range sets {
		if rs := store.Data.NewWriter(v.Key, v.Value).Commit(); !rs.OK() {
			set.Error = types.NewErrorMeta(iamapi.ErrCodeInternalError, "IO Error")
			return
		}
	}

	set.Kind = "AccountFund"
}

func (c AccountMgr) FundSetAction() {

	var set struct {
		types.TypeMeta
		iamapi.AccountFund
	}
	defer c.RenderJson(&set)

	if err := c.Request.JsonDecode(&set.AccountFund); err != nil {
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

	var set_prev iamapi.AccountFund

	if rs := store.Data.NewReader(iamapi.ObjKeyAccFundMgr(set.Id)).Query(); rs.OK() {
		rs.Decode(&set_prev)
	}

	if set_prev.Id == "" || set_prev.Id != set.Id {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeNotFound, "Object Not Found 1")
		return
	}

	set_prev.Updated = uint64(types.MetaTimeNow())
	set_prev.Comment = set.Comment
	set_prev.Type = set.Type
	set_prev.ExpProductLimits = set.ExpProductLimits
	set_prev.ExpProductMax = set.ExpProductMax

	sets := []sko.ClientObjectItem{
		{
			Key:   iamapi.ObjKeyAccFundMgr(set.Id),
			Value: set_prev,
		},
		{
			Key:   iamapi.ObjKeyAccFundUser(set_prev.User, set.Id),
			Value: set_prev,
		},
	}

	for _, v := range sets {
		if rs := store.Data.NewWriter(v.Key, v.Value).Commit(); !rs.OK() {
			set.Error = types.NewErrorMeta(iamapi.ErrCodeInternalError, "IO Error")
			return
		}
	}

	set.Kind = "AccountFund"
}

func (c AccountMgr) ChargeListAction() {

	ls := types.ObjectList{}
	defer c.RenderJson(&ls)

	if !iamclient.SessionAccessAllowed(c.Session, "user.admin", config.Config.InstanceID) {
		ls.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
		return
	}

	k1 := iamapi.ObjKeyAccChargeMgr("zzzzzzzz")
	k2 := iamapi.ObjKeyAccChargeMgr("")
	if rs := store.Data.NewReader(nil).KeyRangeSet(k1, k2).
		ModeRevRangeSet(true).LimitNumSet(1000).Query(); rs.OK() {
		for _, v := range rs.Items {

			var set iamapi.AccountCharge
			if err := v.Decode(&set); err == nil {
				ls.Items = append(ls.Items, set)
			}
		}
	}

	ls.Kind = "AccountChargeList"
}

func (c AccountMgr) ChargeEntryAction() {

	var set struct {
		types.TypeMeta
		iamapi.AccountCharge
	}
	defer c.RenderJson(&set)

	var (
		id   = c.Params.Get("id")
		user = c.Params.Get("user")
	)

	if len(id) < 16 {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "ID Not Found")
		return
	}

	if user == "" {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "user Not Found")
		return
	}

	if !iamclient.SessionAccessAllowed(c.Session, "user.admin", config.Config.InstanceID) {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
		return
	}

	//
	if rs := store.Data.NewReader(iamapi.ObjKeyAccChargeUser(user, id)).Query(); rs.OK() {
		rs.Decode(&set.AccountCharge)
	}
	if set.Id == "" || set.Id != id {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeNotFound, "Object Not Found")
		return
	}

	set.Kind = "AccountCharge"
}

func (c AccountMgr) ChargeSetPayoutAction() {

	var (
		set        types.TypeMeta
		set_charge iamapi.AccountCharge
	)
	defer c.RenderJson(&set)

	if err := c.Request.JsonDecode(&set_charge); err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Bad Request")
		return
	}

	set_charge.Payout = iamapi.AccountFloat64Round(set_charge.Payout, 4)

	if len(set_charge.Id) < 16 {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "ID Not Found")
		return
	}

	if set_charge.Payout < 0 {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "payout Not Found")
		return
	}

	if set_charge.User == "" {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "user Not Found")
		return
	}

	if !iamclient.SessionAccessAllowed(c.Session, "user.admin", config.Config.InstanceID) {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
		return
	}

	var (
		charge   iamapi.AccountCharge
		acc_user iamapi.AccountUser
	)

	//
	if rs := store.Data.NewReader(
		iamapi.ObjKeyAccChargeUser(set_charge.User, set_charge.Id)).Query(); rs.OK() {
		rs.Decode(&charge)
	}
	if charge.Id == "" || charge.Id != set_charge.Id {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeNotFound, "Object Not Found")
		return
	}
	if charge.Payout > 0 {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Payment already Closed")
		return
	}

	//
	if rs := store.Data.NewReader(iamapi.ObjKeyAccUser(set_charge.User)).Query(); rs.OK() {
		rs.Decode(&acc_user)
	} else if !rs.NotFound() {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInternalError, "Server Error")
		return
	}
	if acc_user.User == "" || acc_user.User != set_charge.User {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "User Not Found")
		return
	}

	sets := []sko.ClientObjectItem{}
	updated := uint64(types.MetaTimeNow())

	if charge.Fund != "" {
		var fund iamapi.AccountFund
		if rs := store.Data.NewReader(
			iamapi.ObjKeyAccFundUser(set_charge.User, charge.Fund),
		).Query(); rs.OK() {
			rs.Decode(&fund)
		}
		if fund.Id == "" || fund.Id != charge.Fund {
			set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Fund Not Found")
			return
		}

		//
		fund.Prepay = iamapi.AccountFloat64Round(fund.Prepay-charge.Prepay, 4)
		fund.Payout = iamapi.AccountFloat64Round(fund.Payout+set_charge.Payout, 4)
		fund.ExpProductInpay.Del(charge.Product)
		fund.Updated = updated

		sets = append(sets, sko.ClientObjectItem{
			Key:   iamapi.ObjKeyAccFundUser(set_charge.User, charge.Fund),
			Value: fund,
		})

		sets = append(sets, sko.ClientObjectItem{
			Key:   iamapi.ObjKeyAccFundMgr(charge.Fund),
			Value: fund,
		})
	}

	//
	acc_user.Balance = iamapi.AccountFloat64Round(acc_user.Balance+charge.Prepay-set_charge.Payout, 4)
	acc_user.Prepay = iamapi.AccountFloat64Round(acc_user.Prepay-charge.Prepay, 4)
	acc_user.Updated = updated

	//
	charge.Prepay = 0
	charge.Payout = set_charge.Payout
	charge.Updated = updated

	//
	sets = append(sets, sko.ClientObjectItem{
		Key:   iamapi.ObjKeyAccChargeUser(set_charge.User, set_charge.Id),
		Value: charge,
	})
	sets = append(sets, sko.ClientObjectItem{
		Key:   iamapi.ObjKeyAccUser(set_charge.User),
		Value: acc_user,
	})
	sets = append(sets, sko.ClientObjectItem{
		Key:   iamapi.ObjKeyAccChargeMgr(set_charge.Id),
		Value: charge,
	})

	for _, v := range sets {
		if rs := store.Data.NewWriter(v.Key, v.Value).Commit(); !rs.OK() {
			set.Error = types.NewErrorMeta(iamapi.ErrCodeInternalError, "IO Error")
			return
		}
	}

	set.Kind = "AccountCharge"
}
