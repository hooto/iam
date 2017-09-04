// Copyright 2014 lessos Authors, All rights reserved.
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
	"github.com/lynkdb/iomix/skv"

	"github.com/hooto/iam/config"
	"github.com/hooto/iam/iamapi"
	"github.com/hooto/iam/iamclient"
	"github.com/hooto/iam/store"
)

func (c UserMgr) RoleListAction() {

	ls := iamapi.UserRoleList{}

	defer c.RenderJson(&ls)

	// if !iamclient.SessionAccessAllowed(c.Session, "user.admin", config.Config.InstanceID) {
	// 	ls.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
	// 	return
	// }

	if objs := store.PoScan("role", []byte{}, []byte{}, 10000); objs.OK() {

		rss := objs.KvList()
		for _, obj := range rss {

			var role iamapi.UserRole
			if err := obj.Decode(&role); err == nil {

				if role.Id == 1000 {
					continue
				}

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

	if obj := store.PoGet("role", uint32(c.Params.Uint64("roleid"))); obj.OK() {
		obj.Decode(&set.UserRole)
	}

	if set.Id == 0 {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Role Not Found")
		return
	}

	set.Kind = "UserRole"
}

func (c UserMgr) RoleSetAction() {

	var (
		prev iamapi.UserRole
		set  struct {
			types.TypeMeta
			iamapi.UserRole
		}
		prevVersion uint64
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

	if set.Name == "" {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Bad Request")
		return
	}

	if set.Id == 0 {

		if objs := store.PoRevScan("role", []byte{}, []byte{}, 1); objs.OK() {

			rss := objs.KvList()
			for _, obj := range rss {

				var last_role iamapi.UserRole
				if err := obj.Decode(&last_role); err == nil {
					set.Id = last_role.Id + 1
					break
				}
			}
		}

		if set.Id == 0 {
			set.Error = types.NewErrorMeta("500", "Server Error")
			return
		}

		//
		set.Created = types.MetaTimeNow()
		set.User = "sysadmin"

	} else {

		if obj := store.PoGet("role", set.Id); obj.OK() {
			obj.Decode(&prev)
			prevVersion = obj.Meta().Version
		}

		if prev.Id != set.Id {
			set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "UserRole Not Found")
			return
		}

		prev.Name = set.Name
		prev.Desc = set.Desc
		set.UserRole = prev
	}

	set.Updated = types.MetaTimeNow()
	// roleset["privileges"] = strings.Join(c.Params.Values["privileges"], ",")

	if obj := store.PoPut("role", set.Id, set, &skv.PathWriteOptions{
		PrevVersion: prevVersion,
	}); !obj.OK() {
		set.Error = types.NewErrorMeta("500", obj.Bytex().String())
		return
	}

	set.Kind = "UserRole"
}
