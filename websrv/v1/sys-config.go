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

	"code.hooto.com/lessos/iam/config"
	"code.hooto.com/lessos/iam/iamapi"
	"code.hooto.com/lessos/iam/iamclient"
	"code.hooto.com/lessos/iam/store"
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

	if !iamclient.SessionAccessAllowed(c.Session, "sys.admin", config.Config.InstanceID) {
		ls.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	if objs := store.PvScan("sys-config/", "", "", 1000); objs.OK() {

		rss := objs.KvList()
		for _, obj := range rss {

			switch obj.Meta().Name {
			case "service_name", "webui_banner_title", "user_reg_disable":
				ls.Items.Set(obj.Meta().Name, obj.Bytex().String())
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

	if !iamclient.SessionAccessAllowed(c.Session, "sys.admin", config.Config.InstanceID) {
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

		if obj := store.PvGet("sys-config/" + v.Name); obj.OK() {
			prevVersion = obj.Meta().Version
		}

		if obj := store.PvPut("sys-config/"+v.Name, v.Value, &skv.PvWriteOptions{
			PrevVersion: prevVersion,
		}); !obj.OK() {
			sets.Error = &types.ErrorMeta{"500", obj.Bytex().String()}
			return
		}
	}

	store.SysConfigRefresh()

	sets.Kind = "SysConfigList"
}

func (c SysConfig) MailerAction() {

	ls := iamapi.SysConfigList{}

	defer c.RenderJson(&ls)

	if !iamclient.SessionAccessAllowed(c.Session, "sys.admin", config.Config.InstanceID) {
		ls.Error = &types.ErrorMeta{iamapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	if obj := store.PvGet("sys-config/mailer"); obj.OK() {

		ls.Items.Set("mailer", obj.Data)
	}

	if val, ok := ls.Items.Get("mailer"); val.String() == "" || !ok {
		ls.Items.Set("mailer", "{}")
	}

	ls.Kind = "SysConfigList"
}
