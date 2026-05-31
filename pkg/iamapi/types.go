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

import "github.com/lessos/lessgo/types"

type User struct {
	Name        string        `json:"name" toml:"name"`
	Email       string        `json:"email,omitempty" toml:"email,omitempty"`
	DisplayName string        `json:"display_name,omitempty" toml:"display_name,omitempty"`
	Keys        types.KvPairs `json:"keys,omitempty" toml:"keys,omitempty"`
	Roles       []string      `json:"roles,omitempty" toml:"roles,omitempty"`
	Type        string        `json:"type,omitempty" toml:"type,omitempty"`
	Owners      []string      `json:"owners,omitempty" toml:"owners,omitempty"`
	Members     []string      `json:"members,omitempty" toml:"members,omitempty"`
	Status      uint8         `json:"status,omitempty" toml:"status,omitempty"`
	Created     int64         `json:"created,omitempty" toml:"created,omitempty"`
	Updated     int64         `json:"updated,omitempty" toml:"updated,omitempty"`
}

type UserProfile struct {
	Login       *User  `json:"login,omitempty" toml:"login,omitempty"`
	DisplayName string `json:"display_name,omitempty" toml:"display_name,omitempty"`
	Email       string `json:"email,omitempty" toml:"email,omitempty"`
	Gender      uint8  `json:"gender,omitempty" toml:"gender,omitempty"`
	Birthday    string `json:"birthday,omitempty" toml:"birthday,omitempty"`
	About       string `json:"about,omitempty" toml:"about,omitempty"`
	Photo       string `json:"photo,omitempty" toml:"photo,omitempty"`
	PhotoSource string `json:"photo_source,omitempty" toml:"photo_source,omitempty"`
	Updated     int64  `json:"updated,omitempty" toml:"updated,omitempty"`
}

type UserRole struct {
	Name       string   `json:"name" toml:"name"`
	User       string   `json:"user,omitempty" toml:"user,omitempty"`
	Status     uint8    `json:"status,omitempty" toml:"status,omitempty"`
	Desc       string   `json:"desc,omitempty" toml:"desc,omitempty"`
	Privileges []string `json:"privileges,omitempty" toml:"privileges,omitempty"`
	Created    int64    `json:"created,omitempty" toml:"created,omitempty"`
	Updated    int64    `json:"updated,omitempty" toml:"updated,omitempty"`
}

type AppInstance struct {
	ID         string         `json:"id,omitempty" toml:"id,omitempty"`
	Name       string         `json:"name" toml:"name"`
	User       string         `json:"user,omitempty" toml:"user,omitempty"`
	Version    string         `json:"version,omitempty" toml:"version,omitempty"`
	Status     uint8          `json:"status,omitempty" toml:"status,omitempty"`
	Url        string         `json:"url,omitempty" toml:"url,omitempty"`
	Privileges []AppPrivilege `json:"privileges,omitempty" toml:"privileges,omitempty"`
	SecretKey  string         `json:"secret_key,omitempty" toml:"secret_key,omitempty"`
	Created    int64          `json:"created,omitempty" toml:"created,omitempty"`
	Updated    int64          `json:"updated,omitempty" toml:"updated,omitempty"`
}

type AppPrivilege struct {
	Privilege string   `json:"privilege" toml:"privilege"`
	Desc      string   `json:"desc,omitempty" toml:"desc,omitempty"`
	Roles     []string `json:"roles,omitempty" toml:"roles,omitempty"`
}

type SysConfigMailer struct {
	SmtpHost string `json:"smtp_host" toml:"smtp_host"`
	SmtpPort string `json:"smtp_port" toml:"smtp_port"`
	SmtpUser string `json:"smtp_user" toml:"smtp_user"`
	SmtpPass string `json:"smtp_pass" toml:"smtp_pass"`
}

type UserResetPassword struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Expired  int64  `json:"expired"`
}

type SignInAuthCode struct {
	AccessToken string `json:"access_token"`
	Username    string `json:"username"`
}
