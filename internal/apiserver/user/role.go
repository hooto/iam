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

	"github.com/hooto/iam/v2/internal/data"
	"github.com/hooto/iam/v2/pkg/iamapi"
)

type UserRoleListResponse struct {
	Status inauth.ServiceStatus `json:"status"`
	Items  []iamapi.UserRole    `json:"items,omitempty"`
}

// RoleList returns roles visible to the current user.
func RoleList(ctx httpsrv.Ctx) error {

	u := authCtx(ctx)
	if u == nil {
		return nil
	}

	var rsp UserRoleListResponse
	defer ctx.JSON(&rsp)

	if rs := data.Data.NewRanger(
		iamapi.NsRole(""), iamapi.NsRole("")).SetLimit(1000).Exec(); rs.OK() {

		for _, obj := range rs.Items {

			var role iamapi.UserRole
			if err := obj.JsonDecode(&role); err != nil {
				continue
			}

			if role.Name != "" || role.User == u.Name {
				rsp.Items = append(rsp.Items, role)
			}
		}
	}

	rsp.Status = inauth.NewServiceStatus("200", "ok")
	return nil
}
