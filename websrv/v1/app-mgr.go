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
	"strings"

	"github.com/hooto/httpsrv"
	"github.com/lessos/lessgo/types"

	"github.com/hooto/iam/config"
	"github.com/hooto/iam/data"
	"github.com/hooto/iam/iamapi"
	"github.com/hooto/iam/iamclient"
)

const (
	appMgrInstPageLimit = 100
)

type AppMgr struct {
	*httpsrv.Controller
}

func (c AppMgr) InstListAction() {

	ls := types.ObjectList{}
	defer c.RenderJson(&ls)

	if !iamclient.SessionAccessAllowed(c.Session, "sys.admin", config.Config.InstanceID) {
		ls.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
		return
	}

	var (
		qt = strings.ToLower(c.Params.Value("qry_text"))
	)

	// TODO page
	offset := iamapi.ObjKeyAppInstance("")
	cutset := iamapi.ObjKeyAppInstance("zzzzzzzz")
	if rs := data.Data.NewRanger(offset, cutset).
		SetRevert(true).SetLimit(1000).Exec(); rs.OK() {

		for _, obj := range rs.Items {

			var inst iamapi.AppInstance
			if err := obj.JsonDecode(&inst); err == nil {

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

	var set struct {
		types.TypeMeta
		iamapi.AppInstance
	}
	defer c.RenderJson(&set)

	if !iamclient.SessionAccessAllowed(c.Session, "sys.admin", config.Config.InstanceID) {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
		return
	}

	if obj := data.Data.NewReader(iamapi.ObjKeyAppInstance(c.Params.Value("instid"))).Exec(); obj.OK() {
		obj.Item().JsonDecode(&set.AppInstance)
	}

	if set.Meta.ID == "" {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "App Instance Not Found")
		return
	}

	set.Kind = "AppInstance"
}

func (c AppMgr) InstSetAction() {

	var set struct {
		types.TypeMeta
		iamapi.AppInstance
	}
	defer c.RenderJson(&set)

	if !iamclient.SessionAccessAllowed(c.Session, "sys.admin", config.Config.InstanceID) {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
		return
	}

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

func (c AppMgr) InstDelAction() {

	var set types.TypeMeta
	defer c.RenderJson(&set)

	inst_id := c.Params.Value("inst_id")

	if !iamclient.SessionAccessAllowed(c.Session, "sys.admin", config.Config.InstanceID) {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
		return
	}

	var prev iamapi.AppInstance
	if obj := data.Data.NewReader(iamapi.ObjKeyAppInstance(inst_id)).Exec(); obj.OK() {
		obj.Item().JsonDecode(&prev)
	}

	if prev.Meta.ID == "" {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "App Instance Not Found")
		return
	}

	if obj := data.Data.NewDeleter(iamapi.ObjKeyAppInstance(inst_id)).Exec(); !obj.OK() {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInternalError, obj.ErrorMessage())
		return
	}

	set.Kind = "AppInstance"
}
