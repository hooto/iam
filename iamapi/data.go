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
	"github.com/lynkdb/iomix/skv"
	"github.com/lynkdb/iomix/utils"
)

const (
	dataPrefix        = "iam"
	dataAppInstance   = "ai"  // sko
	dataUser          = "u"   // sko
	dataUserAuth      = "ua"  // skip
	dataUserAuthDeny  = "uad" // skip
	dataPassReset     = "pr"  // skip
	dataAccessKey     = "ak"  // sko
	dataRole          = "r"   // skip
	dataRolePrivilege = "rp"  // skip
	dataUserProfile   = "up"  // sko
	dataAccUser       = "au"  // sko
	dataAccFundUser   = "af"  // sko
	dataAccFundMgr    = "afm" // sko
	dataAccChargeUser = "ac"  // sko
	dataAccChargeMgr  = "acm" // sko
	dataSysConfig     = "sc"  // sko
	dataMsgQueue      = "mq"  // sko
	dataMsgSent       = "ms"  // sko
)

func ObjKeyAppInstance(key string) []byte {
	return []byte(dataPrefix + ":" + dataAppInstance + ":" + key)
}

func ObjKeyUser(uname string) []byte {
	return []byte(dataPrefix + ":" + dataUser + ":" + uname)
}

func ObjKeyAccessKey(uname, id string) []byte {
	return []byte(dataPrefix + ":" + dataAccessKey + ":" + uname + ":" + id)
}

func ObjKeyRole(id uint64) []byte {
	return []byte(dataPrefix + ":" + dataRole + ":" + Uint64ToHexString(id))
}

func ObjKeyUserProfile(uname string) []byte {
	return []byte(dataPrefix + ":" + dataUserProfile + ":" + uname)
}

func ObjKeyAccUser(uname string) []byte {
	return []byte(dataPrefix + ":" + dataAccUser + ":" + uname)
}

func DataAppInstanceKey(id string) skv.KvProgKey {
	if id == "" {
		return skv.NewKvProgKey(dataPrefix, dataAppInstance, []byte{})
	}
	return skv.NewKvProgKey(dataPrefix, dataAppInstance, utils.HexStringToBytes(id))
}

func DataUserKey(uname string) skv.KvProgKey {
	if uname == "" {
		return skv.NewKvProgKey(dataPrefix, dataUser, []byte{})
	}
	return skv.NewKvProgKey(dataPrefix, dataUser, UserIdBytes(uname))
}

func DataUserProfileKey(uname string) skv.KvProgKey {
	if uname == "" {
		return skv.NewKvProgKey(dataPrefix, dataUserProfile, []byte{})
	}
	return skv.NewKvProgKey(dataPrefix, dataUserProfile, UserIdBytes(uname))
}

func DataPasswordResetKey(id string) skv.KvProgKey {
	return skv.NewKvProgKey(dataPrefix, dataPassReset, utils.HexStringToBytes(id))
}

func DataUserAuthDeny(uname, remote_ip string) []byte {
	return []byte(dataPrefix + ":" + dataUserAuthDeny + ":" +
		uname + ":" + remote_ip)
}

func DataUserAuth(uname string, created uint32) []byte {
	return []byte(dataPrefix + ":" + dataUserAuth + ":" +
		uname + ":" + utils.Uint32ToHexString(created))
}

func DataAccessKeyKey(uname, id string) skv.KvProgKey {
	return skv.NewKvProgKey(dataPrefix, dataAccessKey,
		UserIdBytes(uname), utils.HexStringToBytes(id))
}

func DataRoleKey(id uint32) skv.KvProgKey {
	return skv.NewKvProgKey(dataPrefix, dataRole, utils.Uint32ToBytes(id))
}

func DataRolePrivilegeKey(rid uint32, inst string) skv.KvProgKey {
	return skv.NewKvProgKey(dataPrefix, dataRolePrivilege,
		utils.Uint32ToBytes(rid), utils.HexStringToBytes(inst))
}

func DataSysConfigKey(name string) skv.KvProgKey {
	return skv.NewKvProgKey(dataPrefix, dataSysConfig, name)
}

func ObjKeySysConfig(name string) []byte {
	return []byte(dataPrefix + ":" + dataSysConfig + ":" + name)
}

func DataAccUserKey(uname string) skv.KvProgKey {
	if uname == "" {
		return skv.NewKvProgKey(dataPrefix, dataAccUser, []byte{})
	}
	return skv.NewKvProgKey(dataPrefix, dataAccUser, UserIdBytes(uname))
}

func DataAccFundUserKey(uname, id string) skv.KvProgKey {
	return skv.NewKvProgKey(dataPrefix, dataAccFundUser,
		UserIdBytes(uname), utils.HexStringToBytes(id))
}

func ObjKeyAccFundUser(uname, id string) []byte {
	return []byte(dataPrefix + ":" + dataAccFundUser + ":" + uname + ":" + id)
}

func DataAccFundMgrKey(id string) skv.KvProgKey {
	return skv.NewKvProgKey(dataPrefix, dataAccFundMgr, utils.HexStringToBytes(id))
}

func ObjKeyAccFundMgr(id string) []byte {
	return []byte(dataPrefix + ":" + dataAccFundMgr + ":" + id)
}

func DataAccChargeUserKey(uname, id string) skv.KvProgKey {
	return skv.NewKvProgKey(dataPrefix, dataAccChargeUser,
		UserIdBytes(uname), utils.HexStringToBytes(id))
}

func ObjKeyAccChargeUser(uname, id string) []byte {
	return []byte(dataPrefix + ":" + dataAccChargeUser + ":" + uname + ":" + id)
}

func DataAccChargeMgrKey(id string) skv.KvProgKey {
	return skv.NewKvProgKey(dataPrefix, dataAccChargeMgr, utils.HexStringToBytes(id))
}

func ObjKeyAccChargeMgr(id string) []byte {
	return []byte(dataPrefix + ":" + dataAccChargeMgr + ":" + id)
}

func DataAccChargeMgrKeyBytes(id []byte) skv.KvProgKey {
	return skv.NewKvProgKey(dataPrefix, dataAccChargeMgr, id)
}

func DataMsgQueue(id string) []byte {
	return []byte(dataPrefix + ":" + dataMsgQueue + ":" + id)
}

func ObjKeyMsgQueue(id string) []byte {
	return []byte(dataPrefix + ":" + dataMsgQueue + ":" + id)
}

func DataMsgSent(id string) []byte {
	return []byte(dataPrefix + ":" + dataMsgSent + ":" + id)
}

func ObjKeyMsgSent(id string) []byte {
	return []byte(dataPrefix + ":" + dataMsgSent + ":" + id)
}
