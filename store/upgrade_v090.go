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

package store

import (
	"errors"
	"fmt"
	"time"

	"github.com/hooto/hlog4g/hlog"
	"github.com/lynkdb/iomix/sko"
	"github.com/lynkdb/iomix/skv"

	"github.com/hooto/iam/iamapi"
)

func upgrade_v090(dbPrev skv.Connector, dbNext sko.ClientConnector) error {

	if dbPrev == nil || dbNext == nil {
		return errors.New("invalid connect")
	}

	var (
		limit  = 1000
		numAll = 0
		tStart = time.Now()
	)

	// dataAppInstance
	if true {

		var (
			offset = iamapi.DataAppInstanceKey("")
			cutset = iamapi.DataAppInstanceKey("")
			num    = 0
		)

		for {

			rs := dbPrev.KvProgScan(offset, cutset, limit)
			if !rs.OK() {
				return errors.New("server error")
			}

			rss := rs.KvList()
			for _, obj := range rss {

				var item iamapi.AppInstance
				if err := obj.Decode(&item); err != nil {
					return err
				}

				if rs := dbNext.NewWriter(iamapi.ObjKeyAppInstance(item.Meta.ID), item).
					ModeCreateSet(true).Commit(); !rs.OK() {
					return fmt.Errorf("db err %s", rs.Message)
				} else if rs.Meta.Created > 0 && rs.Meta.Updated == 0 {
					continue
				}

				num += 1
				hlog.Printf("warn", "Upgrade App %s", item.Meta.ID)
				offset = iamapi.DataAppInstanceKey(item.Meta.ID)
			}

			if len(rss) < limit {
				break
			}
		}

		if num > 0 {
			hlog.Printf("warn", "Upgrade App %d", num)
		}
		numAll += num
	}

	// dataUser
	users := []string{}
	if true {

		var (
			offset = iamapi.DataUserKey("")
			num    = 0
		)

		for {

			rs := dbPrev.KvProgScan(offset,
				iamapi.DataUserKey(""), limit)
			if !rs.OK() {
				return errors.New("server error")
			}

			rss := rs.KvList()
			for _, obj := range rss {

				var item iamapi.User
				if err := obj.Decode(&item); err != nil {
					return err
				}

				ow := sko.NewObjectWriter(iamapi.ObjKeyUser(item.Name), item).
					ModeCreateSet(true).IncrNamespaceSet("user")
				if item.Name == def_sysadmin {
					ow.Meta.IncrId = 1
				}

				rs := dbNext.Commit(ow)
				if !rs.OK() {
					return fmt.Errorf("db err %s", rs.Message)
				} else if rs.Meta.Created > 0 && rs.Meta.Updated == 0 {
					continue
				}

				num += 1
				hlog.Printf("warn", "Upgrade User %s, Incr %d", item.Name, rs.Meta.IncrId)
				users = append(users, item.Name)

				offset = iamapi.DataUserKey(item.Name)
			}

			if len(rss) < limit {
				break
			}
		}

		if num > 0 {
			hlog.Printf("warn", "Upgrade User %d", num)
		}
		numAll += num
	}

	// dataAccessKey
	for _, user := range users {

		var (
			offset = iamapi.DataAccessKeyKey(user, "")
			cutset = iamapi.DataAccessKeyKey(user, "")
			num    = 0
		)

		for {

			rs := dbPrev.KvProgScan(offset, cutset, limit)
			if !rs.OK() {
				return errors.New("server error")
			}

			rss := rs.KvList()
			for _, obj := range rss {

				var item iamapi.AccessKey
				if err := obj.Decode(&item); err != nil {
					return err
				}

				if rs := dbNext.NewWriter(iamapi.ObjKeyAccessKey(user, item.AccessKey), item).
					ModeCreateSet(true).Commit(); !rs.OK() {
					return fmt.Errorf("db err %s", rs.Message)
				} else if rs.Meta.Created > 0 && rs.Meta.Updated == 0 {
					continue
				}

				num += 1
				hlog.Printf("warn", "Upgrade AK %s %s", user, item.AccessKey)
			}

			break
		}

		if num > 0 {
			hlog.Printf("warn", "Upgrade AK %s %d", user, num)
		}
		numAll += num
	}

	// dataUserProfile
	if true {

		var (
			offset = iamapi.DataUserProfileKey("")
			cutset = iamapi.DataUserProfileKey("")
			num    = 0
		)

		for {

			rs := dbPrev.KvProgScan(offset, cutset, limit)
			if !rs.OK() {
				return errors.New("server error")
			}

			rss := rs.KvList()
			for _, obj := range rss {

				var item iamapi.UserProfile
				if err := obj.Decode(&item); err != nil {
					return err
				}

				if item.Login == nil {
					hlog.Printf("error", "Upgrade UserProfile MIS")
					continue
				}

				if rs := dbNext.NewWriter(iamapi.ObjKeyUserProfile(item.Login.Name), item).
					ModeCreateSet(true).Commit(); !rs.OK() {
					return fmt.Errorf("db err %s", rs.Message)
				} else if rs.Meta.Created > 0 && rs.Meta.Updated == 0 {
					continue
				}

				num += 1
				hlog.Printf("warn", "Upgrade UserProfile %s",
					string(iamapi.ObjKeyUserProfile(item.Login.Name)))

				offset = iamapi.DataUserProfileKey(item.Login.Name)
			}

			if len(rss) < limit {
				break
			}
		}

		if num > 0 {
			hlog.Printf("warn", "Upgrade UserProfile %d", num)
		}
		numAll += num
	}

	// AccUser
	if true {

		var (
			offset = iamapi.DataAccUserKey("")
			num    = 0
		)

		for {

			rs := dbPrev.KvProgScan(offset,
				iamapi.DataAccUserKey(""), limit)
			if !rs.OK() {
				return errors.New("server error")
			}

			rss := rs.KvList()
			for _, obj := range rss {

				var item iamapi.AccountUser
				if err := obj.Decode(&item); err != nil {
					return err
				}

				if rs := dbNext.NewWriter(iamapi.ObjKeyAccUser(item.User), item).
					ModeCreateSet(true).Commit(); !rs.OK() {
					return fmt.Errorf("db err %s", rs.Message)
				} else if rs.Meta.Created > 0 && rs.Meta.Updated == 0 {
					continue
				}

				num += 1
				hlog.Printf("warn", "Upgrade AccUser %s", item.User)

				offset = iamapi.DataAccUserKey(item.User)
			}

			if len(rss) < limit {
				break
			}
		}

		if num > 0 {
			hlog.Printf("warn", "Upgrade AccUser %d", num)
		}
		numAll += num
	}

	// AccFund*
	if true {

		var (
			offset = iamapi.DataAccFundMgrKey("")
			cutset = iamapi.DataAccFundMgrKey("")
			num    = 0
		)

		for {

			rs := dbPrev.KvProgScan(offset, cutset, limit)
			if !rs.OK() {
				return errors.New("server error")
			}

			rss := rs.KvList()
			for _, obj := range rss {

				var item iamapi.AccountFund
				if err := obj.Decode(&item); err != nil {
					return err
				}

				if rs := dbNext.NewWriter(iamapi.ObjKeyAccFundUser(item.User, item.Id), item).
					ModeCreateSet(true).Commit(); !rs.OK() {
					return fmt.Errorf("db err %s", rs.Message)
				} else if rs.Meta.Created > 0 && rs.Meta.Updated == 0 {
					continue
				}

				if rs := dbNext.NewWriter(iamapi.ObjKeyAccFundMgr(item.Id), item).
					ModeCreateSet(true).Commit(); !rs.OK() {
					return fmt.Errorf("db err %s", rs.Message)
				} else if rs.Meta.Created > 0 && rs.Meta.Updated == 0 {
					continue
				}

				num += 1
				hlog.Printf("warn", "Upgrade AccFound %s %s", item.User, item.Id)

				offset = iamapi.DataAccFundMgrKey(item.Id)
			}

			if len(rss) < limit {
				break
			}
		}

		if num > 0 {
			hlog.Printf("warn", "Upgrade AccFund %d", num)
		}
		numAll += num
	}

	// AccCharge*
	if true {

		var (
			offset = iamapi.DataAccChargeMgrKey("")
			cutset = iamapi.DataAccChargeMgrKey("")
			num    = 0
		)

		for {

			rs := dbPrev.KvProgScan(offset, cutset, limit)
			if !rs.OK() {
				return errors.New("server error")
			}

			rss := rs.KvList()
			for _, obj := range rss {

				var item iamapi.AccountCharge
				if err := obj.Decode(&item); err != nil {
					return err
				}

				if rs := dbNext.NewWriter(iamapi.ObjKeyAccChargeUser(item.User, item.Id), item).
					ModeCreateSet(true).Commit(); !rs.OK() {
					return fmt.Errorf("db err %s", rs.Message)
				} else if rs.Meta.Created > 0 && rs.Meta.Updated == 0 {
					continue
				}

				if rs := dbNext.NewWriter(iamapi.ObjKeyAccChargeMgr(item.Id), item).
					ModeCreateSet(true).Commit(); !rs.OK() {
					return fmt.Errorf("db err %s", rs.Message)
				} else if rs.Meta.Created > 0 && rs.Meta.Updated == 0 {
					continue
				}

				num += 1
				hlog.Printf("warn", "Upgrade AccCharge %s %s", item.User, item.Id)

				offset = iamapi.DataAccChargeMgrKey(item.Id)
			}

			if len(rss) < limit {
				break
			}
		}

		if num > 0 {
			hlog.Printf("warn", "Upgrade AccCharge %d", num)
		}
		numAll += num
	}

	// SysConfig
	if true {

		var (
			offset = iamapi.DataSysConfigKey("")
			cutset = iamapi.DataSysConfigKey("")
			num    = 0
		)

		for {

			rs := dbPrev.KvProgScan(offset, cutset, limit)
			if !rs.OK() {
				return errors.New("server error")
			}

			rss := rs.KvList()
			for _, obj := range rss {

				if rs := dbNext.NewWriter(iamapi.ObjKeySysConfig(string(obj.Key)), obj.Bytex().String()).
					ModeCreateSet(true).Commit(); !rs.OK() {
					return fmt.Errorf("db err %s", rs.Message)
				} else if rs.Meta.Created > 0 && rs.Meta.Updated == 0 {
					continue
				}

				num += 1
				hlog.Printf("warn", "Upgrade SysConfig %s, len %d", string(obj.Key), len(obj.Bytex().String()))
			}

			break
		}

		if num > 0 {
			hlog.Printf("warn", "Upgrade SysConfig %d", num)
		}
		numAll += num
	}

	// MsgQueue
	if true {

		var (
			offset = iamapi.PrevDataMsgQueue("")
			cutset = iamapi.PrevDataMsgQueue("")
			num    = 0
		)

		for {

			rs := dbPrev.KvScan(offset, cutset, limit)
			if !rs.OK() {
				return errors.New("server error")
			}

			rss := rs.KvList()
			for _, obj := range rss {

				var item iamapi.MsgItem
				if err := obj.Decode(&item); err != nil {
					return err
				}

				if len(item.Id) < 8 {
					continue
				}

				if rs := dbNext.NewWriter(iamapi.ObjKeyMsgQueue(item.Id), item).
					ModeCreateSet(true).Commit(); !rs.OK() {
					return fmt.Errorf("db err %s", rs.Message)
				} else if rs.Meta.Created > 0 && rs.Meta.Updated == 0 {
					continue
				}

				num += 1
				hlog.Printf("warn", "Upgrade MsgQueue %s", item.Id)

				offset = iamapi.PrevDataMsgQueue(item.Id)
			}

			if len(rss) < limit {
				break
			}
		}

		if num > 0 {
			hlog.Printf("warn", "Upgrade MsgQueue %d", num)
		}
		numAll += num
	}

	// MsgSent
	if true {

		var (
			offset = iamapi.PrevDataMsgSent("")
			cutset = iamapi.PrevDataMsgSent("")
			num    = 0
			tn     = uint32(time.Now().Unix())
		)

		for {

			rs := dbPrev.KvScan(offset, cutset, limit)
			if !rs.OK() {
				return errors.New("server error")
			}

			rss := rs.KvList()
			for _, obj := range rss {

				var item iamapi.MsgItem
				if err := obj.Decode(&item); err != nil {
					return err
				}

				if len(item.Id) < 8 {
					continue
				}

				if item.Created > tn {
					item.Created = tn
				}

				if rs := dbNext.NewWriter(iamapi.ObjKeyMsgSent(item.SentId()), item).
					ModeCreateSet(true).Commit(); !rs.OK() {
					return fmt.Errorf("db err %s", rs.Message)
				} else if rs.Meta.Created > 0 && rs.Meta.Updated == 0 {
					continue
				}

				num += 1
				hlog.Printf("warn", "Upgrade MsgSent %s", item.SentId())

				offset = iamapi.PrevDataMsgSent(item.SentId())
			}

			if len(rss) < limit {
				break
			}
		}

		if num > 0 {
			hlog.Printf("warn", "Upgrade MsgSent %d", num)
		}
		numAll += num
	}

	if numAll > 0 {
		hlog.Printf("warn", "Upgrade %d items, in %v", numAll, time.Since(tStart))
	}

	return nil
}
