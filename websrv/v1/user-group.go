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
	"strings"

	"github.com/hooto/httpsrv"
	"github.com/lessos/lessgo/types"

	"github.com/hooto/iam/config"
	"github.com/hooto/iam/iamapi"
	"github.com/hooto/iam/iamclient"
	"github.com/hooto/iam/store"
)

type UserGroup struct {
	*httpsrv.Controller
	us iamapi.UserSession
}

func (c *UserGroup) Init() int {

	//
	c.us, _ = iamclient.SessionInstance(c.Session)

	if !c.us.IsLogin() {
		c.Response.Out.WriteHeader(401)
		c.RenderJson(types.NewTypeErrorMeta(iamapi.ErrCodeUnauthorized, "Unauthorized"))
		return 1
	}

	return 0
}

func (c UserGroup) ItemAction() {

	var set iamapi.UserGroupItem
	defer c.RenderJson(&set)

	name := iamapi.UserNameFilter(c.Params.Get("name"))
	if err := iamapi.UserNameValid(name); err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, err.Error())
		return
	}

	p := store.UserGet(name)
	if p == nil || p.Type != iamapi.UserTypeGroup {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Item Not Found")
		return
	}

	if !iamapi.ArrayStringHas(p.Owners, c.us.UserName) &&
		!iamapi.ArrayStringHas(p.Members, c.us.UserName) &&
		!iamclient.SessionAccessAllowed(c.Session, "user.admin", config.Config.InstanceID) {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
		return
	}

	set.Name = p.Name
	set.DisplayName = p.DisplayName
	set.Status = p.Status
	set.Owners = p.Owners
	set.Members = p.Members
	set.Created = set.Created
	set.Updated = set.Updated

	set.Kind = "UserGroupItem"
}

func (c UserGroup) ListAction() {

	var ls iamapi.UserGroupList
	defer c.RenderJson(&ls)

	rss := store.UserGroupList()

	for _, v := range rss {

		if v.Type != iamapi.UserTypeGroup {
			continue
		}

		if !iamapi.ArrayStringHas(v.Owners, c.us.UserName) &&
			!iamapi.ArrayStringHas(v.Members, c.us.UserName) &&
			!iamclient.SessionAccessAllowed(c.Session, "user.admin", config.Config.InstanceID) {
			continue
		}

		item := &iamapi.UserGroupItem{
			Name:        v.Name,
			DisplayName: v.DisplayName,
			Status:      v.Status,
			Created:     v.Created,
			Updated:     v.Updated,
			Owners:      v.Owners,
			Members:     v.Members,
		}

		ls.Items = append(ls.Items, item)
	}

	ls.Kind = "UserGroupList"
}

func (c UserGroup) SetAction() {

	var set iamapi.UserGroupItem
	defer c.RenderJson(&set)

	if err := c.Request.JsonDecode(&set); err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Bad Request")
		return
	}

	set.Name = iamapi.UserNameFilter(set.Name)
	if err := iamapi.UserNameValid(set.Name); err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, err.Error())
		return
	}

	var (
		prev = store.UserGet(set.Name)
		chg  = false
	)

	if set.DisplayName == "" {
		set.DisplayName = strings.Title(set.Name)
	}

	if prev == nil {
		prev = &iamapi.User{
			// Id:      iamapi.UserId(set.Name),
			Name:    set.Name,
			Created: types.MetaTimeNow(),
			Type:    iamapi.UserTypeGroup,
		}
		chg = true
	} else {
		if prev.Type != iamapi.UserTypeGroup ||
			(!iamapi.ArrayStringHas(prev.Owners, c.us.UserName) &&
				!iamclient.SessionAccessAllowed(c.Session, "user.admin", config.Config.InstanceID)) {
			set.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
			return
		}
	}

	if prev.DisplayName != set.DisplayName {
		prev.DisplayName = set.DisplayName
		chg = true
	}

	if set.Status > 0 && prev.Status != set.Status {
		prev.Status = set.Status
		chg = true
	}

	if len(set.Owners) > 0 {
		prev.Owners = []string{c.us.UserName}
		for _, v := range set.Owners {

			p := store.UserGet(v)
			if p == nil {
				set.Error = types.NewErrorMeta(
					iamapi.ErrCodeInvalidArgument, "User Not Found ("+v+")")
				return
			}

			iamapi.ArrayStringSet(&prev.Owners, v)
		}
		chg = true
	}

	if len(set.Members) > 0 {
		prev.Members = []string{c.us.UserName}
		for _, v := range set.Members {

			p := store.UserGet(v)
			if p == nil {
				set.Error = types.NewErrorMeta(
					iamapi.ErrCodeInvalidArgument, "User Not Found ("+v+")")
				return
			}

			iamapi.ArrayStringSet(&prev.Members, v)
		}
		chg = true
	}

	if len(set.Members) < 1 || len(set.Owners) < 1 {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Member Not Found")
		return
	}

	if chg {
		prev.Updated = types.MetaTimeNow()
		if !store.UserSet(prev) {
			set.Error = types.NewErrorMeta(iamapi.ErrCodeServerError, "Server Error")
			return
		}
	}

	set.Kind = "UserGropuItem"
}
