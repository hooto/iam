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
	"github.com/lessos/lessgo/utils"

	"github.com/lessos/lessids/config"
	"github.com/lessos/lessids/idsapi"
	"github.com/lessos/lessids/store"
)

type SysConfig struct {
	*httpsrv.Controller
}

var (
	cfgGenKeys = []interface{}{
		"service_name",
		"webui_banner_title",
	}
	cfgMailerKeys = []interface{}{
		"mailer_smtp_host",
		"mailer_smtp_port",
		"mailer_smtp_user",
		"mailer_smtp_pass",
	}
)

func (c SysConfig) GeneralAction() {

	ls := idsapi.SysConfigList{}

	defer c.RenderJson(&ls)

	if !c.Session.AccessAllowed("sys.admin") {
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
			case "service_name", "webui_banner_title":
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

	ls.Kind = "SysConfigList"
}

func (c SysConfig) GeneralSetAction() {

	sets := idsapi.SysConfigList{}

	defer c.RenderJson(&sets)

	if err := c.Request.JsonDecode(&sets); err != nil || len(sets.Items) < 1 {
		sets.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "Bad Request"}
		return
	}

	if !c.Session.AccessAllowed("sys.admin") {
		sets.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	for _, v := range sets.Items {

		mat := false
		for _, vk := range cfgGenKeys {
			if vk.(string) == v.Key {
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

	// config.Config.Refresh()

	sets.Kind = "SysConfigList"
}

func (c SysConfig) MailerAction() {

	set := idsapi.SysConfigMailer{}

	defer c.RenderJson(&set)

	if !c.Session.AccessAllowed("sys.admin") {
		set.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	if obj := store.BtAgent.ObjectGet(btapi.ObjectProposal{
		Meta: btapi.ObjectMeta{
			Path: "/sys-config/mailer",
		},
	}); obj.Error == nil {
		obj.JsonDecode(&set)
	}

	set.Kind = "SysConfigMailer"
}

func (c SysConfig) MailerSetAction() {

	var (
		set  idsapi.SysConfigMailer
		prev idsapi.SysConfigMailer
	)

	defer c.RenderJson(&set)

	if err := c.Request.JsonDecode(&set); err != nil {
		set.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "Bad Request"}
		return
	}

	if !c.Session.AccessAllowed("sys.admin") {
		set.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	var prevVersion uint64
	if obj := store.BtAgent.ObjectGet(btapi.ObjectProposal{
		Meta: btapi.ObjectMeta{
			Path: "/sys-config/mailer",
		},
	}); obj.Error == nil {
		obj.JsonDecode(&prev)
		prevVersion = obj.Meta.Version

		if set.SmtpHost == "" {
			set.SmtpHost = prev.SmtpHost
		}
		if set.SmtpPort == "" {
			set.SmtpPort = prev.SmtpPort
		}
		if set.SmtpUser == "" {
			set.SmtpUser = prev.SmtpUser
		}
		if set.SmtpPass == "" {
			set.SmtpPass = prev.SmtpPass
		}
	}

	setjs, _ := utils.JsonEncode(set)

	if obj := store.BtAgent.ObjectSet(btapi.ObjectProposal{
		Meta: btapi.ObjectMeta{
			Path: "/sys-config/mailer",
		},
		Data:        setjs,
		PrevVersion: prevVersion,
	}); obj.Error != nil {
		set.Error = &types.ErrorMeta{"500", obj.Error.Message}
		return
	}

	config.Config.Refresh()

	set.Kind = "SysConfigMailer"
}
