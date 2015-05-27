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
	"strconv"
	"strings"

	"github.com/lessos/lessgo/data/rdo"
	"github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/types"

	"../../idsapi"
)

const (
	appMgrInstPageLimit = 100
)

type AppMgr struct {
	*httpsrv.Controller
}

func (c AppMgr) InstListAction() {

	rsp := idsapi.AppInstanceList{}

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

	users := []interface{}{}

	q := base.NewQuerySet().From("ids_instance").Limit(appMgrInstPageLimit)

	count, _ := dcn.Base.Count("ids_instance", q.Where)
	page := c.Params.Int64("page")
	if page < 1 {
		page = 1
	}
	if page > 1 {
		q.Offset(int64((page - 1) * appMgrInstPageLimit))
	}

	if rs, err := dcn.Base.Query(q); err == nil && len(rs) > 0 {

		for _, v := range rs {

			item := idsapi.AppInstance{
				Meta: types.ObjectMeta{
					ID:      v.Field("id").String(),
					Name:    v.Field("app_title").String(),
					Created: v.Field("created").TimeFormat("datetime", "atom"),
					Updated: v.Field("updated").TimeFormat("datetime", "atom"),
					UserID:  v.Field("uid").String(),
				},
				AppID:   v.Field("app_id").String(),
				Version: v.Field("version").String(),
				Status:  v.Field("status").Uint8(),
				// Privileges
			}

			uid := v.Field("uid").String()

			inArray := false
			for _, vuid := range users {
				if vuid == uid {
					inArray = true
					break
				}
			}

			if !inArray {
				users = append(users, uid)
			}

			rsp.Items = append(rsp.Items, item)
		}
	}

	rsp.Meta.TotalResults = uint64(count)
	rsp.Meta.StartIndex = uint64((page - 1) * appMgrInstPageLimit)
	rsp.Meta.ItemsPerList = uint64(appMgrInstPageLimit)

	rsp.Kind = "AppInstanceList"

	// // 	//
	// 	q = base.NewQuerySet().From("ids_login").Limit(1000)
	// 	q.Where.And("uid.in", users...)
	// 	rslogin, err := dcn.Base.Query(q)
	// 	if err == nil && len(rslogin) > 0 {
	// 		for _, v := range rslogin {

	// 			for k2, v2 := range ls {
	// 				if v2["uid"] == v.Field("uid").String() {
	// 					ls[k2]["uid_name"] = v.Field("name").String()
	// 				}
	// 			}
	// 		}
	// 	}

	// // 	c.Data["list"] = ls
}

func (c AppMgr) InstEntryAction() {

	rsp := idsapi.AppInstance{}

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

	q := base.NewQuerySet().From("ids_instance").Limit(1)
	q.Where.And("id", c.Params.Get("instid"))

	rsr, err := dcn.Base.Query(q)
	if err != nil || len(rsr) != 1 {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	rsp.Meta.ID = rsr[0].Field("rid").String()
	rsp.Meta.Name = rsr[0].Field("app_title").String()
	rsp.Meta.Created = rsr[0].Field("created").TimeFormat("datetime", "atom")
	rsp.Meta.Updated = rsr[0].Field("updated").TimeFormat("datetime", "atom")

	rsp.AppID = rsr[0].Field("app_id").String()
	rsp.Version = rsr[0].Field("version").String()
	rsp.Status = rsr[0].Field("status").Uint8()

	q = base.NewQuerySet().From("ids_privilege").Limit(500)
	q.Where.And("instance", c.Params.Get("instid"))

	if rspri, err := dcn.Base.Query(q); err == nil && len(rspri) > 0 {

		for _, v := range rspri {

			item := idsapi.AppPrivilege{
				ID:        v.Field("pid").Uint32(),
				Privilege: v.Field("privilege").String(),
				Desc:      v.Field("desc").String(),
			}

			rids := strings.Split(v.Field("roles").String(), ",")
			for _, rv := range rids {

				roleid, _ := strconv.Atoi(rv)
				if roleid > 0 {
					item.Roles = append(item.Roles, uint16(roleid))
				}
			}

			rsp.Privileges = append(rsp.Privileges, item)
		}
	}

	rsp.Kind = "AppInstance"

}

// func (c AppMgr) InstSaveAction() {

// 	c.AutoRender = false

// 	var rsp ResponseJson
// 	rsp.ApiVersion = apiVersion
// 	rsp.Status = 400
// 	rsp.Message = "Bad Request"

// 	defer func() {
// 		if rspj, err := utils.JsonEncode(rsp); err == nil {
// 			io.WriteString(c.Response.Out, rspj)
// 		}
// 	}()

// 	if !c.Session.AccessAllowed("user.admin") {
// 		return
// 	}

// 	dcn, err := rdo.ClientPull("def")
// 	if err != nil {
// 		rsp.Message = "Internal Server Error"
// 		return
// 	}

// 	q := base.NewQuerySet().From("ids_instance").Limit(1)

// 	isNew := true
// 	instset := map[string]interface{}{}

// 	if c.Params.Get("instid") != "" {

// 		q.Where.And("id", c.Params.Get("instid"))

// 		rsinst, err := dcn.Base.Query(q)
// 		if err != nil || len(rsinst) == 0 {
// 			rsp.Status = 400
// 			rsp.Message = http.StatusText(400)
// 			return
// 		}

// 		isNew = false
// 	}

// 	instset["updated"] = base.TimeNow("datetime")
// 	instset["app_title"] = c.Params.Get("app_title")

// 	if isNew {

// 		// TODO

// 	} else {

// 		instset["status"] = c.Params.Get("status")

// 		frupd := base.NewFilter()
// 		frupd.And("id", c.Params.Get("instid"))
// 		if _, err := dcn.Base.Update("ids_instance", instset, frupd); err != nil {
// 			rsp.Status = 500
// 			rsp.Message = "Can not write to database"
// 			return
// 		}
// 	}

// 	rsp.Status = 200
// 	rsp.Message = ""
// }
