package main

import (
    "../deps/lessgo/pagelet"
    "./conf"
    ctrl_def "./controllers"
    "flag"
    "fmt"
    "os"
    "runtime"
    "time"
)

var (
    err        error
    cfg        conf.Config
    flagPrefix = flag.String("prefix", "", "the prefix folder path")
)

func main() {

    runtime.GOMAXPROCS(runtime.NumCPU())

    //
    flag.Parse()
    if cfg, err = conf.NewConfig(*flagPrefix); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    //
    pagelet.Config.InstanceId = "lessids"
    pagelet.Config.UrlBasePath = "ids"
    pagelet.Config.HttpPort = cfg.Port
    pagelet.Config.LessIdsServiceUrl = fmt.Sprintf("http://127.0.0.1:%v/ids", cfg.Port)

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
    pagelet.RegisterController("default", (*ctrl_def.UserMgr)(nil))
    pagelet.RegisterController("default", (*ctrl_def.Status)(nil))
    pagelet.RegisterController("default", (*ctrl_def.SysMgr)(nil))
    pagelet.RegisterController("default", (*ctrl_def.AppAuth)(nil))
    pagelet.RegisterController("default", (*ctrl_def.AppMgr)(nil))

    //
    pagelet.Run()

    //
    for {
        time.Sleep(3e9)
    }
}
