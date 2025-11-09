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

package iamclient

import (
	"encoding/json"
	"errors"
	"sync"

	hauth2 "github.com/hooto/hauth/v2/hauth"

	"github.com/lessos/lessgo/net/httpclient"
)

func NewIdentityAuthService(
	sessionMgr hauth2.SessionTokenManager,
) hauth2.IdentityAuthService {
	return &identityAuthService{
		sessionMgr: sessionMgr,
	}
}

type identityAuthService struct {
	mu         sync.RWMutex
	sessionMgr hauth2.SessionTokenManager
}

func (it *identityAuthService) AuthLogin(req *hauth2.AuthLoginRequest) (*hauth2.AuthLoginResponse, error) {

	js, _ := json.Marshal(req)

	hc := httpclient.Post(ServiceUrl + "/v2/service/auth")
	hc.Header("contentType", "application/json; charset=utf-8")
	hc.SetTimeout(3000)
	hc.Body(js)
	defer hc.Close()

	var rsp hauth2.AuthLoginResponse
	if err := hc.ReplyJson(&rsp); err != nil {
		return nil, err
	} else if rsp.Error != "" {
		return nil, errors.New(rsp.Error)
	}

	accessToken, err := it.sessionMgr.ReSign(rsp.AccessToken, rsp.IdentityToken)
	if err != nil {
		return nil, err
	}
	rsp.AccessToken = accessToken

	return &rsp, nil
}
