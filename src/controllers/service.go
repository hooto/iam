package controllers

import (
    "../../deps/lessgo/data/rdc"
    "../../deps/lessgo/pagelet"
    "../../deps/lessgo/pass"
    "../../deps/lessgo/utils"
    "../models/role"
    "encoding/base64"
    "fmt"
    "io"
    "strconv"
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
        rsp.Message = "Email or Password can not match"
        return
    }

    if !pass.Check(c.Params.Get("passwd"), rsu[0]["pass"].(string)) {
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
        "uid":      rsu[0]["uid"],
        "uname":    rsu[0]["uname"],
        "name":     rsu[0]["name"],
        "roles":    rsu[0]["roles"],
        "timezone": rsu[0]["timezone"],
        "source":   addr,
        "created":  rdc.TimeNow("datetime"),                // TODO
        "expired":  rdc.TimeNowAdd("datetime", "+864000s"), // TODO
    }
    if _, err := dcn.Insert("ids_sessions", session); err != nil {
        rsp.Status = 500
        rsp.Message = "Can not write to database"
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

func (c Service) AuthAction() {

    c.AutoRender = false

    var rsp struct {
        ResponseJson
        Data ResponseSession `json:"data"`
    }
    rsp.ApiVersion = apiVersion
    rsp.Status = 401
    rsp.Message = "Unauthorized"

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

    if c.Params.Get("access_token") == "" {
        return
    }

    q := rdc.NewQuerySet().From("ids_sessions").Limit(1)
    q.Where.And("token", c.Params.Get("access_token"))
    rss, err := dcn.Query(q)
    if err == nil && len(rss) == 0 {
        return
    }

    //
    expired := rss[0]["expired"].(time.Time)
    if expired.Before(time.Now()) {
        return
    }
    //fmt.Println("expired", expired)

    //
    addr := "0.0.0.0"
    if addridx := strings.Index(c.Request.RemoteAddr, ":"); addridx > 0 {
        addr = c.Request.RemoteAddr[:addridx]
    }
    if addr != rss[0]["source"].(string) {
        return
    }

    //
    uid, _ := strconv.Atoi(fmt.Sprintf("%v", rss[0]["uid"]))
    rsp.Data.Uid = uint32(uid)
    rsp.Data.Uname = rss[0]["uname"].(string)
    rsp.Data.Name = rss[0]["name"].(string)
    rsp.Data.Roles = rss[0]["roles"].(string)
    rsp.Data.AccessToken = rss[0]["token"].(string)
    rsp.Data.RefreshToken = rss[0]["refresh"].(string)
    rsp.Data.Timezone = rss[0]["timezone"].(string)
    rsp.Data.Expired = rdc.TimeZoneFormat(expired, rsp.Data.Timezone, "atom")

    rsp.Status = 200
    rsp.Message = ""
}

func (c Service) AccessAllowedAction() {

    c.AutoRender = false

    var rsp struct {
        ResponseJson
    }
    rsp.ApiVersion = apiVersion
    rsp.Status = 401
    rsp.Message = "Unauthorized"

    defer func() {
        if rspj, err := utils.JsonEncode(rsp); err == nil {
            io.WriteString(c.Response.Out, rspj)
        }
    }()

    body := c.Request.RawBodyString()
    if body == "" {
        return
    }

    var req struct {
        AccessToken string `json:"access_token"`
        Data        struct {
            InstanceId string `json:"instanceid"`
            Privilege  string `json:"privilege"`
        }   `json:"data"`
    }
    err := utils.JsonDecode(body, &req)
    if err != nil {
        rsp.Message = err.Error()
        return
    }
    if req.AccessToken == "" {
        return
    }

    dcn, err := rdc.InstancePull("def")
    if err != nil {
        rsp.Message = "Internal Server Error"
        return
    }

    q := rdc.NewQuerySet().From("ids_sessions").Limit(1)
    q.Where.And("token", req.AccessToken)
    rss, err := dcn.Query(q)
    if err == nil && len(rss) == 0 {
        return
    }

    //
    expired := rss[0]["expired"].(time.Time)
    if expired.Before(time.Now()) {
        return
    }

    //
    addr := "0.0.0.0"
    if addridx := strings.Index(c.Request.RemoteAddr, ":"); addridx > 0 {
        addr = c.Request.RemoteAddr[:addridx]
    }

    if addr != rss[0]["source"].(string) {
        return
    }

    if !role.AccessAllowed(rss[0]["roles"].(string), req.Data.InstanceId, req.Data.Privilege) {
        return
    }

    rsp.Status = 200
    rsp.Message = ""
}

func (c Service) PhotoAction() {

    c.AutoRender = false

    photo := "data:image/png;base64,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"

    uid := strings.TrimLeft(c.Request.RequestPath, "service/photo/")
    dcn, err := rdc.InstancePull("def")
    if err != nil {
        return
    }

    q := rdc.NewQuerySet().From("ids_profile").Select("photo").Limit(1)
    q.Where.And("uid", uid)
    rsp, err := dcn.Query(q)
    if err == nil && len(rsp) == 1 {
        if rsp[0]["photo"] != nil && len(rsp[0]["photo"].(string)) > 50 {
            photo = rsp[0]["photo"].(string)
        }
    }

    body64 := strings.SplitAfter(photo, ";base64,")
    if len(body64) != 2 {
        return
    }
    data, err := base64.StdEncoding.DecodeString(body64[1])
    if err != nil {
        return
    }

    io.WriteString(c.Response.Out, string(data))
}
