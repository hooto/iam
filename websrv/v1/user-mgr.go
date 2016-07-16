// Copyright 2014-2016 iam Author, All rights reserved.
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

	"github.com/lessos/iam/base/signup"
	"github.com/lessos/iam/iamapi"
	"github.com/lessos/iam/iamclient"
	"github.com/lessos/iam/store"
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

	ls := iamapi.UserList{}

	defer c.RenderJson(&ls)

	if !iamclient.SessionAccessAllowed(c.Session, "user.admin", "df085c6dc6ff") {
		ls.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	if objs := store.BtAgent.ObjectList("/global/iam/user/"); objs.Error == nil {

		for _, obj := range objs.Items {

			var user iamapi.User
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

	// count, _ := dcn.Base.Count("iam_login", q.Where)

	// if page > 1 {
	// 	q.Offset(int64((page - 1) * userMgrPageLimit))
	// }

	// ls.Meta.TotalResults = uint64(count)
	// ls.Meta.StartIndex = uint64((page - 1) * userMgrPageLimit)
	// ls.Meta.ItemsPerList = uint64(userMgrPageLimit)

	ls.Kind = "UserList"
}

func (c UserMgr) UserEntryAction() {

	set := iamapi.User{}

	defer c.RenderJson(&set)

	if !iamclient.SessionAccessAllowed(c.Session, "user.admin", "df085c6dc6ff") {
		set.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	if obj := store.BtAgent.ObjectGet("/global/iam/user/" + c.Params.Get("userid")); obj.Error == nil {
		obj.JsonDecode(&set)
	}

	// login
	if set.Meta.ID == "" {
		set.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "User Not Found"}
		return
	}
	set.Auth = userMgrPasswdHidden

	//
	var profile iamapi.UserProfile
	if obj := store.BtAgent.ObjectGet("/global/iam/user-profile/" + c.Params.Get("userid")); obj.Error == nil {
		obj.JsonDecode(&profile)
		profile.About = html.EscapeString(profile.About)
	}

	set.Profile = &profile

	set.Kind = "User"
}

func (c UserMgr) UserSetAction() {

	set := iamapi.User{}

	defer c.RenderJson(&set)

	if err := c.Request.JsonDecode(&set); err != nil {
		set.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "Bad Request"}
		return
	}

	if !iamclient.SessionAccessAllowed(c.Session, "user.admin", "df085c6dc6ff") {
		set.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	if err := signup.ValidateEmail(&set); err != nil {
		set.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, err.Error()}
		return
	}

	var prev iamapi.User
	var prevVersion uint64

	if set.Meta.ID == "" {

		if err := signup.ValidateUsername(&set); err != nil {
			set.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, err.Error()}
			return
		}

		set.Meta.ID = utils.StringEncode16(set.Meta.Name, 8)
		set.Meta.Created = utilx.TimeNow("atom")
	}

	if set.Auth != userMgrPasswdHidden && set.Auth != "" {
		set.Auth, _ = pass.HashDefault(set.Auth)
	} else {
		set.Auth = ""
	}

	//
	if obj := store.BtAgent.ObjectGet("/global/iam/user/" + set.Meta.ID); obj.Error == nil {
		obj.JsonDecode(&prev)
		prevVersion = obj.Meta.Version
	}

	//
	if err := signup.ValidateUserID(&set); err != nil {
		set.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, err.Error()}
		return
	}

	//
	if prev.Meta.ID == set.Meta.ID {

		if set.Email != "" {
			prev.Email = set.Email
		}

		if set.Auth != "" {
			prev.Auth = set.Auth
		}

		if set.Timezone != "" {
			prev.Timezone = set.Timezone
		}

		if set.Name != "" {
			prev.Name = set.Name
		}

		prev.Profile = set.Profile

		set = prev
	}

	set.Meta.Updated = utilx.TimeNow("atom")

	if obj := store.BtAgent.ObjectSet("/global/iam/user/"+set.Meta.ID, set, &btapi.ObjectWriteOptions{
		PrevVersion: prevVersion,
	}); obj.Error != nil {
		set.Error = &types.ErrorMeta{"500", obj.Error.Message}
		return
	}

	if set.Profile != nil {

		prevVersion = 0
		var profile iamapi.UserProfile

		if obj := store.BtAgent.ObjectGet("/global/iam/user-profile/" + set.Meta.ID); obj.Error == nil {

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

		if obj := store.BtAgent.ObjectSet("/global/iam/user-profile/"+set.Meta.ID, profile, &btapi.ObjectWriteOptions{
			PrevVersion: prevVersion,
		}); obj.Error != nil {
			set.Error = &types.ErrorMeta{"500", obj.Error.Message}
			return
		}
	}

	set.Kind = "User"
}
