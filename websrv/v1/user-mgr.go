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
	"html"
	"strconv"
	"strings"
	"time"

	"github.com/lessos/lessgo/data/rdo"
	"github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/pass"
	"github.com/lessos/lessgo/types"
	"github.com/lessos/lessgo/utils"

	"../../base/signup"
	"../../idsapi"
)

const (
	userMgrPasswdHidden = "************"
	userMgrPageLimit    = 20
)

var (
	userMgrStatus = map[string]string{
		//0: "Deleted",
		"1": "Active",
		"2": "Banned",
	}
)

type RoleEntry struct {
	Rid, Name, Checked string
}

type UserMgr struct {
	*httpsrv.Controller
}

func (c UserMgr) UserListAction() {

	rsp := idsapi.UserList{}

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

	rdict := map[string]string{}
	q := base.NewQuerySet().From("ids_role").Limit(100)
	rsr, err := dcn.Base.Query(q)
	if err == nil && len(rsr) > 0 {
		for _, v := range rsr {
			rdict[v.Field("rid").String()] = v.Field("name").String()
		}
	}

	page := c.Params.Int64("page")
	if page < 1 {
		page = 1
	}

	// filter: query_text
	q = base.NewQuerySet().From("ids_login").Limit(userMgrPageLimit)
	if query_text := c.Params.Get("query_text"); query_text != "" {
		q.Where.And("name.like", "%"+query_text+"%").
			Or("uname.like", "%"+query_text+"%").
			Or("email.like", "%"+query_text+"%")
	}

	count, _ := dcn.Base.Count("ids_login", q.Where)

	if page > 1 {
		q.Offset(int64((page - 1) * userMgrPageLimit))
	}
	rsl, err := dcn.Base.Query(q)

	if err == nil && len(rsl) > 0 {

		for _, v := range rsl {

			item := idsapi.User{
				Meta: types.ObjectMeta{
					ID:      v.Field("uid").String(),
					Name:    v.Field("uname").String(),
					Created: v.Field("created").String(),
					Updated: v.Field("updated").String(),
				},
				Name:     v.Field("name").String(),
				Email:    v.Field("email").String(),
				Status:   v.Field("status").Uint8(),
				Timezone: v.Field("timezone").String(),
			}

			rids := strings.Split(v.Field("roles").String(), ",")
			for _, rv := range rids {

				roleid, _ := strconv.Atoi(rv)
				if roleid > 0 {

					item.Roles = append(item.Roles, uint16(roleid))
				}
			}

			rsp.Items = append(rsp.Items, item)
		}
	}

	rsp.Meta.TotalResults = uint64(count)
	rsp.Meta.StartIndex = uint64((page - 1) * userMgrPageLimit)
	rsp.Meta.ItemsPerList = uint64(userMgrPageLimit)

	rsp.Kind = "UserList"
}

func (c UserMgr) UserEntryAction() {

	rsp := idsapi.User{}

	defer c.RenderJson(&rsp)

	if !c.Session.AccessAllowed("user.admin") {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	// login
	q := base.NewQuerySet().From("ids_login").Limit(1)
	q.Where.And("uid", c.Params.Get("userid"))
	rslogin, err := dcn.Base.Query(q)
	if err != nil || len(rslogin) != 1 {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	rsp.Meta.ID = rslogin[0].Field("uid").String()
	rsp.Meta.Name = rslogin[0].Field("uname").String()
	rsp.Name = rslogin[0].Field("name").String()
	rsp.Email = rslogin[0].Field("email").String()
	rsp.Auth = userMgrPasswdHidden
	rsp.Status = rslogin[0].Field("status").Uint8()

	rids := strings.Split(rslogin[0].Field("roles").String(), ",")
	for _, rv := range rids {

		roleid, _ := strconv.Atoi(rv)
		if roleid > 0 {

			rsp.Roles = append(rsp.Roles, uint16(roleid))
		}
	}

	//
	rsp.Profile = &idsapi.UserProfile{
		Name: rslogin[0].Field("name").String(),
	}

	//
	q = base.NewQuerySet().From("ids_profile").Limit(1)
	q.Where.And("uid", c.Params.Get("userid"))
	rs, err := dcn.Base.Query(q)
	if err == nil && len(rs) > 0 {
		rsp.Profile.Birthday = rs[0].Field("birthday").String()
		rsp.Profile.About = html.EscapeString(rs[0].Field("aboutme").String())
	}

	rsp.Kind = "User"
}

func (c UserMgr) UserSetAction() {

	rsp := idsapi.User{}

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
		rsp.Error = &types.ErrorMeta{"401", "Access Denied"}
		return
	}

	if err := signup.ValidateEmail(&rsp); err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, err.Error()}
		return
	}

	if rsp.Meta.ID == "" {

		if err := signup.ValidateUsername(&rsp); err != nil {
			rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, err.Error()}
			return
		}

		rsp.Meta.ID = utils.StringEncode16(rsp.Meta.Name, 8)

	}

	if err := signup.ValidateUserID(&rsp); err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, err.Error()}
		return
	}

	isNew := true
	q := base.NewQuerySet().From("ids_login").Limit(1)
	q.Where.And("uid", rsp.Meta.ID)

	rslogin, err := dcn.Base.Query(q)
	if err == nil && len(rslogin) > 0 {
		isNew = false
	}

	loginset := map[string]interface{}{}

	//
	q = base.NewQuerySet().From("ids_login").Limit(1)
	q.Where.And("email", rsp.Email)
	rsu, err := dcn.Base.Query(q)
	if err == nil && len(rsu) > 0 {

		if isNew || rsu[0].Field("uid").String() != rsp.Meta.ID {
			rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "The `Email` already exists, please choose another one"}
			return
		}

	} else {
		loginset["email"] = rsp.Email
	}

	if rsp.Auth != userMgrPasswdHidden {

		pass, _ := pass.HashDefault(rsp.Auth)
		loginset["pass"] = pass
	}

	if isNew {
		loginset["uid"] = rsp.Meta.ID
		loginset["uname"] = rsp.Meta.Name
		loginset["created"] = base.TimeNow("datetime")
		loginset["timezone"] = "UTC"
	}

	loginset["status"] = rsp.Status
	loginset["updated"] = base.TimeNow("datetime")
	loginset["name"] = rsp.Name

	roles := []string{}
	for _, v := range rsp.Roles {
		roles = append(roles, strconv.Itoa(int(v)))
	}

	loginset["roles"] = strings.Join(roles, ",")

	frupd := base.NewFilter()

	if isNew {

		if _, err := dcn.Base.Insert("ids_login", loginset); err != nil {
			rsp.Error = &types.ErrorMeta{idsapi.ErrCodeUnavailable, "Service Unavailable"}
			return
		}

	} else {

		frupd.And("uid", rsp.Meta.ID)
		if _, err := dcn.Base.Update("ids_login", loginset, frupd); err != nil {
			rsp.Error = &types.ErrorMeta{idsapi.ErrCodeUnavailable, "Service Unavailable"}
			return
		}
	}

	if rsp.Profile == nil {
		rsp.Profile = &idsapi.UserProfile{}
	}

	if _, err := time.Parse("2006-01-02", rsp.Profile.Birthday); err != nil {
		rsp.Profile.Birthday = "0000-00-00"
	}

	profile := map[string]interface{}{
		"birthday": rsp.Profile.Birthday,
		"aboutme":  rsp.Profile.About,
		"updated":  base.TimeNow("datetime"),
	}

	if isNew {
		profile["uid"] = rsp.Meta.ID
		profile["gender"] = 0
		profile["created"] = base.TimeNow("datetime")

		dcn.Base.Insert("ids_profile", profile)
	} else {
		dcn.Base.Update("ids_profile", profile, frupd)
	}

	rsp.Kind = "User"
}
