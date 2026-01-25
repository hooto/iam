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

package v2

import (
	"time"

	hauth2 "github.com/hooto/hauth/go/v2"
	"github.com/hooto/httpsrv"
	"github.com/hooto/iam/data"
	"github.com/hooto/iam/iamapi"
)

type Service struct {
	*httpsrv.Controller
}

func (c Service) AuthAction() {

	var (
		req hauth2.AuthLoginRequest
		rsp hauth2.AuthLoginResponse
	)
	defer c.RenderJson(&rsp)

	if err := c.Request.JsonDecode(&req); err != nil {
		rsp.Error = err.Error()
		return
	}

	token, err := hauth2.NewAccessToken(req.LoginToken)
	if err != nil {
		rsp.Error = err.Error()
		return
	}

	ak, err := token.Verify(data.KeyMgr)
	if err != nil {
		rsp.Error = err.Error()
		return
	}

	if ak.User == "" {
		rsp.Error = "access-key not found"
		return
	}

	user := data.UserGet(ak.User)
	if user == nil {
		rsp.Error = "incorrect username or password"
		return
	}

	if user.Type == iamapi.UserTypeGroup {
		rsp.Error = "incorrect username or password"
		return
	}

	iat := time.Now().Unix()

	header := hauth2.TokenHeader{
		Kid: ak.Id,
	}

	claims := hauth2.AccessTokenClaims{
		Jti: hauth2.RandHexString(16),
		Iat: iat,
		Exp: iat + 86400,
		Sub: ak.User,
	}

	at, err := hauth2.Sign(header, claims, []byte(ak.Secret))
	if err != nil {
		rsp.Error = err.Error()
		return
	}
	rsp.AccessToken = at

	rsp.IdentityToken = hauth2.IdentityToken{
		Jti:    claims.Jti,
		Iat:    claims.Iat,
		Exp:    claims.Exp,
		Sub:    claims.Sub,
		Roles:  user.Roles,
		Groups: data.UserGroups(claims.Sub),
		Scopes: ak.Scopes,
	}
}
