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

package auth

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hooto/iam/iamapi"
	"github.com/lessos/lessgo/types"
)

const (
	HttpHeaderKey             = "x-iam-auth"
	ns_auth_type_sha256       = "sha256"
	ns_auth_type_algdef       = ns_auth_type_sha256
	auth_reqtime_range  int64 = 600
	ErrUnAuth                 = "UnAuth"
)

type AuthToken struct {
	ak           iamapi.AccessKey
	token_ar     []string
	typ          string
	User         string
	AccessKey    string
	request_time string
	token_enc    string
}

func NewAuthToken(token string) (*AuthToken, error) {

	t := &AuthToken{
		token_ar: strings.Split(token, "."),
	}

	if len(t.token_ar) < 5 {
		return t, errors.New("Invalid Auth Token")
	}

	switch t.token_ar[0] {

	case ns_auth_type_algdef:
		t.typ = t.token_ar[0]
		t.User = t.token_ar[1]
		t.AccessKey = t.token_ar[2]
		t.request_time = t.token_ar[3]
		t.token_enc = t.token_ar[4]

	default:
		return t, errors.New("Invalid Auth Token")
	}

	return t, nil
}

func (t *AuthToken) Valid(ak iamapi.AccessKey, payload []byte) *types.ErrorMeta {

	if ak.AccessKey != t.AccessKey {
		return types.NewErrorMeta(ErrUnAuth, "Invalid AccessKey")
	}

	rti, err := strconv.ParseInt(t.request_time, 10, 64)
	if err != nil || rti < 1000000000 {
		return types.NewErrorMeta(ErrUnAuth, "Invalid Request Time")
	}
	rtli := time.Now().UTC().Unix()
	if pos := rtli - rti; pos < -auth_reqtime_range || pos > auth_reqtime_range {
		return types.NewErrorMeta(ErrUnAuth, "Invalid Request Time")
	}

	if len(t.token_enc) < 10 {
		return types.NewErrorMeta(ErrUnAuth, "No token Found")
	}

	//
	if t.token_enc != token_encode(t.typ, ak.AccessKey, ak.SecretKey, t.request_time, payload) {
		return types.NewErrorMeta(ErrUnAuth, "Invalid Token")
	}

	return nil
}

type AuthCredentials struct {
	ak      iamapi.AccessKey
	payload []byte
}

func NewAuthCredentials(ak iamapi.AccessKey, payload []byte) AuthCredentials {
	return AuthCredentials{
		ak:      ak,
		payload: payload,
	}
}

func (s *AuthCredentials) HttpHeaderKey() string {
	return HttpHeaderKey
}

func (s *AuthCredentials) HttpHeaderValue() string {

	rt := time_request()

	// TODO
	return fmt.Sprintf("%s.%s.%s.%s.%s",
		ns_auth_type_algdef,
		s.ak.User,
		s.ak.AccessKey,
		rt,
		token_encode(
			ns_auth_type_algdef,
			s.ak.AccessKey,
			s.ak.SecretKey,
			rt,
			[]byte{},
		),
	)
}

func time_request() string {
	return strconv.FormatInt(time.Now().UTC().Unix(), 10)
}

func token_encode(t, access_key, secret_key, rtime string, payload []byte) (rs string) {

	// fmt.Println(time.Now(), t, access_key, secret_key, rtime)

	switch t {

	case ns_auth_type_sha256:
		h := sha256.New()
		h.Write([]byte(fmt.Sprintf("%s.%s.%s", access_key, secret_key, rtime)))
		rs = fmt.Sprintf("%x", h.Sum(nil))

	default:
		rs = "<nil>"
	}

	return
}
