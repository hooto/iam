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

package iamserver

import (
	"errors"

	"github.com/sysinner/incore/v2/pkg/inauth"
)

var AppConfig *AppAuthConfig

type AppAuthConfig struct {
	AppId     string `json:"app_id" yaml:"app_id" toml:"app_id"`
	SecretKey string `json:"secret_key" yaml:"secret_key" toml:"secret_key"`
	Endpoint  string `json:"endpoint" yaml:"endpoint" toml:"endpoint"`

	validTry  int   `json:"-" toml:"-"`
	validFail error `json:"-" toml:"-"`

	SaveFunc func() error `json:"-" toml:"-"`
}

func (it *AppAuthConfig) NewAppCredential() (inauth.AppCredential, error) {
	if err := it.Verify(); err != nil {
		return nil, err
	}
	ak := inauth.NewAppAccessKey()
	ak.Id = it.AppId
	ak.Secret = it.SecretKey

	ac := inauth.NewAppCredential(ak)

	return ac, nil
}

func (it *AppAuthConfig) Verify() error {

	if it == nil || it.AppId == "" || it.Endpoint == "" || it.SecretKey == "" {
		return errors.New("iam app-auth-config not initialized")
	}

	if it.validTry > 0 && it.validFail == nil {
		return nil
	}

	ak := inauth.NewAppAccessKey()
	ak.Id = it.AppId
	ak.Secret = it.SecretKey

	ac := inauth.NewAppCredential(ak)
	at := ac.AuthToken()

	it.validTry += 1

	var rsp struct {
		Status inauth.ServiceStatus `json:"status"`
	}
	if err := iamPost(
		it.Endpoint, "/v2/open/app-auth/verify",
		at,
		map[string]string{},
		&rsp,
	); err != nil {
		return err
	}

	if rsp.Status.Code != "200" {
		return errors.New("Verification failed: " + rsp.Status.Message)
	}

	return nil
}
