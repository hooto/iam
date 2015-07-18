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

package ctrl

import (
	"net/http"

	"github.com/lessos/bigtree/btapi"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessids/store"
)

type Service struct {
	*httpsrv.Controller
}

func (c Service) LoginAction() {

	c.Data["continue"] = c.Params.Get("continue")
	if c.Params.Get("persistent") == "1" {
		c.Data["persistentChecked"] = "checked"
	}
}

func (c Service) SignOutAction() {

	c.Data["continue"] = "/ids"

	token := c.Params.Get("access_token")

	if token == "" {
		session, _ := c.Session.SessionFetch()
		token = session.AccessToken
	}

	if len(c.Params.Get("continue")) > 0 {
		c.Data["continue"] = c.Params.Get("continue")
	}

	obj := store.BtAgent.ObjectGet(btapi.ObjectProposal{
		Meta: btapi.ObjectMeta{
			Path: "/session/" + token,
		},
	})

	store.BtAgent.ObjectDel(btapi.ObjectProposal{
		Meta: btapi.ObjectMeta{
			Path: "/session/" + token,
		},
		PrevVersion: obj.Meta.Version,
	})

	ck := &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	}
	http.SetCookie(c.Response.Out, ck)
}
