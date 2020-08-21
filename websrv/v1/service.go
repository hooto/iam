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
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hooto/hlog4g/hlog"
	"github.com/hooto/httpsrv"
	"github.com/lessos/lessgo/pass"
	"github.com/lessos/lessgo/types"

	"github.com/hooto/hauth/go/hauth/v1"
	"github.com/hooto/iam/base/role"
	"github.com/hooto/iam/iamapi"
	"github.com/hooto/iam/store"
)

type Service struct {
	*httpsrv.Controller
}

func urlHost(requrl string) string {

	u, err := url.Parse(requrl)

	if err != nil {
		return "localhost"
	}

	if i := strings.Index(u.Host, ":"); i > 0 {
		return u.Host[:i]
	}

	return u.Host
}

func (c Service) LoginAuthAction() {

	rsp := iamapi.ServiceLoginAuth{
		RedirectUri: "/iam",
	}
	defer c.RenderJson(&rsp)

	if c.Params.Get("passwd") == "" {
		rsp.Error = types.NewErrorMeta("400", "Username or Password can not be empty")
		return
	}

	uname := iamapi.UserNameFilter(c.Params.Get("uname"))
	if err := iamapi.UserNameValid(uname); err != nil {
		rsp.Error = types.NewErrorMeta("400", "incorrect username or password")
		hlog.Printf("info", "service/login-auth fail user %s", uname)
		return
	}

	user := store.UserGet(uname)
	if user == nil {
		rsp.Error = types.NewErrorMeta("400", "incorrect username or password")
		hlog.Printf("info", "service/login-auth fail user %s", uname)
		return
	}

	if user.Type == iamapi.UserTypeGroup {
		rsp.Error = types.NewErrorMeta("400", "incorrect username or password")
		hlog.Printf("info", "service/login-auth fail user %s", uname)
		return
	}

	addr := "127.0.0.1"
	if addridx := strings.Index(c.Request.RemoteAddr, ":"); addridx > 0 {
		addr = c.Request.RemoteAddr[:addridx]
	}

	err_num := 0
	err_key := iamapi.ObjKeyUserAuthDeny(uname, addr)
	if rs := store.Data.NewReader(err_key).Query(); rs.OK() {
		err_num = rs.DataValue().Int()
		if err_num >= 20 {
			rsp.Error = types.NewErrorMeta("400",
				fmt.Sprintf("more than %d times failed to verify this signin, please try again in 1 day later", err_num))
			return
		}
	}

	if auth := user.Keys.Get(iamapi.UserKeyDefault); auth == nil ||
		!pass.Check(c.Params.Get("passwd"), auth.String()) {
		err_num++
		store.Data.NewWriter(err_key, err_num).ExpireSet(86400000).Commit()
		rsp.Error = types.NewErrorMeta("400", "incorrect username or password")
		hlog.Printf("info", "service/login-auth fail user %s", uname)
		return
	}

	var (
		ttl = int64(864000)
		ap  = hauth.NewUserPayload(
			user.Name,
			user.DisplayName,
			user.Roles,
			store.UserGroups(uname),
			ttl,
		)
	)

	if rs := store.Data.NewWriter(iamapi.ObjKeyUserAuth(ap.Id, uint32(ap.Expired)), ap).
		ExpireSet(ttl * 1000).Commit(); !rs.OK() {
		rsp.Error = types.NewErrorMeta("500", rs.Message)
		return
	}

	rsp.AccessToken = ap.SignToken(store.KeyMgr)

	if len(c.Params.Get("redirect_token")) > 20 {

		rt := iamapi.ServiceRedirectTokenDecode(c.Params.Get("redirect_token"))

		if len(rt.RedirectUri) > 0 {

			rsp.RedirectUri = rt.RedirectUri

			if urlHost(rsp.RedirectUri) != urlHost(c.Request.URL.Host) {

				if strings.Index(rsp.RedirectUri, "?") == -1 {
					rsp.RedirectUri += "?"
				} else {
					rsp.RedirectUri += "&"
				}

				rsp.RedirectUri += iamapi.AccessTokenKey + "=" + rsp.AccessToken +
					"&expires_in=" + strconv.Itoa(int(ttl))

				if len(rt.State) > 0 {
					rsp.RedirectUri += "&state=" + url.QueryEscape(rt.State)
				}
			}
		}
	}

	http.SetCookie(c.Response.Out, &http.Cookie{
		Name:     iamapi.AccessTokenKey,
		Value:    rsp.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(time.Duration(ttl) * time.Second),
	})

	rsp.Kind = "ServiceLoginAuth"

	hlog.Printf("info", "service/login-auth user %s", user.Name)
}

func (c Service) AuthAction() {

	var set types.TypeMeta
	defer c.RenderJson(&set)

	token := c.Params.Get(iamapi.AccessTokenKey)

	if _, err := hauth.UserValid(token, store.KeyMgr); err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "Unauthorized")
		return
	}

	//
	// addr := "0.0.0.0"
	// if addridx := strings.Index(c.Request.RemoteAddr, ":"); addridx > 0 {
	// 	addr = c.Request.RemoteAddr[:addridx]
	// }
	// if addr != set.ClientAddr {
	// 	set.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "Unauthorized")
	// 	return
	// }

	set.Kind = "UserSession"
}

func (c Service) AccessAllowedAction() {

	var (
		req iamapi.UserAccessEntry
		rsp iamapi.UserAccessEntry
	)
	defer c.RenderJson(&rsp)

	if len(c.Request.RawBody) == 0 {
		rsp.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "Unauthorized")
		return
	}

	// fmt.Println("AccessAllowedAction", string(c.Request.RawBody))

	if err := c.Request.JsonDecode(&req); err != nil {
		rsp.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, err.Error())
		return
	}

	ap, err := hauth.UserValid(req.AccessToken, store.KeyMgr)
	if err != nil {
		rsp.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "Unauthorized")
		return
	}

	//
	// addr := "0.0.0.0"
	// if addridx := strings.Index(c.Request.RemoteAddr, ":"); addridx > 0 {
	// 	addr = c.Request.RemoteAddr[:addridx]
	// }
	// if addr != session.ClientAddr {
	// 	rsp.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "Unauthorized")
	// 	return
	// }

	if !role.AccessAllowed(ap.Id, ap.Roles, req.InstanceID, req.Privilege) {
		rsp.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "Unauthorized")
		return
	}

	rsp.Kind = "UserAccessEntry"
}

var (
	iam_v1_service_len = len("iam/v1/service/photo/")
	iam_v1_service_def = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAGAAAABgCAYAAADimHc4AAAAAXNSR0IArs4c6QAAAAZiS0dEAP8A/wD/oL2nkwAAAAlwSFlzAAALEwAACxMBAJqcGAAAAAd0SU1FB94EBRAHIE63lHIAAAAZdEVYdENvbW1lbnQAQ3JlYXRlZCB3aXRoIEdJTVBXgQ4XAAAFW0lEQVR42u2cwWtUVxSHv5moJBHEmBasMc0mdpGFRaHRRZMuBDe6qEVadaN0IVKIIoKC+jeopV10VWhRpF0FQSXZVpRsivgHhLGGlJYWQ+LCaJzpYs4j05d73ptJZua+ue988JhhZt49v3vOm3Pvu+/eC4ZhGIZhGIZhGIZhGIZhGIZhhE6hg7RuA/qB/cCnwBjwHlCU78vAP8BvwCPgd+BfYNHCvDEOAd8DM0ClwWNGzj1kbmycYeAZsFTj0LLyPu37JSlr2NyaTg9wNcGxjR7xc6+KDcPBx8ATxXHlBoJRTinjidgyathdh5Pjn78C5uR4Vee/ptbGbnN7lcE6cvwy8By4Iz0gjTH5zXM5J62NGMy780cTrth38noPOLGOsk/IubVluQIymlfnbwNKKanmc6BrAza6pIyk1FQSLbnjRo1T4vn/7yb3VnqkzHg7EL2/kTfnH0tw/jww0AKbA1K2FoRjeQrAbELO726h3e6ENmE2L86fcDggugqH2mB/SGkP3om2tlH04PzNwEGH7QLwq6SIVjMvtgoOfxwUjcGyS+n1LFAd6WwX+8Wmq1e0K+QAHFXuUJ960PJU0XI05ADcUXL/SQ9aTiptwZ12CfDxQKYSey141OJdT9HjP6FQU8l5jzrmHXqC7QUNKFdfyWMASjEtmtYgAtCjpMAHHgPwQEk5PSEGoKJ8/tZjAN42qLWjA1BUKulzTH5QcbjP9rFl9Cn9bp9jMLOKpr5Q7wNc/e4Fj3oWlPuSYHmpVLjLg5Yu5YJ4GXIArrB2HL4C3PWg5W7sQoj0XAk5AAdwTx3xcS9QUrQcCDkAO4AXjjT0FjjVRh2nxGY8/bwQjUFzU8m9c22q/A6x5WqLbpIDPqM6X9M1T+dcG+yfwz1PaEm05YKH6FMJT7fQ7mn06YsPyRGbcT8cL8tnZ1pg84yU7XoWXCHwR5EuxnE/nowcNAb0NsFOr5Tl6v5Gxzg5pEj14bhrKCC6KqeAfRuwsU/K0KaiVERDkRyzqDTI5Zqhij+AkQbKHJFzFhJyfhlbwvS/ICTN8a913M/ATuDD2LFTvqunjIo5fy2PUlJFIytlyimp7ZG5ey1bcM8ZbcYKmfgc0C3mbp29rH9NWNqxN0sV7cqY4z+RvvoPrD6TbdZMhWiM/0tgE/Aav7MxMkdtj6UdR9SzyjXdwLUmNrzrbZCv0dop8So+tyr4CvgaOBxLE4XY+0pM55/ANNVtCFbEiVE63UR1O4PDwAeOcl3lR0wDPwK/5OHKvwi8of4lqY+p7g8xBLwvjtbYJL8ZknMeU//S1TeiLVi2At+iLxmtff0LuNxE25elzHJKV7UiGreG5vyPqO7XkOb8B8CFFuq4IDbSgvBMNAdBL6srFJMaw2/a1Bh2i62kxj9aqdkbgvOXSR4KztK8oLjO5U4OgrYstPb1VgZ03krR2Kplsy0f17mvVCj622fpIch4TFtc8/1OGz+6ntLQZXHzpOGUDsL1TnH+JfR1uGWyvTnGKKvPpV33Cpey7vwxqnv3aJsldcK2YcPoGz69Inm7HO9jO1MJ3bojHZRCjyTUY8rX2FEaEwl5/yc6a5vMgmjW2oOJrAnuTxhnmevg+5i5hDGk/iwJLaFvOba9gwOwHX3Ls1JWRH7huNuNxJ4PYCjlvFK3Zam794Z3UvmLTme1sVpHHaeVOk76ruOI0ki9Bo4TDselTq66jvgUNqn0mWcIjxnl3mbSp6gVJT+GuMyzT6nrik9Ri468eI9wueeor9cpjmdZO7Yf8m60g6x9hnDW9x3jEeA28B2wh/DZI3W9LXUvYBiGYRiGYRiGYRiGYRiGYSTwHy11zABJLMguAAAAAElFTkSuQmCC"
)

func (c Service) PhotoAction() {

	c.AutoRender = false

	var photo string

	if len(c.Request.RequestPath) > iam_v1_service_len {

		uname := c.Request.RequestPath[iam_v1_service_len:]

		var profile iamapi.UserProfile

		if obj := store.Data.NewReader(iamapi.ObjKeyUserProfile(uname)).Query(); obj.OK() {
			if err := obj.Decode(&profile); err == nil && len(profile.Photo) > 50 {
				photo = profile.Photo
			}
		}
	}

	if photo == "" {
		photo = iam_v1_service_def
	}

	body64 := strings.SplitAfter(photo, ";base64,")
	if len(body64) != 2 {
		return
	}
	data, err := base64.StdEncoding.DecodeString(body64[1])
	if err != nil {
		return
	}

	c.Response.Out.Write(data)
}
