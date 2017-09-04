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
	"net/url"
	"sync"
	"time"

	"github.com/hooto/httpsrv"
	"github.com/hooto/iam/iamapi"
	"github.com/lessos/lessgo/encoding/json"
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
)

var (
	sessions     = map[string]iamapi.UserSession{}
	sessions_aks = map[string]iamapi.AccessKeySession{}
	sessions_tto = time.Now()
	locker       sync.RWMutex
	c_roles            = iamapi.UserRoleList{}
	c_roles_tto  int64 = 0
)

func Expired(ttl int) time.Time {
	return time.Now().Add(time.Second * time.Duration(ttl))
}

func innerExpiredClean() {

	if sessions_tto.After(time.Now()) {
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

	for k, v := range sessions_aks {

		if v.Expired > tnow {
			continue
		}

		delete(sessions_aks, k)
	}

	sessions_tto = time.Now().Add(time.Second * 60)
}

func service_prefix() string {

	if ServiceUrlFrontend != "" {
		return ServiceUrlFrontend
	}

	return ServiceUrl
}

func LoginUrl(backurl string) string {
	return service_prefix() + "/service/login?continue=" + backurl
}

func AuthServiceUrl(client_id, redirect_uri, state string) string {
	return auth_service_url(service_prefix(), client_id, redirect_uri, state)
}

func auth_service_url(service_url, client_id, redirect_uri, state string) string {
	return fmt.Sprintf(
		"%s/service/login?response_type=token&client_id=%s&redirect_uri=%s&state=%s",
		service_url, client_id, url.QueryEscape(redirect_uri), url.QueryEscape(state),
	)
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
	if err != nil || us.UserName == "" {
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

func AppRoleList(s *httpsrv.Session, appid string) (*iamapi.UserRoleList, error) {

	token := s.Get(AccessTokenKey)
	if ServiceUrl == "" || token == "" {
		return nil, errors.New("Unauthorized")
	}

	tnu := time.Now().UTC().Unix()
	if c_roles_tto > tnu {
		return &c_roles, nil
	}

	hc := httpclient.Get(fmt.Sprintf(
		"%s/v1/app-auth/role-list?%s=%s&appid=%s",
		ServiceUrl,
		AccessTokenKey,
		token,
		appid,
	))
	defer hc.Close()

	var rls iamapi.UserRoleList
	if err := hc.ReplyJson(&rls); err != nil {
		return nil, err
	} else if rls.Error != nil || rls.Kind != "UserRoleList" {
		return nil, errors.New("Network is unreachable, Please try again later")
	}

	c_roles_tto = tnu
	c_roles = rls

	return &c_roles, nil
}
