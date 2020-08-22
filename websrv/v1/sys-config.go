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

	"github.com/hooto/iam/config"
	"github.com/hooto/iam/data"
	"github.com/hooto/iam/iamapi"
	"github.com/hooto/iam/iamclient"
)

type SysConfig struct {
	*httpsrv.Controller
}

var (
	cfgGenKeys = []string{
		"service_name",
		"webui_banner_title",
		"user_reg_disable",
		"service_login_form_alert_msg",
		"mailer",
	}
)

func (c SysConfig) GeneralAction() {

	ls := iamapi.SysConfigList{}
	defer c.RenderJson(&ls)

	if !iamclient.SessionAccessAllowed(c.Session, "sys.admin", config.Config.InstanceID) {
		ls.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
		return
	}

	if rs := data.Data.NewReader(nil).KeyRangeSet(
		iamapi.ObjKeySysConfig(""), iamapi.ObjKeySysConfig("")).LimitNumSet(1000).Query(); rs.OK() {

		for _, obj := range rs.Items {

			switch string(obj.Meta.Key) {
			case "service_name", "webui_banner_title", "user_reg_disable", "service_login_form_alert_msg":
				ls.Items.Set(string(obj.Meta.Key), obj.DataValue().String())
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

	if _, ok := ls.Items.Get("service_login_form_alert_msg"); !ok {
		ls.Items.Set("service_login_form_alert_msg", config.Config.ServiceLoginFormAlertMsg)
	}

	ls.Kind = "SysConfigList"
}

func (c SysConfig) GeneralSetAction() {

	sets := iamapi.SysConfigList{}

	defer c.RenderJson(&sets)

	if err := c.Request.JsonDecode(&sets); err != nil || len(sets.Items) < 1 {
		sets.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Bad Request")
		return
	}

	if !iamclient.SessionAccessAllowed(c.Session, "sys.admin", config.Config.InstanceID) {
		sets.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
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

		if v.Value == "" {
			data.Data.NewWriter(iamapi.ObjKeySysConfig(v.Name), nil).
				ModeDeleteSet(true).Commit()
			if v.Name == "service_login_form_alert_msg" {
				config.Config.ServiceLoginFormAlertMsg = ""
			}
		} else {

			if obj := data.Data.NewWriter(iamapi.ObjKeySysConfig(v.Name), v.Value).
				Commit(); !obj.OK() {
				sets.Error = types.NewErrorMeta("500", "DB ERROR #2 "+obj.Message)
				return
			}
		}
	}

	data.SysConfigRefresh()

	sets.Kind = "SysConfigList"
}

func (c SysConfig) MailerAction() {

	ls := iamapi.SysConfigList{}

	defer c.RenderJson(&ls)

	if !iamclient.SessionAccessAllowed(c.Session, "sys.admin", config.Config.InstanceID) {
		ls.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Access Denied")
		return
	}

	if obj := data.Data.NewReader(iamapi.ObjKeySysConfig("mailer")).Query(); obj.OK() {
		ls.Items.Set("mailer", obj.DataValue().String())
	}

	if val, ok := ls.Items.Get("mailer"); val.String() == "" || !ok {
		ls.Items.Set("mailer", "{}")
	}

	ls.Kind = "SysConfigList"
}
