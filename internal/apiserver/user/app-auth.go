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
	"time"

	"github.com/google/uuid"
	"github.com/hooto/httpsrv"
	"github.com/sysinner/innerstack/v2/pkg/inauth"

	"github.com/hooto/iam/v2/internal/data"
	"github.com/hooto/iam/v2/pkg/iamapi"
)

func accessTokenFromRequest(ctx httpsrv.Ctx, bodyToken string) string {
	if bodyToken != "" {
		return bodyToken
	}
	cookie, err := ctx.Request().Cookie(inauth.AppHttpHeaderKey)
	if err != nil || cookie.Value == "" {
		return ""
	}
	return cookie.Value
}

func verifyAccessToken(accessToken string) (*inauth.AccessToken, error) {
	token, err := inauth.ParseAccessToken(accessToken)
	if err != nil {
		return nil, err
	}
	if _, err = token.Verify(data.KeyMgr); err != nil {
		return nil, err
	}
	return token, nil
}

// AppAuth_Register creates a new AppInstance with auto-generated ID and SecretKey.
func AppAuth_Register(ctx httpsrv.Ctx) error {
	u := authCtx(ctx)
	if u == nil {
		return nil
	}
	var (
		req struct {
			AccessToken string `json:"access_token,omitempty"`
			Name        string `json:"name"`
			Url         string `json:"url"`
		}
		rsp struct {
			Status inauth.ServiceStatus `json:"status"`
			App    *iamapi.AppInstance  `json:"app,omitempty"`
		}
	)
	defer ctx.JSON(&rsp)

	if err := ctx.Request().JsonDecode(&req); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", "Bad Argument")
		return nil
	}

	accessToken := accessTokenFromRequest(ctx, req.AccessToken)
	if accessToken == "" {
		rsp.Status = inauth.NewServiceStatus("401", "Unauthorized")
		return nil
	}

	_, err := verifyAccessToken(accessToken)
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("401", err.Error())
		return nil
	}

	tn := time.Now().Unix()

	secretKey, err := inauth.GenerateSecretKeyBase62(40)
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("500", "Failed to generate secret key")
		return nil
	}

	app := iamapi.AppInstance{
		ID:        uuid.New().String(),
		Name:      req.Name,
		User:      u.Name,
		Url:       req.Url,
		SecretKey: secretKey,
		Status:    1,
		Created:   tn,
		Updated:   tn,
	}

	if rs := data.Data.NewWriter(iamapi.NsAppInstance(app.ID), nil).SetJsonValue(app).Exec(); !rs.OK() {
		rsp.Status = inauth.NewServiceStatus("500", rs.ErrorMessage())
		return nil
	}

	rsp.App = &app
	rsp.Status = inauth.NewServiceStatus("200", "ok")
	return nil
}

// AppAuth_List returns all apps belonging to the authenticated user.
func AppAuth_List(ctx httpsrv.Ctx) error {
	u := authCtx(ctx)
	if u == nil {
		return nil
	}
	var (
		req struct {
			AccessToken string `json:"access_token,omitempty"`
		}
		rsp struct {
			Status inauth.ServiceStatus `json:"status"`
			Items  []iamapi.AppInstance `json:"items,omitempty"`
		}
	)
	defer ctx.JSON(&rsp)

	ctx.Request().JsonDecode(&req)

	accessToken := accessTokenFromRequest(ctx, req.AccessToken)
	if accessToken == "" {
		rsp.Status = inauth.NewServiceStatus("401", "Unauthorized")
		return nil
	}

	_, err := verifyAccessToken(accessToken)
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("401", err.Error())
		return nil
	}

	if rs := data.Data.NewRanger(
		iamapi.NsAppInstance(""),
		iamapi.NsAppInstance("~"),
	).Exec(); rs.OK() {
		for _, item := range rs.Items {
			var app iamapi.AppInstance
			item.JsonDecode(&app)
			if app.User == u.Name {
				app.SecretKey = maskSecret(app.SecretKey)
				rsp.Items = append(rsp.Items, app)
			}
		}
	}

	rsp.Status = inauth.NewServiceStatus("200", "ok")
	return nil
}

// AppAuth_Update updates an existing AppInstance's name and callback URL.
func AppAuth_Update(ctx httpsrv.Ctx) error {
	u := authCtx(ctx)
	if u == nil {
		return nil
	}
	var (
		req struct {
			AccessToken string `json:"access_token,omitempty"`
			ID          string `json:"id"`
			Name        string `json:"name"`
			Url         string `json:"url"`
		}
		rsp struct {
			Status inauth.ServiceStatus `json:"status"`
			App    *iamapi.AppInstance  `json:"app,omitempty"`
		}
	)
	defer ctx.JSON(&rsp)

	if err := ctx.Request().JsonDecode(&req); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", "Bad Argument")
		return nil
	}

	accessToken := accessTokenFromRequest(ctx, req.AccessToken)
	if accessToken == "" {
		rsp.Status = inauth.NewServiceStatus("401", "Unauthorized")
		return nil
	}

	_, err := verifyAccessToken(accessToken)
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("401", err.Error())
		return nil
	}

	var app iamapi.AppInstance
	if rs := data.Data.NewReader(iamapi.NsAppInstance(req.ID)).Exec(); rs.OK() {
		rs.Item().JsonDecode(&app)
	}

	if app.ID != req.ID {
		rsp.Status = inauth.NewServiceStatus("404", "App not found")
		return nil
	}

	if app.User != u.Name {
		rsp.Status = inauth.NewServiceStatus("401", "Unauthorized")
		return nil
	}

	app.Name = req.Name
	app.Url = req.Url
	app.Updated = time.Now().Unix()

	if rs := data.Data.NewWriter(iamapi.NsAppInstance(app.ID), nil).SetJsonValue(app).Exec(); !rs.OK() {
		rsp.Status = inauth.NewServiceStatus("500", rs.ErrorMessage())
		return nil
	}

	app.SecretKey = ""
	rsp.App = &app
	rsp.Status = inauth.NewServiceStatus("200", "ok")
	return nil
}

// AppAuth_Delete removes an AppInstance.
func AppAuth_Delete(ctx httpsrv.Ctx) error {
	var (
		req struct {
			AccessToken string `json:"access_token,omitempty"`
			ID          string `json:"id"`
		}
		rsp struct {
			Status inauth.ServiceStatus `json:"status"`
		}
	)
	defer ctx.JSON(&rsp)

	if err := ctx.Request().JsonDecode(&req); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", "Bad Argument")
		return nil
	}

	accessToken := accessTokenFromRequest(ctx, req.AccessToken)
	if accessToken == "" {
		rsp.Status = inauth.NewServiceStatus("401", "Unauthorized")
		return nil
	}

	token, err := verifyAccessToken(accessToken)
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("401", err.Error())
		return nil
	}

	var app iamapi.AppInstance
	if rs := data.Data.NewReader(iamapi.NsAppInstance(req.ID)).Exec(); rs.OK() {
		rs.Item().JsonDecode(&app)
	}

	if app.ID != req.ID {
		rsp.Status = inauth.NewServiceStatus("404", "App not found")
		return nil
	}

	if app.User != token.Claims.Sub {
		rsp.Status = inauth.NewServiceStatus("401", "Unauthorized")
		return nil
	}

	data.Data.NewDeleter(iamapi.NsAppInstance(req.ID)).Exec()
	rsp.Status = inauth.NewServiceStatus("200", "ok")
	return nil
}

func maskSecret(s string) string {
	if len(s) > 8 {
		return s[:4] + "****" + s[len(s)-4:]
	}
	return "****"
}
