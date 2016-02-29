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
	"time"

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
	err                      error
	Prefix                   = "/opt/lessos/ids"
	UserRegistrationDisabled = false
	Config                   = ConfigCommon{
		ServiceName:      "lessOS Identity Service",
		WebUiBannerTitle: "Account Center",
	}
)

type ConfigCommon struct {
	Port             uint16 `json:"port"`
	ServiceName      string `json:"service_name"`
	WebUiBannerTitle string
	// Mailer           idsapi.SysConfigMailer `json:"mailer"`
	// BtData btapi.DataAccessConfig `json:"btdata,omitempty"`
}

func Init(prefix string) error {

	if prefix == "" {
		prefix, err = filepath.Abs(filepath.Dir(os.Args[0]) + "/..")
		if err == nil {
			Prefix = prefix
		}
	}
	reg, _ := regexp.Compile("/+")
	Prefix = "/" + strings.Trim(reg.ReplaceAllString(Prefix, "/"), "/")

	file := Prefix + "/etc/main.json"
	if _, err := os.Stat(file); err != nil && os.IsNotExist(err) {
		file = Prefix + "/etc/main.json.dev"
	}
	if _, err := os.Stat(file); err != nil && os.IsNotExist(err) {
		return fmt.Errorf("Error: config file is not exists %s", Prefix+"/etc/main.json")
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

	InitConfig()

	return InitConfig()
}

func InitConfig() error {

	go func() {

		for {

			if store.Ready {
				Config.Refresh()
				break
			}

			time.Sleep(1e9)
		}
	}()

	return nil
}

func (c *ConfigCommon) Refresh() {

	//
	var mailer idsapi.SysConfigMailer
	if obj := store.BtAgent.ObjectGet("/global/ids/sys-config/mailer"); obj.Error == nil {

		obj.JsonDecode(&mailer)

		if mailer.SmtpHost != "" {

			email.MailerRegister("def",
				mailer.SmtpHost,
				mailer.SmtpPort,
				mailer.SmtpUser,
				mailer.SmtpPass)
		}
	}

	if obj := store.BtAgent.ObjectGet("/global/ids/sys-config/service_name"); obj.Error == nil {

		if obj.Data != "" {
			c.ServiceName = obj.Data
		}
	}

	if obj := store.BtAgent.ObjectGet("/global/ids/sys-config/webui_banner_title"); obj.Error == nil {

		if obj.Data != "" {
			c.WebUiBannerTitle = obj.Data
		}
	}

	if obj := store.BtAgent.ObjectGet("/global/ids/sys-config/user_reg_disable"); obj.Error == nil {

		if obj.Data == "1" {
			UserRegistrationDisabled = true
		} else {
			UserRegistrationDisabled = false
		}
	}
}
