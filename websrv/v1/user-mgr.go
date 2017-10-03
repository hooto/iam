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
	"html"
	"sort"
	"strings"
	"time"

	"github.com/hooto/httpsrv"
	"github.com/lessos/lessgo/pass"
	"github.com/lessos/lessgo/types"

	"github.com/hooto/iam/base/signup"
	"github.com/hooto/iam/config"
	"github.com/hooto/iam/iamapi"
	"github.com/hooto/iam/iamclient"
	"github.com/hooto/iam/store"
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

	if !iamclient.SessionAccessAllowed(c.Session, "user.admin", config.Config.InstanceID) {
		ls.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
		return
	}

	var (
		qt = strings.ToLower(c.Params.Get("qry_text"))
	)

	// TODO page
	if rs := store.PoScan("user", []byte{}, []byte{}, 1000); rs.OK() {

		rss := rs.KvList()
		for _, obj := range rss {

			var user iamapi.User
			if err := obj.Decode(&user); err == nil {

				if qt != "" && (!strings.Contains(user.Name, qt) &&
					!strings.Contains(user.Email, qt)) {
					continue
				}

				user.Id = ""
				user.Keys = nil

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

	set := iamapi.UserEntry{}
	defer c.RenderJson(&set)

	if !iamclient.SessionAccessAllowed(c.Session, "user.admin", config.Config.InstanceID) {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
		return
	}

	userid := iamapi.UserId(c.Params.Get("username"))

	if obj := store.PoGet("user", userid); obj.OK() {
		obj.Decode(&set.Login)
	}

	// login
	if set.Login.Name == "" {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "User Not Found")
		return
	}
	set.Login.Keys = types.KvPairs{}
	set.Login.Keys.Set(iamapi.UserKeyDefault, userMgrPasswdHidden)

	//
	var profile iamapi.UserProfile
	if obj := store.PoGet("user-profile", userid); obj.OK() {
		obj.Decode(&profile)

		profile.About = html.EscapeString(profile.About)
		profile.Photo = ""
		profile.PhotoSource = ""
		profile.Login = nil
	}

	set.Profile = &profile

	set.Kind = "User"
}

func (c UserMgr) UserSetAction() {

	set := iamapi.UserEntry{}
	defer c.RenderJson(&set)

	if err := c.Request.JsonDecode(&set); err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Bad Request")
		return
	}

	if !iamclient.SessionAccessAllowed(c.Session, "user.admin", config.Config.InstanceID) {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
		return
	}

	if err := signup.ValidateUsername(&set.Login); err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, err.Error())
		return
	}

	if err := signup.ValidateEmail(&set.Login); err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, err.Error())
		return
	}

	set.Login.Id = iamapi.UserId(set.Login.Name)

	if auth := set.Login.Keys.Get(iamapi.UserKeyDefault); auth != nil {

		if auth.String() != userMgrPasswdHidden && auth.String() != "" {
			authenc, _ := pass.HashDefault(auth.String())
			set.Login.Keys.Set(iamapi.UserKeyDefault, authenc)
		} else {
			set.Login.Keys.Del(iamapi.UserKeyDefault)
		}
	} else {
		set.Login.Keys.Del(iamapi.UserKeyDefault)
	}

	var prev iamapi.UserEntry

	//
	if obj := store.PoGet("user", set.Login.Id); obj.OK() {
		obj.Decode(&prev.Login)
	}

	//
	if prev.Login.Id == set.Login.Id {

		if set.Login.Email != "" {
			prev.Login.Email = set.Login.Email
		}

		if auth := set.Login.Keys.Get(iamapi.UserKeyDefault); auth != nil {
			prev.Login.Keys.Set(iamapi.UserKeyDefault, auth.String())
		}

		if set.Login.DisplayName != "" {
			prev.Login.DisplayName = set.Login.DisplayName
		}

		if len(set.Login.Roles) > 0 && !set.Login.Roles.Equal(prev.Login.Roles) {
			prev.Login.Roles = set.Login.Roles
		}

		// prev.Profile = set.Profile

		set.Login = prev.Login
	} else {
		set.Login.Created = types.MetaTimeNow()
	}

	set.Login.Updated = types.MetaTimeNow()
	sort.Sort(set.Login.Roles)

	if obj := store.PoPut("user", set.Login.Id, set.Login, nil); !obj.OK() {
		set.Error = types.NewErrorMeta("500", obj.Bytex().String())
		return
	}

	if set.Profile != nil {

		var profile iamapi.UserProfile

		if obj := store.PoGet("user-profile", set.Login.Id); obj.OK() {

			obj.Decode(&profile)

			if _, err := time.Parse("2006-01-02", set.Profile.Birthday); err == nil {
				profile.Birthday = set.Profile.Birthday
			}

			if set.Profile.About != "" {
				profile.About = set.Profile.About
			}
		}

		if obj := store.PoPut("user-profile", set.Login.Id, profile, nil); !obj.OK() {
			set.Error = types.NewErrorMeta("500", obj.Bytex().String())
			return
		}
	}

	set.Kind = "User"
}
