package conf

import (
	"../../deps/lessgo/data/rdo"
	"../../deps/lessgo/data/rdo/base"
	"../../deps/lessgo/net/email"
	"../../deps/lessgo/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	Version       = "0.2.1dev"
	GroupMember   = 100
	GroupSysAdmin = 1
)

var cfg Config

type ConfigMailer struct {
	SmtpHost string `json:"smtp_host"`
	SmtpPort string `json:"smtp_port"`
	SmtpUser string `json:"smtp_user"`
	SmtpPass string `json:"smtp_pass"`
}

type Config struct {
	ServiceName      string `json:"service_name"`
	Port             int    `json:"port"`
	DomainDef        string `json:"domaindef"`
	Version          string
	Prefix           string
	KeeperAgent      string
	WebServer        string
	WebPort          string
	WebDaemon        string
	WebConfig        string
	DatabasePath     string
	WebUiBannerTitle string
	Mailer           ConfigMailer `json:"mailer"`
	Database         base.Config  `json:"database"`
}

func NewConfig(prefix string) (Config, error) {

	var err error

	cfg.Version = Version

	if prefix == "" {
		prefix, err = filepath.Abs(filepath.Dir(os.Args[0]) + "/..")
		if err != nil {
			prefix = "/opt/lessids"
		}
	}
	reg, _ := regexp.Compile("/+")
	cfg.Prefix = "/" + strings.Trim(reg.ReplaceAllString(prefix, "/"), "/")

	file := cfg.Prefix + "/etc/lessids.conf"
	if _, err := os.Stat(file); err != nil && os.IsNotExist(err) {
		file = cfg.Prefix + "/etc/lessids.conf.dev"
	}
	if _, err := os.Stat(file); err != nil && os.IsNotExist(err) {
		return cfg, errors.New("Error: config file is not exists")
	}

	fp, err := os.Open(file)
	if err != nil {
		return cfg, errors.New(fmt.Sprintf("Error: Can not open (%s)", file))
	}
	defer fp.Close()

	cfgstr, err := ioutil.ReadAll(fp)
	if err != nil {
		return cfg, errors.New(fmt.Sprintf("Error: Can not read (%s)", file))
	}

	if err = json.Unmarshal(cfgstr, &cfg); err != nil {
		return cfg, errors.New(fmt.Sprintf("Error: "+
			"config file invalid. (%s)", err.Error()))
	}

	if cfg.ServiceName == "" {
		cfg.ServiceName = "less Identity"
	}

	cfg.WebUiBannerTitle = "Account Center"

	// if cfg.DatabasePath == "" {
	//     cfg.DatabasePath = cfg.Prefix + "/var/lessids.sqlite"
	// }

	_, err = cfg.DatabaseInstance()
	if err != nil {
		return cfg, err
	}

	cfg.Refresh()

	return cfg, nil
}

func (c *Config) Refresh() {

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		return
	}

	q := base.NewQuerySet().From("ids_sysconfig").Limit(1000)
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
		case "mailer":
			var mailer ConfigMailer
			err := utils.JsonDecode(val, &mailer)

			if err == nil && mailer.SmtpHost != "" {
				c.Mailer = mailer

				email.MailerRegister("def",
					c.Mailer.SmtpHost,
					c.Mailer.SmtpPort,
					c.Mailer.SmtpUser,
					c.Mailer.SmtpPass)

			}
		}
	}
}

func ConfigFetch() *Config {
	return &cfg
}
