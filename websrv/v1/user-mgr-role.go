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

	if !iamclient.SessionAccessAllowed(c.Session, "user.admin", config.Config.InstanceID) {
		ls.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	if objs := store.PvScan("role/", "", "", 10000); objs.OK() {

		rss := objs.KvList()
		for _, obj := range rss {

			var role iamapi.UserRole
			if err := obj.Decode(&role); err == nil {

				ls.Items = append(ls.Items, role)
			}
		}
	}

	ls.Kind = "UserRoleList"
}

func (c UserMgr) RoleEntryAction() {

	set := iamapi.UserRole{}

	defer c.RenderJson(&set)

	if !iamclient.SessionAccessAllowed(c.Session, "user.admin", config.Config.InstanceID) {
		set.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	if obj := store.PvGet("role/" + c.Params.Get("roleid")); obj.OK() {
		obj.Decode(&set)
	}

	if set.Meta.ID == "" {
		set.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "Role Not Found"}
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
		set.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "Bad Request"}
		return
	}

	if !iamclient.SessionAccessAllowed(c.Session, "user.admin", config.Config.InstanceID) {
		set.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	if set.Meta.ID == "" {

		//
		set.Meta.ID = utils.StringNewRand(8)
		set.Meta.Created = utilx.TimeNow("atom")
		set.Meta.UserID = utils.StringEncode16("sysadmin", 8)

	} else {

		if obj := store.PvGet("role/" + set.Meta.ID); obj.OK() {
			obj.Decode(&prev)
			prevVersion = obj.Meta().Version
		}

		if prev.Meta.ID != set.Meta.ID {
			set.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "UserRole Not Found"}
			return
		}

		set.Meta.Created = prev.Meta.Updated
		set.Meta.UserID = prev.Meta.UserID
	}

	set.Meta.Updated = utilx.TimeNow("atom")
	// roleset["privileges"] = strings.Join(c.Params.Values["privileges"], ",")

	if obj := store.PvPut("role/"+set.Meta.ID, set, &skv.PvWriteOptions{
		PrevVersion: prevVersion,
	}); !obj.OK() {
		set.Error = &types.ErrorMeta{"500", obj.Bytex().String()}
		return
	}

	set.Kind = "UserRole"
}
