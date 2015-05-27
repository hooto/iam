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
	"encoding/base64"
	"io"
	"strings"
	"time"

	"github.com/lessos/lessgo/data/rdo"
	"github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/pass"
	"github.com/lessos/lessgo/types"
	"github.com/lessos/lessgo/utils"

	"../../base/role"
	"../../idsapi"
)

type Service struct {
	*httpsrv.Controller
}

func (c Service) LoginAuthAction() {

	rsp := idsapi.ServiceLoginAuth{
		Continue: "/ids",
	}

	defer c.RenderJson(&rsp)

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Error = &types.ErrorMeta{"400", "Internal Server Error"}
		return
	}

	if c.Params.Get("email") == "" || c.Params.Get("passwd") == "" {
		rsp.Error = &types.ErrorMeta{"400", "Bad Request"}
		return
	}

	email := strings.ToLower(c.Params.Get("email"))

	q := base.NewQuerySet().From("ids_login").Limit(1)
	q.Where.And("email", email)
	rsu, err := dcn.Base.Query(q)
	if err == nil && len(rsu) == 0 {
		rsp.Error = &types.ErrorMeta{"400", "Email or Password can not match"}
		return
	}

	if !pass.Check(c.Params.Get("passwd"), rsu[0].Field("pass").String()) {
		rsp.Error = &types.ErrorMeta{"400", "Bad Email or Password can not match"}
		return
	}

	rsp.AccessToken = utils.StringNewRand(24)

	addr := "127.0.0.1"
	if addridx := strings.Index(c.Request.RemoteAddr, ":"); addridx > 0 {
		addr = c.Request.RemoteAddr[:addridx]
	}

	session := map[string]interface{}{
		"token":    rsp.AccessToken,
		"refresh":  utils.StringNewRand(24),
		"status":   1,
		"uid":      rsu[0].Field("uid").String(),
		"uname":    rsu[0].Field("uname").String(),
		"name":     rsu[0].Field("name").String(),
		"roles":    rsu[0].Field("roles").String(),
		"timezone": rsu[0].Field("timezone").String(),
		"source":   addr,
		"created":  base.TimeNow("datetime"),                // TODO
		"expired":  base.TimeNowAdd("datetime", "+864000s"), // TODO
	}
	if _, err := dcn.Base.Insert("ids_sessions", session); err != nil {
		rsp.Error = &types.ErrorMeta{"500", "Can not write to database" + err.Error()}
		return
	}

	if len(c.Params.Get("continue")) > 0 {
		rsp.Continue = c.Params.Get("continue")
		if strings.Index(rsp.Continue, "?") == -1 {
			rsp.Continue += "?"
		} else {
			rsp.Continue += "&"
		}
		rsp.Continue += "access_token=" + rsp.AccessToken
	}

	rsp.Kind = "ServiceLoginAuth"
}

func (c Service) AuthAction() {

	rsp := idsapi.UserSession{}

	defer c.RenderJson(&rsp)

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Error = &types.ErrorMeta{"500", "Internal Server Error"}
		return
	}

	if c.Params.Get("access_token") == "" {
		rsp.Error = &types.ErrorMeta{"401", "Unauthorized"}
		return
	}

	q := base.NewQuerySet().From("ids_sessions").Limit(1)
	q.Where.And("token", c.Params.Get("access_token"))
	rss, err := dcn.Base.Query(q)
	if err == nil && len(rss) == 0 {
		rsp.Error = &types.ErrorMeta{"401", "Unauthorized"}
		return
	}

	//
	expired := rss[0].Field("expired").TimeParse("datetime")
	if expired.Before(time.Now()) {
		rsp.Error = &types.ErrorMeta{"401", "Unauthorized"}
		return
	}
	//fmt.Println("expired", expired)

	//
	addr := "0.0.0.0"
	if addridx := strings.Index(c.Request.RemoteAddr, ":"); addridx > 0 {
		addr = c.Request.RemoteAddr[:addridx]
	}
	if addr != rss[0].Field("source").String() {
		rsp.Error = &types.ErrorMeta{"401", "Unauthorized"}
		return
	}

	//
	rsp.AccessToken = rss[0].Field("token").String()
	rsp.RefreshToken = rss[0].Field("refresh").String()

	rsp.UserID = rss[0].Field("uid").String()
	rsp.UserName = rss[0].Field("uname").String()
	rsp.Name = rss[0].Field("name").String()
	rsp.Roles = rss[0].Field("roles").String()
	rsp.Timezone = rss[0].Field("timezone").String()
	rsp.Expired = base.TimeZoneFormat(expired, rsp.Timezone, "atom")

	rsp.Kind = "UserSession"
}

func (c Service) AccessAllowedAction() {

	var (
		req idsapi.UserAccessEntry
		rsp idsapi.UserAccessEntry
	)

	defer c.RenderJson(&rsp)

	if len(c.Request.RawBody) == 0 {
		rsp.Error = &types.ErrorMeta{"401", "Unauthorized"}
		return
	}

	if err := c.Request.JsonDecode(&req); err != nil {
		rsp.Error = &types.ErrorMeta{"401", err.Error()}
		return
	}
	if req.AccessToken == "" {
		rsp.Error = &types.ErrorMeta{"401", "Unauthorized"}
		return
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Error = &types.ErrorMeta{"500", "Internal Server Error"}
		return
	}

	q := base.NewQuerySet().From("ids_sessions").Limit(1)
	q.Where.And("token", req.AccessToken)
	rss, err := dcn.Base.Query(q)
	if err == nil && len(rss) == 0 {
		rsp.Error = &types.ErrorMeta{"401", "Unauthorized"}
		return
	}

	//
	expired := rss[0].Field("expired").TimeParse("datetime")
	if expired.Before(time.Now()) {
		rsp.Error = &types.ErrorMeta{"401", "Unauthorized"}
		return
	}

	//
	addr := "0.0.0.0"
	if addridx := strings.Index(c.Request.RemoteAddr, ":"); addridx > 0 {
		addr = c.Request.RemoteAddr[:addridx]
	}

	if addr != rss[0].Field("source").String() {
		rsp.Error = &types.ErrorMeta{"401", "Unauthorized"}
		return
	}

	if !role.AccessAllowed(rss[0].Field("roles").String(), req.InstanceID, req.Privilege) {
		rsp.Error = &types.ErrorMeta{"401", "Unauthorized"}
		return
	}

	rsp.Kind = "UserAccessEntry"
}

func (c Service) PhotoAction() {

	c.AutoRender = false

	photo := "data:image/png;base64,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"

	if len(c.Request.RequestPath) > 14 {

		uid := c.Request.RequestPath[14:]
		dcn, err := rdo.ClientPull("def")
		if err != nil {
			return
		}

		q := base.NewQuerySet().From("ids_profile").Select("photo").Limit(1)
		q.Where.And("uid", uid)
		rsp, err := dcn.Base.Query(q)
		if err == nil && len(rsp) == 1 {
			if len(rsp[0].Field("photo").String()) > 50 {
				photo = rsp[0].Field("photo").String()
			}
		}
	}

	body64 := strings.SplitAfter(photo, ";base64,")
	if len(body64) != 2 {
		return
	}
	data, err := base64.StdEncoding.DecodeString(body64[1])
	if err != nil {
		return
	}

	io.WriteString(c.Response.Out, string(data))
}
