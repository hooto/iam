// Copyright 2014 lessos Authors, All rights reserved.
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

	"code.hooto.com/lessos/iam/iamapi"
	"code.hooto.com/lessos/iam/iamclient"
	"code.hooto.com/lessos/iam/store"
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
	c.Data["access_token_key"] = iamclient.AccessTokenKey

	if len(c.Params.Get("continue")) > 0 {
		c.Data["continue"] = c.Params.Get("continue")
	}

	token := iamapi.AccessTokenFrontend(c.Params.Get(iamclient.AccessTokenKey))
	if !token.Valid() {
		session, _ := iamclient.SessionInstance(c.Session)
		token = iamapi.AccessTokenFrontend(session.FullToken())
	}

	if token.Valid() {
		store.PvDel("session/"+token.SessionPath(), nil)
	}

	http.SetCookie(c.Response.Out, &http.Cookie{
		Name:   iamclient.AccessTokenKey,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}
