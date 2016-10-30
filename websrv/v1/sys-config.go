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

	"github.com/lessos/iam/config"
	"github.com/lessos/iam/iamapi"
	"github.com/lessos/iam/iamclient"
	"github.com/lessos/iam/store"
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

	ls := iamapi.SysConfigList{}

	defer c.RenderJson(&ls)

	if !iamclient.SessionAccessAllowed(c.Session, "sys.admin", "df085c6dc6ff") {
		ls.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	if objs := store.BtAgent.ObjectList("/global/iam/sys-config/"); objs.Error == nil {

		for _, obj := range objs.Items {

			switch obj.Meta.Name {
			case "service_name", "webui_banner_title", "user_reg_disable":
				ls.Items.Set(obj.Meta.Name, obj.Data)
			}
		}
	}

	if val, ok := ls.Items.Get("service_name"); val.String() == "" || !ok {
		ls.Items.Set("service_name", config.Config.ServiceName)
	}

	if val, ok := ls.Items.Get("webui_banner_title"); val.String() == "" || !ok {
		ls.Items.Set("webui_banner_title", config.Config.WebUiBannerTitle)
	}

	if val, ok := ls.Items.Get("user_reg_disable"); !ok || val.String() == "" {
		ls.Items.Set("user_reg_disable", "0")
	}

	ls.Kind = "SysConfigList"
}

func (c SysConfig) GeneralSetAction() {

	sets := iamapi.SysConfigList{}

	defer c.RenderJson(&sets)

	if err := c.Request.JsonDecode(&sets); err != nil || len(sets.Items) < 1 {
		sets.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "Bad Request"}
		return
	}

	if !iamclient.SessionAccessAllowed(c.Session, "sys.admin", "df085c6dc6ff") {
		sets.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	for _, v := range sets.Items {

		mat := false
		for _, vk := range cfgGenKeys {
			if vk == v.Name {
				mat = true
				break
			}
		}
		if !mat {
			continue
		}

		var prevVersion uint64

		if obj := store.BtAgent.ObjectGet("/global/iam/sys-config/" + v.Name); obj.Error == nil {
			prevVersion = obj.Meta.Version
		}

		if obj := store.BtAgent.ObjectSet("/global/iam/sys-config/"+v.Name, v.Value, &btapi.ObjectWriteOptions{
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

	ls := iamapi.SysConfigList{}

	defer c.RenderJson(&ls)

	if !iamclient.SessionAccessAllowed(c.Session, "sys.admin", "df085c6dc6ff") {
		ls.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	if obj := store.BtAgent.ObjectGet("/global/iam/sys-config/mailer"); obj.Error == nil {

		ls.Items.Set("mailer", obj.Data)
	}

	if val, ok := ls.Items.Get("mailer"); val.String() == "" || !ok {
		ls.Items.Set("mailer", "{}")
	}

	ls.Kind = "SysConfigList"
}
