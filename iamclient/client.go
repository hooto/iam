// Copyright 2014 lessos Authors, All rights reserved.
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
	"errors"
	"fmt"
	"sync"
	"time"

	"code.hooto.com/lessos/iam/iamapi"
	"github.com/lessos/lessgo/encoding/json"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/net/httpclient"
	"github.com/lessos/lessgo/types"
)

const (
	AccessTokenKey = iamapi.AccessTokenKey
)

var (
	InstanceID         = ""
	InstanceOwner      = ""
	ServiceUrl         = "http://127.0.0.1:9528/iam"
	ServiceUrlFrontend = ""
	sessions           = map[string]iamapi.UserSession{}
	nextClean          = time.Now()
	locker             sync.Mutex
)

func Expired(ttl int) time.Time {
	return time.Now().Add(time.Second * time.Duration(ttl))
}

func innerExpiredClean() {

	if nextClean.After(time.Now()) {
		return
	}

	locker.Lock()
	defer locker.Unlock()

	tnow := types.MetaTimeNow()

	for k, v := range sessions {

		if v.Expired > tnow {
			continue
		}

		delete(sessions, k)
	}

	nextClean = time.Now().Add(time.Second * 60)
}

func LoginUrl(backurl string) string {

	if ServiceUrlFrontend != "" {
		return ServiceUrlFrontend + "/service/login?continue=" + backurl
	}

	return ServiceUrl + "/service/login?continue=" + backurl
}

func AuthServiceUrl(client_id, redirect_uri, state string) string {

	if ServiceUrlFrontend != "" {
		return auth_service_url(ServiceUrlFrontend, client_id, redirect_uri, state)
	}

	return auth_service_url(ServiceUrl, client_id, redirect_uri, state)
}

func auth_service_url(service_url, client_id, redirect_uri, state string) string {
	return fmt.Sprintf("%s/service/login?response_type=token&client_id=%s&redirect_uri=%s&state=%s",
		service_url, client_id, redirect_uri, state)
}

func SessionAccessToken(s *httpsrv.Session) string {
	return s.Get(AccessTokenKey)
}

func SesionSet(s *httpsrv.Session) error {

	return nil
}

func SessionIsLogin(s *httpsrv.Session) bool {

	if s == nil {
		return false
	}

	return _is_login(s.Get(AccessTokenKey))
}

func SessionAccessAllowed(s *httpsrv.Session, privilege, client_id string) bool {

	if s == nil {
		return false
	}

	return _access_allowed(privilege, s.Get(AccessTokenKey), client_id)
}

func SessionInstance(s *httpsrv.Session) (session iamapi.UserSession, err error) {

	if s == nil {
		return iamapi.UserSession{}, errors.New("No Session Found")
	}

	return Instance(s.Get(AccessTokenKey))
}

func Instance(token string) (session iamapi.UserSession, err error) {

	if ServiceUrl == "" || token == "" {
		return session, errors.New("Unauthorized")
	}

	if session, ok := sessions[token]; ok {
		return session, nil
	}

	hc := httpclient.Get(fmt.Sprintf(
		"%s/v1/service/auth?%s=%s",
		ServiceUrl,
		AccessTokenKey,
		token,
	))
	defer hc.Close()

	var us iamapi.UserSession

	err = hc.ReplyJson(&us)
	if err != nil || us.Error != nil || us.Kind != "UserSession" {
		return session, errors.New("Unauthorized")
	}

	if types.MetaTimeNow() > us.Expired {
		return session, errors.New("Unauthorized")
	}

	locker.Lock()
	sessions[token] = us // TODO Cache API
	locker.Unlock()

	return us, nil
}

func _is_login(token string) bool {

	if _, err := Instance(token); err != nil {
		return false
	}

	return true
}

func _access_allowed(privilege, token, instanceid string) bool {

	if !_is_login(token) {
		return false
	}

	req := iamapi.UserAccessEntry{
		AccessToken: token,
		InstanceID:  instanceid,
		Privilege:   privilege,
	}

	js, _ := json.Encode(req, "")
	hc := httpclient.Post(ServiceUrl + "/v1/service/access-allowed")
	hc.Header("contentType", "application/json; charset=utf-8")
	hc.Body(js)
	defer hc.Close()

	var us iamapi.UserAccessEntry
	if err := hc.ReplyJson(&us); err != nil || us.Kind != "UserAccessEntry" {
		return false
	}

	return true
}
