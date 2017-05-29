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

package v1

import (
	"code.hooto.com/lynkdb/iomix/skv"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/types"

	"code.hooto.com/lessos/iam/iamapi"
	"code.hooto.com/lessos/iam/iamclient"
	"code.hooto.com/lessos/iam/store"
)

const (
	myAppInstPageLimit = 100
)

type MyApp struct {
	*httpsrv.Controller
	us iamapi.UserSession
}

func (c *MyApp) Init() int {

	//
	c.us, _ = iamclient.SessionInstance(c.Session)

	if !c.us.IsLogin() {
		c.Response.Out.WriteHeader(401)
		c.RenderJson(types.NewTypeErrorMeta(iamapi.ErrCodeUnauthorized, "Unauthorized"))
		return 1
	}

	return 0
}

func (c MyApp) InstListAction() {

	ls := iamapi.AppInstanceList{}

	defer c.RenderJson(&ls)

	if objs := store.PoScan("app-instance", []byte{}, []byte{}, 1000); objs.OK() {

		rss := objs.KvList()
		for _, obj := range rss {

			var inst iamapi.AppInstance
			if err := obj.Decode(&inst); err == nil {

				if inst.Meta.User == c.us.UserName {
					ls.Items = append(ls.Items, inst)
				}
			}
		}
	}

	ls.Kind = "AppInstanceList"
}

func (c MyApp) InstEntryAction() {

	set := iamapi.AppInstance{}

	defer c.RenderJson(&set)

	if obj := store.PoGet("app-instance", c.Params.Get("instid")); obj.OK() {
		obj.Decode(&set)
	}

	if set.Meta.ID == "" {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "App Instance Not Found")
		return
	}

	if set.Meta.User != c.us.UserName {
		set = iamapi.AppInstance{}
		set.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "Access Denied")
		return
	}

	set.Kind = "AppInstance"
}

func (c MyApp) InstSetAction() {

	set := iamapi.AppInstance{}
	defer c.RenderJson(&set)

	if err := c.Request.JsonDecode(&set); err != nil || set.Meta.ID == "" {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "InvalidArgument")
		return
	}

	var prev iamapi.AppInstance
	var prevVersion uint64
	if obj := store.PoGet("app-instance", set.Meta.ID); obj.OK() {
		obj.Decode(&prev)
		prevVersion = obj.Meta().Version
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

		if obj := store.PoPut("app-instance", set.Meta.ID, prev, &skv.PathWriteOptions{
			PrevVersion: prevVersion,
		}); !obj.OK() {
			set.Error = types.NewErrorMeta(iamapi.ErrCodeInternalError, obj.Bytex().String())
			return
		}
	}

	set.Kind = "AppInstance"
}
