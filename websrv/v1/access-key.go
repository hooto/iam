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
	"fmt"
	"strings"
	"time"

	"github.com/hooto/hauth/go/hauth/v1"
	"github.com/hooto/httpsrv"
	"github.com/lessos/lessgo/crypto/idhash"
	"github.com/lessos/lessgo/types"
	iox_utils "github.com/lynkdb/iomix/utils"

	"github.com/hooto/iam/data"
	"github.com/hooto/iam/iamapi"
	"github.com/hooto/iam/iamclient"
)

var (
	ak_limit = 20
)

type AccessKey struct {
	*httpsrv.Controller
	us iamapi.UserSession
}

func (c *AccessKey) Init() int {

	//
	c.us, _ = iamclient.SessionInstance(c.Session)

	if !c.us.IsLogin() {
		c.Response.Out.WriteHeader(401)
		c.RenderJson(types.NewTypeErrorMeta(iamapi.ErrCodeUnauthorized, "Unauthorized"))
		return 1
	}

	return 0
}

func (c AccessKey) EntryAction() {

	var set types.WebServiceResult
	defer c.RenderJson(&set)

	id := c.Params.Value("access_key_id")
	if id == "" {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeNotFound, "Access Key Not Found")
		return
	}

	var ak hauth.AccessKey
	if rs := data.Data.NewReader(iamapi.NsAccessKey(c.us.UserName, id)).Query(); rs.OK() {
		rs.Decode(&ak)
	}

	if ak.Id != "" && ak.Id == id {
		set.Kind = "AccessKey"
		set.Item = ak
	} else {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeNotFound, "Access Key Not Found")
	}
}

func (c AccessKey) ListAction() {

	var ls types.WebServiceResult
	defer c.RenderJson(&ls)

	k1 := iamapi.NsAccessKey(c.us.UserName, "zzzzzzzz")
	k2 := iamapi.NsAccessKey(c.us.UserName, "")
	if rs := data.Data.NewReader(nil).KeyRangeSet(k1, k2).
		ModeRevRangeSet(true).LimitNumSet(int64(ak_limit)).Query(); rs.OK() {

		for _, v := range rs.Items {
			var ak hauth.AccessKey
			if err := v.Decode(&ak); err == nil {
				ls.Items = append(ls.Items, ak)
			}
		}
	}

	ls.Kind = "AccessKeyList"
}

func (c AccessKey) SetAction() {

	var set struct {
		types.TypeMeta
		hauth.AccessKey
	}
	defer c.RenderJson(&set)

	if err := c.Request.JsonDecode(&set.AccessKey); err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Bad Request")
		return
	}

	var prev hauth.AccessKey
	if len(set.AccessKey.Id) < 16 {
		set.AccessKey.Id = iox_utils.Uint32ToHexString(uint32(time.Now().Unix())) + idhash.RandHexString(8)
	} else {

		if rs := data.Data.NewReader(
			iamapi.NsAccessKey(c.us.UserName, set.AccessKey.Id)).Query(); rs.OK() {
			rs.Decode(&prev)
		}
	}

	if rs := data.Data.NewReader(nil).KeyRangeSet(
		iamapi.NsAccessKey(c.us.UserName, ""), iamapi.NsAccessKey(c.us.UserName, "")).
		LimitNumSet(int64(ak_limit + 1)).Query(); rs.OK() {
		if len(rs.Items) > ak_limit && prev.Id == "" {
			set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, fmt.Sprintf("Num Out Range (%d)", ak_limit))
			return
		}
	}

	if prev.Id == "" {
		prev = set.AccessKey
	} else {

		prev.Status = set.AccessKey.Status
		prev.Description = set.AccessKey.Description

		for _, v := range set.AccessKey.Scopes {
			prev.ScopeSet(v)
		}
	}

	if len(prev.Secret) < 40 {
		prev.Secret = idhash.RandBase64String(40)
	}

	if len(prev.User) < 1 {
		prev.User = c.us.UserName
	}

	if rs := data.Data.NewWriter(iamapi.NsAccessKey(c.us.UserName, prev.Id), prev).
		Commit(); rs.OK() {
		set.Kind = "AccessKey"
	} else {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInternalError, "IO Error")
	}
}

func (c AccessKey) DelAction() {

	var set types.TypeMeta
	defer c.RenderJson(&set)

	id := c.Params.Value("access_key_id")
	if id == "" {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeNotFound, "Access Key Not Found")
		return
	}

	if rs := data.Data.NewWriter(iamapi.NsAccessKey(c.us.UserName, id), nil).
		ModeDeleteSet(true).Commit(); rs.OK() {
		set.Kind = "AccessKey"
	} else {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInternalError, "IO Error")
	}
}

func (c AccessKey) BindAction() {

	var set types.TypeMeta
	defer c.RenderJson(&set)

	var (
		id    = c.Params.Value("access_key_id")
		bname = c.Params.Value("scope_content")
	)
	if id == "" && bname == "" {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeNotFound, "Access Key Not Found")
		return
	}

	var ak hauth.AccessKey
	if rs := data.Data.NewReader(iamapi.NsAccessKey(c.us.UserName, id)).Query(); rs.OK() {
		rs.Decode(&ak)
	}

	if id != ak.Id {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeNotFound, "Access Key Not Found")
		return
	}

	ar := strings.Split(bname, "=")
	if len(ar) != 2 {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeNotFound, "Invalid Bound Value")
		return
	}

	ak.ScopeSet(hauth.NewScopeFilter(
		strings.TrimSpace(ar[0]),
		strings.TrimSpace(ar[1])))

	if rs := data.Data.NewWriter(iamapi.NsAccessKey(c.us.UserName, ak.Id), ak).
		Commit(); rs.OK() {
		set.Kind = "AccessKey"
	} else {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInternalError, "IO Error")
	}
}

func (c AccessKey) UnbindAction() {

	var set types.TypeMeta
	defer c.RenderJson(&set)

	var (
		id    = c.Params.Value("access_key_id")
		bname = c.Params.Value("scope_content")
	)
	if id == "" && bname == "" {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeNotFound, "Access Key Not Found")
		return
	}

	var ak hauth.AccessKey
	if rs := data.Data.NewReader(iamapi.NsAccessKey(c.us.UserName, id)).Query(); rs.OK() {
		rs.Decode(&ak)
	}

	if id != ak.Id {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeNotFound, "Access Key Not Found")
		return
	}

	ar := strings.Split(bname, "=")
	if len(ar) > 2 {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeNotFound, "Invalid Bound Value")
		return
	}

	bname = strings.TrimSpace(ar[0])
	if bname == "" {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeNotFound, "Invalid Bound Value")
		return
	}
	ak.ScopeDel(bname)

	if rs := data.Data.NewWriter(iamapi.NsAccessKey(c.us.UserName, ak.Id), ak).
		Commit(); rs.OK() {
		set.Kind = "AccessKey"
	} else {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInternalError, "IO Error")
	}
}
