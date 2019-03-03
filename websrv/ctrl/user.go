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

package ctrl

import (
	"github.com/hooto/hlang4g/hlang"
	"github.com/hooto/httpsrv"
	"github.com/lessos/lessgo/types"

	"github.com/hooto/iam/config"
	"github.com/hooto/iam/iamapi"
	"github.com/hooto/iam/iamclient"
	"github.com/hooto/iam/store"
)

type User struct {
	*httpsrv.Controller
}

func (c User) PanelInfoAction() {

	rsp := map[string]interface{}{}
	//
	nav := []map[string]string{
		{"path": "app/index",
			"title": hlang.StdLangFeed.Translate(c.Request.Locale, "Authorized Apps")},
		{"path": "access-key/index",
			"title": hlang.StdLangFeed.Translate(c.Request.Locale, "Keys")},
		{"path": "account/index",
			"title": hlang.StdLangFeed.Translate(c.Request.Locale, "Account")},
	}

	if iamclient.SessionAccessAllowed(c.Session, "user.admin", config.Config.InstanceID) {
		nav = append(nav, map[string]string{
			"path":  "user-mgr/index",
			"title": hlang.StdLangFeed.Translate(c.Request.Locale, "Users"),
		})
		nav = append(nav, map[string]string{
			"path":  "acc-mgr/index",
			"title": hlang.StdLangFeed.Translate(c.Request.Locale, "Accounts"),
		})
	}

	if iamclient.SessionAccessAllowed(c.Session, "sys.admin", config.Config.InstanceID) {
		nav = append(nav, map[string]string{
			"path":  "app-mgr/index",
			"title": hlang.StdLangFeed.Translate(c.Request.Locale, "Apps"),
		})
		nav = append(nav, map[string]string{
			"path":  "sys-mgr/index",
			"title": hlang.StdLangFeed.Translate(c.Request.Locale, "System"),
		})
	}

	rsp["topnav"] = nav
	rsp["webui_banner_title"] = config.Config.WebUiBannerTitle

	c.RenderJson(rsp)
}

func (c User) EntryAction() {

	var set iamapi.UserEntry
	defer c.RenderJson(&set)

	user := c.Params.Get("user")

	if user == "" {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeNotFound, "User Not Found")
		return
	}

	// profile
	var profile iamapi.UserProfile
	if obj := store.Data.KvProgGet(iamapi.DataUserProfileKey(user)); obj.OK() {
		obj.Decode(&profile)
		set.Login = iamapi.User{
			Name: user,
		}
		set.Kind = "UserEntry"
	} else {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeNotFound, "User Not Found")
	}
}
