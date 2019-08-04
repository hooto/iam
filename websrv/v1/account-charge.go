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
	"github.com/hooto/httpsrv"
	"github.com/hooto/iam/iamapi"
	"github.com/hooto/iam/iamauth"
	"github.com/hooto/iam/store"
	"github.com/lessos/lessgo/types"
	"github.com/lynkdb/iomix/skv"
)

type AccountCharge struct {
	*httpsrv.Controller
	us iamapi.UserSession
}

func (c AccountCharge) PreValidAction() {

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
	authValidator, err := iamauth.NewAppValidatorWithHttpRequest(c.Request.Request)
	if err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, err.Error())
		return
	}

	var ak iamapi.AccessKey
	if rs := store.Data.KvProgGet(iamapi.DataAccessKeyKey(authValidator.User, authValidator.AccessKey)); rs.OK() {
		rs.Decode(&ak)
	}
	if ak.AccessKey == "" || ak.AccessKey != authValidator.AccessKey {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "No Auth Found, AK "+authValidator.AccessKey)
		return
	}
	if terr := authValidator.SignValid(c.Request.RawBody, ak.AuthKey()); terr != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "Auth Sign Failed")
		return
	}

	set.Prepay = iamapi.AccountFloat64Round(set.Prepay, 2)

	var acc_user iamapi.AccountUser
	if rs := store.Data.KvProgGet(iamapi.DataAccUserKey(set.User)); rs.OK() {
		rs.Decode(&acc_user)
	} else if !rs.NotFound() {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInternalError, "Server Error")
		return
	}

	if acc_user.User == "" || acc_user.User != set.User {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeAccChargeOut, "Out of balance")
		return
	}

	actives := []iamapi.AccountFund{}
	ka := iamapi.DataAccFundUserKey(set.User, "")
	if rs := store.Data.KvProgScan(ka, ka, 1000); rs.OK() {
		rss := rs.KvList()
		for _, v := range rss {
			var v2 iamapi.AccountFund
			if err := v.Decode(&v2); err == nil {
				if (v2.Amount - v2.Payout - v2.Prepay) > 0 {
					actives = append(actives, v2)
				}
			}
		}
	}

	var active iamapi.AccountFund
	for _, v := range actives {

		if v.ExpProductMax > 0 &&
			len(v.ExpProductInpay) >= v.ExpProductMax &&
			!v.ExpProductInpay.Has(set.Product) {
			continue
		}

		balance := v.Amount - v.Prepay - v.Payout
		if set.Prepay > balance {
			continue
		}

		active = v
		break
	}

	if active.Id == "" {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeAccChargeOut, "Out of balance")
		return
	}

	set.Kind = "AccountCharge"
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
	authValidator, err := iamauth.NewAppValidatorWithHttpRequest(c.Request.Request)
	if err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, err.Error())
		return
	}

	var ak iamapi.AccessKey
	if rs := store.Data.KvProgGet(iamapi.DataAccessKeyKey(authValidator.User, authValidator.AccessKey)); rs.OK() {
		rs.Decode(&ak)
	}
	if ak.AccessKey == "" || ak.AccessKey != authValidator.AccessKey {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "No Auth Found, AK "+authValidator.AccessKey)
		return
	}
	if terr := authValidator.SignValid(c.Request.RawBody, ak.AuthKey()); terr != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "Auth Sign Failed")
		return
	}

	var (
		_, charge_id = iamapi.AccountChargeId(set.Product, set.TimeStart)
		charge       iamapi.AccountCharge
	)

	if rs := store.Data.KvProgGet(
		iamapi.DataAccChargeUserKey(set.User, charge_id),
	); rs.OK() {
		if err := rs.Decode(&charge); err == nil {
			if charge.Prepay == set.Prepay {
				set.Kind = "AccountChargePrepay"
				return
			}
		}
	}

	set.Prepay = iamapi.AccountFloat64Round(set.Prepay, 2)

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
	charge.Comment = set.Comment

	var acc_user iamapi.AccountUser
	if rs := store.Data.KvProgGet(iamapi.DataAccUserKey(charge.User)); rs.OK() {
		rs.Decode(&acc_user)
	} else if !rs.NotFound() {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInternalError, "Server Error")
		return
	}

	if acc_user.User == "" || acc_user.User != charge.User {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "User Not Found")
		return
	}

	var active iamapi.AccountFund

	if charge.Fund == "" {

		actives := []iamapi.AccountFund{}
		ka := iamapi.DataAccFundUserKey(charge.User, "")
		if rs := store.Data.KvProgScan(ka, ka, 1000); rs.OK() {
			rss := rs.KvList()
			for _, v := range rss {
				var v2 iamapi.AccountFund
				if err := v.Decode(&v2); err == nil {
					if (v2.Amount - v2.Payout - v2.Prepay) > 0 {
						actives = append(actives, v2)
					}
				}
			}
		}

		for _, v := range actives {

			if v.ExpProductMax > 0 &&
				len(v.ExpProductInpay) >= v.ExpProductMax &&
				!v.ExpProductInpay.Has(charge.Product) {
				continue
			}

			balance := v.Amount - v.Prepay - v.Payout
			if charge.Prepay > balance {
				continue
			}

			active = v
			charge.Fund = v.Id
			break
		}
	}

	if active.Id == "" || active.Id != charge.Fund {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeAccChargeOut, "Out of balance")
		return
	}

	active.Prepay = iamapi.AccountFloat64Round(active.Prepay+charge.Prepay, 2)
	active.Updated = uint64(types.MetaTimeNow())
	active.ExpProductInpay.Set(charge.Product)

	acc_user.Balance = iamapi.AccountFloat64Round(acc_user.Balance-charge.Prepay, 2)
	acc_user.Prepay = iamapi.AccountFloat64Round(acc_user.Prepay+charge.Prepay, 2)

	sets := []skv.KvProgKeyValue{
		{
			Key: iamapi.DataAccFundUserKey(charge.User, active.Id),
			Val: skv.NewKvEntry(active),
		},
		{
			Key: iamapi.DataAccChargeUserKey(charge.User, charge_id),
			Val: skv.NewKvEntry(charge),
		},
		{
			Key: iamapi.DataAccUserKey(charge.User),
			Val: skv.NewKvEntry(acc_user),
		},
		{
			Key: iamapi.DataAccFundMgrKey(active.Id),
			Val: skv.NewKvEntry(active),
		},
		{
			Key: iamapi.DataAccChargeMgrKey(charge_id),
			Val: skv.NewKvEntry(charge),
		},
	}

	for _, v := range sets {
		if rs := store.Data.KvProgPut(v.Key, v.Val, nil); !rs.OK() {
			set.Error = types.NewErrorMeta(iamapi.ErrCodeInternalError, "IO Error")
			return
		}
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
	authValidator, err := iamauth.NewAppValidatorWithHttpRequest(c.Request.Request)
	if err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, err.Error())
		return
	}

	var ak iamapi.AccessKey
	if rs := store.Data.KvProgGet(iamapi.DataAccessKeyKey(authValidator.User, authValidator.AccessKey)); rs.OK() {
		rs.Decode(&ak)
	}
	if ak.AccessKey == "" || ak.AccessKey != authValidator.AccessKey {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "No Auth Found, AK "+authValidator.AccessKey)
		return
	}
	if terr := authValidator.SignValid(c.Request.RawBody, ak.AuthKey()); terr != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "Auth Sign Failed")
		return
	}

	//
	var acc_user iamapi.AccountUser
	// hlog.Printf("info", "%s %s %d %d", set.User, userid, set.TimeStart, set.TimeClose)
	if rs := store.Data.KvProgGet(iamapi.DataAccUserKey(set.User)); rs.OK() {
		rs.Decode(&acc_user)
	} else if !rs.NotFound() {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInternalError, "Server Error")
		return
	}
	if acc_user.User == "" || acc_user.User != set.User {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "User Not Found")
		return
	}

	var (
		_, charge_id = iamapi.AccountChargeId(set.Product, set.TimeStart)
		charge       iamapi.AccountCharge
	)
	if rs := store.Data.KvProgGet(
		iamapi.DataAccChargeUserKey(set.User, charge_id),
	); rs.OK() {
		rs.Decode(&charge)
	}

	set.Payout = iamapi.AccountFloat64Round(set.Payout, 2)

	if charge_id != charge.Id {

		charge.Id = charge_id
		charge.Created = uint64(types.MetaTimeNow())
		charge.User = set.User

		charge.Product = set.Product
		charge.TimeStart = set.TimeStart
		charge.TimeClose = set.TimeClose
	}

	if set.TimeClose > 0 && charge.TimeClose != set.TimeClose && set.TimeClose > charge.TimeStart {
		charge.TimeClose = set.TimeClose
	}

	charge.Payout = set.Payout
	charge.Updated = uint64(types.MetaTimeNow())
	charge.Comment = set.Comment

	var (
		active  iamapi.AccountFund
		actives = []iamapi.AccountFund{}
	)

	ka := iamapi.DataAccFundUserKey(set.User, "")
	if rs := store.Data.KvProgScan(ka, ka, 1000); rs.OK() {
		rss := rs.KvList()
		for _, v := range rss {
			var v2 iamapi.AccountFund
			if err := v.Decode(&v2); err == nil {
				actives = append(actives, v2)
			}
		}
	}

	for _, v := range actives {

		balance := v.Amount - v.Payout - v.Prepay

		if (charge.Fund == "" && set.Payout <= balance) ||
			(charge.Fund != "" && charge.Fund == v.Id) {

			if charge.Fund == "" {
				if v.ExpProductMax > 0 &&
					len(v.ExpProductInpay) >= v.ExpProductMax &&
					!v.ExpProductInpay.Has(charge.Product) {
					continue
				}
				charge.Fund = v.Id
			}

			active = v

			break
		}
	}

	if charge.Fund == "" || active.Id == "" {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "InvalidArgument")
		return
	}

	if charge.Prepay > 0 {
		active.Prepay = iamapi.AccountFloat64Round(active.Prepay-charge.Prepay, 2)
		acc_user.Prepay = iamapi.AccountFloat64Round(acc_user.Prepay-charge.Prepay, 2)
	}

	active.Payout = iamapi.AccountFloat64Round(active.Payout+charge.Payout, 2)
	active.Updated = uint64(types.MetaTimeNow())
	active.ExpProductInpay.Del(charge.Product)

	acc_user.Balance = iamapi.AccountFloat64Round(acc_user.Balance-charge.Payout, 2)
	acc_user.Updated = active.Updated

	sets := []skv.KvProgKeyValue{
		{
			Key: iamapi.DataAccFundUserKey(set.User, active.Id),
			Val: skv.NewKvEntry(active),
		},
		{
			Key: iamapi.DataAccChargeUserKey(set.User, charge_id),
			Val: skv.NewKvEntry(charge),
		},
		{
			Key: iamapi.DataAccUserKey(set.User),
			Val: skv.NewKvEntry(acc_user),
		},
		{
			Key: iamapi.DataAccFundMgrKey(active.Id),
			Val: skv.NewKvEntry(active),
		},
		{
			Key: iamapi.DataAccChargeMgrKey(charge_id),
			Val: skv.NewKvEntry(charge),
		},
	}

	for _, v := range sets {
		if rs := store.Data.KvProgPut(v.Key, v.Val, nil); !rs.OK() {
			set.Error = types.NewErrorMeta(iamapi.ErrCodeInternalError, "IO Error")
			return
		}
	}

	set.Kind = "AccountChargePayout"
}
