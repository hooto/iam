package controllers

import (
    "../../deps/lessgo/data/rdc"
    "../../deps/lessgo/pagelet"
    //"../../deps/lessgo/pass"
    //"../../deps/lessgo/utils"
    //"../models/login"
    //"../models/profile"
    //"../models/session"
    //"encoding/base64"
    "fmt"
    //"html"
    //"io"
    "../reg/signup"
    "strings"

    "../../deps/lessgo/pass"
    "../../deps/lessgo/utils"
    "io"
)

type UserMgr struct {
    *pagelet.Controller
}

func (c UserMgr) IndexAction() {

    if !c.Session.AccessAllowed("user.admin") {
        c.RenderError(401, "Access Denied")
        return
    }
}

func (c UserMgr) ListAction() {

    if !c.Session.AccessAllowed("user.admin") {
        c.RenderError(200, "Access Denied")
        return
    }

    dcn, err := rdc.InstancePull("def")
    if err != nil {
        return
    }

    rdict := map[string]string{}
    q := rdc.NewQuerySet().From("ids_role").Limit(100)
    rsr, err := dcn.Query(q)
    if err == nil && len(rsr) > 0 {
        //c.ViewData["roles"] = rsr
        for _, v := range rsr {
            rdict[fmt.Sprintf("%v", v["rid"])] = v["name"].(string)
        }
    }

    // filter: query_text
    q = rdc.NewQuerySet().From("ids_login").Limit(20)
    if query_text := c.Params.Get("query_text"); query_text != "" {
        q.Where.And("name.like", "%"+query_text+"%").
            Or("uname.like", "%"+query_text+"%").
            Or("email.like", "%"+query_text+"%")
        c.ViewData["query_text"] = query_text
    }

    rsl, err := dcn.Query(q)
    if err == nil && len(rsl) > 0 {

        for k, v := range rsl {

            rids := strings.Split(v["roles"].(string), ",")
            for rk, rv := range rids {

                rname, ok := rdict[rv]
                if !ok {
                    continue
                }

                rids[rk] = rname
            }

            rsl[k]["roles_display"] = rids
        }

        c.ViewData["list"] = rsl
    }

    c.ViewData["query_role"] = c.Params.Get("query_role")
}

func (c UserMgr) NewAction() {

    if !c.Session.AccessAllowed("user.admin") {
        c.RenderError(200, "Access Denied")
        return
    }

    dcn, err := rdc.InstancePull("def")
    if err != nil {
        return
    }

    q := rdc.NewQuerySet().From("ids_role").Limit(100)
    q.Where.And("status", 1)
    rsr, err := dcn.Query(q)
    if err == nil && len(rsr) > 0 {
        c.ViewData["roles"] = rsr
    }
}

func (c UserMgr) NewSaveAction() {

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

    if err := signup.Validate(c.Params); err != nil {
        rsp.Message = err.Error()
        return
    }

    dcn, err := rdc.InstancePull("def")
    if err != nil {
        rsp.Message = "Internal Server Error"
        return
    }

    q := rdc.NewQuerySet().From("ids_login").Limit(1)
    q.Where.And("email", c.Params.Get("email"))
    rsu, err := dcn.Query(q)
    if err == nil && len(rsu) == 1 {
        rsp.Message = "The `Email` already exists, please choose another one"
        return
    }

    q = rdc.NewQuerySet().From("ids_login").Limit(1)
    q.Where.And("uname", c.Params.Get("uname"))
    rsu, err = dcn.Query(q)
    if err == nil && len(rsu) == 1 {
        rsp.Message = "The `Username` already exists, please choose another one"
        return
    }

    pass, err := pass.HashDefault(c.Params.Get("passwd"))
    if err != nil {
        return
    }

    item := map[string]interface{}{
        "uname":    c.Params.Get("uname"),
        "email":    c.Params.Get("email"),
        "pass":     pass,
        "name":     c.Params.Get("name"),
        "status":   1,
        "roles":    strings.Join(c.Params.Values["roles"], ","),
        "timezone": "UTC",                   // TODO
        "created":  rdc.TimeNow("datetime"), // TODO
        "updated":  rdc.TimeNow("datetime"), // TODO
    }
    //fmt.Println(item)
    rst, err := dcn.Insert("ids_login", item)
    if err != nil {
        rsp.Status = 500
        rsp.Message = "Can not write to database"
        return
    }
    lastid, err := rst.LastInsertId()
    if lastid > 0 {
        profile := map[string]interface{}{
            "uid":      lastid,
            "gender":   0,
            "birthday": c.Params.Get("birthday"),
            "aboutme":  c.Params.Get("aboutme"),
            "created":  rdc.TimeNow("datetime"),
            "updated":  rdc.TimeNow("datetime"),
        }
        dcn.Insert("ids_profile", profile)
    }

    rsp.Status = 200
    rsp.Message = ""
}
