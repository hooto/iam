package controllers

import (
    "../../deps/lessgo/data/rdc"
    "../../deps/lessgo/pagelet"
    "../../deps/lessgo/utils"
    "../models/profile"
    "../models/session"
    "encoding/base64"
    "fmt"
    "html"
    "io"
    "strings"
    "time"
)

type User struct {
    *pagelet.Controller
}

func (c User) IndexAction() {

    s := session.GetSession(c.Request)
    if s.Uid == 0 {
        c.RenderRedirect("/ids/service/login")
        return
    }

    dcn, err := rdc.InstancePull("def")
    if err != nil {
        return
    }

    // login
    q := rdc.NewQuerySet().From("ids_login").Limit(1)
    q.Where.And("uid", s.Uid)
    rslogin, err := dcn.Query(q)
    if err != nil || len(rslogin) != 1 {
        c.RenderRedirect("/ids/service/login")
        return
    }
    c.ViewData["login_uid"] = fmt.Sprintf("%v", rslogin[0]["uid"])
    c.ViewData["login_name"] = rslogin[0]["name"].(string)
    c.ViewData["login_email"] = rslogin[0]["email"].(string)

    //
    q = rdc.NewQuerySet().From("ids_profile").Limit(1)
    q.Where.And("uid", s.Uid)
    rsp, err := dcn.Query(q)
    if err != nil || len(rsp) != 1 {

        item := map[string]interface{}{
            "uid":     s.Uid,
            "gender":  0,
            "created": time.Now().Format("2006-01-02 15:04:05"), // TODO
            "updated": time.Now().Format("2006-01-02 15:04:05"), // TODO
        }
        dcn.Insert("ids_profile", item)
    } else {
        if rslogin[0]["photo"] != nil && len(rslogin[0]["photo"].(string)) > 0 {
            c.ViewData["photo"] = rslogin[0]["photo"].(string)
        }
    }
}

func (c User) ProfileSetAction() {

    s := session.GetSession(c.Request)
    if s.Uid == 0 {
        return
    }

    dcn, err := rdc.InstancePull("def")
    if err != nil {
        return
    }

    q := rdc.NewQuerySet().From("ids_login").Limit(1)
    q.Where.And("uid", s.Uid)
    rslogin, err := dcn.Query(q)
    if err != nil || len(rslogin) != 1 {
        return
    }

    q = rdc.NewQuerySet().From("ids_profile").Limit(1)
    q.Where.And("uid", s.Uid)
    rsp, err := dcn.Query(q)
    if err != nil || len(rsp) != 1 {
        return
    }

    c.ViewData["login_uid"] = s.Uid
    c.ViewData["login_name"] = rslogin[0]["name"].(string)
    if rsp[0]["birthday"] != nil {
        c.ViewData["profile_birthday"] = rsp[0]["birthday"].(string)
    }
    if rsp[0]["aboutme"] != nil {
        c.ViewData["profile_aboutme"] = html.EscapeString(rsp[0]["aboutme"].(string))
    }
}

func (c User) ProfilePutAction() {

    c.AutoRender = false

    var rsp ResponseJson

    rsp.ApiVersion = apiVersion
    rsp.Status = 400
    rsp.Message = "Bad Request"

    defer func() {
        if rspj, err := utils.JsonEncode(rsp); err == nil {
            io.WriteString(c.Response.Out, rspj)
        }
    }()

    if err := profile.PutValidate(c.Params); err != nil {
        rsp.Message = err.Error()
        return
    }

    s := session.GetSession(c.Request)
    if s.Uid == 0 {
        return
    }

    dcn, err := rdc.InstancePull("def")
    if err != nil {
        return
    }

    itemlogin := map[string]interface{}{
        "name":    c.Params.Get("name"),
        "updated": time.Now().Format("2006-01-02 15:04:05"),
    }
    ft := rdc.NewFilter()
    ft.And("uid", s.Uid)
    dcn.Update("ids_login", itemlogin, ft)

    itemprofile := map[string]interface{}{
        "birthday": c.Params.Get("birthday"),
        "aboutme":  c.Params.Get("aboutme"),
        "updated":  time.Now().Format("2006-01-02 15:04:05"), // TODO
    }
    dcn.Update("ids_profile", itemprofile, ft)

    rsp.Status = 200
    rsp.Message = "Successfully Updated"
}

func (c User) PhotoSetAction() {

    s := session.GetSession(c.Request)
    if s.Uid == 0 {
        return
    }

    c.ViewData["login_uid"] = s.Uid
}

func (c User) PhotoPutAction() {

    c.AutoRender = false

    var rsp ResponseJson

    rsp.ApiVersion = apiVersion
    rsp.Status = 400
    rsp.Message = "Bad Request"

    defer func() {
        if rspj, err := utils.JsonEncode(rsp); err == nil {
            io.WriteString(c.Response.Out, rspj)
        }
    }()

    s := session.GetSession(c.Request)
    if s.Uid == 0 {
        rsp.Message = "E02"
        return
    }

    body := c.Request.RawBodyString()
    if body == "" {
        rsp.Message = "E0"
        return
    }

    var req struct {
        //AccessToken string `json:"access_token"`
        Data struct {
            Name string `json:"name"`
            Size int64  `json:"size"`
            Data string `json:"data"`
        } `json:"data"`
    }
    err := utils.JsonDecode(body, &req)
    if err != nil {
        rsp.Message = err.Error()
        return
    }

    body64 := strings.SplitAfter(req.Data.Data, ";base64,")
    if len(body64) != 2 {
        rsp.Message = "E1"
        return
    }
    _, err = base64.StdEncoding.DecodeString(body64[1])
    if err != nil {
        rsp.Message = err.Error()
        rsp.Message = "E2"
        return
    }

    dcn, err := rdc.InstancePull("def")
    if err != nil {
        rsp.Status = 500
        rsp.Message = "Can not pull database instance"
        return
    }

    itemprofile := map[string]interface{}{
        "photo":   req.Data.Data,
        "updated": time.Now().Format("2006-01-02 15:04:05"),
    }
    ft := rdc.NewFilter()
    ft.And("uid", s.Uid)
    dcn.Update("ids_profile", itemprofile, ft)

    rsp.Status = 200
    rsp.Message = "Successfully changed, Page redirecting"
}
