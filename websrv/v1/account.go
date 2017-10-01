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
	"github.com/lynkdb/iomix/skv"
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

func (c Account) RechargeListAction() {

	ls := types.ObjectList{}
	defer c.RenderJson(&ls)

	k := skv.NewProgKey("iam", "acc_recharge", iamapi.UserId(c.us.UserName), "")
	if rs := store.Data.ProgScan(k, k, 100); rs.OK() {
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

func (c Account) ActiveListAction() {

	ls := types.ObjectList{}
	defer c.RenderJson(&ls)

	k := skv.NewProgKey("iam", "acc_active", iamapi.UserId(c.us.UserName), "")
	if rs := store.Data.ProgScan(k, k, 100); rs.OK() {
		rss := rs.KvList()
		for _, v := range rss {

			var set iamapi.AccountActive
			if err := v.Decode(&set); err == nil {
				ls.Items = append(ls.Items, set)
			}
		}
	}

	ls.Kind = "AccountActiveList"
}

func (c Account) ChargeListAction() {

	ls := types.ObjectList{}
	defer c.RenderJson(&ls)

	k := skv.NewProgKey("iam", "acc_charge", iamapi.UserId(c.us.UserName), "")
	if rs := store.Data.ProgRevScan(k, k, 100); rs.OK() {
		rss := rs.KvList()
		for _, v := range rss {

			var set iamapi.AccountChargeEntry
			if err := v.Decode(&set); err == nil {
				ls.Items = append(ls.Items, set)
			}
		}
	}

	ls.Kind = "AccountChargeList"
}

func (c Account) ChargePayoutListAction() {

	ls := types.ObjectList{}
	defer c.RenderJson(&ls)

	k := skv.NewProgKey("iam", "acc_charge", iamapi.UserId(c.us.UserName), "")
	if rs := store.Data.ProgRevScan(k, k, 100); rs.OK() {
		rss := rs.KvList()
		for _, v := range rss {

			var set iamapi.AccountChargeEntry
			if err := v.Decode(&set); err == nil {
				if set.Payout > 0 {
					ls.Items = append(ls.Items, set)
				}
			}
		}
	}

	ls.Kind = "AccountChargeList"
}
