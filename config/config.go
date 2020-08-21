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

	"github.com/hooto/hauth/go/hauth/v1"
	"github.com/hooto/htoml4g/htoml"
	"github.com/lessos/lessgo/crypto/idhash"
	"github.com/lessos/lessgo/encoding/json"
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
	Config                   = &ConfigCommon{}
	configFilePath           = ""
)

type ConfigCommon struct {
	InstanceID               string             `json:"instance_id" toml:"instance_id"`
	HttpPort                 uint16             `json:"http_port,omitempty" toml:"http_port,omitempty"`
	ServiceName              string             `json:"service_name" toml:"service_name"`
	AccessKeys               []*hauth.AccessKey `json:"access_keys" toml:"access_keys"`
	WebUiBannerTitle         string             `json:"-" toml:"-"`
	ServiceLoginFormAlertMsg string             `json:"-" toml:"-"`
}

func (it *ConfigCommon) reset() {

	if it.InstanceID == "" {
		it.InstanceID = idhash.RandHexString(16)
	}

	if it.ServiceName == "" {
		it.ServiceName = "hooto IAM Service"
	}

	if it.WebUiBannerTitle == "" {
		it.WebUiBannerTitle = "Account Panel"
	}

	if len(it.AccessKeys) < 1 {
		it.AccessKeys = append(it.AccessKeys, &hauth.AccessKey{
			Id:     idhash.RandHexString(8),
			Secret: idhash.RandBase64String(40),
		})
	}
}

func Setup(prefix string) error {

	if prefix != "" {
		Prefix = filepath.Clean(prefix)
	} else {
		Prefix, _ = filepath.Abs(filepath.Dir(os.Args[0]) + "/..")
	}

	configFilePath = Prefix + "/etc/iam_config.toml"

	if err := htoml.DecodeFromFile(Config, configFilePath); err != nil {

		if !os.IsNotExist(err) {
			return err
		}

		if err := json.DecodeFile(Prefix+"/etc/iam_config.json", &Config); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	Config.reset()

	return Flush()
}

func SetupConfig(prefix string, cfg *ConfigCommon) error {

	if prefix != "" {
		Prefix = filepath.Clean(prefix)
	} else {
		Prefix, _ = filepath.Abs(filepath.Dir(os.Args[0]) + "/..")
	}

	if cfg == nil {
		cfg = &ConfigCommon{}
	}
	cfg.reset()

	Config = cfg
	return nil
}

func Flush() error {
	if configFilePath != "" {
		return htoml.EncodeToFile(Config, configFilePath, nil)
	}
	return nil
}
