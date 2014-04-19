package controllers

import (
    "../../deps/lessgo/data/rdc"
    "../../deps/lessgo/net/email"
    "../../deps/lessgo/pagelet"
    "../../deps/lessgo/pass"
    "../../deps/lessgo/utils"
    "../conf"
    "../models/login"
    "../reg/signup"
    "fmt"
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
    if _, err := dcn.Insert("ids_login", item); err != nil {
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

    q := rdc.NewQuerySet().From("ids_login").Limit(1)
    q.Where.And("email", c.Params.Get("email"))
    rsl, err := dcn.Query(q)
    if err != nil || len(rsl) != 1 {
        rsp.Message = "Email can not found"
        return
    }

    id := utils.StringNewRand36(24)
    taf, _ := time.ParseDuration("+3600s")
    item := map[string]interface{}{
        "id":      id,
        "status":  0,
        "email":   c.Params.Get("email"),                             // TODO
        "expired": time.Now().Add(taf).Format("2006-01-02 15:04:05"), // TODO
    }
    if _, err := dcn.Insert("ids_resetpass", item); err != nil {
        rsp.Status = 500
        rsp.Message = "Can not write to database"
        return
    }

    mr, err := email.MailerPull("def")
    if err != nil {
        rsp.Message = "Internal Server Error"
        return
    }

    cfg := conf.ConfigFetch()

    // TODO tempate
    body := fmt.Sprintf(`<html>
<body>
<div>You recently requested a password reset for your %s account. To create a new password, click on the link below:</div>
<br>
<a href="http://%s/ids/reg/pass-reset?id=%s">Reset My Password</a>
<br>
<div>This request was made on %s.</div>
<br>
<div>Regards,</div>
<div>%s Account Services</div>

<div>********************************************************</div>
<div>Please do not reply to this message. Mail sent to this address cannot be answered.</div>
</body>
</html>`, cfg.ServiceName, c.Request.Host, id, time.Now().Format("2006-01-02 15:04:05"), cfg.ServiceName)

    err = mr.SendMail(c.Params.Get("email"), c.T("Reset your password"), body)

    rsp.Status = 200
    rsp.Message = ""
}

func (c Reg) PassResetAction() {

    if c.Params.Get("id") == "" {
        return
    }

    dcn, err := rdc.InstancePull("def")
    if err != nil {
        return
    }

    q := rdc.NewQuerySet().From("ids_resetpass").Limit(1)
    q.Where.And("id", c.Params.Get("id"))
    rsr, err := dcn.Query(q)
    if err != nil || len(rsr) != 1 {
        return
    }

    if rsr[0]["expired"].(time.Time).Before(time.Now()) {
        return
    }

    c.ViewData["pass_reset_id"] = c.Params.Get("id")
}

func (c Reg) PassResetPutAction() {

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

    if c.Params.Get("id") == "" {
        rsp.Message = "Token can not be null"
        return
    }

    if err := login.PassSetValidate(c.Params); err != nil {
        rsp.Message = err.Error()
        return
    }

    dcn, err := rdc.InstancePull("def")
    if err != nil {
        rsp.Message = "Internal Server Error"
        return
    }

    q := rdc.NewQuerySet().From("ids_resetpass").Limit(1)
    q.Where.And("id", c.Params.Get("id"))
    rsr, err := dcn.Query(q)
    if err != nil || len(rsr) != 1 {
        rsp.Message = "Token not found"
        return
    }

    if rsr[0]["expired"].(time.Time).Before(time.Now()) {
        rsp.Message = "Token expired"
        return
    }

    if rsr[0]["email"].(string) != c.Params.Get("email") {
        rsp.Message = "Email or Birthday is not valid"
        return
    }

    q = rdc.NewQuerySet().From("ids_login").Limit(1)
    q.Where.And("email", c.Params.Get("email"))
    rsl, err := dcn.Query(q)
    if err != nil || len(rsl) != 1 {
        rsp.Message = "User can not found"
        return
    }

    q = rdc.NewQuerySet().From("ids_profile").Limit(1)
    q.Where.And("uid", rsl[0]["uid"])
    rspf, err := dcn.Query(q)
    if err != nil || len(rspf) != 1 {
        rsp.Message = "User can not found"
        return
    }
    if fmt.Sprintf("%v", rspf[0]["birthday"]) != c.Params.Get("birthday") {
        rsp.Message = "Email or Birthday is not valid"
        return
    }

    pass, err := pass.HashDefault(c.Params.Get("passwd"))
    if err != nil {
        rsp.Message = "Internal Server Error"
        return
    }

    itemlogin := map[string]interface{}{
        "pass":    pass,
        "updated": time.Now().Format("2006-01-02 15:04:05"),
    }
    ft := rdc.NewFilter()
    ft.And("uid", rsl[0]["uid"])
    dcn.Update("ids_login", itemlogin, ft)

    //
    delfr := rdc.NewFilter()
    delfr.And("id", c.Params.Get("id"))
    dcn.Delete("ids_resetpass", delfr)

    rsp.Status = 200
    rsp.Message = "Successfully Updated"
}
