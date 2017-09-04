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
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/hooto/iam/iamapi"
	"github.com/lessos/lessgo/net/httpclient"
	"github.com/lessos/lessgo/types"
)

const (
	ns_auth_type              = "la"
	ns_auth_type_sha256       = "sha256"
	ns_auth_type_algdef       = ns_auth_type_sha256
	ns_auth_token             = "tk"
	ns_auth_client_id         = "c"
	ns_auth_reqtime           = "rt"
	auth_reqtime_range  int64 = 600
)

func AccessKeySession(app_aka, user_aka iamapi.AccessKeyAuth) (iamapi.AccessKeySession, error) {

	if session, ok := sessions_aks[user_aka.Key]; ok {
		return session, nil
	}

	hc := httpclient.Get(fmt.Sprintf(
		"%s/v1/app-auth/user-access-key?user=%s&access_key=%s",
		ServiceUrl,
		user_aka.User,
		user_aka.Key,
	))
	defer hc.Close()

	hc.Header("Auth", app_aka.Encode())

	var session iamapi.AccessKeySession

	err := hc.ReplyJson(&session)
	if err != nil || session.SecretKey == "" {
		return session, errors.New("Unauthorized")
	}

	if types.MetaTimeNow() > session.Expired {
		return session, errors.New("Unauthorized")
	}

	locker.Lock()
	sessions_aks[session.AccessKey] = session // TODO Cache API
	locker.Unlock()

	return session, nil
}

/*
func AccessKeyAuthValidString(auth, secret_key string) error {

	aka, err := iamapi.AccessKeyAuthDecode(auth)
	if err != nil {
		return err
	}

	return AccessKeyAuthValid(aka, secret_key)
}
*/

func AccessKeyAuthValid(aka iamapi.AccessKeyAuth, secret_key string) error {

	if err := aka.Valid(); err != nil {
		return err
	}

	rtli := auth_time()
	if pos := rtli - aka.Time; pos < -auth_reqtime_range || pos > auth_reqtime_range {
		return errors.New("Invalid Request Time")
	}

	if aka.Token != token_encode(
		aka.Type, aka.Key, aka.Time, secret_key) {
		return errors.New("Invalid Token")
	}

	return nil
}

func NewAccessKeyAuth(user, access_key, secret_key, data string) (iamapi.AccessKeyAuth, error) {

	rt := auth_time()

	return iamapi.AccessKeyAuth{
		Type: ns_auth_type_algdef,
		User: user,
		Key:  access_key,
		Time: rt,
		Token: token_encode(
			ns_auth_type_algdef,
			access_key,
			rt,
			secret_key,
		),
	}, nil
}

func auth_time() int64 {
	return time.Now().UTC().Unix()
}

func token_encode(t string, client_id string, rtime int64, secret_key string) (rs string) {

	switch t {

	case ns_auth_type_sha256:
		h := sha256.New()
		h.Write([]byte(fmt.Sprintf("%s.%d.%s", client_id, rtime, secret_key)))
		rs = fmt.Sprintf("%x", h.Sum(nil))

	default:
		rs = "<nil>"
	}

	return
}
