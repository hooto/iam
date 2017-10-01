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

package iamapi

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"regexp"
	"time"

	"github.com/lessos/lessgo/crypto/idhash"
	"github.com/lessos/lessgo/types"
)

var (
	AccGenIdReg           = regexp.MustCompile("^[0-9a-f]{16,24}$")
	accPayTimeZero uint32 = 1500000000
)

const (
	AccountCurrencyTypeCash    uint8 = 1
	AccountCurrencyTypeVirtual uint8 = 32
	AccountCurrencyTypeCard    uint8 = 33
)

func AccountCurrencyTypeValid(v uint8) bool {
	if v == AccountCurrencyTypeCash || v == AccountCurrencyTypeVirtual || v == AccountCurrencyTypeCard {
		return true
	}
	return false
}

// iam/acc_recharge/user-id/rand-id
// iam/acc_recharge_ls/rand-id
type AccountRecharge struct {
	Id               string                    `json:"id"`
	Type             uint8                     `json:"type"`
	User             string                    `json:"user"`
	UserOpr          string                    `json:"user_opr"`
	Amount           float64                   `json:"amount"`
	CashType         uint16                    `json:"cash_type,omitempty"`
	CashAmount       float32                   `json:"cash_amount,omitempty"`
	Priority         uint8                     `json:"priority"`
	Options          types.Labels              `json:"options,omitempty"`
	Comment          string                    `json:"comment,omitempty"`
	Created          uint64                    `json:"created"`
	Updated          uint64                    `json:"updated"`
	ExpProductLimits types.ArrayNameIdentifier `json:"exp_product_limits,omitempty"`
}

type AccountCurrencyOption struct {
	Name  string       `json:"name"`
	Items types.Labels `json:"items,omitempty"`
}

// iam/acc_active/user-id/rand-id
type AccountActive struct {
	Id               string                    `json:"id"`
	Type             uint8                     `json:"type"`
	User             string                    `json:"user"`
	Amount           float64                   `json:"amount"`
	Prepay           float64                   `json:"prepay"`
	Payout           float64                   `json:"payout"`
	Priority         uint8                     `json:"priority"`
	Options          types.Labels              `json:"options,emitempty"`
	Created          uint64                    `json:"created"`
	Updated          uint64                    `json:"updated"`
	ExpProductLimits types.ArrayNameIdentifier `json:"exp_product_limits,omitempty"`
}

// iam/acc_charge/user-id/hash-id
// iam/acc_payout/user-id/hash-id
type AccountChargeEntry struct {
	types.TypeMeta `json:",inline"`
	Id             string               `json:"id"`
	RcId           string               `json:"rcid"`
	User           string               `json:"user"`
	Product        types.NameIdentifier `json:"product"`
	Prepay         float64              `json:"prepay"`
	Payout         float64              `json:"payout"`
	TimeStart      uint32               `json:"time_start"`
	TimeClose      uint32               `json:"time_close"`
	Created        uint64               `json:"created"`
	Updated        uint64               `json:"updated"`
}

func AccountChargeEntryId(prod types.NameIdentifier, start, close uint32) string {

	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, close)

	return hex.EncodeToString(bs) + idhash.HashToHexString([]byte(fmt.Sprintf(
		"acc.charge.%s.%d.%d",
		prod.String(),
		start, close,
	)), 16)
}

const (
	PayTypeLease uint8 = 1
	PayTypeOrder uint8 = 2
)

type AccountPayout struct {
	Id        string               `json:"id"`
	Type      uint8                `json:"type"`
	RcId      string               `json:"rcid"`
	User      string               `json:"user"`
	Product   types.NameIdentifier `json:"product"`
	Payout    float64              `json:"payout"`
	Comments  types.Labels         `json:"comments,omitempty"`
	PayOpened uint64               `json:"pay_opened"`
	PayClosed uint64               `json:"pay_closed"`
	Created   uint64               `json:"created"`
}

type AccountChargePrepay struct {
	types.TypeMeta `json:",inline"`
	User           string               `json:"user"`
	Product        types.NameIdentifier `json:"product"`
	Prepay         float64              `json:"prepay"`
	TimeStart      uint32               `json:"time_start"`
	TimeClose      uint32               `json:"time_close"`
}

func (this *AccountChargePrepay) Valid() error {

	if !UserNameRe2.MatchString(this.User) {
		return errors.New("Invalid User")
	}

	if err := this.Product.Valid(); err != nil {
		return errors.New("Invalid Product : " + err.Error())
	}

	if this.Prepay == 0 {
		return errors.New("Invalid Prepay")
	}

	if this.TimeStart < accPayTimeZero {
		return errors.New("Invalid TimeStart")
	}

	if this.TimeClose < accPayTimeZero {
		return errors.New("Invalid TimeClose")
	}

	if this.TimeStart >= this.TimeClose {
		return errors.New("Invalid TimeStart or TimeClose")
	}

	return nil
}

type AccountChargePayout struct {
	types.TypeMeta `json:",inline"`
	Id             string               `json:"id"`
	User           string               `json:"user"`
	Product        types.NameIdentifier `json:"product"`
	Payout         float64              `json:"payout"`
	TimeStart      uint32               `json:"time_start"`
	TimeClose      uint32               `json:"time_close"`
}

func (this *AccountChargePayout) Valid() error {

	if !UserNameRe2.MatchString(this.User) {
		return errors.New("Invalid User")
	}

	if err := this.Product.Valid(); err != nil {
		return errors.New("Invalid Product : " + err.Error())
	}

	if this.Payout <= 0 {
		return errors.New("Invalid Payout")
	}

	if this.TimeStart < accPayTimeZero {
		return errors.New("Invalid TimeStart")
	}

	if this.TimeClose < accPayTimeZero {
		return errors.New("Invalid TimeClose")
	}

	if this.TimeStart >= this.TimeClose {
		return errors.New("Invalid TimeStart or TimeClose")
	}

	return nil
}

type AccountChargeStatus struct {
	PayOpened uint64 `json:"pay_opened"`
	PayTime   uint64 `json:"pay_time"`
}

const (
	AccountChargeTypePrepay uint8 = 1
	AccountChargeTypePayout uint8 = 2

	AccountChargeCycleHour  uint32 = 3600
	AccountChargeCycleDay   uint32 = 86400
	AccountChargeCycleMonth uint32 = 2592000
)

func account_charge_cycle_fix(cycle uint32) uint32 {
	if cycle < AccountChargeCycleHour {
		cycle = AccountChargeCycleHour
	} else if cycle > AccountChargeCycleMonth {
		cycle = AccountChargeCycleMonth
	}
	if fix := cycle % 3600; fix > 0 {
		cycle -= fix
	}
	return cycle
}

func AccountChargeCycleTimeCloseNow(cycle uint32) uint32 {

	cycle = account_charge_cycle_fix(cycle)
	var (
		tm  = time.Now()
		ctm = uint32(tm.Unix())
	)

	if cycle >= AccountChargeCycleMonth {

		if tm.Month() == 12 {
			ctm = uint32(time.Date(tm.Year()+1, 1, 1,
				0, 0, 0, 0, time.Local).Unix())
		} else {
			ctm = uint32(time.Date(tm.Year(), tm.Month()+1, 1,
				0, 0, 0, 0, time.Local).Unix())
		}

	} else if cycle >= AccountChargeCycleDay {

		ctm = uint32(time.Date(tm.Year(), tm.Month(), tm.Day(),
			0, 0, 0, 0, time.Local).AddDate(0, 0, 1).Unix())

	} else {

		offset := uint32(tm.Hour()*3600 + tm.Minute()*60 + tm.Second())
		if fix := offset % cycle; fix > 0 {
			ctm -= fix
		}

		ctm += cycle
	}

	return ctm

	// if cycle > AccountChargeCycleHour {
	// 	offset += uint32(tm.Hour()) * 3600
	// }
	// if cycle > AccountChargeCycleDay {
	// 	offset += uint32(tm.Day()) * 86400
	// }
	// if fix := uint32(tm.Minute()*60+tm.Second()) % cycle; fix > 0 {
	// 	ctm -= fix
	// }
	// return (ctm + cycle)
}

func AccountChargeCycleTimeClose(cycle, ctc uint32) uint32 {

	cycle = account_charge_cycle_fix(cycle)
	var (
		tm  = time.Unix(int64(ctc), 0)
		ctm = ctc
	)

	if cycle >= AccountChargeCycleMonth {

		ctm = uint32(time.Date(tm.Year(), tm.Month(), 1,
			0, 0, 0, 0, time.Local).Unix())

		if ctm < ctc {

			if tm.Month() == 12 {
				ctm = uint32(time.Date(tm.Year()+1, 1, 1,
					0, 0, 0, 0, time.Local).Unix())
			} else {
				ctm = uint32(time.Date(tm.Year(), tm.Month()+1, 1,
					0, 0, 0, 0, time.Local).Unix())
			}
		}

	} else if cycle >= AccountChargeCycleDay {

		ctm = uint32(time.Date(tm.Year(), tm.Month(), tm.Day(),
			0, 0, 0, 0, time.Local).Unix())

		if ctm < ctc {
			ctm = uint32(time.Date(tm.Year(), tm.Month(), tm.Day(),
				0, 0, 0, 0, time.Local).AddDate(0, 0, 1).Unix())
		}

	} else {

		offset := uint32(tm.Hour()*3600 + tm.Minute()*60 + tm.Second())
		if fix := offset % cycle; fix > 0 {
			ctm = ctm - fix + cycle
		}
	}

	return ctm
	// cycle = account_charge_cycle_fix(cycle)
	// tm := time.Unix(int64(ctm), 0)
	// if fix := uint32(tm.Minute()*60+tm.Second()) % cycle; fix > 0 {
	// 	ctm = ctm - fix + cycle
	// }
}

// deposit
// type AccoutPaymentStatus struct {
// 	PayOpened uint64 `json:"pay_opened"`
// 	PayTime   uint32 `json:"pay_time"`
// }

func AccountFloat64Round(f float64) float64 {
	return math.Trunc(f*1e4+0.5) * 1e-4
}
