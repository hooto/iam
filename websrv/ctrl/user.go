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

package ctrl

import (
	"github.com/lessos/iam/config"
	"github.com/lessos/iam/iamclient"
	"github.com/lessos/lessgo/httpsrv"
)

type User struct {
	*httpsrv.Controller
}

func (c User) PanelInfoAction() {

	rsp := map[string]interface{}{}
	//
	nav := []map[string]string{
		{"path": "#my-app/index", "title": "My Applications"},
	}

	if iamclient.SessionAccessAllowed(c.Session, "user.admin", "df085c6dc6ff") {
		nav = append(nav, map[string]string{
			"path":  "#user-mgr/index",
			"title": "User Manage",
		})
	}

	if iamclient.SessionAccessAllowed(c.Session, "sys.admin", "df085c6dc6ff") {

		nav = append(nav, map[string]string{
			"path":  "#app-mgr/index",
			"title": "Applications",
		})
		nav = append(nav, map[string]string{
			"path":  "#sys-mgr/index",
			"title": "System Settings",
		})
	}

	rsp["topnav"] = nav
	rsp["webui_banner_title"] = config.Config.WebUiBannerTitle

	c.RenderJson(rsp)
}
