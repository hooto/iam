// Copyright 2015 lessOS.com, All rights reserved.
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

package idclient

import (
	"errors"
	"sync"
	"time"

	"github.com/lessos/lessgo/net/httpclient"
	"github.com/lessos/lessgo/utils"
	"github.com/lessos/lessgo/utilx"
	"github.com/lessos/lessids/idsapi"
)

var (
	ServiceUrl        = "http://127.0.0.1:50101/ids"
	sessions          = map[string]idsapi.UserSession{}
	nextClean         = time.Now()
	innerExpiredRange = time.Second * 1800
	locker            sync.Mutex
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

	for k, v := range sessions {

		if v.InnerExpired.Before(time.Now()) {
			continue
		}

		delete(sessions, k)
	}

	nextClean = time.Now().Add(time.Second * 60)
}

func LoginUrl(backurl string) string {
	return ServiceUrl + "/service/login?continue=" + backurl
}

func Instance(token string) (session idsapi.UserSession, err error) {

	if ServiceUrl == "" || token == "" {
		return session, errors.New("Unauthorized")
	}

	if session, ok := sessions[token]; ok {
		return session, nil
	}

	hc := httpclient.Get(ServiceUrl + "/v1/service/auth?access_token=" + token)

	var us idsapi.UserSession

	err = hc.ReplyJson(&us)
	if err != nil || us.Error != nil || us.Kind != "UserSession" {
		return session, errors.New("Unauthorized")
	}

	us.InnerExpired = time.Now().Add(innerExpiredRange)

	exp := utilx.TimeParse(us.Expired, "atom")
	if us.InnerExpired.After(exp) {
		us.InnerExpired = exp
	}

	locker.Lock()
	sessions[token] = us // TODO Cache API
	locker.Unlock()

	return us, nil
}

func IsLogin(token string) bool {

	if _, err := Instance(token); err != nil {
		return false
	}

	return true
}

func AccessAllowed(privilege, instanceid, token string) bool {

	if !IsLogin(token) {
		return false
	}

	req := idsapi.UserAccessEntry{
		AccessToken: token,
		InstanceID:  instanceid,
		Privilege:   privilege,
	}

	js, _ := utils.JsonEncode(req)
	hc := httpclient.Post(ServiceUrl + "/v1/service/access-allowed")
	hc.Header("contentType", "application/json; charset=utf-8")
	hc.Body(js)

	var us idsapi.UserAccessEntry
	if err := hc.ReplyJson(&us); err != nil || us.Kind != "UserAccessEntry" {
		return false
	}

	return true
}
