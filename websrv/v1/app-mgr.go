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

package v1

import (
	"github.com/lessos/bigtree/btapi"

	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/types"
	"github.com/lessos/lessgo/utilx"

	"github.com/lessos/iam/iamapi"
	"github.com/lessos/iam/iamclient"
	"github.com/lessos/iam/store"
)

const (
	appMgrInstPageLimit = 100
)

type AppMgr struct {
	*httpsrv.Controller
}

func (c AppMgr) InstListAction() {

	ls := iamapi.AppInstanceList{}

	defer c.RenderJson(&ls)

	if !iamclient.SessionAccessAllowed(c.Session, "sys.admin", "df085c6dc6ff") {
		ls.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	if objs := store.BtAgent.ObjectList("/global/iam/app-instance/"); objs.Error == nil {

		for _, obj := range objs.Items {

			var inst iamapi.AppInstance
			if err := obj.JsonDecode(&inst); err == nil {

				ls.Items = append(ls.Items, inst)
			}
		}
	}

	// TODO Query

	ls.Kind = "AppInstanceList"
}

func (c AppMgr) InstEntryAction() {

	set := iamapi.AppInstance{}

	defer c.RenderJson(&set)

	if !iamclient.SessionAccessAllowed(c.Session, "sys.admin", "df085c6dc6ff") {
		set.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	if obj := store.BtAgent.ObjectGet("/global/iam/app-instance/" + c.Params.Get("instid")); obj.Error == nil {
		obj.JsonDecode(&set)
	}

	if set.Meta.ID == "" {
		set.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "App Instance Not Found"}
		return
	}

	set.Kind = "AppInstance"
}

func (c AppMgr) InstSetAction() {

	set := iamapi.AppInstance{}

	defer c.RenderJson(&set)

	if !iamclient.SessionAccessAllowed(c.Session, "sys.admin", "df085c6dc6ff") {
		set.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	if err := c.Request.JsonDecode(&set); err != nil || set.Meta.ID == "" {
		set.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "InvalidArgument"}
		return
	}

	var prev iamapi.AppInstance
	var prevVersion uint64
	if obj := store.BtAgent.ObjectGet("/global/iam/app-instance/" + set.Meta.ID); obj.Error == nil {
		obj.JsonDecode(&prev)
		prevVersion = obj.Meta.Version
	}

	if prev.Meta.ID == "" || prevVersion < 1 {
		set.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "App Instance Not Found"}
		return
	}

	if set.AppTitle != prev.AppTitle || set.Url != prev.Url {
		prev.Meta.Updated = utilx.TimeNow("atom")
		prev.AppTitle = set.AppTitle
		prev.Url = set.Url

		if obj := store.BtAgent.ObjectSet("/global/iam/app-instance/"+set.Meta.ID, prev, &btapi.ObjectWriteOptions{
			PrevVersion: prevVersion,
		}); obj.Error != nil {
			set.Error = &types.ErrorMeta{iamapi.ErrCodeInternalError, obj.Error.Message}
			return
		}
	}

	set.Kind = "AppInstance"
}
