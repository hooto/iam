package controllers

import (
    "../../deps/lessgo/data/rdo"
    "../../deps/lessgo/data/rdo/base"
    "../../deps/lessgo/pagelet"
    "../../deps/lessgo/pass"
    "../../deps/lessgo/utils"
    "io"
    "strings"
)

type AppAuth struct {
    *pagelet.Controller
}

func (c AppAuth) RegisterAction() {

    c.AutoRender = false

    var rsp struct {
        ResponseJson
        Data struct {
            InstanceId string `json:"instance_id"`
            AppId      string `json:"app_id"`
            AppTitle   string `json:"app_title"`
            Version    string `json:"version"`
            Privileges []struct {
                Key  string `json:"key"`
                Desc string `json:"desc"`
            }   `json:"privileges"`
            Continue    string `json:"continue"`
            AccessToken string `json:"access_token"`
        }   `json:"data"`
    }
    rsp.ApiVersion = apiVersion
    rsp.Status = 400
    rsp.Message = "Bad Request"

    rsp.Data.Continue = "/ids"

    defer func() {
        if rspj, err := utils.JsonEncode(rsp); err == nil {
            io.WriteString(c.Response.Out, rspj)
        }
    }()

    dcn, err := rdo.ClientPull("def")
    if err != nil {
        rsp.Message = "Internal Server Error"
        return
    }

    if c.Params.Get("email") == "" || c.Params.Get("passwd") == "" {
        return
    }

    email := strings.ToLower(c.Params.Get("email"))

    q := base.NewQuerySet().From("ids_login").Limit(1)
    q.Where.And("email", email)
    rsu, err := dcn.Base.Query(q)
    if err == nil && len(rsu) == 0 {
        rsp.Message = "Email or Password can not match"
        return
    }

    if !pass.Check(c.Params.Get("passwd"), rsu[0].Field("pass").String()) {
        rsp.Message = "Email or Password can not match"
        return
    }

    rsp.Data.AccessToken = utils.StringNewRand36(24)

    addr := "127.0.0.1"
    if addridx := strings.Index(c.Request.RemoteAddr, ":"); addridx > 0 {
        addr = c.Request.RemoteAddr[:addridx]
    }
    //fmt.Println(c.Request.RemoteAddr, addr, c.Request.Request)

    session := map[string]interface{}{
        "token":    rsp.Data.AccessToken,
        "refresh":  utils.StringNewRand36(24),
        "status":   1,
        "uid":      rsu[0].Field("uid").Int(),
        "uuid":     rsu[0].Field("uuid").String(),
        "uname":    rsu[0].Field("uname").String(),
        "name":     rsu[0].Field("name").String(),
        "roles":    rsu[0].Field("roles").String(),
        "timezone": rsu[0].Field("timezone").String(),
        "source":   addr,
        "created":  base.TimeNow("datetime"),                // TODO
        "expired":  base.TimeNowAdd("datetime", "+864000s"), // TODO
    }
    if _, err := dcn.Base.Insert("ids_sessions", session); err != nil {
        rsp.Status = 500
        rsp.Message = "Can not write to database" + err.Error()
        return
    }

    if len(c.Params.Get("continue")) > 0 {
        rsp.Data.Continue = c.Params.Get("continue")
        if strings.Index(rsp.Data.Continue, "?") == -1 {
            rsp.Data.Continue += "?"
        } else {
            rsp.Data.Continue += "&"
        }
        rsp.Data.Continue += "access_token=" + rsp.Data.AccessToken
    }

    rsp.Status = 200
    rsp.Message = ""
}
