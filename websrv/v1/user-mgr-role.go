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

package v1

import (
	"github.com/lessos/lessgo/types"

	"github.com/hooto/iam/config"
	"github.com/hooto/iam/data"
	"github.com/hooto/iam/iamapi"
	"github.com/hooto/iam/iamclient"
)

func (c UserMgr) RoleListAction() {

	ls := iamapi.UserRoleList{}

	defer c.RenderJson(&ls)

	// if !iamclient.SessionAccessAllowed(c.Session, "user.admin", config.Config.InstanceID) {
	// 	ls.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
	// 	return
	// }

	if rs := data.Data.NewReader(nil).KeyRangeSet(
		iamapi.ObjKeyRole(""), iamapi.ObjKeyRole("")).LimitNumSet(1000).Query(); rs.OK() {

		for _, obj := range rs.Items {

			var role iamapi.UserRole
			if err := obj.DataValue().Decode(&role, nil); err == nil {

				if obj.Meta.IncrId == 1000 {
					continue
				}

				role.Id = uint32(obj.Meta.IncrId)
				ls.Items = append(ls.Items, role)
			}
		}
	}

	ls.Kind = "UserRoleList"
}

func (c UserMgr) RoleEntryAction() {

	var set struct {
		types.TypeMeta
		iamapi.UserRole
	}
	defer c.RenderJson(&set)

	// if !iamclient.SessionAccessAllowed(c.Session, "user.admin", config.Config.InstanceID) {
	// 	set.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
	// 	return
	// }

	// TODO roleid
	name := c.Params.Get("role_name")
	if rs := data.Data.NewReader(iamapi.ObjKeyRole(name)).Query(); rs.OK() {
		rs.DataValue().Decode(&set.UserRole, nil)
	}

	if set.Name == "" {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Role Not Found")
		return
	}

	set.Kind = "UserRole"
}

func (c UserMgr) RoleSetAction() {

	var (
		set struct {
			types.TypeMeta
			iamapi.UserRole
		}
	)
	defer c.RenderJson(&set)

	if err := c.Request.JsonDecode(&set.UserRole); err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Bad Request")
		return
	}

	if !iamclient.SessionAccessAllowed(c.Session, "user.admin", config.Config.InstanceID) {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
		return
	}

	if !iamapi.UsernameRE.MatchString(set.Name) {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Invalid Role Name")
		return
	}

	rsp := data.Data.NewReader(iamapi.ObjKeyRole(set.Name)).Query()

	if rsp.NotFound() {

		set.Created = types.MetaTimeNow()
		set.User = "sysadmin"

	} else if rsp.OK() {

		var prev iamapi.UserRole
		rsp.DataValue().Decode(&prev, nil)

		if prev.Created > 0 {
			prev.Desc = set.Desc
			set.UserRole = prev
		}

	} else {
		set.Error = types.NewErrorMeta("500", rsp.Message)
		return
	}

	set.Updated = types.MetaTimeNow()
	// roleset["privileges"] = strings.Join(c.Params.Values["privileges"], ",")

	if rs := data.Data.NewWriter(iamapi.ObjKeyRole(set.Name), set.UserRole).
		IncrNamespaceSet("role").Commit(); !rs.OK() {
		set.Error = types.NewErrorMeta("500", rs.Message)
		return
	}

	set.Kind = "UserRole"
}
