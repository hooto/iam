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
	"strings"

	"code.hooto.com/lynkdb/iomix/skv"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/types"
	"github.com/lessos/lessgo/utilx"

	"code.hooto.com/lessos/iam/config"
	"code.hooto.com/lessos/iam/iamapi"
	"code.hooto.com/lessos/iam/iamclient"
	"code.hooto.com/lessos/iam/store"
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

	if !iamclient.SessionAccessAllowed(c.Session, "sys.admin", config.Config.InstanceID) {
		ls.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
		return
	}

	var (
		qt = strings.ToLower(c.Params.Get("qry_text"))
	)

	// TODO page
	if objs := store.PvScan("app-instance", "", "", 1000); objs.OK() {

		rss := objs.KvList()
		for _, obj := range rss {

			var inst iamapi.AppInstance
			if err := obj.Decode(&inst); err == nil {

				if qt != "" && (!strings.Contains(inst.AppID, qt) &&
					!strings.Contains(inst.AppTitle, qt)) {
					continue
				}

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

	if !iamclient.SessionAccessAllowed(c.Session, "sys.admin", config.Config.InstanceID) {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
		return
	}

	if obj := store.PvGet("app-instance/" + c.Params.Get("instid")); obj.OK() {
		obj.Decode(&set)
	}

	if set.Meta.ID == "" {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "App Instance Not Found")
		return
	}

	set.Kind = "AppInstance"
}

func (c AppMgr) InstSetAction() {

	set := iamapi.AppInstance{}

	defer c.RenderJson(&set)

	if !iamclient.SessionAccessAllowed(c.Session, "sys.admin", config.Config.InstanceID) {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
		return
	}

	if err := c.Request.JsonDecode(&set); err != nil || set.Meta.ID == "" {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "InvalidArgument")
		return
	}

	var prev iamapi.AppInstance
	var prevVersion uint64
	if obj := store.PvGet("app-instance/" + set.Meta.ID); obj.OK() {
		obj.Decode(&prev)
		prevVersion = obj.Meta().Version
	}

	if prev.Meta.ID == "" || prevVersion < 1 {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "App Instance Not Found")
		return
	}

	if set.AppTitle != prev.AppTitle || set.Url != prev.Url {
		prev.Meta.Updated = utilx.TimeNow("atom")
		prev.AppTitle = set.AppTitle
		prev.Url = set.Url

		if obj := store.PvPut("app-instance/"+set.Meta.ID, prev, &skv.PvWriteOptions{
			PrevVersion: prevVersion,
		}); !obj.OK() {
			set.Error = types.NewErrorMeta(iamapi.ErrCodeInternalError, obj.Bytex().String())
			return
		}
	}

	set.Kind = "AppInstance"
}
