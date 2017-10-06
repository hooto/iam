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
	dataPassReset     = "pr"
	dataSession       = "s"
	dataAccessKey     = "ak"
	dataRole          = "r"
	dataRolePrivilege = "rp"
	dataSysConfig     = "sc"
	dataUserProfile   = "up"
	dataAccUser       = "au"
	dataAccFundUser   = "af"
	dataAccFundMgr    = "afm"
	dataAccChargeUser = "ac"
	dataAccChargeMgr  = "acm"
)

func DataAppInstanceKey(id string) skv.ProgKey {
	if id == "" {
		return skv.NewProgKey(dataPrefix, dataAppInstance, []byte{})
	}
	return skv.NewProgKey(dataPrefix, dataAppInstance, utils.HexStringToBytes(id))
}

func DataUserKey(uname string) skv.ProgKey {
	if uname == "" {
		return skv.NewProgKey(dataPrefix, dataUser, []byte{})
	}
	return skv.NewProgKey(dataPrefix, dataUser, UserIdBytes(uname))
}

func DataUserProfileKey(uname string) skv.ProgKey {
	if uname == "" {
		return skv.NewProgKey(dataPrefix, dataUserProfile, []byte{})
	}
	return skv.NewProgKey(dataPrefix, dataUserProfile, UserIdBytes(uname))
}

func DataPasswordResetKey(id string) skv.ProgKey {
	return skv.NewProgKey(dataPrefix, dataPassReset, utils.HexStringToBytes(id))
}

func DataSessionKey(uname, id string) skv.ProgKey {
	return skv.NewProgKey(dataPrefix, dataSession,
		UserIdBytes(uname), utils.HexStringToBytes(id))
}

func DataAccessKeyKey(uname, id string) skv.ProgKey {
	return skv.NewProgKey(dataPrefix, dataAccessKey,
		UserIdBytes(uname), utils.HexStringToBytes(id))
}

func DataRoleKey(id uint32) skv.ProgKey {
	return skv.NewProgKey(dataPrefix, dataRole, utils.Uint32ToBytes(id))
}

func DataRolePrivilegeKey(rid uint32, inst string) skv.ProgKey {
	return skv.NewProgKey(dataPrefix, dataRolePrivilege,
		utils.Uint32ToBytes(rid), utils.HexStringToBytes(inst))
}

func DataSysConfigKey(name string) skv.ProgKey {
	return skv.NewProgKey(dataPrefix, dataSysConfig, name)
}

func DataAccUserKey(uname string) skv.ProgKey {
	if uname == "" {
		return skv.NewProgKey(dataPrefix, dataAccUser, []byte{})
	}
	return skv.NewProgKey(dataPrefix, dataAccUser, UserIdBytes(uname))
}

func DataAccFundUserKey(uname, id string) skv.ProgKey {
	return skv.NewProgKey(dataPrefix, dataAccFundUser,
		UserIdBytes(uname), utils.HexStringToBytes(id))
}

func DataAccFundMgrKey(id string) skv.ProgKey {
	return skv.NewProgKey(dataPrefix, dataAccFundMgr, utils.HexStringToBytes(id))
}

func DataAccChargeUserKey(uname, id string) skv.ProgKey {
	return skv.NewProgKey(dataPrefix, dataAccChargeUser,
		UserIdBytes(uname), utils.HexStringToBytes(id))
}

func DataAccChargeMgrKey(id string) skv.ProgKey {
	return skv.NewProgKey(dataPrefix, dataAccChargeMgr, utils.HexStringToBytes(id))
}
