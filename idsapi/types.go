// Copyright 2015 lessOS.com, All rights reserved.
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

package idsapi

import (
	"time"

	"github.com/lessos/lessgo/types"
)

const (
	ErrCodeAccessDenied    = "AccessDenied"
	ErrCodeUnauthorized    = "Unauthorized" // Need to login and fetch a new access-token
	ErrCodeInvalidArgument = "InvalidArgument"
	ErrCodeUnavailable     = "Unavailable"
	ErrCodeInternalError   = "InternalError"
	ErrCodeNotFound        = "NotFound"
)

type ServiceLoginAuth struct {
	types.TypeMeta `json:",inline"`
	Continue       string `json:"continue,omitempty"`
	AccessToken    string `json:"access_token,omitempty"`
}

type UserSession struct {
	types.TypeMeta `json:",inline"`
	AccessToken    string    `json:"access_token"`
	RefreshToken   string    `json:"refresh_token,omitempty"`
	UserID         string    `json:"userid"`
	UserName       string    `json:"username"`
	ClientAddr     string    `json:"client_addr,omitempty"`
	Name           string    `json:"name"`
	Data           string    `json:"data,omitempty"`
	Roles          []uint32  `json:"roles,omitempty"`
	Groups         []uint32  `json:"groups,omitempty"`
	InnerExpired   time.Time `json:"inner_expired,omitempty"`
	Timezone       string    `json:"timezone"`
	Source         string    `json:"source,omitempty"`
	Created        string    `json:"created"`
	Expired        string    `json:"expired"`
}

func (s *UserSession) IsLogin() bool {
	return (s.UserID != "")
}

type UserAccessEntry struct {
	types.TypeMeta `json:",inline"`
	AccessToken    string `json:"access_token"`
	Privilege      string `json:"privilege"`
	InstanceID     string `json:"instance_id,omitempty"`
}

type User struct {
	types.TypeMeta `json:",inline"`
	Meta           types.ObjectMeta `json:"meta,omitempty"`
	Email          string           `json:"email,omitempty"`
	Name           string           `json:"name,omitempty"`
	Auth           string           `json:"auth,omitempty"`
	Timezone       string           `json:"timezone,omitempty"`
	Roles          []uint32         `json:"roles,omitempty"`
	Groups         []uint32         `json:"groups,omitempty"`
	Status         uint8            `json:"status"`
	Profile        *UserProfile     `json:"profile,omitempty"`
}

type UserList struct {
	types.TypeMeta `json:",inline"`
	Meta           types.ListMeta `json:"meta,omitempty"`
	Items          []User         `json:"items,omitempty"`
}

type UserProfile struct {
	types.TypeMeta `json:",inline"`
	Login          User   `json:"login,omitempty"`
	Gender         uint8  `json:"gender,omitempty"`
	Photo          string `json:"photo,omitempty"`
	PhotoSource    string `json:"photoSource,omitempty"`
	Name           string `json:"name,omitempty"`
	Birthday       string `json:"birthday,omitempty"`
	About          string `json:"about,omitempty"`
}

type UserPasswordSet struct {
	types.TypeMeta  `json:",inline"`
	CurrentPassword string `json:"currentPassword,omitempty"`
	NewPassword     string `json:"newPassword,omitempty"`
}

type UserPasswordReset struct {
	types.TypeMeta `json:",inline"`
	ID             string `json:"id,omitempty"`
	UserID         string `json:"userid,omitempty"`
	Email          string `json:"email,omitempty"`
	Expired        string `json:"expired,omitempty"`
}

type UserEmailSet struct {
	types.TypeMeta `json:",inline"`
	Auth           string `json:"auth,omitempty"`
	Email          string `json:"email,omitempty"`
}

type UserPhotoSet struct {
	types.TypeMeta `json:",inline"`
	Name           string `json:"name,omitempty"`
	Size           int    `json:"size,omitempty"`
	Data           string `json:"data,omitempty"`
}

type UserRole struct {
	types.TypeMeta `json:",inline"`
	Meta           types.ObjectMeta `json:"meta,omitempty"`
	IdxID          uint32           `json:"idxid"`
	Status         uint8            `json:"status"`
	Desc           string           `json:"desc,omitempty"`
	Privileges     []string         `json:"privileges,omitempty"`
}

type UserRoleList struct {
	types.TypeMeta `json:",inline"`
	Items          []UserRole `json:"items,omitempty"`
}

type UserPrivilege struct {
	types.TypeMeta `json:",inline"`
	Meta           types.ObjectMeta `json:"meta,omitempty"`
	Instance       string           `json:"instance"`
	Desc           string           `json:"desc,omitempty"`
}

type UserPrivilegeList struct {
	types.TypeMeta `json:",inline"`
	Items          []UserPrivilege `json:"items,omitempty"`
}

type AppPrivilege struct {
	// ID        uint32   `json:"id,omitempty"`
	Privilege string   `json:"privilege"`
	Desc      string   `json:"desc,omitempty"`
	Roles     []uint32 `json:"roles,omitempty"`
	// ExtRoles  []uint32 `json:"extroles,omitempty"`
}

type AppInstance struct {
	types.TypeMeta `json:",inline"`
	Meta           types.ObjectMeta `json:"meta,omitempty"`
	AppID          string           `json:"app_id,omitempty"`
	AppTitle       string           `json:"app_title,omitempty"`
	Version        string           `json:"version,omitempty"`
	Status         uint8            `json:"status,omitempty"`
	Url            string           `json:"url,omitempty"`
	Privileges     []AppPrivilege   `json:"privileges,omitempty"`
}

type AppInstanceList struct {
	types.TypeMeta `json:",inline"`
	Meta           types.ListMeta `json:"meta,omitempty"`
	Items          []AppInstance  `json:"items,omitempty"`
}

type AppAuthInfo struct {
	types.TypeMeta `json:",inline"`
	InstanceID     string `json:"instance_id"`
	AppID          string `json:"app_id"`
	// Version    string `json:"version"`
}

type AppInstanceRegister struct {
	types.TypeMeta `json:",inline"`
	AccessToken    string      `json:"access_token,omitempty"`
	AccessKey      AccessKey   `json:"access_key,omitempty"`
	Instance       AppInstance `json:"instance"`
}

type SysConfigList struct {
	types.TypeMeta `json:",inline"`
	Items          types.LabelListMeta `json:"items,omitempty"`
}

type SysConfigMailer struct {
	SmtpHost string `json:"smtp_host"`
	SmtpPort string `json:"smtp_port"`
	SmtpUser string `json:"smtp_user"`
	SmtpPass string `json:"smtp_pass"`
}

type AccessKey struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key,omitempty"`
}
