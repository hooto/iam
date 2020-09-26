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
	"github.com/lynkdb/iomix/utils"
)

const (
	dataPrefix        = "iam"
	dataAppInstance   = "ai"  // kv2
	dataUser          = "u"   // kv2
	dataUserAuth      = "ua"  // skip
	dataUserAuthDeny  = "uad" // skip
	dataPassReset     = "pr"  // skip
	dataAccessKey     = "ak2" // hauth
	dataRole          = "r"   // skip
	dataRolePrivilege = "rp"  // skip
	dataUserProfile   = "up"  // kv2
	dataAccUser       = "au"  // kv2
	dataAccFundUser   = "af"  // kv2
	dataAccFundMgr    = "afm" // kv2
	dataAccChargeUser = "ac"  // kv2
	dataAccChargeMgr  = "acm" // kv2
	dataSysConfig     = "sc"  // kv2
	dataMsgQueue      = "mq"  // kv2
	dataMsgSent       = "ms"  // kv2
)

func ObjKeyAppInstance(key string) []byte {
	return []byte(dataPrefix + ":" + dataAppInstance + ":" + key)
}

func ObjKeyUser(uname string) []byte {
	return []byte(dataPrefix + ":" + dataUser + ":" +
		UserNameFilter(uname))
}

func ObjKeyUserProfile(uname string) []byte {
	return []byte(dataPrefix + ":" + dataUserProfile + ":" +
		UserNameFilter(uname))
}

func ObjKeyPasswordReset(id string) []byte {
	return []byte(dataPrefix + ":" + dataPassReset + ":" + id)
}

func DataUserAuthDeny(uname, remote_ip string) []byte {
	return []byte(dataPrefix + ":" + dataUserAuthDeny + ":" +
		uname + ":" + remote_ip)
}

func ObjKeyUserAuthDeny(uname, remote_ip string) []byte {
	return []byte(dataPrefix + ":" + dataUserAuthDeny + ":" +
		UserNameFilter(uname) + ":" + remote_ip)
}

func DataUserAuth(uname string, created uint32) []byte {
	return []byte(dataPrefix + ":" + dataUserAuth + ":" +
		uname + ":" + utils.Uint32ToHexString(created))
}

func ObjKeyUserAuth(uname string, created uint32) []byte {
	return []byte(dataPrefix + ":" + dataUserAuth + ":" +
		UserNameFilter(uname) + ":" + utils.Uint32ToHexString(created))
}

func NsAccessKey(uname, id string) []byte {
	if uname == "" {
		return []byte(dataPrefix + ":" + dataAccessKey + ":")
	}
	return []byte(dataPrefix + ":" + dataAccessKey + ":" +
		UserNameFilter(uname) + ":" + id)
}

func ObjKeySysConfig(name string) []byte {
	return []byte(dataPrefix + ":" + dataSysConfig + ":" + name)
}

func ObjKeyAccUser(uname string) []byte {
	return []byte(dataPrefix + ":" + dataAccUser + ":" +
		UserNameFilter(uname))
}

func ObjKeyAccFundUser(uname, id string) []byte {
	return []byte(dataPrefix + ":" + dataAccFundUser + ":" +
		UserNameFilter(uname) + ":" + id)
}

func ObjKeyAccFundMgr(id string) []byte {
	return []byte(dataPrefix + ":" + dataAccFundMgr + ":" + id)
}

func ObjKeyAccChargeUser(uname, id string) []byte {
	return []byte(dataPrefix + ":" + dataAccChargeUser + ":" +
		UserNameFilter(uname) + ":" + id)
}

func ObjKeyAccChargeMgr(id string) []byte {
	return []byte(dataPrefix + ":" + dataAccChargeMgr + ":" + id)
}

func PrevDataMsgQueue(id string) []byte {
	return []byte(dataPrefix + ":" + dataMsgQueue + ":" + id)
}

func PrevDataMsgSent(id string) []byte {
	return []byte(dataPrefix + ":" + dataMsgSent + ":" + id)
}

func ObjKeyRole(name string) []byte {
	return []byte(dataPrefix + ":" + dataRole + ":" +
		UserNameFilter(name))
}

func ObjKeyRolePrivilege(rid uint32, appid string) []byte {
	return []byte(dataPrefix + ":" + dataRolePrivilege + ":" +
		utils.Uint32ToHexString(rid) + ":" + appid)
}

func ObjKeyMsgQueue(id string) []byte {
	return []byte(dataPrefix + ":" + dataMsgQueue + ":" + id)
}

func ObjKeyMsgSent(id string) []byte {
	return []byte(dataPrefix + ":" + dataMsgSent + ":" + id)
}
