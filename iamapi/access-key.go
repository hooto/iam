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

package iamapi

import (
	"encoding/base64"
	"errors"

	"github.com/hooto/iam/iamauth"
	"github.com/lessos/lessgo/encoding/json"
	"github.com/lessos/lessgo/types"
)

type AccessKey struct {
	User        string           `json:"user,omitempty"`
	AccessKey   string           `json:"access_key"`
	SecretKey   string           `json:"secret_key,omitempty"`
	Created     uint64           `json:"created,omitempty"`
	Action      uint16           `json:"action,omitempty"`
	Description string           `json:"desc,omitempty"`
	Bounds      []AccessKeyBound `json:"bounds,omitempty"`
}

func (it *AccessKey) AuthKey() *iamauth.AuthKey {
	return &iamauth.AuthKey{
		User:      it.User,
		AccessKey: it.AccessKey,
		SecretKey: it.SecretKey,
	}
}

type AccessKeyBound struct {
	Name    string `json:"name"`
	Created uint64 `json:"created,omitempty"`
}

func (it AccessKeyBound) IterKey() string {
	return it.Name
}

type AccessKeyBounds []AccessKeyBound

type AccessKeyAuth struct {
	Type  string `json:"t"`
	User  string `json:"u"`
	Key   string `json:"k"`
	Time  int64  `json:"rt"`
	Token string `json:"tk"`
}

func (t AccessKeyAuth) Encode() string {
	bs, _ := json.Encode(t, "")
	return base64.StdEncoding.EncodeToString(bs)
}

func (t AccessKeyAuth) Valid() error {

	//
	if len(t.Type) == 0 {
		return errors.New("No Auth Type Found")
	}

	//
	if len(t.User) == 0 {
		return errors.New("No Auth User Found")
	}

	//
	if len(t.Key) == 0 {
		return errors.New("No Auth AccessKey Found")
	}

	if t.Time < 1000000000 {
		return errors.New("Invalid Request Time")
	}

	//
	if len(t.Token) < 30 {
		return errors.New("No Auth Token Found")
	}

	return nil
}

func AccessKeyAuthDecode(auth string) (AccessKeyAuth, error) {

	var t AccessKeyAuth
	if len(auth) < 30 {
		return t, errors.New("Unauthorized")
	}

	bs, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		return t, err
	}

	if err = json.Decode(bs, &t); err != nil {
		return t, err
	}

	err = t.Valid()

	return t, err
}

// Access Key SESSION
// K1(2)VERIFY-SIGNATURE(36)PAYLOAD-DATA
type AccessKeySession struct {
	AccessKey string            `json:"ak"`
	SecretKey string            `json:"sk"`
	User      string            `json:"ur"`
	Roles     types.ArrayUint32 `json:"rs"`
	Expired   int64             `json:"ex"` // unix seconds
}
