// Copyright 2014 Eryx <evorui at gmail dot com>, All rights reserved.
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

import "fmt"

const (
	Role_Sysadmin  = "sa"
	Role_User      = "user"
	Role_Guest     = "guest"
	Role_Developer = "dev"
)

const (
	UserSysadmin    = "sysadmin"
	DefaultPassword = "changeme"
)

const (
	UserKeyDefault = "std"

	UserTypeGroup = "group"
)

var (
	ErrCodeAccessDenied    = "AccessDenied"
	ErrCodeUnauthorized    = "Unauthorized" // Need to login and fetch a new access-token
	ErrCodeInvalidArgument = "InvalidArgument"
	ErrCodeUnavailable     = "Unavailable"
	ErrCodeServerError     = "ServerError"
	ErrCodeInternalError   = "InternalError"
	ErrCodeNotFound        = "NotFound"
)

func NsUser(uname string) []byte {
	return []byte(fmt.Sprintf("iam/v2/u/%s", uname))
}

func NsUserProfile(uname string) []byte {
	return []byte(fmt.Sprintf("iam/v2/up/%s", uname))
}

func NsUserAuthDeny(uname, remoteIp string) []byte {
	return []byte(fmt.Sprintf("iam/v2/uad/%s/%s", uname, remoteIp))
}

func NsUserSession(uname string, created uint32) []byte {
	return []byte(fmt.Sprintf("iam/v2/us/%s/%012d", uname, created))
}

func NsUserResetPassword(id string) []byte {
	return []byte(fmt.Sprintf("iam/v2/upr/%s", id))
}

func NsAccessKey(uname, id string) []byte {
	if id == "" {
		return []byte(fmt.Sprintf("iam/v2/ak/%s", uname))
	}
	return []byte(fmt.Sprintf("iam/v2/ak/%s/%s", uname, id))
}

func NsAppInstance(id string) []byte {
	return []byte(fmt.Sprintf("iam/v2/app/%s", id))
}

func NsRole(name string) []byte {
	return []byte(fmt.Sprintf("iam/v2/role/%s", name))
}

func NsRolePrivilege(role, appid string) []byte {
	return []byte(fmt.Sprintf("iam/role/%s/%s", role, appid))
}

func NsSysConfig(name string) []byte {
	return []byte(fmt.Sprintf("iam/v2/sysconfig/%s", name))
}

func NsAuthCode(code string) []byte {
	return []byte(fmt.Sprintf("iam/v2/ac/%s", code))
}

func NsMsgQueue(id string) []byte {
	return []byte(fmt.Sprintf("iam/v2/msg/queue/%s", id))
}

func NsMsgSent(id string) []byte {
	return []byte(fmt.Sprintf("iam/v2/msg/done/%s", id))
}
