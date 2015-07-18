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

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/lessos/bigtree/btapi"
	"github.com/lessos/lessgo/net/email"
	"github.com/lessos/lessids/idsapi"
	"github.com/lessos/lessids/store"
)

const (
	Version       = "0.2.1dev"
	GroupMember   = 100
	GroupSysAdmin = 1
)

var (
	err    error
	Config = ConfigCommon{
		ServiceName:      "lessOS Identity Service",
		WebUiBannerTitle: "Account Center",
	}
)

type ConfigCommon struct {
	Port             uint16 `json:"port"`
	Prefix           string
	ServiceName      string `json:"service_name"`
	WebUiBannerTitle string
	Mailer           idsapi.SysConfigMailer `json:"mailer"`
	BtData           btapi.DataAccessConfig `json:"btdata,omitempty"`
}

func Init(prefix string) error {

	if prefix == "" {
		prefix, err = filepath.Abs(filepath.Dir(os.Args[0]) + "/..")
		if err != nil {
			prefix = "/opt/lessos/ids"
		}
	}
	reg, _ := regexp.Compile("/+")
	Config.Prefix = "/" + strings.Trim(reg.ReplaceAllString(prefix, "/"), "/")

	file := Config.Prefix + "/etc/main.json"
	if _, err := os.Stat(file); err != nil && os.IsNotExist(err) {
		file = Config.Prefix + "/etc/main.json.dev"
	}
	if _, err := os.Stat(file); err != nil && os.IsNotExist(err) {
		return fmt.Errorf("Error: config file is not exists %s", Config.Prefix+"/etc/main.json")
	}

	fp, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("Error: Can not open (%s)", file)
	}
	defer fp.Close()

	cfgstr, err := ioutil.ReadAll(fp)
	if err != nil {
		return fmt.Errorf("Error: Can not read (%s)", file)
	}

	if err = json.Unmarshal(cfgstr, &Config); err != nil {
		return fmt.Errorf("Error: config file invalid. (%s)", err.Error())
	}

	Config.Refresh()

	return nil
}

func (c *ConfigCommon) Refresh() {

	//
	var mailer idsapi.SysConfigMailer
	if obj := store.BtAgent.ObjectGet(btapi.ObjectProposal{
		Meta: btapi.ObjectMeta{
			Path: "/sys-config/mailer",
		},
	}); obj.Error == nil {
		obj.JsonDecode(&mailer)

		if c.Mailer.SmtpHost != "" {

			email.MailerRegister("def",
				c.Mailer.SmtpHost,
				c.Mailer.SmtpPort,
				c.Mailer.SmtpUser,
				c.Mailer.SmtpPass)
		}
	}

	if obj := store.BtAgent.ObjectGet(btapi.ObjectProposal{
		Meta: btapi.ObjectMeta{
			Path: "/sys-config/service_name",
		},
	}); obj.Error == nil {

		if obj.Data != "" {
			c.ServiceName = obj.Data
		}
	}

	if obj := store.BtAgent.ObjectGet(btapi.ObjectProposal{
		Meta: btapi.ObjectMeta{
			Path: "/sys-config/webui_banner_title",
		},
	}); obj.Error == nil {

		if obj.Data != "" {
			c.WebUiBannerTitle = obj.Data
		}
	}
}
