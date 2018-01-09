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
	"github.com/hooto/iam/iamapi"
	"github.com/hooto/iam/iamclient"
	"github.com/hooto/iam/store"
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

	if rs := store.Data.ProgGet(
		iamapi.DataAccUserKey(c.us.UserName),
	); rs.OK() {
		rs.Decode(&set.AccountUser)
	}

	if set.AccountUser.User == c.us.UserName {
		set.Kind = "AccountUser"
	}
}

func (c Account) FundListAction() {

	ls := types.ObjectList{}
	defer c.RenderJson(&ls)

	k := iamapi.DataAccFundUserKey(c.us.UserName, "")
	if rs := store.Data.ProgRevScan(k, k, 1000); rs.OK() {
		rss := rs.KvList()
		for _, v := range rss {

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

	k := iamapi.DataAccChargeUserKey(c.us.UserName, "")
	if rs := store.Data.ProgRevScan(k, k, 1000); rs.OK() {
		rss := rs.KvList()
		for _, v := range rss {

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

	k := iamapi.DataAccChargeUserKey(c.us.UserName, "")
	if rs := store.Data.ProgRevScan(k, k, 1000); rs.OK() {
		rss := rs.KvList()
		for _, v := range rss {

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
