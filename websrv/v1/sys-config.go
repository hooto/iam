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
	"github.com/lessos/lessgo/data/rdo"
	"github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/types"

	"../../config"
	"../../idsapi"
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

	rsp := idsapi.SysConfigList{}

	defer c.RenderJson(&rsp)

	if !c.Session.AccessAllowed("sys.admin") {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeUnavailable, "Service Unavailable"}
		return
	}

	q := base.NewQuerySet().From("ids_sysconfig").Limit(10)
	q.Where.And("key", "service_name").Or("key", "webui_banner_title")
	rs, err := dcn.Base.Query(q)
	if err == nil && len(rs) > 0 {

		for _, v := range rs {
			rsp.Items = rsp.Items.Insert(v.Field("key").String(), v.Field("value").String())
		}
	}

	if val, ok := rsp.Items.Fetch("service_name"); val == "" || !ok {
		rsp.Items = rsp.Items.Insert("service_name", config.Config.ServiceName)
	}

	if val, ok := rsp.Items.Fetch("webui_banner_title"); val == "" || !ok {
		rsp.Items = rsp.Items.Insert("webui_banner_title", config.Config.WebUiBannerTitle)
	}

	rsp.Kind = "SysConfigList"
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

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		sets.Error = &types.ErrorMeta{idsapi.ErrCodeUnavailable, "Service Unavailable"}
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

		set := map[string]interface{}{
			"value":   v.Val,
			"updated": base.TimeNow("datetime"),
		}
		ft := base.NewFilter()
		ft.And("key", v.Key)

		qry := base.NewQuerySet()
		qry.From("ids_sysconfig")
		qry.Where = ft

		if _, err := dcn.Base.Fetch(qry); err == nil {
			_, err = dcn.Base.Update("ids_sysconfig", set, ft)
		} else {
			set["key"] = v.Key
			set["created"] = base.TimeNow("datetime")
			_, err = dcn.Base.Insert("ids_sysconfig", set)
		}
	}

	// config.Config.Refresh()

	sets.Kind = "SysConfigList"
}

func (c SysConfig) MailerAction() {

	rsp := idsapi.SysConfigList{}

	defer c.RenderJson(&rsp)

	if !c.Session.AccessAllowed("sys.admin") {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeUnavailable, "Service Unavailable"}
		return
	}

	q := base.NewQuerySet().From("ids_sysconfig").Limit(10)
	q.Where.And("key.in", cfgMailerKeys...)

	rs, err := dcn.Base.Query(q)
	if err == nil && len(rs) > 0 {
		for _, v := range rs {
			rsp.Items = rsp.Items.Insert(v.Field("key").String(), v.Field("value").String())
		}
	}

	for _, val := range cfgMailerKeys {
		if _, ok := rsp.Items.Fetch(val.(string)); !ok {
			rsp.Items = rsp.Items.Insert(val.(string), "")
		}
	}

	rsp.Kind = "SysConfigList"
}

func (c SysConfig) MailerSetAction() {

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

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		sets.Error = &types.ErrorMeta{idsapi.ErrCodeUnavailable, "Service Unavailable"}
		return
	}

	for _, v := range sets.Items {

		mat := false
		for _, vk := range cfgMailerKeys {
			if vk.(string) == v.Key {
				mat = true
				break
			}
		}

		if !mat {
			continue
		}

		set := map[string]interface{}{
			"value":   v.Val,
			"updated": base.TimeNow("datetime"),
		}

		ft := base.NewFilter()
		ft.And("key", v.Key)

		qry := base.NewQuerySet()
		qry.From("ids_sysconfig")
		qry.Where = ft

		if _, err := dcn.Base.Fetch(qry); err == nil {
			dcn.Base.Update("ids_sysconfig", set, ft)
		} else {
			set["created"] = base.TimeNow("datetime")
			set["key"] = v.Key
			dcn.Base.Insert("ids_sysconfig", set)
		}
	}

	config.Config.Refresh()

	sets.Kind = "SysConfigList"
}
