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

package user

import (
	"github.com/hooto/httpsrv"
	"github.com/sysinner/innerstack/v2/pkg/inauth"
)

// NewModule creates an httpsrv.Module for user, access-key and app-auth API endpoints.
func NewModule() *httpsrv.Module {
	mod := httpsrv.NewModule()

	// user profile
	mod.RegisterAction("/profile", ProfileEntry)
	mod.RegisterAction("/profile-set", ProfileSet)
	mod.RegisterAction("/pass-set", PassSet)
	mod.RegisterAction("/email-set", EmailSet)
	mod.RegisterAction("/photo-set", ProfilePhotoSet)

	// access key
	mod.RegisterAction("/keys/list", AccessKeyList)
	mod.RegisterAction("/keys/entry", AccessKeyEntry)
	mod.RegisterAction("/keys/set", AccessKeySet)
	mod.RegisterAction("/keys/delete", AccessKeyDelete)
	mod.RegisterAction("/keys/bind", AccessKeyBind)
	mod.RegisterAction("/keys/unbind", AccessKeyUnbind)

	// app auth
	mod.RegisterAction("/apps/register", AppAuth_Register)
	mod.RegisterAction("/apps/list", AppAuth_List)
	mod.RegisterAction("/apps/update", AppAuth_Update)
	mod.RegisterAction("/apps/delete", AppAuth_Delete)

	return mod
}

type StatusResponse struct {
	Status inauth.ServiceStatus `json:"status"`
}

func NewStatusResponse(code, message string) *StatusResponse {
	return &StatusResponse{
		Status: inauth.NewServiceStatus(code, message),
	}
}
