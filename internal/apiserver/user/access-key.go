// Copyright 2014 Eryx <evorui at gmail dot com>, All rights reserved.
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

package user

import (
	"fmt"
	"strings"

	"github.com/hooto/httpsrv"
	"github.com/sysinner/innerstack/v2/pkg/inauth"

	"github.com/hooto/iam/v2/internal/data"
	"github.com/hooto/iam/v2/pkg/iamapi"
)

const accessKeyLimit = 20

type AccessKeyEntryResponse struct {
	Status inauth.ServiceStatus `json:"status"`
	Item   *inauth.AccessKey    `json:"item,omitempty"`
}

// AccessKeyEntry returns a single access key by ID.
func AccessKeyEntry(ctx httpsrv.Ctx) error {

	u := authCtx(ctx)
	if u == nil {
		return nil
	}

	var rsp AccessKeyEntryResponse
	defer ctx.JSON(&rsp)

	id := ctx.Params().Value("access_key_id")
	if id == "" {
		rsp.Status = inauth.NewServiceStatus("404", "Access Key Not Found")
		return nil
	}

	var ak inauth.AccessKey
	if rs := data.Data.NewReader(iamapi.NsAccessKey(u.Name, id)).Exec(); rs.OK() {
		rs.Item().JsonDecode(&ak)
	}

	if ak.GetId() != id {
		rsp.Status = inauth.NewServiceStatus("404", "Access Key Not Found")
		return nil
	}

	ak.Secret = "" // do not return secret in frontend

	rsp.Status = inauth.NewServiceStatus("200", "ok")
	rsp.Item = &ak
	return nil
}

type AccessKeyListResponse struct {
	Status inauth.ServiceStatus `json:"status"`
	Items  []*inauth.AccessKey  `json:"items,omitempty"`
}

// AccessKeyList returns all access keys for the current user.
func AccessKeyList(ctx httpsrv.Ctx) error {

	u := authCtx(ctx)
	if u == nil {
		return nil
	}

	var rsp AccessKeyListResponse
	defer ctx.JSON(&rsp)

	k1 := iamapi.NsAccessKey(u.Name, "")
	k2 := iamapi.NsAccessKey(u.Name, "zzzzzzzz")
	if rs := data.Data.NewRanger(k1, k2).
		SetLimit(int64(accessKeyLimit)).Exec(); rs.OK() {
		for _, v := range rs.Items {
			var ak inauth.AccessKey
			if err := v.JsonDecode(&ak); err == nil &&
				ak.State != inauth.AccessKey_State_Disable {
				ak.Secret = "" // do not return secret in list
				rsp.Items = append(rsp.Items, &ak)
			}
		}
	}

	rsp.Status = inauth.NewServiceStatus("200", "ok")
	return nil
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

// AccessKeySet creates or updates an access key.
func AccessKeySet(ctx httpsrv.Ctx) error {

	u := authCtx(ctx)
	if u == nil {
		return nil
	}

	var (
		req AccessKeySetRequest
		rsp AccessKeySetResponse
	)
	defer ctx.JSON(&rsp)

	if err := ctx.Request().JsonDecode(&req); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", "Bad Request")
		return nil
	}

	var prev inauth.AccessKey
	if req.Id == "" {
		// create new access key using inauth factory
		newAk := inauth.NewUserAccessKey()
		newAk.User = u.Name
		newAk.State = inauth.AccessKey_State_Active
		newAk.Description = req.Description
		newAk.Scopes = req.Scopes
		prev = *newAk
	} else {
		// update existing access key
		if rs := data.Data.NewReader(
			iamapi.NsAccessKey(u.Name, req.Id)).Exec(); rs.OK() {
			rs.Item().JsonDecode(&prev)
		}
		if prev.GetId() != req.Id {
			rsp.Status = inauth.NewServiceStatus("404", "Access Key Not Found")
			return nil
		}
	}

	// enforce access key count limit on create
	if req.Id == "" {
		if rs := data.Data.NewRanger(
			iamapi.NsAccessKey(u.Name, ""), iamapi.NsAccessKey(u.Name, "")).
			SetLimit(int64(accessKeyLimit + 1)).Exec(); rs.OK() {
			if len(rs.Items) > accessKeyLimit {
				rsp.Status = inauth.NewServiceStatus("400",
					fmt.Sprintf("Num Out Range (%d)", accessKeyLimit))
				return nil
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

	if rs := data.Data.NewWriter(iamapi.NsAccessKey(u.Name, prev.Id), nil).SetJsonValue(prev).
		Exec(); rs.OK() {
		rsp.Status = inauth.NewServiceStatus("200", "ok")
		rsp.Item = &prev
		data.KeyMgr.Set(&prev)
	} else {
		rsp.Status = inauth.NewServiceStatus("500", "IO Error "+rs.ErrorMessage())
	}
	return nil
}

type AccessKeyDeleteResponse struct {
	Status inauth.ServiceStatus `json:"status"`
}

// AccessKeyDelete removes an access key by ID.
func AccessKeyDelete(ctx httpsrv.Ctx) error {

	u := authCtx(ctx)
	if u == nil {
		return nil
	}

	var rsp AccessKeyDeleteResponse
	defer ctx.JSON(&rsp)

	id := ctx.Params().Value("access_key_id")
	if id == "" {
		rsp.Status = inauth.NewServiceStatus("404", "Access Key Not Found")
		return nil
	}

	if rs := data.Data.NewDeleter(iamapi.NsAccessKey(u.Name, id)).Exec(); rs.OK() {
		rsp.Status = inauth.NewServiceStatus("200", "ok")
		data.KeyMgr.Del(id)
	} else {
		rsp.Status = inauth.NewServiceStatus("500", "IO Error")
	}
	return nil
}

type AccessKeyBindResponse struct {
	Status inauth.ServiceStatus `json:"status"`
}

// AccessKeyBind adds a scope binding to an access key.
func AccessKeyBind(ctx httpsrv.Ctx) error {

	u := authCtx(ctx)
	if u == nil {
		return nil
	}

	var rsp AccessKeyBindResponse
	defer ctx.JSON(&rsp)

	var (
		id    = ctx.Params().Value("access_key_id")
		bname = ctx.Params().Value("scope_content")
	)
	if id == "" || bname == "" {
		rsp.Status = inauth.NewServiceStatus("404", "Access Key Not Found")
		return nil
	}

	var ak inauth.AccessKey
	if rs := data.Data.NewReader(iamapi.NsAccessKey(u.Name, id)).Exec(); rs.OK() {
		rs.Item().JsonDecode(&ak)
	}

	if id != ak.GetId() {
		rsp.Status = inauth.NewServiceStatus("404", "Access Key Not Found")
		return nil
	}

	ar := strings.Split(bname, "=")
	if len(ar) != 2 {
		rsp.Status = inauth.NewServiceStatus("400", "Invalid Bound Value")
		return nil
	}

	scope := strings.TrimSpace(ar[0]) + "=" + strings.TrimSpace(ar[1])
	scopeSet(&ak.Scopes, scope)

	if rs := data.Data.NewWriter(iamapi.NsAccessKey(u.Name, ak.Id), nil).SetJsonValue(ak).
		Exec(); rs.OK() {
		rsp.Status = inauth.NewServiceStatus("200", "ok")
		data.KeyMgr.Set(&ak)
	} else {
		rsp.Status = inauth.NewServiceStatus("500", "IO Error")
	}
	return nil
}

// AccessKeyUnbind removes a scope binding from an access key.
func AccessKeyUnbind(ctx httpsrv.Ctx) error {

	u := authCtx(ctx)
	if u == nil {
		return nil
	}

	var rsp AccessKeyBindResponse
	defer ctx.JSON(&rsp)

	var (
		id    = ctx.Params().Value("access_key_id")
		bname = ctx.Params().Value("scope_content")
	)
	if id == "" || bname == "" {
		rsp.Status = inauth.NewServiceStatus("404", "Access Key Not Found")
		return nil
	}

	var ak inauth.AccessKey
	if rs := data.Data.NewReader(iamapi.NsAccessKey(u.Name, id)).Exec(); rs.OK() {
		rs.Item().JsonDecode(&ak)
	}

	if id != ak.GetId() {
		rsp.Status = inauth.NewServiceStatus("404", "Access Key Not Found")
		return nil
	}

	ar := strings.Split(bname, "=")
	if len(ar) > 2 {
		rsp.Status = inauth.NewServiceStatus("400", "Invalid Bound Value")
		return nil
	}

	name := strings.TrimSpace(ar[0])
	if name == "" {
		rsp.Status = inauth.NewServiceStatus("400", "Invalid Bound Value")
		return nil
	}
	scopeDel(&ak.Scopes, name)

	if rs := data.Data.NewWriter(iamapi.NsAccessKey(u.Name, ak.Id), nil).SetJsonValue(ak).
		Exec(); rs.OK() {
		rsp.Status = inauth.NewServiceStatus("200", "ok")
		data.KeyMgr.Set(&ak)
	} else {
		rsp.Status = inauth.NewServiceStatus("500", "IO Error")
	}
	return nil
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
