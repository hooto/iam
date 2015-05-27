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
	"github.com/lessos/lessgo/data/rdo"
	"github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/types"

	"../../idsapi"
)

func (c UserMgr) RoleListAction() {

	rsp := idsapi.UserRoleList{}

	defer c.RenderJson(&rsp)

	if !c.Session.AccessAllowed("user.admin") {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeUnavailable, "Service Unavailable"}
		return
	}

	q := base.NewQuerySet().From("ids_role").Limit(1000)
	if c.Params.Get("status") != "0" {
		q.Where.And("status", 1)
	}

	rsr, err := dcn.Base.Query(q)

	if err == nil && len(rsr) > 0 {
		for _, v := range rsr {
			rsp.Items = append(rsp.Items, idsapi.UserRole{
				Meta: types.ObjectMeta{
					ID:      v.Field("rid").String(),
					Name:    v.Field("name").String(),
					Created: v.Field("created").TimeFormat("datetime", "atom"),
					Updated: v.Field("updated").TimeFormat("datetime", "atom"),
				},
				Status: v.Field("status").Uint8(),
				Desc:   v.Field("desc").String(),
			})
		}
	}

	rsp.Kind = "UserRoleList"
}

func (c UserMgr) RoleEntryAction() {

	rsp := idsapi.UserRole{}

	defer c.RenderJson(&rsp)

	if !c.Session.AccessAllowed("user.admin") {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeUnavailable, "Service Unavailable"}
		return
	}

	q := base.NewQuerySet().From("ids_role").Limit(1)
	q.Where.And("rid", c.Params.Get("roleid"))

	rsr, err := dcn.Base.Query(q)
	if err != nil || len(rsr) != 1 {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	rsp.Meta.ID = rsr[0].Field("rid").String()
	rsp.Meta.Name = rsr[0].Field("name").String()
	rsp.Meta.Created = rsr[0].Field("created").TimeFormat("datetime", "atom")
	rsp.Meta.Updated = rsr[0].Field("updated").TimeFormat("datetime", "atom")
	rsp.Desc = rsr[0].Field("desc").String()
	rsp.Status = rsr[0].Field("status").Uint8()

	rsp.Kind = "UserRole"
}

func (c UserMgr) RoleSetAction() {

	rsp := idsapi.UserRole{}

	defer c.RenderJson(&rsp)

	if err := c.Request.JsonDecode(&rsp); err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "Bad Request"}
		return
	}

	if !c.Session.AccessAllowed("user.admin") {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeUnavailable, "Service Unavailable"}
		return
	}

	q := base.NewQuerySet().From("ids_role").Limit(1)

	isNew := true
	roleset := map[string]interface{}{}

	if rsp.Meta.ID != "" {

		q.Where.And("rid", rsp.Meta.ID)

		rsrole, err := dcn.Base.Query(q)
		if err != nil || len(rsrole) == 0 {
			rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "UserRole Not Found"}
			return
		}

		isNew = false
	}

	roleset["updated"] = base.TimeNow("datetime")
	roleset["name"] = rsp.Meta.Name
	roleset["desc"] = rsp.Desc
	roleset["status"] = rsp.Status
	// roleset["privileges"] = strings.Join(c.Params.Values["privileges"], ",")

	if isNew {

		si, err := c.Session.SessionFetch()
		if err != nil {
			return
		}

		roleset["created"] = base.TimeNow("datetime")
		roleset["uid"] = si.UserID

		_, err = dcn.Base.Insert("ids_role", roleset)
		if err != nil {
			rsp.Error = &types.ErrorMeta{idsapi.ErrCodeUnavailable, "Service Unavailable"}
			return
		}

	} else {

		frupd := base.NewFilter()
		frupd.And("rid", rsp.Meta.ID)
		if _, err := dcn.Base.Update("ids_role", roleset, frupd); err != nil {
			rsp.Error = &types.ErrorMeta{idsapi.ErrCodeUnavailable, "Service Unavailable"}
			return
		}
	}

	rsp.Kind = "UserRole"
}
