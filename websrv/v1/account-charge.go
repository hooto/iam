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
	"github.com/hooto/iam/auth"
	"github.com/hooto/iam/iamapi"
	"github.com/hooto/iam/store"
	"github.com/lessos/lessgo/types"
	"github.com/lynkdb/iomix/skv"
)

type AccountCharge struct {
	*httpsrv.Controller
	us iamapi.UserSession
}

func (c AccountCharge) PrepayAction() {

	set := iamapi.AccountChargePrepay{}
	defer c.RenderJson(&set)

	if err := c.Request.JsonDecode(&set); err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "InvalidArgument")
		return
	}

	if err := set.Valid(); err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, err.Error())
		return
	}

	//
	auth_token, err := auth.NewAuthToken(c.Request.Header.Get(auth.HttpHeaderKey))
	if err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "No Auth Found #01")
		return
	}

	var ak iamapi.AccessKey
	if rs := store.PoGet("ak/"+auth_token.User, auth_token.AccessKey); rs.OK() {
		rs.Decode(&ak)
	}
	if ak.AccessKey == "" || ak.AccessKey != auth_token.AccessKey {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "No Auth Found #02")
		return
	}
	if terr := auth_token.Valid(ak, c.Request.RawBody); terr != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "No Auth Found #03 "+terr.Message)
		return
	}

	var (
		charge_id = iamapi.AccountChargeEntryId(set.Product, set.TimeStart, set.TimeClose)
		charge    iamapi.AccountChargeEntry
	)

	if rs := store.Data.ProgGet(
		skv.NewProgKey("iam", "acc_charge", iamapi.UserId(set.User), charge_id),
	); rs.OK() {
		if err := rs.Decode(&charge); err == nil {
			if charge.Prepay == set.Prepay {
				set.Kind = "AccountChargePrepay"
				return
			}
		}
	}

	set.Prepay = iamapi.AccountFloat64Round(set.Prepay)

	if charge_id != charge.Id {
		charge.Id = charge_id
		charge.Created = uint64(types.MetaTimeNow())
		charge.User = set.User
	}

	charge.Product = set.Product
	charge.TimeStart = set.TimeStart
	charge.TimeClose = set.TimeClose

	charge.Prepay = set.Prepay
	charge.Updated = uint64(types.MetaTimeNow())

	var (
		userid = iamapi.UserId(charge.User)
		login  iamapi.User
	)
	if obj := store.PoGet("user", userid); obj.OK() {
		obj.Decode(&login)
	}
	if login.Id == "" || login.Id != userid {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "User Not Found")
		return
	}

	if charge.RcId == "" {

		actives := []iamapi.AccountActive{}
		ka := skv.NewProgKey("iam", "acc_active", userid, "")
		if rs := store.Data.ProgScan(ka, ka, 1000); rs.OK() {
			rss := rs.KvList()
			for _, v := range rss {
				var active iamapi.AccountActive
				if err := v.Decode(&active); err == nil {
					if active.Amount > 0 {
						actives = append(actives, active)
					}
				}
			}
		}

		for _, active := range actives {

			balance := active.Amount - active.Prepay - active.Payout

			if charge.Prepay > balance {
				continue
			}

			charge.RcId = active.Id

			active.Prepay += charge.Prepay
			active.Amount -= charge.Prepay
			active.Updated = uint64(types.MetaTimeNow())

			login.EcoinAmount -= charge.Prepay
			login.EcoinPrepay += charge.Prepay

			if rs := store.Data.ProgPut(
				skv.NewProgKey("iam", "acc_active", userid, active.Id),
				skv.NewProgValue(active),
				nil,
			); !rs.OK() {
				set.Error = types.NewErrorMeta(iamapi.ErrCodeInternalError, rs.Bytex().String())
				return
			}

			break
		}
	}

	if charge.RcId == "" {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeAccChargeOut, "")
		return
	}

	if rs := store.Data.ProgPut(
		skv.NewProgKey("iam", "acc_charge", userid, charge.Id),
		skv.NewProgValue(charge),
		nil,
	); !rs.OK() {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInternalError, rs.Bytex().String())
		return
	}

	if rs := store.PoPut("user", userid, login, &skv.PathWriteOptions{
		Force: true,
	}); !rs.OK() {
		set.Error = types.NewErrorMeta("500", rs.Bytex().String())
		return
	}

	set.Kind = "AccountChargePrepay"
}

func (c AccountCharge) PayoutAction() {

	set := iamapi.AccountChargePayout{}
	defer c.RenderJson(&set)

	if err := c.Request.JsonDecode(&set); err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "InvalidArgument")
		return
	}

	if err := set.Valid(); err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, err.Error())
		return
	}

	//
	auth_token, err := auth.NewAuthToken(c.Request.Header.Get(auth.HttpHeaderKey))
	if err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "No Auth Found #01")
		return
	}

	var ak iamapi.AccessKey
	if rs := store.PoGet("ak/"+auth_token.User, auth_token.AccessKey); rs.OK() {
		rs.Decode(&ak)
	}
	if ak.AccessKey == "" || ak.AccessKey != auth_token.AccessKey {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "No Auth Found #02")
		return
	}

	if terr := auth_token.Valid(ak, c.Request.RawBody); terr != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "No Auth Found #03")
		return
	}

	//
	var (
		userid = iamapi.UserId(set.User)
		login  iamapi.User
	)
	// hlog.Printf("info", "%s %s %d %d", set.User, userid, set.TimeStart, set.TimeClose)
	if obj := store.PoGet("user", userid); obj.OK() {
		obj.Decode(&login)
	}
	if login.Id == "" || login.Id != userid {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "User Not Found")
		return
	}

	var (
		charge_id = iamapi.AccountChargeEntryId(set.Product, set.TimeStart, set.TimeClose)
		charge    iamapi.AccountChargeEntry
	)
	if rs := store.Data.ProgGet(
		skv.NewProgKey("iam", "acc_charge", iamapi.UserId(set.User), charge_id),
	); rs.OK() {
		rs.Decode(&charge)
	}

	set.Payout = iamapi.AccountFloat64Round(set.Payout)

	if charge_id != charge.Id {

		charge.Id = charge_id
		charge.Created = uint64(types.MetaTimeNow())
		charge.User = set.User

		charge.Product = set.Product
		charge.TimeStart = set.TimeStart
		charge.TimeClose = set.TimeClose
	}

	charge.Payout = set.Payout
	charge.Updated = uint64(types.MetaTimeNow())

	var (
		active  iamapi.AccountActive
		actives = []iamapi.AccountActive{}
	)

	ka := skv.NewProgKey("iam", "acc_active", userid, "")
	if rs := store.Data.ProgScan(ka, ka, 1000); rs.OK() {
		rss := rs.KvList()
		for _, v := range rss {
			var v2 iamapi.AccountActive
			if err := v.Decode(&v2); err == nil {
				actives = append(actives, v2)
			}
		}
	}

	for _, v := range actives {

		if (charge.RcId == "" && set.Payout <= v.Amount) ||
			(charge.RcId != "" && charge.RcId == v.Id) {

			active = v
			charge.RcId = v.Id
			break
		}
	}

	if charge.RcId == "" || active.Id == "" {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "InvalidArgument")
		return
	}

	if charge.Prepay > 0 {
		active.Amount += charge.Prepay
		active.Prepay -= charge.Prepay
		login.EcoinPrepay -= charge.Prepay
	}

	active.Amount -= charge.Payout
	active.Updated = uint64(types.MetaTimeNow())

	sets := []skv.ProgKeyValue{}
	sets = append(sets, skv.ProgKeyValue{
		Key: skv.NewProgKey("iam", "acc_active", userid, active.Id),
		Val: skv.NewProgValue(active),
	})
	sets = append(sets, skv.ProgKeyValue{
		Key: skv.NewProgKey("iam", "acc_charge", userid, charge.Id),
		Val: skv.NewProgValue(charge),
	})

	for _, v := range sets {
		if rs := store.Data.ProgPut(v.Key, v.Val, nil); !rs.OK() {
			set.Error = types.NewErrorMeta(iamapi.ErrCodeInternalError, "IO Error")
			return
		}
	}

	login.EcoinAmount -= charge.Payout
	if rs := store.PoPut("user", userid, login, &skv.PathWriteOptions{
		Force: true,
	}); !rs.OK() {
		set.Error = types.NewErrorMeta("500", rs.Bytex().String())
		return
	}

	set.Kind = "AccountChargePayout"
}
