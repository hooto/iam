package controllers

import (
    "../../deps/lessgo/data/rdc"
    "../../deps/lessgo/pagelet"
    "../../deps/lessgo/pass"
    "../../deps/lessgo/utils"
    "io"
    "strings"
    "time"
)

type Service struct {
    *pagelet.Controller
}

func (c Service) IndexAction() {

}

func (c Service) LoginAction() {
    c.ViewData["continue"] = c.Params.Get("continue")
    if c.Params.Get("persistent") == "1" {
        c.ViewData["persistentChecked"] = "checked"
    }
}

func (c Service) LoginAuthAction() {

    c.AutoRender = false

    var rsp struct {
        ResponseJson
        Data struct {
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

    dcn, err := rdc.InstancePull("def")
    if err != nil {
        rsp.Message = "Internal Server Error"
        return
    }

    if c.Params.Get("email") == "" || c.Params.Get("passwd") == "" {
        return
    }

    email := strings.ToLower(c.Params.Get("email"))

    q := rdc.NewQuerySet().From("ids_login").Limit(1)
    q.Where.And("email", email)
    rsu, err := dcn.Query(q)
    if err == nil && len(rsu) == 0 {
        rsp.Message = "Email or Password can not match 1"
        return
    }

    if !pass.Check(c.Params.Get("passwd"), rsu[0]["pass"].(string)) {
        rsp.Message = "Email or Password can not match"
        return
    }

    rsp.Data.AccessToken = utils.StringNewRand36(24)

    addr := "0.0.0.0"
    if addridx := strings.Index(c.Request.RemoteAddr, ":"); addridx > 0 {
        addr = c.Request.RemoteAddr[:addridx]
    }

    session := map[string]interface{}{
        "token":   rsp.Data.AccessToken,
        "status":  1,
        "uid":     rsu[0]["uid"],
        "uname":   rsu[0]["uname"],
        "source":  addr,
        "created": time.Now().Format("2006-01-02 15:04:05"), // TODO
        "timeout": 10 * 24 * 3600,
    }
    if err := dcn.Insert("ids_sessions", session); err != nil {
        rsp.Status = 500
        rsp.Message = "Can not write to database"
        return
    }

    if len(c.Params.Get("continue")) > 0 {
        rsp.Data.Continue = c.Params.Get("continue")
    }
    if strings.Index(rsp.Data.Continue, "?") == -1 {
        rsp.Data.Continue += "?"
    }
    rsp.Data.Continue += "&access_token=" + rsp.Data.AccessToken

    rsp.Status = 200
    rsp.Message = ""
}
