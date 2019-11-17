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

	"github.com/lessos/lessgo/crypto/idhash"
	"github.com/lessos/lessgo/encoding/json"

	"github.com/hooto/iam/iamauth"
)

const (
	Version       = "0.9.0"
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
)

type ConfigCommon struct {
	filepath                 string             `json:"-"`
	InstanceID               string             `json:"instance_id"`
	HttpPort                 uint16             `json:"http_port,omitempty"`
	ServiceName              string             `json:"service_name"`
	AuthKeys                 []*iamauth.AuthKey `json:"auth_keys"`
	WebUiBannerTitle         string             `json:"-"`
	ServiceLoginFormAlertMsg string             `json:"-"`
}

func (it *ConfigCommon) Flush() error {
	if it.filepath == "" {
		return nil
	}
	return json.EncodeToFile(it, it.filepath, "  ")
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

	if err := json.DecodeFile(Prefix+"/etc/iam_config.json", &Config); err != nil && !os.IsNotExist(err) {
		return err
	}
	Config.filepath = Prefix + "/etc/iam_config.json"

	if Config.InstanceID == "" {
		Config.InstanceID = idhash.RandHexString(16)
	}

	if Config.ServiceName == "" {
		Config.ServiceName = "hooto IAM Service"
	}

	if len(Config.AuthKeys) < 1 {
		Config.AuthKeys = append(Config.AuthKeys, &iamauth.AuthKey{
			AccessKey: idhash.RandHexString(8),
			SecretKey: idhash.RandBase64String(40),
		})
	}

	return Config.Flush()
}
