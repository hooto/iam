package controllers

import (
    "../../deps/lessgo/data/rdc"
    "../../deps/lessgo/pagelet"
    "../../deps/lessgo/pass"
    "../../deps/lessgo/utils"
    "../models/login"
    "../reg/signup"
    "io"
    "time"
)

type Reg struct {
    *pagelet.Controller
}

func (c Reg) IndexAction() {

}

func (c Reg) SignUpAction() {
    c.ViewData["continue"] = c.Params.Get("continue")
}

func (c Reg) SignUpRegAction() {

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

    uname := utils.StringNewRand36(8)
    pass, err := pass.HashDefault(c.Params.Get("passwd"))
    if err != nil {
        return
    }

    item := map[string]interface{}{
        "uname":   uname,
        "email":   c.Params.Get("email"),
        "pass":    pass,
        "name":    c.Params.Get("name"),
        "status":  1,
        "created": time.Now().Format("2006-01-02 15:04:05"), // TODO
        "updated": time.Now().Format("2006-01-02 15:04:05"), // TODO
    }
    if err := dcn.Insert("ids_login", item); err != nil {
        rsp.Status = 500
        rsp.Message = "Can not write to database"
        return
    }

    rsp.Status = 200
    rsp.Message = ""
}

func (c Reg) ForgotPassAction() {

}

func (c Reg) ForgotPassPutAction() {

    c.AutoRender = false

    var rsp struct {
        ResponseJson
        Data struct {
            Continue string `json:"continue"`
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

    if err := login.EmailSetValidate(c.Params); err != nil {
        rsp.Message = err.Error()
        return
    }

    dcn, err := rdc.InstancePull("def")
    if err != nil {
        rsp.Message = "Internal Server Error"
        return
    }

    id := utils.StringNewRand36(24)
    taf, _ := time.ParseDuration("+3600s")
    item := map[string]interface{}{
        "id":      id,
        "status": 0,
        "email":   c.Params.Get("email"),                             // TODO
        "expired": time.Now().Add(taf).Format("2006-01-02 15:04:05"), // TODO
    }
    if err := dcn.Insert("ids_resetpass", item); err != nil {
        rsp.Status = 500
        rsp.Message = "Can not write to database"
        return
    }

    rsp.Status = 200
    rsp.Message = ""
}
