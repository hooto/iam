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

package apiserver

import (
	"github.com/hooto/httpsrv"
	"github.com/sysinner/incore/v2/pkg/inauth"

	"github.com/hooto/iam/v2/internal/data"
	"github.com/hooto/iam/v2/pkg/iamapi"
)

type AdminAccessKey struct {
	*httpsrv.Controller
}

type AdminAccessKeyListRequest struct {
	AccessToken string `json:"access_token,omitempty"`
}

type AdminAccessKeyListResponse struct {
	Status inauth.ServiceStatus `json:"status"`
	Items  []*inauth.AccessKey  `json:"items,omitempty"`
}

func (c AdminAccessKey) ListAction() {

	// user := userAuth(c.Controller)
	// if user == nil {
	// 	return
	// }

	var rsp AdminAccessKeyListResponse
	defer c.RenderJson(&rsp)

	k1 := iamapi.NsAccessKey("", "")
	k2 := iamapi.NsAccessKey("", "")
	if rs := data.Data.NewRanger(k1, k2).
		SetLimit(int64(accessKeyLimit)).Exec(); rs.OK() {
		for _, v := range rs.Items {
			var ak inauth.AccessKey
			if err := v.JsonDecode(&ak); err == nil &&
				ak.State != inauth.AccessKey_State_Disable {
				rsp.Items = append(rsp.Items, &ak)
			}
		}
	}

	rsp.Status = inauth.NewServiceStatus("200", "ok")
}
