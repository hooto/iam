package conf

import (
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
    Version       = "0.1.0dev"
    GroupMember   = 100
    GroupSysAdmin = 1
)

var cfg Config

type Config struct {
    ServiceName  string `json:"service_name"`
    Port         int    `json:"port"`
    DomainDef    string `json:"domaindef"`
    MailerHost   string `json:"mailer_host"`
    MailerUser   string `json:"mailer_user"`
    MailerPass   string `json:"mailer_pass"`
    Version      string
    Prefix       string
    KeeperAgent  string
    WebServer    string
    WebPort      string
    WebDaemon    string
    WebConfig    string
    DatabasePath string
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

    file := cfg.Prefix + "/etc/lessids.json"
    if _, err := os.Stat(file); err != nil && os.IsNotExist(err) {
        file = cfg.Prefix + "/etc/lessids.json.dev"
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

    if cfg.DatabasePath == "" {
        cfg.DatabasePath = cfg.Prefix + "/var/lessids.sqlite"
    }

    return cfg, nil
}

func ConfigFetch() *Config {
    return &cfg
}
