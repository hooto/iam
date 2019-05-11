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
	dataAppInstance   = "ai"
	dataUser          = "u"
	dataUserAuth      = "ua"
	dataUserAuthDeny  = "uad"
	dataPassReset     = "pr"
	dataAccessKey     = "ak"
	dataRole          = "r"
	dataRolePrivilege = "rp"
	dataUserProfile   = "up"
	dataAccUser       = "au"
	dataAccFundUser   = "af"
	dataAccFundMgr    = "afm"
	dataAccChargeUser = "ac"
	dataAccChargeMgr  = "acm"
	dataSysConfig     = "sc"
)

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

func DataAccFundMgrKey(id string) skv.KvProgKey {
	return skv.NewKvProgKey(dataPrefix, dataAccFundMgr, utils.HexStringToBytes(id))
}

func DataAccChargeUserKey(uname, id string) skv.KvProgKey {
	return skv.NewKvProgKey(dataPrefix, dataAccChargeUser,
		UserIdBytes(uname), utils.HexStringToBytes(id))
}

func DataAccChargeMgrKey(id string) skv.KvProgKey {
	return skv.NewKvProgKey(dataPrefix, dataAccChargeMgr, utils.HexStringToBytes(id))
}

func DataAccChargeMgrKeyBytes(id []byte) skv.KvProgKey {
	return skv.NewKvProgKey(dataPrefix, dataAccChargeMgr, id)
}
