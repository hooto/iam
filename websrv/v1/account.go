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
	"github.com/hooto/iam/data"
	"github.com/hooto/iam/iamapi"
	"github.com/hooto/iam/iamclient"
	"github.com/lessos/lessgo/types"
)

type Account struct {
	*httpsrv.Controller
	us iamapi.UserSession
}

func (c *Account) Init() int {

	//
	c.us, _ = iamclient.SessionInstance(c.Session)

	if !c.us.IsLogin() {
		c.Response.Out.WriteHeader(401)
		c.RenderJson(types.NewTypeErrorMeta(iamapi.ErrCodeUnauthorized, "Unauthorized"))
		return 1
	}

	return 0
}

func (c Account) UserEntryAction() {

	var set struct {
		types.TypeMeta
		iamapi.AccountUser
	}
	defer c.RenderJson(&set)

	if rs := data.Data.NewReader(
		iamapi.ObjKeyAccUser(c.us.UserName),
	).Query(); rs.OK() {
		rs.Decode(&set.AccountUser)
	}

	set.Balance = iamapi.AccountFloat64Round(set.Balance, 2)
	set.Prepay = iamapi.AccountFloat64Round(set.Prepay, 2)

	if set.AccountUser.User == c.us.UserName {
		set.Kind = "AccountUser"
	}
}

func (c Account) FundListAction() {

	ls := types.ObjectList{}
	defer c.RenderJson(&ls)

	k1 := iamapi.ObjKeyAccFundUser(c.us.UserName, "zzzzzzzz")
	k2 := iamapi.ObjKeyAccFundUser(c.us.UserName, "")
	if rs := data.Data.NewReader(nil).KeyRangeSet(k1, k2).
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

func (c Account) ChargeListAction() {

	ls := types.ObjectList{}
	defer c.RenderJson(&ls)

	k1 := iamapi.ObjKeyAccChargeUser(c.us.UserName, "zzzzzzzz")
	k2 := iamapi.ObjKeyAccChargeUser(c.us.UserName, "")
	if rs := data.Data.NewReader(nil).KeyRangeSet(k1, k2).
		ModeRevRangeSet(true).LimitNumSet(1000).Query(); rs.OK() {
		for _, v := range rs.Items {

			var set iamapi.AccountCharge
			if err := v.Decode(&set); err == nil {
				if set.Prepay > 0 && set.Payout == 0 {
					ls.Items = append(ls.Items, set)
				}
			}
		}
	}

	ls.Kind = "AccountChargeList"
}

func (c Account) ChargePayoutListAction() {

	ls := types.ObjectList{}
	defer c.RenderJson(&ls)

	k1 := iamapi.ObjKeyAccChargeUser(c.us.UserName, "zzzzzzzz")
	k2 := iamapi.ObjKeyAccChargeUser(c.us.UserName, "")
	if rs := data.Data.NewReader(nil).KeyRangeSet(k1, k2).
		ModeRevRangeSet(true).LimitNumSet(1000).Query(); rs.OK() {
		for _, v := range rs.Items {

			var set iamapi.AccountCharge
			if err := v.Decode(&set); err == nil {
				if set.Payout > 0 {
					ls.Items = append(ls.Items, set)
				}
			}
		}
	}

	ls.Kind = "AccountChargeList"
}
