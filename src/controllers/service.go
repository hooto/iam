package controllers

import (
    "../../deps/lessgo/pagelet"
    "../../deps/lessgo/pass"
    "../../deps/lessgo/utils"
    "../../deps/lessgo/data/rdc"
    "strings"
    "io"
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
            Continue string `json:"continue"`
            AccessToken string `json:"access_token"`
        } `json:"data"`
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

    if c.Params.Get("userid") == "" || c.Params.Get("passwd") == "" {
        return
    }

    userid := strings.ToLower(c.Params.Get("userid"))

    q := rdc.NewQuerySet().From("ids_login").Limit(1)
    q.Where.And("uid", userid)
    rsu, err := dcn.Query(q)
    if err == nil && len(rsu) == 0 {
        rsp.Message = "User ID or Password can not match"
        return
    }
    
    if pass.Check(c.Params.Get("passwd"), rsu[0]["pass"].(string)) {
        rsp.Message = "User ID or Password can not match"
        return
    }

    rsp.Data.AccessToken = utils.StringNewRand(32)

    //session.Set(rsp.AccessToken, userid, "7200")

    rsp.Status = 200
    rsp.Message = ""
}