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

package v1

import (
	"github.com/lessos/bigtree/btapi"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/types"
	"github.com/lessos/lessgo/utilx"

	"github.com/lessos/lessids/idclient"
	"github.com/lessos/lessids/idsapi"
	"github.com/lessos/lessids/store"
)

const (
	myAppInstPageLimit = 100
)

type MyApp struct {
	*httpsrv.Controller
}

func (c MyApp) InstListAction() {

	ls := idsapi.AppInstanceList{}

	defer c.RenderJson(&ls)

	session, err := idclient.SessionInstance(c.Session)

	if err != nil || !session.IsLogin() {
		ls.Error = &types.ErrorMeta{idsapi.ErrCodeUnauthorized, "Access Denied"}
		return
	}

	if objs := store.BtAgent.ObjectList("/global/ids/app-instance/"); objs.Error == nil {

		for _, obj := range objs.Items {

			var inst idsapi.AppInstance
			if err := obj.JsonDecode(&inst); err == nil {

				if inst.Meta.UserID == session.UserID {
					ls.Items = append(ls.Items, inst)
				}
			}
		}
	}

	ls.Kind = "AppInstanceList"
}

func (c MyApp) InstEntryAction() {

	set := idsapi.AppInstance{}

	defer c.RenderJson(&set)

	session, err := idclient.SessionInstance(c.Session)

	if err != nil || !session.IsLogin() {
		set.Error = &types.ErrorMeta{idsapi.ErrCodeUnauthorized, "Access Denied"}
		return
	}

	if obj := store.BtAgent.ObjectGet("/global/ids/app-instance/" + c.Params.Get("instid")); obj.Error == nil {
		obj.JsonDecode(&set)
	}

	if set.Meta.ID == "" {
		set.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "App Instance Not Found"}
		return
	}

	if set.Meta.UserID != session.UserID {
		set = idsapi.AppInstance{}
		set.Error = &types.ErrorMeta{idsapi.ErrCodeUnauthorized, "Access Denied"}
		return
	}

	set.Kind = "AppInstance"
}

func (c MyApp) InstSetAction() {

	set := idsapi.AppInstance{}

	defer c.RenderJson(&set)

	session, err := idclient.SessionInstance(c.Session)

	if err != nil || !session.IsLogin() {
		set.Error = &types.ErrorMeta{idsapi.ErrCodeUnauthorized, "Access Denied"}
		return
	}

	if err := c.Request.JsonDecode(&set); err != nil || set.Meta.ID == "" {
		set.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "InvalidArgument"}
		return
	}

	var prev idsapi.AppInstance
	var prevVersion uint64
	if obj := store.BtAgent.ObjectGet("/global/ids/app-instance/" + set.Meta.ID); obj.Error == nil {
		obj.JsonDecode(&prev)
		prevVersion = obj.Meta.Version
	}

	if prev.Meta.ID == "" || prevVersion < 1 {
		set.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "App Instance Not Found"}
		return
	}

	if prev.Meta.UserID != session.UserID {
		set.Error = &types.ErrorMeta{idsapi.ErrCodeUnauthorized, "Access Denied"}
		return
	}

	if set.AppTitle != prev.AppTitle || set.Url != prev.Url {

		prev.Meta.Updated = utilx.TimeNow("atom")
		prev.AppTitle = set.AppTitle
		prev.Url = set.Url

		if obj := store.BtAgent.ObjectSet("/global/ids/app-instance/"+set.Meta.ID, prev, &btapi.ObjectWriteOptions{
			PrevVersion: prevVersion,
		}); obj.Error != nil {
			set.Error = &types.ErrorMeta{idsapi.ErrCodeInternalError, obj.Error.Message}
			return
		}
	}

	set.Kind = "AppInstance"
}
