package main

import (
    "../deps/lessgo/data/rdc"
    "../deps/lessgo/net/email"
    "../deps/lessgo/pagelet"
    "./conf"
    ctrl_def "./controllers"
    "flag"
    "fmt"
    "log"
    "os"
    "runtime"
    "time"
)

const (
    VERSION string = "0.1.0dev"
)

var cfg conf.Config

var flagPrefix = flag.String("prefix", "", "the prefix folder path")

func main() {

    var err error

    runtime.GOMAXPROCS(runtime.NumCPU())

    //
    flag.Parse()
    if cfg, err = conf.NewConfig(*flagPrefix); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    cfg.Version = VERSION

    if cn, err := cfg.DatabaseInstance(); err == nil {
        rdc.InstanceRegister("def", cn)
    } else {
        log.Fatal("Database Error:", err)
    }

    if cfg.MailerHost != "" && cfg.MailerUser != "" {
        email.MailerRegister("def", cfg.MailerHost, cfg.MailerUser, cfg.MailerPass)
    }

    //
    pagelet.Config.UrlBasePath = "ids"
    pagelet.Config.HttpPort = cfg.Port

    //
    pagelet.Config.ViewPath("default", cfg.Prefix+"/src/views")
    // TODO auto config
    pagelet.Config.I18n(cfg.Prefix + "/src/i18n/en.json")
    pagelet.Config.I18n(cfg.Prefix + "/src/i18n/zh_CN.json")
    //
    pagelet.Config.RouteStaticAppend("default", "/~", cfg.Prefix+"/static")
    pagelet.Config.RouteAppend("default", "/:controller/:action")

    //
    pagelet.RegisterController("default", (*ctrl_def.Index)(nil))
    pagelet.RegisterController("default", (*ctrl_def.Error)(nil))
    pagelet.RegisterController("default", (*ctrl_def.Service)(nil))
    pagelet.RegisterController("default", (*ctrl_def.Reg)(nil))
    pagelet.RegisterController("default", (*ctrl_def.User)(nil))

    //
    pagelet.Run()

    //
    for {
        time.Sleep(3e9)
    }
}
