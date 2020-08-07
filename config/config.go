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

package config

import (
	"os"
	"path/filepath"

	"github.com/hooto/htoml4g/htoml"
	"github.com/lessos/lessgo/crypto/idhash"
	"github.com/lessos/lessgo/encoding/json"

	"github.com/hooto/hauth/go/hauth/v1"
)

const (
	Version       = "0.9.1"
	GroupMember   = 100
	GroupSysAdmin = 1
)

var (
	err                      error
	VersionHash              = Version
	Prefix                   = "/opt/hooto/iam"
	UserRegistrationDisabled = false
	Config                   = ConfigCommon{
		InstanceID:               "",
		ServiceName:              "hooto IAM Service",
		WebUiBannerTitle:         "Account Panel",
		ServiceLoginFormAlertMsg: "",
	}
	AuthKeyMgr     = hauth.NewAuthKeyManager()
	configFilePath = ""
)

type ConfigCommon struct {
	filepath                 string           `json:"-" toml:"-"`
	InstanceID               string           `json:"instance_id" toml:"instance_id"`
	HttpPort                 uint16           `json:"http_port,omitempty" toml:"http_port,omitempty"`
	ServiceName              string           `json:"service_name" toml:"service_name"`
	AuthKeys                 []*hauth.AuthKey `json:"auth_keys" toml:"auth_keys"`
	WebUiBannerTitle         string           `json:"-" toml:"-"`
	ServiceLoginFormAlertMsg string           `json:"-" toml:"-"`
}

func Setup(prefix string) error {

	if prefix == "" {
		if prefix, err = filepath.Abs(filepath.Dir(os.Args[0]) + "/.."); err == nil {
			Prefix = prefix
		}
	} else {
		Prefix = prefix
	}
	Prefix = filepath.Clean(Prefix)

	configFilePath = Prefix + "/etc/iam_config.toml"

	if err := htoml.DecodeFromFile(&Config, configFilePath); err != nil {

		if !os.IsNotExist(err) {
			return err
		}

		if err := json.DecodeFile(Prefix+"/etc/iam_config.json", &Config); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	if Config.InstanceID == "" {
		Config.InstanceID = idhash.RandHexString(16)
	}

	if Config.ServiceName == "" {
		Config.ServiceName = "hooto IAM Service"
	}

	if len(Config.AuthKeys) < 1 {
		Config.AuthKeys = append(Config.AuthKeys, &hauth.AuthKey{
			AccessKey: idhash.RandHexString(8),
			SecretKey: idhash.RandBase64String(40),
		})
	}
	for _, v := range Config.AuthKeys {
		AuthKeyMgr.KeySet(v)
	}

	return Flush()
}

func Flush() error {
	return htoml.EncodeToFile(Config, configFilePath, nil)
}
