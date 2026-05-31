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

package open

import (
	"encoding/json"
	"log/slog"

	"github.com/hooto/httpsrv"
	"github.com/sysinner/incore/v2/pkg/inauth"

	"github.com/hooto/iam/v2/internal/data"
	"github.com/hooto/iam/v2/pkg/iamapi"
)

// appAuth verifies the access token from cookie and returns the
// parsed AccessKey. It writes a 401 JSON response and returns nil on failure.
func appAuth(ctx *httpsrv.Context) *inauth.AccessKey {

	tokenStr := ""
	cookie, err := ctx.Request().Cookie(inauth.AppHttpHeaderKey)
	if err == nil && cookie.Value != "" {
		tokenStr = cookie.Value
	}

	if tokenStr == "" {
		ctx.RenderJson(struct {
			Status inauth.ServiceStatus `json:"status"`
		}{Status: inauth.NewServiceStatus("401", "Unauthorized")})
		return nil
	}

	token, err := inauth.ParseAccessToken(tokenStr)
	if err != nil || token.Header.Kid == "" {
		ctx.RenderJson(struct {
			Status inauth.ServiceStatus `json:"status"`
		}{Status: inauth.NewServiceStatus("401", "Unauthorized")})
		return nil
	}

	ak, err := token.Verify(data.KeyMgr)
	if err != nil {
		ctx.RenderJson(struct {
			Status inauth.ServiceStatus `json:"status"`
		}{Status: inauth.NewServiceStatus("401", "Unauthorized")})
		return nil
	}

	return ak
}

// AppAuth_Verify verifies app credentials (app_id + secret_key).
// Used by third-party apps during setup to validate their IAM configuration.
func AppAuth_Verify(ctx *httpsrv.Context) error {

	ak := appAuth(ctx)
	if ak == nil {
		return nil
	}

	var (
		req struct {
			AppId string `json:"app_id"`
		}
		rsp struct {
			Status inauth.ServiceStatus `json:"status"`
			App    *iamapi.AppInstance  `json:"app,omitempty"`
		}
	)
	defer ctx.RenderJson(&rsp)

	if err := ctx.Request().JsonDecode(&req); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", "Bad Argument")
		return nil
	}

	if req.AppId == "" || req.AppId != ak.Id {
		rsp.Status = inauth.NewServiceStatus("400", "app_id required")
		return nil
	}

	var app iamapi.AppInstance
	if rs := data.Data.NewReader(iamapi.NsAppInstance(req.AppId)).Exec(); rs.OK() {
		rs.Item().JsonDecode(&app)
	}

	if app.ID != req.AppId {
		rsp.Status = inauth.NewServiceStatus("404", "App not found")
		return nil
	}

	app.SecretKey = ""

	rsp.App = &app
	rsp.Status = inauth.NewServiceStatus("200", "ok")

	return nil
}

type AppAuth_TokenExchangeRequest struct {
	Code string `json:"code"`
}

type AppAuth_TokenExchangeResponse struct {
	Status        inauth.ServiceStatus  `json:"status"`
	AccessToken   string                `json:"access_token,omitempty"`
	IdentityToken *inauth.IdentityToken `json:"identity_token,omitempty"`
}

// AppAuth_TokenExchange exchanges a one-time auth code for access_token + identity_token.
// Called by third-party app backend using app credentials.
func AppAuth_TokenExchange(ctx *httpsrv.Context) error {

	var (
		req AppAuth_TokenExchangeRequest
		rsp AppAuth_TokenExchangeResponse
	)
	defer ctx.RenderJson(&rsp)

	if err := ctx.Request().JsonDecode(&req); err != nil {
		// rsp.Status = inauth.NewServiceStatus("400", "Invalid request format")
		// return nil
	}

	if req.Code == "" {
		if ctx.Params().Value("code") == "" {
			rsp.Status = inauth.NewServiceStatus("400", "code required")
			return nil
		}
		req.Code = ctx.Params().Value("code")
	}

	// look up auth code
	var entry iamapi.SignInAuthCode
	if rs := data.Data.NewReader(iamapi.NsAuthCode(req.Code)).Exec(); !rs.OK() {
		rsp.Status = inauth.NewServiceStatus("404", "Invalid or expired auth code")
		return nil
	} else {
		rs.Item().JsonDecode(&entry)
	}

	// one-time use: delete the code
	data.Data.NewDeleter(iamapi.NsAuthCode(req.Code)).Exec()

	if entry.AccessToken == "" {
		rsp.Status = inauth.NewServiceStatus("404", "Invalid auth code")
		return nil
	}

	// parse and verify the stored access token
	token, err := inauth.ParseAccessToken(entry.AccessToken)
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("401", "Invalid access token")
		return nil
	}

	if _, err := token.Verify(data.KeyMgr); err != nil {
		rsp.Status = inauth.NewServiceStatus("401", "Invalid access token")
		return nil
	}

	// lookup session for identity token
	var st inauth.SessionToken
	if rs := data.Data.NewReader(iamapi.NsUserSession(token.Claims.Jti, uint32(token.Claims.Exp))).Exec(); rs.OK() {
		rs.Item().JsonDecode(&st)
	}

	rsp.Status = inauth.NewServiceStatus("200", "ok")
	rsp.AccessToken = entry.AccessToken
	if st.IdentityToken != nil {
		rsp.IdentityToken = st.IdentityToken
	}

	js, _ := json.Marshal(rsp)
	slog.Info("service/auth-token-exchange response", "body", string(js))

	return nil
}
