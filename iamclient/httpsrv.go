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

package iamclient

import (
	"net/http"
	"strings"
	"time"

	"code.hooto.com/lessos/iam/iamapi"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/types"
)

type AuthSession struct {
	types.TypeMeta `json:",inline"`
	UserID         string            `json:"userid"`
	UserName       string            `json:"username"`
	Name           string            `json:"name"`
	IamUrl         string            `json:"iam_url"`
	PhotoUrl       string            `json:"photo_url"`
	InstanceOwner  bool              `json:"instance_owner,omitempty"`
	Roles          types.ArrayUint32 `json:"roles,omitempty"`
}

type Auth struct {
	*httpsrv.Controller
}

func (c Auth) CbAction() {

	http.SetCookie(c.Response.Out, &http.Cookie{
		Name:     AccessTokenKey,
		Value:    c.Params.Get(AccessTokenKey),
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(time.Second * time.Duration(c.Params.Int64("expires_in"))),
	})

	if c.Params.Get("state") != "" {
		c.Redirect(c.Params.Get("state"))
	} else {
		c.Redirect(c.UrlModuleBase(""))
	}
}

func (c Auth) LoginAction() {

	referer := c.UrlModuleBase("")
	if len(c.Request.Referer()) > 10 {
		referer = c.Request.Referer()
	}

	c.Redirect(AuthServiceUrl(
		InstanceID,
		c.UrlModuleBase("auth/cb"),
		referer,
	))
}

func (c Auth) SessionAction() {

	set := AuthSession{
		IamUrl:   service_prefix(),
		PhotoUrl: service_prefix() + "/v1/service/photo/guest",
	}

	if session, err := SessionInstance(c.Session); err == nil {

		set.UserID = session.UserID
		set.UserName = session.UserName
		set.Name = session.Name
		set.PhotoUrl = service_prefix() + "/v1/service/photo/" + session.UserID
		set.Roles = session.Roles

		if InstanceOwner == set.UserID {
			set.InstanceOwner = true
		}

		set.Kind = "AuthSession"

	} else {
		set.Error = types.NewErrorMeta("401", "Unauthorized")
	}

	c.RenderJson(set)
}

func (c Auth) SignOutAction() {

	http.SetCookie(c.Response.Out, &http.Cookie{
		Name:    AccessTokenKey,
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(-86400),
	})

	referer := ""
	if len(c.Request.Referer()) > 10 {
		referer = strings.TrimRight(c.Request.Referer(), "/")
	}

	if referer == "" {
		referer = c.UrlModuleBase("")
	}

	if strings.Contains(referer, "?") {
		referer += "&"
	} else {
		referer += "?"
	}

	c.Redirect(referer + "_iam_out=1")
}

func (c Auth) AppRoleListAction() {

	sets, err := AppRoleList(c.Session, "") // TODO appid
	if err == nil {
		c.RenderJson(sets)
	} else if err != nil {
		c.RenderJson(types.NewTypeErrorMeta("500", err.Error()))
	} else {
		c.RenderJson(types.NewTypeErrorMeta(iamapi.ErrCodeUnauthorized, "Unauthorized"))
	}
}
