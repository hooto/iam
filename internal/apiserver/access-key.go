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

package apiserver

import (
	"fmt"
	"strings"

	"github.com/hooto/httpsrv"
	"github.com/sysinner/incore/v2/pkg/inauth"

	"github.com/hooto/iam/v2/internal/data"
	"github.com/hooto/iam/v2/pkg/iamapi"
)

const (
	accessKeyLimit = 20
)

type AccessKey struct {
	*httpsrv.Controller
}

type AccessKeyEntryResponse struct {
	Status inauth.ServiceStatus `json:"status"`
	Item   *inauth.AccessKey    `json:"item,omitempty"`
}

func (c AccessKey) EntryAction() {

	user := userAuth(c.Controller)
	if user == nil {
		return
	}

	var rsp AccessKeyEntryResponse
	defer c.RenderJson(&rsp)

	id := c.Params.Value("access_key_id")
	if id == "" {
		rsp.Status = inauth.NewServiceStatus("404", "Access Key Not Found")
		return
	}

	var ak inauth.AccessKey
	if rs := data.Data.NewReader(iamapi.NsAccessKey(user.Name, id)).Exec(); rs.OK() {
		rs.Item().JsonDecode(&ak)
	}

	if ak.GetId() != id {
		rsp.Status = inauth.NewServiceStatus("404", "Access Key Not Found")
		return
	}

	rsp.Status = inauth.NewServiceStatus("200", "ok")
	rsp.Item = &ak
}

type AccessKeyListRequest struct {
	AccessToken string `json:"access_token,omitempty"`
}

type AccessKeyListResponse struct {
	Status inauth.ServiceStatus `json:"status"`
	Items  []*inauth.AccessKey  `json:"items,omitempty"`
}

func (c AccessKey) ListAction() {

	user := userAuth(c.Controller)
	if user == nil {
		return
	}

	var rsp AccessKeyListResponse
	defer c.RenderJson(&rsp)

	k1 := iamapi.NsAccessKey(user.Name, "")
	k2 := iamapi.NsAccessKey(user.Name, "zzzzzzzz")
	if rs := data.Data.NewRanger(k1, k2).
		SetLimit(int64(accessKeyLimit)).Exec(); rs.OK() {
		for _, v := range rs.Items {
			var ak inauth.AccessKey
			if err := v.JsonDecode(&ak); err == nil &&
				ak.State != inauth.AccessKey_State_Disable {
				rsp.Items = append(rsp.Items, &ak)
			}
		}
	}

	rsp.Status = inauth.NewServiceStatus("200", "ok")
}

type AccessKeySetRequest struct {
	Id          string   `json:"id,omitempty"`
	State       string   `json:"state,omitempty"`
	Description string   `json:"description,omitempty"`
	Scopes      []string `json:"scopes,omitempty"`
}

type AccessKeySetResponse struct {
	Status inauth.ServiceStatus `json:"status"`
	Item   *inauth.AccessKey    `json:"item,omitempty"`
}

func (c AccessKey) SetAction() {

	user := userAuth(c.Controller)
	if user == nil {
		return
	}

	var (
		req AccessKeySetRequest
		rsp AccessKeySetResponse
	)
	defer c.RenderJson(&rsp)

	if err := c.Request.JsonDecode(&req); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", "Bad Request")
		return
	}

	var prev inauth.AccessKey
	if req.Id == "" {
		// create new access key using inauth factory
		newAk := inauth.NewUserAccessKey()
		newAk.User = user.Name
		newAk.State = inauth.AccessKey_State_Active
		newAk.Description = req.Description
		newAk.Scopes = req.Scopes
		prev = *newAk
	} else {
		// update existing access key
		if rs := data.Data.NewReader(
			iamapi.NsAccessKey(user.Name, req.Id)).Exec(); rs.OK() {
			rs.Item().JsonDecode(&prev)
		}
		if prev.GetId() != req.Id {
			rsp.Status = inauth.NewServiceStatus("404", "Access Key Not Found")
			return
		}
	}

	// enforce access key count limit on create
	if req.Id == "" {
		if rs := data.Data.NewRanger(
			iamapi.NsAccessKey(user.Name, ""), iamapi.NsAccessKey(user.Name, "")).
			SetLimit(int64(accessKeyLimit + 1)).Exec(); rs.OK() {
			if len(rs.Items) > accessKeyLimit {
				rsp.Status = inauth.NewServiceStatus("400",
					fmt.Sprintf("Num Out Range (%d)", accessKeyLimit))
				return
			}
		}
	}

	if prev.GetId() == req.Id {
		// update existing
		if req.State != "" {
			prev.State = req.State
		}
		if req.Description != "" {
			prev.Description = req.Description
		}
		for _, s := range req.Scopes {
			scopeSet(&prev.Scopes, s)
		}
	}

	if rs := data.Data.NewWriter(iamapi.NsAccessKey(user.Name, prev.Id), nil).SetJsonValue(prev).
		Exec(); rs.OK() {
		rsp.Status = inauth.NewServiceStatus("200", "ok")
		rsp.Item = &prev
		data.KeyMgr.Set(&prev)
	} else {
		rsp.Status = inauth.NewServiceStatus("500", "IO Error "+rs.ErrorMessage())
	}
}

type AccessKeyDeleteResponse struct {
	Status inauth.ServiceStatus `json:"status"`
}

func (c AccessKey) DeleteAction() {

	user := userAuth(c.Controller)
	if user == nil {
		return
	}

	var rsp AccessKeyDeleteResponse
	defer c.RenderJson(&rsp)

	id := c.Params.Value("access_key_id")
	if id == "" {
		rsp.Status = inauth.NewServiceStatus("404", "Access Key Not Found")
		return
	}

	if rs := data.Data.NewDeleter(iamapi.NsAccessKey(user.Name, id)).Exec(); rs.OK() {
		rsp.Status = inauth.NewServiceStatus("200", "ok")
		data.KeyMgr.Del(id)
	} else {
		rsp.Status = inauth.NewServiceStatus("500", "IO Error")
	}
}

type AccessKeyBindResponse struct {
	Status inauth.ServiceStatus `json:"status"`
}

func (c AccessKey) BindAction() {

	user := userAuth(c.Controller)
	if user == nil {
		return
	}

	var rsp AccessKeyBindResponse
	defer c.RenderJson(&rsp)

	var (
		id    = c.Params.Value("access_key_id")
		bname = c.Params.Value("scope_content")
	)
	if id == "" || bname == "" {
		rsp.Status = inauth.NewServiceStatus("404", "Access Key Not Found")
		return
	}

	var ak inauth.AccessKey
	if rs := data.Data.NewReader(iamapi.NsAccessKey(user.Name, id)).Exec(); rs.OK() {
		rs.Item().JsonDecode(&ak)
	}

	if id != ak.GetId() {
		rsp.Status = inauth.NewServiceStatus("404", "Access Key Not Found")
		return
	}

	ar := strings.Split(bname, "=")
	if len(ar) != 2 {
		rsp.Status = inauth.NewServiceStatus("400", "Invalid Bound Value")
		return
	}

	scope := strings.TrimSpace(ar[0]) + "=" + strings.TrimSpace(ar[1])
	scopeSet(&ak.Scopes, scope)

	if rs := data.Data.NewWriter(iamapi.NsAccessKey(user.Name, ak.Id), nil).SetJsonValue(ak).
		Exec(); rs.OK() {
		rsp.Status = inauth.NewServiceStatus("200", "ok")
		data.KeyMgr.Set(&ak)
	} else {
		rsp.Status = inauth.NewServiceStatus("500", "IO Error")
	}
}

func (c AccessKey) UnbindAction() {

	user := userAuth(c.Controller)
	if user == nil {
		return
	}

	var rsp AccessKeyBindResponse
	defer c.RenderJson(&rsp)

	var (
		id    = c.Params.Value("access_key_id")
		bname = c.Params.Value("scope_content")
	)
	if id == "" || bname == "" {
		rsp.Status = inauth.NewServiceStatus("404", "Access Key Not Found")
		return
	}

	var ak inauth.AccessKey
	if rs := data.Data.NewReader(iamapi.NsAccessKey(user.Name, id)).Exec(); rs.OK() {
		rs.Item().JsonDecode(&ak)
	}

	if id != ak.GetId() {
		rsp.Status = inauth.NewServiceStatus("404", "Access Key Not Found")
		return
	}

	ar := strings.Split(bname, "=")
	if len(ar) > 2 {
		rsp.Status = inauth.NewServiceStatus("400", "Invalid Bound Value")
		return
	}

	name := strings.TrimSpace(ar[0])
	if name == "" {
		rsp.Status = inauth.NewServiceStatus("400", "Invalid Bound Value")
		return
	}
	scopeDel(&ak.Scopes, name)

	if rs := data.Data.NewWriter(iamapi.NsAccessKey(user.Name, ak.Id), nil).SetJsonValue(ak).
		Exec(); rs.OK() {
		rsp.Status = inauth.NewServiceStatus("200", "ok")
		data.KeyMgr.Set(&ak)
	} else {
		rsp.Status = inauth.NewServiceStatus("500", "IO Error")
	}
}

// scopeSet adds a scope to the list if not already present.
func scopeSet(scopes *[]string, s string) {
	if scopes == nil || s == "" {
		return
	}
	for _, v := range *scopes {
		if v == s {
			return
		}
	}
	*scopes = append(*scopes, s)
}

// scopeDel removes scopes matching the given prefix (before '=' if present).
func scopeDel(scopes *[]string, prefix string) {
	if scopes == nil || prefix == "" {
		return
	}
	n := 0
	for _, v := range *scopes {
		if strings.HasPrefix(v, prefix+"=") || v == prefix {
			continue
		}
		(*scopes)[n] = v
		n++
	}
	*scopes = (*scopes)[:n]
}
