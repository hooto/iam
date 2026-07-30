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
	"github.com/hooto/httpsrv/v2"
	"github.com/sysinner/innerstack/v2/pkg/inauth"
)

// Register wires the user, access-key and app-auth API endpoints onto the
// given router.
func Register(r httpsrv.Router) {

	// user profile
	r.All("/profile", ProfileEntry)
	r.All("/profile-set", ProfileSet)
	r.All("/pass-set", PassSet)
	r.All("/email-set", EmailSet)
	r.All("/photo-set", ProfilePhotoSet)

	// access key
	r.All("/keys/list", AccessKeyList)
	r.All("/keys/entry", AccessKeyEntry)
	r.All("/keys/set", AccessKeySet)
	r.All("/keys/delete", AccessKeyDelete)
	r.All("/keys/bind", AccessKeyBind)
	r.All("/keys/unbind", AccessKeyUnbind)

	// app auth
	r.All("/apps/register", AppAuth_Register)
	r.All("/apps/list", AppAuth_List)
	r.All("/apps/update", AppAuth_Update)
	r.All("/apps/delete", AppAuth_Delete)
}

type StatusResponse struct {
	Status inauth.ServiceStatus `json:"status"`
}

func NewStatusResponse(code, message string) *StatusResponse {
	return &StatusResponse{
		Status: inauth.NewServiceStatus(code, message),
	}
}
