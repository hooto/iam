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
	"errors"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/lessos/lessgo/encoding/json"
	"github.com/lessos/lessgo/net/httpclient"
	"github.com/lessos/lessgo/types"

	"github.com/hooto/hauth/go/hauth/v1"
	"github.com/hooto/httpsrv"
	"github.com/hooto/iam/iamapi"
)

const (
	AccessTokenKey = iamapi.AccessTokenKey
)

var (
	InstanceID         = ""
	InstanceOwner      = ""
	ServiceUrl         = "http://127.0.0.1:9528/iam"
	ServiceUrlFrontend = ""
	ServiceUrlGlobal   = ""
)

var (
	locker      sync.RWMutex
	c_roles           = iamapi.UserRoleList{}
	c_roles_tto int64 = 0
)

func Expired(ttl int) time.Time {
	return time.Now().Add(time.Second * time.Duration(ttl))
}

func service_url_global() string {
	if ServiceUrlGlobal != "" {
		return ServiceUrlGlobal
	}
	return ServiceUrl
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

func SessionIsLogin(s *httpsrv.Session) bool {

	if s == nil {
		return false
	}

	return tokenIsLogin(s.Get(AccessTokenKey))
}

func SessionAccessAllowed(s *httpsrv.Session, privilege, client_id string) bool {

	if s == nil {
		return false
	}

	return tokenAccessAllowed(privilege, s.Get(AccessTokenKey), client_id)
}

func SessionInstance(s *httpsrv.Session) (iamapi.UserSession, error) {

	if s == nil {
		return iamapi.UserSession{}, errors.New("No Session Found")
	}

	sess, err := Instance(s.Get(AccessTokenKey))
	if err != nil {
		return iamapi.UserSession{}, err
	}

	return *sess, nil
}

func Instance(token string) (*iamapi.UserSession, error) {

	if ServiceUrl == "" || token == "" {
		return nil, errors.New("Unauthorized")
	}

	ap, err := hauth.NewUserValidator(token)
	if err != nil {
		return nil, err
	}

	if ap.IsExpired() {
		return nil, errors.New("auth expired")
	}

	tn := time.Now().Unix()

	// cache
	session := SessionCache(ap.Id)
	if session != nil && !sessionCacheRefresh(session, tn) {
		return session, nil
	}

	hc := httpclient.Get(fmt.Sprintf(
		"%s/v1/service/auth?%s=%s",
		ServiceUrl,
		AccessTokenKey,
		token,
	))
	hc.SetTimeout(3000)
	defer hc.Close()

	var us types.TypeMeta
	if err = hc.ReplyJson(&us); err != nil {
		err = errors.New("Network error, please try again later")
	} else if us.Kind != "UserSession" {
		if us.Error != nil {
			err = errors.New("Unauthorized " + us.Error.Message)
		} else {
			err = errors.New("Unauthorized")
		}
	} else {

		session = &iamapi.UserSession{
			AccessToken: token,
			UserName:    ap.Id,
			DisplayName: ap.Name,
			Roles:       ap.Roles,
			Groups:      ap.Groups,
			Expired:     ap.Expired,
		}
		SessionSync(session, tn)
	}

	if err != nil && session == nil {
		return nil, err
	}

	return session, nil
}

func tokenIsLogin(token string) bool {

	if _, err := Instance(token); err != nil {
		return false
	}

	return true
}

func tokenAccessAllowed(privilege, token, instanceid string) bool {

	if !tokenIsLogin(token) {
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
	hc.SetTimeout(3000)
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
	hc.SetTimeout(3000)
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
