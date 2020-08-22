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

package worker

import (
	"sync"
	"time"

	"github.com/hooto/hlog4g/hlog"
	"github.com/hooto/iam/data"
	"github.com/hooto/iam/iamapi"
	"github.com/lessos/lessgo/types"
	kv2 "github.com/lynkdb/kvspec/go/kvspec/v2"
)

var (
	mu                          sync.Mutex
	accountChargeCloseRefreshed uint32 = 0
	accountChargeClosePending          = false
	accountChargeCloseTimeout   uint32 = 864000
)

func AccountChargeCloseRefresh() {

	tn := uint32(time.Now().Unix())
	if accountChargeCloseRefreshed+7200 > tn {
		return
	}

	mu.Lock()
	if accountChargeClosePending {
		mu.Unlock()
		return
	}
	accountChargeClosePending = true
	mu.Unlock()

	defer func() {
		accountChargeClosePending = false
	}()

	var (
		offset = iamapi.ObjKeyAccChargeMgr("zzzzzzzz")
		cutset = iamapi.ObjKeyAccChargeMgr("")
		limit  = 100
		num    = 10000 // TODO
	)

	for {

		rs := data.Data.NewReader(nil).KeyRangeSet(offset, cutset).
			ModeRevRangeSet(true).
			LimitNumSet(int64(limit)).Query()
		if !rs.OK() {
			break
		}

		for _, v := range rs.Items {

			var set iamapi.AccountCharge
			if err := v.Decode(&set); err != nil {
				continue
			}

			offset = iamapi.ObjKeyAccChargeMgr(set.Id)

			if set.Prepay == 0 || set.Payout > 0 {
				continue
			}

			if (set.TimeClose + accountChargeCloseTimeout) > tn {
				continue
			}

			//

			var (
				charge   iamapi.AccountCharge
				acc_user iamapi.AccountUser
			)

			//
			if rs := data.Data.NewReader(iamapi.ObjKeyAccChargeUser(set.User, set.Id)).Query(); rs.OK() {
				rs.Decode(&charge)
			}
			if charge.Id == "" || charge.Id != set.Id {
				continue
			}
			if charge.Payout > 0 {
				continue
			}

			if rs := data.Data.NewReader(iamapi.ObjKeyAccUser(set.User)).Query(); rs.OK() {
				rs.Decode(&acc_user)
			} else if !rs.NotFound() {
				continue
			}
			if acc_user.User == "" || acc_user.User != set.User {
				continue
			}

			sets := []kv2.ClientObjectItem{}
			updated := types.MetaTimeNow()

			if charge.Fund != "" {
				var fund iamapi.AccountFund
				if rs := data.Data.NewReader(
					iamapi.ObjKeyAccFundUser(set.User, charge.Fund)).Query(); rs.OK() {
					rs.Decode(&fund)
				}
				if fund.Id == "" || fund.Id != charge.Fund {
					continue
				}

				//
				fund.Prepay = iamapi.AccountFloat64Round(fund.Prepay-charge.Prepay, 4)
				fund.Payout = iamapi.AccountFloat64Round(fund.Payout+set.Payout, 4)
				fund.ExpProductInpay.Del(charge.Product)
				fund.Updated = updated

				sets = append(sets, kv2.ClientObjectItem{
					Key:   iamapi.ObjKeyAccFundUser(set.User, charge.Fund),
					Value: fund,
				})

				sets = append(sets, kv2.ClientObjectItem{
					Key:   iamapi.ObjKeyAccFundMgr(charge.Fund),
					Value: fund,
				})
			}

			//
			acc_user.Balance = iamapi.AccountFloat64Round(acc_user.Balance+charge.Prepay-set.Payout, 4)
			acc_user.Prepay = iamapi.AccountFloat64Round(acc_user.Prepay-charge.Prepay, 4)
			acc_user.Updated = updated

			//
			charge.Prepay = 0
			charge.Payout = set.Payout
			charge.Updated = updated

			//
			sets = append(sets, kv2.ClientObjectItem{
				Key:   iamapi.ObjKeyAccChargeUser(set.User, set.Id),
				Value: charge,
			})
			sets = append(sets, kv2.ClientObjectItem{
				Key:   iamapi.ObjKeyAccUser(set.User),
				Value: acc_user,
			})
			sets = append(sets, kv2.ClientObjectItem{
				Key:   iamapi.ObjKeyAccChargeMgr(set.Id),
				Value: charge,
			})

			hlog.Printf("warn", "iam/worker charge-payout force close, user %s, charge id %s",
				set.User, set.Id)

			for _, v := range sets {
				if rs := data.Data.NewWriter(v.Key, v.Value).Commit(); !rs.OK() {
					return
				}
			}
		}

		num -= len(rs.Items)

		if len(rs.Items) < 1 || num < 1 || !rs.Next {
			break
		}

	}

	accountChargeCloseRefreshed = tn
}
