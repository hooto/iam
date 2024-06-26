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

	keyPrefix := iamapi.ObjKeySysConfig("")

	if rs := data.Data.NewRanger(keyPrefix, iamapi.ObjKeySysConfig("zzzz")).SetLimit(1000).Exec(); rs.OK() {

		for _, obj := range rs.Items {

			if len(obj.Key) <= len(keyPrefix) {
				continue
			}

			k := string(obj.Key[len(keyPrefix):])

			switch k {
			case "service_name", "webui_banner_title", "user_reg_disable", "service_login_form_alert_msg":
				ls.Items.Set(k, obj.StringValue())
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
			data.Data.NewDeleter(iamapi.ObjKeySysConfig(v.Name)).Exec()
			if v.Name == "service_login_form_alert_msg" {
				config.Config.ServiceLoginFormAlertMsg = ""
			}
		} else {

			if obj := data.Data.NewWriter(iamapi.ObjKeySysConfig(v.Name), nil).SetJsonValue(v.Value).
				Exec(); !obj.OK() {
				sets.Error = types.NewErrorMeta("500", "DB ERROR #2 "+obj.ErrorMessage())
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

	if obj := data.Data.NewReader(iamapi.ObjKeySysConfig("mailer")).Exec(); obj.OK() {
		ls.Items.Set("mailer", obj.Item().StringValue())
	}

	if val, ok := ls.Items.Get("mailer"); val.String() == "" || !ok {
		ls.Items.Set("mailer", "{}")
	}

	ls.Kind = "SysConfigList"
}
