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

//go:generate protoc --go_out=plugins=grpc:. types.proto
//go:generate protobuf_slice "*.proto"

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"hash/crc32"
	"regexp"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/lessos/lessgo/crypto/idhash"
	"github.com/lessos/lessgo/encoding/json"
	"github.com/lessos/lessgo/types"
)

const (
	AccessTokenKey = "_iam_at"
	UserKeyDefault = "std"

	ErrCodeAccessDenied    = "AccessDenied"
	ErrCodeUnauthorized    = "Unauthorized" // Need to login and fetch a new access-token
	ErrCodeInvalidArgument = "InvalidArgument"
	ErrCodeUnavailable     = "Unavailable"
	ErrCodeServerError     = "ServerError"
	ErrCodeInternalError   = "InternalError"
	ErrCodeNotFound        = "NotFound"
	ErrCodeAccChargeOut    = "AccChargeOut"
)

var (
	UserNameRe2            = regexp.MustCompile("^[a-z]{1}[a-z0-9]{3,29}$")
	UserRoleNameRe2        = regexp.MustCompile("^[a-z]{1}[a-z0-9]{3,19}$")
	UserEmailRe2           = regexp.MustCompile("^[_a-z0-9-]+(\\.[_a-z0-9-]+)*@[a-z0-9-]+(\\.[a-z0-9-]+)*(\\.[a-z]{2,10})$")
	accessTokenFrontendRe2 = regexp.MustCompile("^[a-z0-9]{4,30}\\/[a-f0-9]{20,40}$")
	AppInstanceIdReg       = regexp.MustCompile("^[a-f0-9]{16,24}$")
)

func UserNameFilter(name string) string {
	name = strings.ToLower(name)
	name2 := ""
	for _, v := range name {
		if (v >= 'a' && v <= 'z') || (v >= '0' || v <= '9') {
			name2 += string(v)
		}
	}
	return name2
}

func UserIdBytes(name string) []byte {
	return idhash.Hash([]byte(name), 4)
}

func UserId(name string) string {
	return hex.EncodeToString(UserIdBytes(name))
}

func Hash32(v string) uint32 {
	u32 := crc32.ChecksumIEEE([]byte(v))
	if u32 < 200 {
		u32 = 200
	}
	return u32
}

func UserNameValid(user string) error {
	if UserNameRe2.MatchString(user) {
		return nil
	}
	return errors.New("Invalid UserName")
}

type UserSession struct {
	AccessToken  string            `json:"access_token" toml:"access_token"`
	RefreshToken string            `json:"refresh_token,omitempty" toml:"refresh_token,omitempty"`
	UserName     string            `json:"username" toml:"username"`
	DisplayName  string            `json:"display_name,omitempty" toml:"display_name,omitempty"`
	Roles        types.ArrayUint32 `json:"roles,omitempty" toml:"roles,omitempty"`
	Groups       []string          `json:"groups,omitempty" toml:"groups,omitempty"`
	ClientAddr   string            `json:"client_addr,omitempty" toml:"client_addr,omitempty"`
	Created      int64             `json:"created" toml:"created"`
	Expired      int64             `json:"expired" toml:"expired"`
	Cached       int64             `json:"cached,omitempty" toml:"cached,omitempty"`
}

func (s *UserSession) IsLogin() bool {
	return (s.AccessToken != "")
}

func (s *UserSession) UserId() string {
	return UserId(s.UserName)
}

func (s *UserSession) AccessAllow(name string) bool {
	if name != "" {
		if name == s.UserName {
			return true
		}
		for _, v := range s.Groups {
			if name == v {
				return true
			}
		}
	}
	return false
}

type UserAccessEntry struct {
	types.TypeMeta `json:",inline" toml:",inline"`
	AccessToken    string `json:"access_token" toml:"access_token"`
	Privilege      string `json:"privilege" toml:"privilege"`
	InstanceID     string `json:"instance_id,omitempty" toml:"instance_id,omitempty"`
}

const (
	UserTypeGroup uint32 = 1 << 1
)

type User struct {
	// Id          string            `json:"id,omitempty" toml:"id,omitempty"`
	Name        string            `json:"name" toml:"name"`
	Email       string            `json:"email,omitempty" toml:"email,omitempty"`
	DisplayName string            `json:"display_name,omitempty" toml:"display_name,omitempty"`
	Keys        types.KvPairs     `json:"keys,omitempty" toml:"keys,omitempty"`
	Roles       types.ArrayUint32 `json:"roles,omitempty" toml:"roles,omitempty"`
	Type        uint32            `json:"type,omitempty" toml:"type,omitempty"`
	Owners      []string          `json:"owners,omitempty" toml:"owners,omitempty"`
	Members     []string          `json:"members,omitempty" toml:"members,omitempty"`
	Status      uint8             `json:"status,omitempty" toml:"status,omitempty"`
	Created     types.MetaTime    `json:"created,omitempty" toml:"created,omitempty"`
	Updated     types.MetaTime    `json:"updated,omitempty" toml:"updated,omitempty"`
}

type UserEntry struct {
	types.TypeMeta `json:",inline" toml:",inline"`
	Login          User         `json:"login,omitempty" toml:"login,omitempty"`
	Profile        *UserProfile `json:"profile,omitempty" toml:"profile,omitempty"`
}

type UserList struct {
	types.TypeMeta `json:",inline" toml:",inline"`
	Meta           types.ListMeta `json:"meta,omitempty" toml:"meta,omitempty"`
	Items          []User         `json:"items,omitempty" toml:"items,omitempty"`
}

type UserProfile struct {
	Login       *User          `json:"login,omitempty" toml:"login,omitempty"`
	Gender      uint8          `json:"gender,omitempty" toml:"gender,omitempty"`
	Photo       string         `json:"photo,omitempty" toml:"photo,omitempty"`
	PhotoSource string         `json:"photo_source,omitempty" toml:"photo_source,omitempty"`
	Birthday    string         `json:"birthday,omitempty" toml:"birthday,omitempty"`
	About       string         `json:"about,omitempty" toml:"about,omitempty"`
	Updated     types.MetaTime `json:"updated,omitempty" toml:"updated,omitempty"`
}

type UserGroupItem struct {
	types.TypeMeta `json:",inline" toml:",inline"`
	Name           string         `json:"name" toml:"name"`
	DisplayName    string         `json:"display_name,omitempty" toml:"display_name,omitempty"`
	Owners         []string       `json:"owners,omitempty" toml:"owners,omitempty"`
	Members        []string       `json:"members,omitempty" toml:"members,omitempty"`
	Status         uint8          `json:"status" toml:"status"`
	Created        types.MetaTime `json:"created" toml:"created"`
	Updated        types.MetaTime `json:"updated" toml:"updated"`
}

type UserGroupList struct {
	types.TypeMeta `json:",inline" toml:",inline"`
	Items          []*UserGroupItem `json:"items,omitempty" toml:"items,omitempty"`
}

type UserPasswordSet struct {
	types.TypeMeta  `json:",inline" toml:",inline"`
	CurrentPassword string `json:"current_password,omitempty" toml:"current_password,omitempty"`
	NewPassword     string `json:"new_password,omitempty" toml:"new_password,omitempty"`
}

type UserPasswordReset struct {
	types.TypeMeta `json:",inline" toml:",inline"`
	Id             string `json:"id,omitempty" toml:"id,omitempty"`
	UserName       string `json:"username,omitempty" toml:"username,omitempty"`
	Email          string `json:"email,omitempty" toml:"email,omitempty"`
	Expired        string `json:"expired,omitempty" toml:"expired,omitempty"`
}

type UserEmailSet struct {
	types.TypeMeta `json:",inline" toml:",inline"`
	Auth           string `json:"auth,omitempty" toml:"auth,omitempty"`
	Email          string `json:"email,omitempty" toml:"email,omitempty"`
}

type UserPhotoSet struct {
	types.TypeMeta `json:",inline" toml:",inline"`
	Name           string `json:"name,omitempty" toml:"name,omitempty"`
	Size           int    `json:"size,omitempty" toml:"size,omitempty"`
	Data           string `json:"data,omitempty" toml:"data,omitempty"`
}

type UserRole struct {
	Id         uint32         `json:"id" toml:"id"`
	Name       string         `json:"name" toml:"name"`
	User       string         `json:"user,omitempty" toml:"user,omitempty"`
	Status     uint8          `json:"status,omitempty" toml:"status,omitempty"`
	Desc       string         `json:"desc,omitempty" toml:"desc,omitempty"`
	Privileges []string       `json:"privileges,omitempty" toml:"privileges,omitempty"`
	Created    types.MetaTime `json:"created,omitempty" toml:"created,omitempty"`
	Updated    types.MetaTime `json:"updated,omitempty" toml:"updated,omitempty"`
}

type UserRoleList struct {
	types.TypeMeta `json:",inline" toml:",inline"`
	Items          []UserRole `json:"items,omitempty" toml:"items,omitempty"`
}

type UserPrivilege struct {
	types.TypeMeta `json:",inline" toml:",inline"`
	Meta           types.InnerObjectMeta `json:"meta,omitempty" toml:"meta,omitempty"`
	Instance       string                `json:"instance" toml:"instance"`
	Desc           string                `json:"desc,omitempty" toml:"desc,omitempty"`
}

type UserPrivilegeList struct {
	types.TypeMeta `json:",inline" toml:",inline"`
	Items          []UserPrivilege `json:"items,omitempty" toml:"items,omitempty"`
}

type AppPrivilege struct {
	// ID        uint32   `json:"id,omitempty" toml:"id,omitempty"`
	Privilege string            `json:"privilege" toml:"privilege"`
	Desc      string            `json:"desc,omitempty" toml:"desc,omitempty"`
	Roles     types.ArrayUint32 `json:"roles,omitempty" toml:"roles,omitempty"`
	// ExtRoles  types.ArrayUint32 `json:"extroles,omitempty" toml:"extroles,omitempty"`
}

type AppInstance struct {
	Meta       types.InnerObjectMeta `json:"meta,omitempty" toml:"meta,omitempty"`
	AppID      string                `json:"app_id,omitempty" toml:"app_id,omitempty"`
	AppTitle   string                `json:"app_title,omitempty" toml:"app_title,omitempty"`
	Version    string                `json:"version,omitempty" toml:"version,omitempty"`
	Status     uint8                 `json:"status,omitempty" toml:"status,omitempty"`
	Url        string                `json:"url,omitempty" toml:"url,omitempty"`
	Privileges []AppPrivilege        `json:"privileges,omitempty" toml:"privileges,omitempty"`
	SecretKey  string                `json:"secret_key,omitempty" toml:"secret_key,omitempty"`
}

type AppInstanceList struct {
	types.TypeMeta `json:",inline" toml:",inline"`
	Meta           types.ListMeta `json:"meta,omitempty" toml:"meta,omitempty"`
	Items          []AppInstance  `json:"items,omitempty" toml:"items,omitempty"`
}

type AppAuthInfo struct {
	types.TypeMeta `json:",inline" toml:",inline"`
	InstanceID     string `json:"instance_id" toml:"instance_id"`
	AppID          string `json:"app_id" toml:"app_id"`
	// Version        string `json:"version,omitempty" toml:"version,omitempty"`
}

type AppInstanceRegister struct {
	types.TypeMeta `json:",inline" toml:",inline"`
	AccessToken    string      `json:"access_token,omitempty" toml:"access_token,omitempty"`
	Instance       AppInstance `json:"instance" toml:"instance"`
}

type SysConfigList struct {
	types.TypeMeta `json:",inline" toml:",inline"`
	Items          types.Labels `json:"items,omitempty" toml:"items,omitempty"`
}

type SysConfigMailer struct {
	SmtpHost string `json:"smtp_host" toml:"smtp_host"`
	SmtpPort string `json:"smtp_port" toml:"smtp_port"`
	SmtpUser string `json:"smtp_user" toml:"smtp_user"`
	SmtpPass string `json:"smtp_pass" toml:"smtp_pass"`
}

type ServiceLoginAuth struct {
	types.TypeMeta `json:",inline" toml:",inline"`
	RedirectUri    string `json:"redirect_uri,omitempty" toml:"redirect_uri,omitempty"`
	AccessToken    string `json:"access_token,omitempty" toml:"access_token,omitempty"`
}

type ServiceRedirectToken struct {
	RedirectUri string `json:"uri,omitempty" toml:"uri,omitempty"`
	State       string `json:"state,omitempty" toml:"state,omitempty"`
	ClientId    string `json:"cid,omitempty" toml:"cid,omitempty"`
	Persistent  int    `json:"pt,omitempty" toml:"pt,omitempty"`
}

func ServiceRedirectTokenValid(tokenstr string) bool {
	if _, err := base64.StdEncoding.DecodeString(tokenstr); err == nil {
		return true
	}
	return false
}

func (s *ServiceRedirectToken) Encode() string {

	if len(s.RedirectUri) > 200 {
		s.RedirectUri = s.RedirectUri[:200]
	}
	if len(s.State) > 100 {
		s.State = s.State[:100]
	}
	if len(s.ClientId) > 40 {
		s.ClientId = s.ClientId[:40]
	}

	js, _ := json.Encode(s, "")
	return base64.StdEncoding.EncodeToString(js)
}

func ServiceRedirectTokenDecode(tokenstr string) ServiceRedirectToken {
	var token ServiceRedirectToken
	if jsb, err := base64.StdEncoding.DecodeString(tokenstr); err == nil {
		json.Decode(jsb, &token)
	}
	return token
}

func Uint32ToHexString(v uint32) string {
	return BytesToHexString(Uint32ToBytes(v))
}

func Uint64ToHexString(v uint64) string {
	return BytesToHexString(Uint64ToBytes(v))
}

func Uint32ToBytes(v uint32) []byte {
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, v)
	return bs
}

func Uint64ToBytes(v uint64) []byte {
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, v)
	return bs
}

func BytesToHexString(bs []byte) string {
	return hex.EncodeToString(bs)
}

func HexStringToBytes(s string) []byte {
	dec, err := hex.DecodeString(s)
	if err != nil {
		return []byte{}
	}
	return dec
}

func OpActionAllow(opbase, op uint32) bool {
	return (op & opbase) == op
}

func OpActionRemove(opbase, op uint32) uint32 {
	return (opbase | op) - (op)
}

func OpActionAppend(opbase, op uint32) uint32 {
	return (opbase | op)
}

type WebServiceKind struct {
	Kind  string           `json:"kind" toml:"kind"`
	Error *types.ErrorMeta `json:"error,omitempty" toml:"error,omitempty"`
	Data  proto.Message    `json:"data,omitempty" toml:"data,omitempty"`
}
