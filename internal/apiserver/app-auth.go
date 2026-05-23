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
	"time"

	"github.com/hooto/httpsrv"
	"github.com/sysinner/incore/v2/pkg/inauth"

	"github.com/hooto/iam/v2/internal/data"
	"github.com/hooto/iam/v2/pkg/iamapi"
)

type AppAuth struct {
	*httpsrv.Controller
}

type AppAuthRegisterRequest struct {
	AccessToken string             `json:"access_token,omitempty" toml:"access_token,omitempty"`
	App         iamapi.AppInstance `json:"app" toml:"app"`
}

type AppAuthRegisterResponse struct {
	Status inauth.ServiceStatus `json:"status"`
	App    *iamapi.AppInstance  `json:"app,omitempty"`
}

func (c AppAuth) RegisterAction() {

	var (
		req AppAuthRegisterRequest
		rsp AppAuthRegisterResponse
	)
	defer c.RenderJson(&rsp)

	if err := c.Request.JsonDecode(&req); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", "Bad Argument")
		return
	}

	tn := time.Now().Unix()

	if err := iamapi.AppIdValid(req.App.ID); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", err.Error())
		return
	}

	// parse and verify the access token
	token, err := inauth.ParseAccessToken(req.AccessToken)
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("401", err.Error())
		return
	}

	_, err = token.Verify(data.KeyMgr)
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("401", "Unauthorized")
		return
	}

	var prev iamapi.AppInstance
	if rs := data.Data.NewReader(iamapi.NsAppInstance(req.App.ID)).Exec(); rs.OK() {
		rs.Item().JsonDecode(&prev)
	}

	if prev.ID == "" {

		req.App.Created = tn
		req.App.Updated = tn
		req.App.Status = 1
		req.App.User = token.Claims.Sub

	} else {

		if prev.User != token.Claims.Sub {
			rsp.Status = inauth.NewServiceStatus("401", "Unauthorized")
			return
		}

		req.App.Created = prev.Created
		req.App.User = prev.User
		req.App.Status = prev.Status
	}

	if rs := data.Data.NewWriter(iamapi.NsAppInstance(req.App.ID), nil).SetJsonValue(req.App).
		Exec(); !rs.OK() {
		rsp.Status = inauth.NewServiceStatus("500", rs.ErrorMessage())
		return
	}

	rsp.Status = inauth.NewServiceStatus("200", "ok")
	rsp.App = &req.App
}

type AppAuthInfoRequest struct {
	AccessToken string `json:"access_token,omitempty" toml:"access_token,omitempty"`
	AppID       string `json:"app_id,omitempty"`
}

type AppAuthInfoResponse struct {
	Status inauth.ServiceStatus `json:"status"`
	App    *iamapi.AppInstance  `json:"app,omitempty"`
}

func (c AppAuth) InfoAction() {

	var (
		req AppAuthInfoRequest
		rsp AppAuthInfoResponse
	)
	defer c.RenderJson(&rsp)

	c.Request.JsonDecode(&req)

	if req.AppID == "" {
		rsp.Status = inauth.NewServiceStatus("404", "App Instance Not Found")
		return
	}

	// parse and verify the access token
	token, err := inauth.ParseAccessToken(req.AccessToken)
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("401", err.Error())
		return
	}

	_, err = token.Verify(data.KeyMgr)
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("401", "Unauthorized")
		return
	}
	user := data.UserGet(token.Claims.Sub)
	if user == nil {
		rsp.Status = inauth.NewServiceStatus("401", "User not found")
		return
	}

	var inst iamapi.AppInstance
	if rs := data.Data.NewReader(iamapi.NsAppInstance(req.AppID)).Exec(); rs.OK() {
		rs.Item().JsonDecode(&inst)
	}

	if inst.ID == req.AppID {
		if inst.User != token.Claims.Sub &&
			inst.User != iamapi.UserSysadmin {
			rsp.Status = inauth.NewServiceStatus("401", "Unauthorized")
		} else {
			rsp.Status = inauth.NewServiceStatus("200", "ok")
			rsp.App = &inst
		}
	} else {
		rsp.Status = inauth.NewServiceStatus("404", "App Instance Not Found")
	}
}

/**
type AppAuthRoleListRequest struct {
	AccessToken string `json:"access_token,omitempty" toml:"access_token,omitempty"`
}

type AppAuthRoleListResponse struct {
	Status inauth.ServiceStatus `json:"status"`
	Items  []iamapi.UserRole    `json:"items,omitempty"`
}

func (c AppAuth) RoleListAction() {

	var (
		req AppAuthRoleListRequest
		rsp AppAuthRoleListResponse
	)
	defer c.RenderJson(&rsp)

	c.Request.JsonDecode(&req)

	// parse and verify the access token
	token, err := inauth.ParseAccessToken(req.AccessToken)
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("401", err.Error())
		return
	}

	_, err = token.Verify(data.KeyMgr)
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("401", "Unauthorized")
		return
	}

	if rs := data.Data.NewRanger(
		iamapi.NsRole(""), iamapi.NsRole("")).SetLimit(100).Exec(); rs.OK() {

		for _, v := range rs.Items {

			var role iamapi.UserRole
			if err := v.JsonDecode(&role); err == nil {

				if role.Status == 0 {
					continue
				}

				rsp.Items = append(rsp.Items, iamapi.UserRole{
					Name: role.Name,
					Desc: role.Desc,
				})
			}
		}
	}

	rsp.Status = inauth.NewServiceStatus("200", "ok")
}

type AppAuthUserAccessKeyRequest struct {
	AccessToken string `json:"access_token,omitempty" toml:"access_token,omitempty"`
}

type AppAuthUserAccessKeyResponse struct {
	Status        inauth.ServiceStatus `json:"status"`
	IdentityToken inauth.IdentityToken `json:"identity_token"`
}

func (c AppAuth) UserAccessKeyAction() {

	var (
		req AppAuthUserAccessKeyRequest
		rsp AppAuthUserAccessKeyResponse
	)
	defer c.RenderJson(&rsp)

	// parse and verify the access token
	token, err := inauth.ParseAccessToken(req.AccessToken)
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("401", err.Error())
		return
	}

	_, err = token.Verify(data.KeyMgr)
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("401", "Unauthorized")
		return
	}

	var (
		username  = c.Params.Value("user")
		accessKey = c.Params.Value("access_key")
	)
	if username == "" || accessKey == "" {
		rsp.Status = inauth.NewServiceStatus("400", "Bad Argument")
		return
	}

	appAuth := c.Request.Header.Get("Auth")
	if appAuth == "" {
		rsp.Status = inauth.NewServiceStatus("401", "Unauthorized")
		return
	}

	appAka, err := iamapi.AccessKeyAuthDecode(appAuth)
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("401", "Unauthorized")
		return
	}

	if err := appAka.Valid(); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", err.Error())
		return
	}

	var app iamapi.AppInstance
	if rs := data.Data.NewReader(iamapi.NsAppInstance(appAka.Key)).Exec(); rs.OK() {
		rs.Item().JsonDecode(&app)
	}

	if app.ID != appAka.Key {
		rsp.Status = inauth.NewServiceStatus("400", "Bad Argument")
		return
	}

	if err := accessKeyAuthValid(appAka, app.SecretKey); err != nil {
		rsp.Status = inauth.NewServiceStatus("401", "Unauthorized")
		return
	}

	var userAk inauth.AccessKey
	if rs := data.Data.NewReader(iamapi.NsAccessKey(username, accessKey)).Exec(); rs.OK() {
		rs.Item().JsonDecode(&userAk)
	}

	if userAk.GetId() != accessKey ||
		userAk.GetState() != inauth.AccessKey_State_Active {
		rsp.Status = inauth.NewServiceStatus("400", "Access Key Not Found")
		return
	}

	// check if the access key allows access to the app
	if !userAk.Allow("app:" + appAka.Key) {
		rsp.Status = inauth.NewServiceStatus("400", "Access Key scope not allowed")
		return
	}

	var user iamapi.User
	if rs := data.Data.NewReader(iamapi.NsUser(username)).Exec(); rs.OK() {
		rs.Item().JsonDecode(&user)
	}

	rsp.Status = inauth.NewServiceStatus("200", "ok")
	rsp.AccessKeySession = iamapi.AccessKeySession{
		User:      username,
		AccessKey: userAk.GetId(),
		SecretKey: userAk.GetSecret(),
		Roles:     user.Roles,
		Expired:   time.Now().Unix() + 864000,
	}

	slog.Info("app-auth UserAccessKey ok", "user", username, "access_key", accessKey)
}
*/
