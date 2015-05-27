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

	"github.com/lessos/lessgo/data/rdo"
	"github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/net/email"
)

const (
	Version       = "0.2.1dev"
	GroupMember   = 100
	GroupSysAdmin = 1
)

var (
	err    error
	Config ConfigCommon
)

type ConfigMailer struct {
	SmtpHost string `json:"smtp_host"`
	SmtpPort string `json:"smtp_port"`
	SmtpUser string `json:"smtp_user"`
	SmtpPass string `json:"smtp_pass"`
}

type ConfigCommon struct {
	Port             uint16 `json:"port"`
	Prefix           string
	ServiceName      string `json:"service_name"`
	WebUiBannerTitle string
	Mailer           ConfigMailer `json:"mailer"`
	Database         base.Config  `json:"database"`
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

	if Config.ServiceName == "" {
		Config.ServiceName = "lessOS Identity Service"
	}

	Config.WebUiBannerTitle = "Account Center"

	_, err = Config.DatabaseInstance()
	if err != nil {
		return err
	}

	Config.Refresh()

	return nil
}

func (c *ConfigCommon) Refresh() {

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		return
	}

	q := base.NewQuerySet().From("ids_sysconfig").Limit(100)
	rs, err := dcn.Base.Query(q)
	if err != nil || len(rs) < 1 {
		return
	}

	for _, v := range rs {

		val := v.Field("value").String()
		if val == "" {
			continue
		}

		switch v.Field("key").String() {
		case "service_name":
			c.ServiceName = val
		case "webui_banner_title":
			c.WebUiBannerTitle = val
		case "mailer_smtp_host":
			c.Mailer.SmtpHost = val
		case "mailer_smtp_port":
			c.Mailer.SmtpPort = val
		case "mailer_smtp_user":
			c.Mailer.SmtpUser = val
		case "mailer_smtp_pass":
			c.Mailer.SmtpPass = val
		}
	}

	if c.Mailer.SmtpHost != "" {

		email.MailerRegister("def",
			c.Mailer.SmtpHost,
			c.Mailer.SmtpPort,
			c.Mailer.SmtpUser,
			c.Mailer.SmtpPass)
	}
}
