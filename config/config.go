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

package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lessos/lessgo/encoding/json"
)

const (
	Version       = "0.2.1dev"
	GroupMember   = 100
	GroupSysAdmin = 1
)

var (
	err                      error
	Prefix                   = "/opt/lessos/iam"
	UserRegistrationDisabled = false
	Config                   = ConfigCommon{
		InstanceID:       "",
		ServiceName:      "lessOS IAM Service",
		WebUiBannerTitle: "Account Center",
	}
)

type ConfigCommon struct {
	InstanceID       string `json:"instance_id"`
	HttpPort         uint16 `json:"http_port,omitempty"`
	ServiceName      string `json:"service_name"`
	WebUiBannerTitle string
}

func Init(prefix string) error {

	if prefix == "" {
		if prefix, err = filepath.Abs(filepath.Dir(os.Args[0]) + "/.."); err == nil {
			Prefix = prefix
		}
	}
	Prefix = filepath.Clean(Prefix)

	file := Prefix + "/etc/iam_config.json"
	if _, err := os.Stat(file); err != nil && os.IsNotExist(err) {
		file = Prefix + "/etc/iam_config.json.dev"
	}
	if _, err := os.Stat(file); err != nil && os.IsNotExist(err) {
		return fmt.Errorf("Error: config file is not exists %s", Prefix+"/etc/iam_config.json")
	}

	if Config.InstanceID == "" {
		return fmt.Errorf("No InstanceID Found")
	}

	if err := json.DecodeFile(file, &Config); err != nil {
		return err
	}

	return InitConfig()
}

func InitConfig() error {

	if Config.InstanceID == "" {
		return fmt.Errorf("No InstanceID Found")
	}

	return nil
}
