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

type AppAuthConfig struct {
	AppId     string `json:"app_id" yaml:"app_id" toml:"app_id"`
	SecretKey string `json:"secret_key" yaml:"secret_key" toml:"secret_key"`
	BaseURL   string `json:"base_url" yaml:"base_url" toml:"base_url"`

	Flush func() error `json:"-" yaml:"-" toml:"-"`
}

func (it *AppAuthConfig) Valid() error {
	if it == nil || it.AppId == "" || it.BaseURL == "" || it.SecretKey == "" {
		return errors.New("iam app-auth-config not initialized")
	}
	return nil
}

func (it *AppAuthConfig) NewAppCredential() (inauth.AppCredential, error) {
	if err := it.Valid(); err != nil {
		return nil, err
	}
	ak := inauth.NewAppAccessKey()
	ak.Id = it.AppId
	ak.Secret = it.SecretKey

	ac := inauth.NewAppCredential(ak)

	return ac, nil
}
