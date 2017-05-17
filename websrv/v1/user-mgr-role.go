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
	"code.hooto.com/lynkdb/iomix/skv"
	iox_utils "code.hooto.com/lynkdb/iomix/utils"
	"github.com/lessos/lessgo/types"
	"github.com/lessos/lessgo/utils"
	"github.com/lessos/lessgo/utilx"

	"code.hooto.com/lessos/iam/config"
	"code.hooto.com/lessos/iam/iamapi"
	"code.hooto.com/lessos/iam/iamclient"
	"code.hooto.com/lessos/iam/store"
)

func (c UserMgr) RoleListAction() {

	ls := iamapi.UserRoleList{}

	defer c.RenderJson(&ls)

	// if !iamclient.SessionAccessAllowed(c.Session, "user.admin", config.Config.InstanceID) {
	// 	ls.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
	// 	return
	// }

	if objs := store.PvScan("role/", "", "", 10000); objs.OK() {

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

	set := iamapi.UserRole{}

	defer c.RenderJson(&set)

	// if !iamclient.SessionAccessAllowed(c.Session, "user.admin", config.Config.InstanceID) {
	// 	set.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
	// 	return
	// }

	if obj := store.PvGet("role/" + c.Params.Get("roleid")); obj.OK() {
		obj.Decode(&set)
	}

	if set.Meta.ID == "" {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Role Not Found")
		return
	}

	set.Kind = "UserRole"
}

func (c UserMgr) RoleSetAction() {

	var (
		prev        iamapi.UserRole
		set         iamapi.UserRole
		prevVersion uint64
	)

	defer c.RenderJson(&set)

	if err := c.Request.JsonDecode(&set); err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Bad Request")
		return
	}

	if !iamclient.SessionAccessAllowed(c.Session, "user.admin", config.Config.InstanceID) {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
		return
	}

	if set.Meta.Name == "" {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Bad Request")
		return
	}

	if set.Meta.ID == "" {

		last_id := uint32(0)

		if objs := store.PvRevScan("role/", "", "", 1); objs.OK() {

			rss := objs.KvList()
			for _, obj := range rss {

				var last_role iamapi.UserRole
				if err := obj.Decode(&last_role); err == nil {
					last_id = last_role.Id
				}
			}
		}

		if last_id < 1000 {
			set.Error = types.NewErrorMeta("500", "Server Error")
			return
		}

		last_id++

		//
		set.Meta.ID = iox_utils.BytesToHexString(iox_utils.Uint32ToBytes(last_id))
		set.Id = last_id
		set.Meta.Created = utilx.TimeNow("atom")
		set.Meta.UserID = utils.StringEncode16("sysadmin", 8)

	} else {

		if obj := store.PvGet("role/" + set.Meta.ID); obj.OK() {
			obj.Decode(&prev)
			prevVersion = obj.Meta().Version
		}

		if prev.Meta.ID != set.Meta.ID {
			set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "UserRole Not Found")
			return
		}

		prev.Meta.Name = set.Meta.Name
		prev.Desc = set.Desc
		set = prev
	}

	set.Meta.Updated = utilx.TimeNow("atom")
	// roleset["privileges"] = strings.Join(c.Params.Values["privileges"], ",")

	if obj := store.PvPut("role/"+set.Meta.ID, set, &skv.PvWriteOptions{
		PrevVersion: prevVersion,
	}); !obj.OK() {
		set.Error = types.NewErrorMeta("500", obj.Bytex().String())
		return
	}

	set.Kind = "UserRole"
}
