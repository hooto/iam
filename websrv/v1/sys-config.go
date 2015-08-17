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

	"github.com/lessos/lessids/config"
	"github.com/lessos/lessids/idsapi"
	"github.com/lessos/lessids/store"
	"github.com/lessos/lessids/idclient"
)

type SysConfig struct {
	*httpsrv.Controller
}

var (
	cfgGenKeys = []string{
		"service_name",
		"webui_banner_title",
		"user_reg_disable",
		"mailer",
	}
)

func (c SysConfig) GeneralAction() {

	ls := idsapi.SysConfigList{}

	defer c.RenderJson(&ls)

	if !idclient.SessionAccessAllowed(c.Session, "sys.admin", "df085c6dc6ff") {
		ls.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	if objs := store.BtAgent.ObjectList(btapi.ObjectProposal{
		Meta: btapi.ObjectMeta{
			Path: "/sys-config/",
		},
	}); objs.Error == nil {

		for _, obj := range objs.Items {

			switch obj.Name() {
			case "service_name", "webui_banner_title", "user_reg_disable":
				ls.Items = ls.Items.Insert(obj.Name(), obj.Data)
			}
		}
	}

	if val, ok := ls.Items.Fetch("service_name"); val == "" || !ok {
		ls.Items = ls.Items.Insert("service_name", config.Config.ServiceName)
	}

	if val, ok := ls.Items.Fetch("webui_banner_title"); val == "" || !ok {
		ls.Items = ls.Items.Insert("webui_banner_title", config.Config.WebUiBannerTitle)
	}

	if val, ok := ls.Items.Fetch("user_reg_disable"); !ok || val == "" {
		ls.Items = ls.Items.Insert("user_reg_disable", "0")
	}

	ls.Kind = "SysConfigList"
}

func (c SysConfig) GeneralSetAction() {

	sets := idsapi.SysConfigList{}

	defer c.RenderJson(&sets)

	if err := c.Request.JsonDecode(&sets); err != nil || len(sets.Items) < 1 {
		sets.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "Bad Request"}
		return
	}

	if !idclient.SessionAccessAllowed(c.Session, "sys.admin","df085c6dc6ff") {
		sets.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	for _, v := range sets.Items {

		mat := false
		for _, vk := range cfgGenKeys {
			if vk == v.Key {
				mat = true
				break
			}
		}
		if !mat {
			continue
		}

		var prevVersion uint64

		if obj := store.BtAgent.ObjectGet(btapi.ObjectProposal{
			Meta: btapi.ObjectMeta{
				Path: "/sys-config/" + v.Key,
			},
		}); obj.Error == nil {
			prevVersion = obj.Meta.Version
		}

		if obj := store.BtAgent.ObjectSet(btapi.ObjectProposal{
			Meta: btapi.ObjectMeta{
				Path: "/sys-config/" + v.Key,
			},
			Data:        v.Val,
			PrevVersion: prevVersion,
		}); obj.Error != nil {
			sets.Error = &types.ErrorMeta{"500", obj.Error.Message}
			return
		}
	}

	config.Config.Refresh()

	sets.Kind = "SysConfigList"
}

func (c SysConfig) MailerAction() {

	ls := idsapi.SysConfigList{}

	defer c.RenderJson(&ls)

	if !idclient.SessionAccessAllowed(c.Session, "sys.admin","df085c6dc6ff") {
		ls.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	if obj := store.BtAgent.ObjectGet(btapi.ObjectProposal{
		Meta: btapi.ObjectMeta{
			Path: "/sys-config/mailer",
		},
	}); obj.Error == nil {

		ls.Items = ls.Items.Insert("mailer", obj.Data)
	}

	if val, ok := ls.Items.Fetch("mailer"); val == "" || !ok {
		ls.Items = ls.Items.Insert("mailer", "{}")
	}

	ls.Kind = "SysConfigList"
}
