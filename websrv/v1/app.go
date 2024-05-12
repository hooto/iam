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
	"github.com/hooto/httpsrv"
	"github.com/lessos/lessgo/types"

	"github.com/hooto/iam/data"
	"github.com/hooto/iam/iamapi"
	"github.com/hooto/iam/iamclient"
)

const (
	myAppInstPageLimit = 100
)

type App struct {
	*httpsrv.Controller
	us iamapi.UserSession
}

func (c *App) Init() int {

	//
	c.us, _ = iamclient.SessionInstance(c.Session)

	if !c.us.IsLogin() {
		c.Response.Out.WriteHeader(401)
		c.RenderJson(types.NewTypeErrorMeta(iamapi.ErrCodeUnauthorized, "Unauthorized"))
		return 1
	}

	return 0
}

func (c App) InstListAction() {

	ls := types.ObjectList{}
	defer c.RenderJson(&ls)

	if rs := data.Data.NewRanger(
		iamapi.ObjKeyAppInstance(""), iamapi.ObjKeyAppInstance("zzzzzzzz")).
		SetRevert(true).SetLimit(1000).Exec(); rs.OK() {

		for _, obj := range rs.Items {

			var inst iamapi.AppInstance
			if err := obj.JsonDecode(&inst); err == nil {

				if inst.Meta.User == c.us.UserName {
					ls.Items = append(ls.Items, inst)
				}
			}
		}
	}

	ls.Kind = "AppInstanceList"
}

func (c App) InstEntryAction() {

	var set struct {
		types.TypeMeta
		iamapi.AppInstance
	}
	defer c.RenderJson(&set)

	if obj := data.Data.NewReader(iamapi.ObjKeyAppInstance(c.Params.Value("instid"))).Exec(); obj.OK() {
		obj.Item().JsonDecode(&set.AppInstance)
	}

	if set.Meta.ID == "" {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "App Instance Not Found")
		return
	}

	if set.Meta.User != c.us.UserName {
		set.AppInstance = iamapi.AppInstance{}
		set.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "Access Denied")
		return
	}

	set.Kind = "AppInstance"
}

func (c App) InstSetAction() {

	var set struct {
		types.TypeMeta
		iamapi.AppInstance
	}
	defer c.RenderJson(&set)

	if err := c.Request.JsonDecode(&set.AppInstance); err != nil || set.Meta.ID == "" {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "InvalidArgument")
		return
	}

	var prev iamapi.AppInstance
	if obj := data.Data.NewReader(iamapi.ObjKeyAppInstance(set.Meta.ID)).Exec(); obj.OK() {
		obj.Item().JsonDecode(&prev)
	}

	if prev.Meta.ID == "" {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "App Instance Not Found")
		return
	}

	if prev.Meta.User != c.us.UserName {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "Access Denied")
		return
	}

	if set.AppTitle != prev.AppTitle || set.Url != prev.Url {

		prev.Meta.Updated = types.MetaTimeNow()
		prev.AppTitle = set.AppTitle
		prev.Url = set.Url

		if obj := data.Data.NewWriter(iamapi.ObjKeyAppInstance(set.Meta.ID), nil).SetJsonValue(prev).
			Exec(); !obj.OK() {
			set.Error = types.NewErrorMeta(iamapi.ErrCodeInternalError, obj.ErrorMessage())
			return
		}
	}

	set.Kind = "AppInstance"
}

func (c App) InstDelAction() {

	var set types.TypeMeta
	defer c.RenderJson(&set)

	inst_id := c.Params.Value("inst_id")

	var prev iamapi.AppInstance
	if obj := data.Data.NewReader(iamapi.ObjKeyAppInstance(inst_id)).Exec(); obj.OK() {
		obj.Item().JsonDecode(&prev)
	}

	if prev.Meta.ID == "" {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "App Instance Not Found")
		return
	}

	if prev.Meta.User != c.us.UserName {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "Access Denied")
		return
	}

	if obj := data.Data.NewDeleter(iamapi.ObjKeyAppInstance(inst_id)).Exec(); !obj.OK() {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInternalError, obj.ErrorMessage())
		return
	}

	set.Kind = "AppInstance"
}
