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

	if !c.Session.AccessAllowed("user.admin") {
		ls.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	if objs := store.BtAgent.ObjectList(btapi.ObjectProposal{
		Meta: btapi.ObjectMeta{
			Path: "/app-instance/",
		},
	}); objs.Error == nil {

		for _, obj := range objs.Items {

			var inst idsapi.AppInstance
			if err := obj.JsonDecode(&inst); err == nil {

				ls.Items = append(ls.Items, inst)
			}
		}
	}

	// TODO Query

	ls.Kind = "AppInstanceList"
}

func (c MyApp) InstEntryAction() {

	set := idsapi.AppInstance{}

	defer c.RenderJson(&set)

	if !c.Session.AccessAllowed("user.admin") {
		set.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	if obj := store.BtAgent.ObjectGet(btapi.ObjectProposal{
		Meta: btapi.ObjectMeta{
			Path: "/app-instance/" + c.Params.Get("instid"),
		},
	}); obj.Error == nil {
		obj.JsonDecode(&set)
	}

	if set.Meta.ID == "" {
		set.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "App Instance Not Found"}
		return
	}

	// TODO set.Privileges

	set.Kind = "AppInstance"
}

// func (c MyApp) InstSaveAction() {

// 	c.AutoRender = false

// 	var set ResponseJson
// 	set.ApiVersion = apiVersion
// 	set.Status = 400
// 	set.Message = "Bad Request"

// 	defer func() {
// 		if setj, err := utils.JsonEncode(set); err == nil {
// 			io.WriteString(c.Response.Out, setj)
// 		}
// 	}()

// 	if !c.Session.AccessAllowed("user.admin") {
// 		return
// 	}

// 	dcn, err := rdo.ClientPull("def")
// 	if err != nil {
// 		set.Message = "Internal Server Error"
// 		return
// 	}

// 	q := base.NewQuerySet().From("ids_instance").Limit(1)

// 	isNew := true
// 	instset := map[string]interface{}{}

// 	if c.Params.Get("instid") != "" {

// 		q.Where.And("id", c.Params.Get("instid"))

// 		rsinst, err := dcn.Base.Query(q)
// 		if err != nil || len(rsinst) == 0 {
// 			set.Status = 400
// 			set.Message = http.StatusText(400)
// 			return
// 		}

// 		isNew = false
// 	}

// 	instset["updated"] = base.TimeNow("datetime")
// 	instset["app_title"] = c.Params.Get("app_title")

// 	if isNew {

// 		// TODO

// 	} else {

// 		instset["status"] = c.Params.Get("status")

// 		frupd := base.NewFilter()
// 		frupd.And("id", c.Params.Get("instid"))
// 		if _, err := dcn.Base.Update("ids_instance", instset, frupd); err != nil {
// 			set.Status = 500
// 			set.Message = "Can not write to database"
// 			return
// 		}
// 	}

// 	set.Status = 200
// 	set.Message = ""
// }
