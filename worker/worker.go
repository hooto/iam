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
	"github.com/hooto/iam/iamapi"
	"github.com/hooto/iam/store"
	"github.com/lessos/lessgo/types"
	"github.com/lynkdb/iomix/skv"
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
		offset = iamapi.DataAccChargeMgrKey("")
		cutset = iamapi.DataAccChargeMgrKey("")
		limit  = 100
		num    = 10000 // TODO
	)

	for {
		rs := store.Data.KvProgRevScan(offset, cutset, limit)
		if !rs.OK() {
			break
		}

		rss := rs.KvList()

		for _, v := range rss {

			var set iamapi.AccountCharge
			if err := v.Decode(&set); err != nil {
				continue
			}

			if set.Prepay == 0 || set.Payout > 0 {
				continue
			}

			if set.TimeClose+accountChargeCloseTimeout > tn {
				continue
			}

			//

			var (
				charge   iamapi.AccountCharge
				acc_user iamapi.AccountUser
			)

			//
			if rs := store.Data.KvProgGet(iamapi.DataAccChargeUserKey(set.User, set.Id)); rs.OK() {
				rs.Decode(&charge)
			}
			if charge.Id == "" || charge.Id != set.Id {
				continue
			}
			if charge.Payout > 0 {
				continue
			}

			if rs := store.Data.KvProgGet(iamapi.DataAccUserKey(set.User)); rs.OK() {
				rs.Decode(&acc_user)
			} else if !rs.NotFound() {
				continue
			}
			if acc_user.User == "" || acc_user.User != set.User {
				continue
			}

			sets := []skv.KvProgKeyValue{}
			updated := uint64(types.MetaTimeNow())

			if charge.Fund != "" {
				var fund iamapi.AccountFund
				if rs := store.Data.KvProgGet(
					iamapi.DataAccFundUserKey(set.User, charge.Fund),
				); rs.OK() {
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

				sets = append(sets, skv.KvProgKeyValue{
					Key: iamapi.DataAccFundUserKey(set.User, charge.Fund),
					Val: skv.NewKvEntry(fund),
				})

				sets = append(sets, skv.KvProgKeyValue{
					Key: iamapi.DataAccFundMgrKey(charge.Fund),
					Val: skv.NewKvEntry(fund),
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
			sets = append(sets, skv.KvProgKeyValue{
				Key: iamapi.DataAccChargeUserKey(set.User, set.Id),
				Val: skv.NewKvEntry(charge),
			})
			sets = append(sets, skv.KvProgKeyValue{
				Key: iamapi.DataAccUserKey(set.User),
				Val: skv.NewKvEntry(acc_user),
			})
			sets = append(sets, skv.KvProgKeyValue{
				Key: iamapi.DataAccChargeMgrKey(set.Id),
				Val: skv.NewKvEntry(charge),
			})

			hlog.Printf("warn", "iam/worker charge-payout force close, user %s, charge id %s",
				set.User, set.Id)

			for _, v := range sets {
				if rs := store.Data.KvProgPut(v.Key, v.Val, nil); !rs.OK() {
					return
				}
			}
		}

		num -= len(rss)

		if num < 1 || len(rss) < limit {
			break
		}

		offset = iamapi.DataAccChargeMgrKeyBytes(rss[len(rss)-1].Key)
	}

	accountChargeCloseRefreshed = tn
}
