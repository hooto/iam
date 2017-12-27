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
	"net/http"

	"github.com/hooto/httpsrv"
	"github.com/hooto/iam/config"
	"github.com/hooto/iam/iamapi"
	"github.com/hooto/iam/iamclient"
	"github.com/hooto/iam/store"
)

type Service struct {
	*httpsrv.Controller
}

func (c Service) LoginAction() {

	var rt iamapi.ServiceRedirectToken

	if len(c.Params.Get("redirect_token")) > 20 {
		rt = iamapi.ServiceRedirectTokenDecode(c.Params.Get("redirect_token"))
	}

	if rt.RedirectUri == "" && c.Params.Get("redirect_uri") != "" {
		rt = iamapi.ServiceRedirectToken{
			RedirectUri: c.Params.Get("redirect_uri"),
			State:       c.Params.Get("state"),
			ClientId:    c.Params.Get("client_id"),
			Persistent:  int(c.Params.Int64("persistent")),
		}
	}

	if rt.RedirectUri != "" {
		c.Data["redirect_token"] = rt.Encode()
	}

	if rt.Persistent == 1 {
		c.Data["persistent_checked"] = "checked"
	}

	if config.Config.ServiceLoginFormAlertMsg != "" {
		c.Data["alert_msg"] = config.Config.ServiceLoginFormAlertMsg
	}
}

func (c Service) SignOutAction() {

	c.Data["access_token_key"] = iamclient.AccessTokenKey

	token := iamapi.AccessTokenFrontend(c.Params.Get(iamclient.AccessTokenKey))
	if !token.Valid() {
		session, _ := iamclient.SessionInstance(c.Session)
		token = iamapi.AccessTokenFrontend(session.FullToken())
	}

	if token.Valid() {
		store.Data.ProgDel(iamapi.DataSessionKey(token.User(), token.Id()), nil)
	}

	http.SetCookie(c.Response.Out, &http.Cookie{
		Name:   iamclient.AccessTokenKey,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}
