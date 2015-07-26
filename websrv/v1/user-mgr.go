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
	"time"

	"github.com/lessos/bigtree/btapi"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/pass"
	"github.com/lessos/lessgo/types"
	"github.com/lessos/lessgo/utils"
	"github.com/lessos/lessgo/utilx"

	"github.com/lessos/lessids/base/signup"
	"github.com/lessos/lessids/idsapi"
	"github.com/lessos/lessids/store"
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

type UserMgr struct {
	*httpsrv.Controller
}

func (c UserMgr) UserListAction() {

	ls := idsapi.UserList{}

	defer c.RenderJson(&ls)

	if !c.Session.AccessAllowed("user.admin") {
		ls.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	if objs := store.BtAgent.ObjectList(btapi.ObjectProposal{
		Meta: btapi.ObjectMeta{
			Path: "/user/",
		},
	}); objs.Error == nil {

		for _, obj := range objs.Items {

			var user idsapi.User
			if err := obj.JsonDecode(&user); err == nil {
				user.Auth = ""
				ls.Items = append(ls.Items, user)
			}
		}
	}

	// page := c.Params.Int64("page")
	// if page < 1 {
	// 	page = 1
	// }

	// filter: query_text

	// count, _ := dcn.Base.Count("ids_login", q.Where)

	// if page > 1 {
	// 	q.Offset(int64((page - 1) * userMgrPageLimit))
	// }

	// ls.Meta.TotalResults = uint64(count)
	// ls.Meta.StartIndex = uint64((page - 1) * userMgrPageLimit)
	// ls.Meta.ItemsPerList = uint64(userMgrPageLimit)

	ls.Kind = "UserList"
}

func (c UserMgr) UserEntryAction() {

	set := idsapi.User{}

	defer c.RenderJson(&set)

	if !c.Session.AccessAllowed("user.admin") {
		set.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	if obj := store.BtAgent.ObjectGet(btapi.ObjectProposal{
		Meta: btapi.ObjectMeta{
			Path: "/user/" + c.Params.Get("userid"),
		},
	}); obj.Error == nil {
		obj.JsonDecode(&set)
	}

	// login
	if set.Meta.ID == "" {
		set.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "User Not Found"}
		return
	}
	set.Auth = userMgrPasswdHidden

	//
	var profile idsapi.UserProfile
	if obj := store.BtAgent.ObjectGet(btapi.ObjectProposal{
		Meta: btapi.ObjectMeta{
			Path: "/user-profile/" + c.Params.Get("userid"),
		},
	}); obj.Error == nil {
		obj.JsonDecode(&profile)
		profile.About = html.EscapeString(profile.About)
	}

	set.Profile = &profile

	set.Kind = "User"
}

func (c UserMgr) UserSetAction() {

	set := idsapi.User{}

	defer c.RenderJson(&set)

	if err := c.Request.JsonDecode(&set); err != nil {
		set.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "Bad Request"}
		return
	}

	if !c.Session.AccessAllowed("user.admin") {
		set.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	if err := signup.ValidateEmail(&set); err != nil {
		set.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, err.Error()}
		return
	}

	var prevVersion uint64

	if set.Meta.ID == "" {

		if err := signup.ValidateUsername(&set); err != nil {
			set.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, err.Error()}
			return
		}

		set.Meta.Name = set.Meta.Name
		set.Meta.ID = utils.StringEncode16(set.Meta.Name, 8)
		set.Meta.Created = utilx.TimeNow("atom")

		set.Auth, _ = pass.HashDefault(set.Auth)

	} else {

		var prev idsapi.User
		if obj := store.BtAgent.ObjectGet(btapi.ObjectProposal{
			Meta: btapi.ObjectMeta{
				Path: "/user/" + set.Meta.ID,
			},
		}); obj.Error == nil {
			obj.JsonDecode(&prev)
			prevVersion = obj.Meta.Version
		}

		if prev.Meta.ID != set.Meta.ID {
			set.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "User Not Found"}
			return
		}

		if set.Email != "" {
			prev.Email = set.Email
		}

		if set.Auth != userMgrPasswdHidden || set.Auth != "" {
			prev.Auth = set.Auth
		}

		if set.Timezone != "" {
			prev.Timezone = set.Timezone
		}

		if set.Name != "" {
			prev.Name = set.Name
		}

		set = prev
	}

	if err := signup.ValidateUserID(&set); err != nil {
		set.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, err.Error()}
		return
	}

	set.Meta.Updated = utilx.TimeNow("atom")

	setjs, _ := utils.JsonEncode(set)

	if obj := store.BtAgent.ObjectSet(btapi.ObjectProposal{
		Meta: btapi.ObjectMeta{
			Path: "/user/" + set.Meta.ID,
		},
		Data:        setjs,
		PrevVersion: prevVersion,
	}); obj.Error != nil {
		set.Error = &types.ErrorMeta{"500", obj.Error.Message}
		return
	}

	prevVersion = 0
	var profile idsapi.UserProfile

	if obj := store.BtAgent.ObjectGet(btapi.ObjectProposal{
		Meta: btapi.ObjectMeta{
			Path: "/user-profile/" + set.Meta.ID,
		},
	}); obj.Error == nil {

		obj.JsonDecode(&profile)
		prevVersion = obj.Meta.Version

		if _, err := time.Parse("2006-01-02", set.Profile.Birthday); err == nil {
			profile.Birthday = set.Profile.Birthday
		}

		if set.Profile.About != "" {
			profile.About = set.Profile.About
		}

		if set.Name != "" {
			profile.Name = set.Name
		}
	}

	profjs, _ := utils.JsonEncode(profile)

	if obj := store.BtAgent.ObjectSet(btapi.ObjectProposal{
		Meta: btapi.ObjectMeta{
			Path: "/user-profile/" + set.Meta.ID,
		},
		Data:        profjs,
		PrevVersion: prevVersion,
	}); obj.Error != nil {
		set.Error = &types.ErrorMeta{"500", obj.Error.Message}
		return
	}

	set.Kind = "User"
}
