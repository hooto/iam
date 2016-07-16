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
	"net/http"

	"github.com/lessos/bigtree/btapi"
	"github.com/lessos/iam/iamclient"
	"github.com/lessos/iam/store"
	"github.com/lessos/lessgo/httpsrv"
)

type Service struct {
	*httpsrv.Controller
}

func (c Service) LoginAction() {

	if c.Params.Get("persistent") == "1" {
		c.Data["persistent_checked"] = "checked"
	}

	if c.Params.Get("client_id") != "" {
		c.Data["client_id"] = c.Params.Get("client_id")
	}

	if c.Params.Get("redirect_uri") != "" {
		c.Data["redirect_uri"] = c.Params.Get("redirect_uri")
	}

	if c.Params.Get("state") != "" {
		c.Data["state"] = c.Params.Get("state")
	}
}

func (c Service) SignOutAction() {

	c.Data["continue"] = "/iam"

	if len(c.Params.Get("continue")) > 0 {
		c.Data["continue"] = c.Params.Get("continue")
	}

	token := c.Params.Get("access_token")

	if len(token) < 30 {
		session, _ := iamclient.SessionInstance(c.Session)
		token = session.FullToken()
	}

	if len(token) > 30 {

		token = token[:8] + "/" + token[9:]

		store.BtAgent.ObjectDel("/global/iam/session/"+token, &btapi.ObjectWriteOptions{
			Force: true,
		})
	}

	http.SetCookie(c.Response.Out, &http.Cookie{
		Name:   iamclient.AccessTokenKey,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}
