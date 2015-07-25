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

package v1

import (
	"github.com/lessos/bigtree/btapi"
	"github.com/lessos/lessgo/types"
	"github.com/lessos/lessgo/utils"
	"github.com/lessos/lessgo/utilx"

	"github.com/lessos/lessids/idsapi"
	"github.com/lessos/lessids/store"
)

func (c UserMgr) RoleListAction() {

	ls := idsapi.UserRoleList{}

	defer c.RenderJson(&ls)

	if !c.Session.AccessAllowed("user.admin") {
		ls.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	if objs := store.BtAgent.ObjectList(btapi.ObjectProposal{
		Meta: btapi.ObjectMeta{
			Path: "/role/",
		},
	}); objs.Error == nil {

		for _, obj := range objs.Items {

			var role idsapi.UserRole
			if err := obj.JsonDecode(&role); err == nil {

				ls.Items = append(ls.Items, role)
			}
		}
	}

	ls.Kind = "UserRoleList"
}

func (c UserMgr) RoleEntryAction() {

	set := idsapi.UserRole{}

	defer c.RenderJson(&set)

	if !c.Session.AccessAllowed("user.admin") {
		set.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	if obj := store.BtAgent.ObjectGet(btapi.ObjectProposal{
		Meta: btapi.ObjectMeta{
			Path: "/role/" + c.Params.Get("roleid"),
		},
	}); obj.Error == nil {
		obj.JsonDecode(&set)
	}

	if set.Meta.ID == "" {
		set.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "Role Not Found"}
		return
	}

	set.Kind = "UserRole"
}

func (c UserMgr) RoleSetAction() {

	var (
		prev        idsapi.UserRole
		set         idsapi.UserRole
		prevVersion uint64
	)

	defer c.RenderJson(&set)

	if err := c.Request.JsonDecode(&set); err != nil {
		set.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "Bad Request"}
		return
	}

	if !c.Session.AccessAllowed("user.admin") {
		set.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	if set.Meta.ID == "" {

		//
		set.Meta.ID = utils.StringNewRand(8)
		set.Meta.Created = utilx.TimeNow("atom")
		set.Meta.UserID = utils.StringEncode16("sysadmin", 8)

	} else {

		if obj := store.BtAgent.ObjectGet(btapi.ObjectProposal{
			Meta: btapi.ObjectMeta{
				Path: "/role/" + set.Meta.ID,
			},
		}); obj.Error == nil {
			obj.JsonDecode(&prev)
			prevVersion = obj.Meta.Version
		}

		if prev.Meta.ID != set.Meta.ID {
			set.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "UserRole Not Found"}
			return
		}

		set.Meta.Created = prev.Meta.Updated
		set.Meta.UserID = prev.Meta.UserID
	}

	set.Meta.Updated = utilx.TimeNow("atom")
	// roleset["privileges"] = strings.Join(c.Params.Values["privileges"], ",")

	setjs, _ := utils.JsonEncode(set)

	if obj := store.BtAgent.ObjectSet(btapi.ObjectProposal{
		Meta: btapi.ObjectMeta{
			Path: "/role/" + set.Meta.ID,
		},
		Data:        setjs,
		PrevVersion: prevVersion,
	}); obj.Error != nil {
		set.Error = &types.ErrorMeta{"500", obj.Error.Message}
		return
	}

	set.Kind = "UserRole"
}
