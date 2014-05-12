package controllers

import (
    "../../deps/lessgo/data/rdc"
    "../../deps/lessgo/pagelet"
    "../../deps/lessgo/pass"
    "../../deps/lessgo/utils"
    "../conf"
    "../models/login"
    "../models/profile"
    "../models/session"
    "bytes"
    "encoding/base64"
    "fmt"
    "github.com/eryx/imaging"
    "html"
    "image"
    "image/png"
    "io"
    "strings"
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

    //
    menus := []map[string]string{
        {"path": "#user/my", "title": "My Account"},
    }
    if c.Session.AccessAllowed("user.admin") {
        menus = append(menus, map[string]string{
            "path":  "#user-mgr/index",
            "title": "User Manage",
        })
    }
    if c.Session.AccessAllowed("sys.admin") {
        menus = append(menus, map[string]string{
            "path":  "#sys-mgr/index",
            "title": "System Settings",
        })
    }

    c.ViewData["menus"] = menus

    c.ViewData["webui_banner_title"] = conf.ConfigFetch().WebUiBannerTitle
}

func (c User) MyAction() {

    s := session.GetSession(c.Request)
    if s.Uid == 0 {
        c.RenderError(401, "Access Denied")
        return
    }

    dcn, err := rdc.InstancePull("def")
    if err != nil {
        c.RenderError(401, "Access Denied")
        return
    }

    // login
    q := rdc.NewQuerySet().From("ids_login").Limit(1)
    q.Where.And("uid", s.Uid)
    rslogin, err := dcn.Query(q)
    if err != nil || len(rslogin) != 1 {
        c.RenderError(401, "Access Denied")
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
            "created": rdc.TimeNow("datetime"), // TODO
            "updated": rdc.TimeNow("datetime"), // TODO
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
        "updated": rdc.TimeNow("datetime"),
    }
    ft := rdc.NewFilter()
    ft.And("uid", s.Uid)
    dcn.Update("ids_login", itemlogin, ft)

    itemprofile := map[string]interface{}{
        "birthday": c.Params.Get("birthday"),
        "aboutme":  c.Params.Get("aboutme"),
        "updated":  rdc.TimeNow("datetime"), // TODO
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
        return
    }

    body := c.Request.RawBodyString()
    if body == "" {
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

    //
    img64 := strings.SplitAfter(req.Data.Data, ";base64,")
    if len(img64) != 2 {
        return
    }
    imgreader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(img64[1]))
    imgsrc, _, err := image.Decode(imgreader)
    if err != nil {
        rsp.Message = err.Error()
        return
    }
    imgnew := imaging.Thumbnail(imgsrc, 96, 96, imaging.CatmullRom)

    var imgbuf bytes.Buffer
    err = png.Encode(&imgbuf, imgnew)
    if err != nil {
        rsp.Message = err.Error()
        return
    }
    imgphoto := base64.StdEncoding.EncodeToString(imgbuf.Bytes())

    dcn, err := rdc.InstancePull("def")
    if err != nil {
        rsp.Status = 500
        rsp.Message = "Can not pull database instance"
        return
    }

    itemprofile := map[string]interface{}{
        "photo":    "data:image/png;base64," + imgphoto,
        "photosrc": req.Data.Data,
        "updated":  rdc.TimeNow("datetime"),
    }
    ft := rdc.NewFilter()
    ft.And("uid", s.Uid)
    dcn.Update("ids_profile", itemprofile, ft)

    rsp.Status = 200
    rsp.Message = "Successfully changed, Page redirecting"
}

func (c User) PassSetAction() {
}

func (c User) PassPutAction() {

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

    if err := login.PassSetValidate(c.Params); err != nil {
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

    q := rdc.NewQuerySet().From("ids_login").Limit(1)
    q.Where.And("uid", s.Uid)
    rsu, err := dcn.Query(q)
    if err == nil && len(rsu) == 0 {
        rsp.Message = "User can not found"
        return
    }

    if !pass.Check(c.Params.Get("passwd_current"), rsu[0]["pass"].(string)) {
        rsp.Message = "Current Password can not match"
        return
    }

    pass, err := pass.HashDefault(c.Params.Get("passwd"))
    if err != nil {
        return
    }

    itemlogin := map[string]interface{}{
        "pass":    pass,
        "updated": rdc.TimeNow("datetime"),
    }
    ft := rdc.NewFilter()
    ft.And("uid", s.Uid)
    dcn.Update("ids_login", itemlogin, ft)

    rsp.Status = 200
    rsp.Message = "Successfully Updated"
}

func (c User) EmailSetAction() {

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
    rsu, err := dcn.Query(q)
    if err == nil && len(rsu) == 1 {
        c.ViewData["login_email"] = rsu[0]["email"].(string)
    }
}

func (c User) EmailPutAction() {

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

    if err := login.EmailSetValidate(c.Params); err != nil {
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

    q := rdc.NewQuerySet().From("ids_login").Limit(1)
    q.Where.And("uid", s.Uid)
    rsu, err := dcn.Query(q)
    if err == nil && len(rsu) == 0 {
        rsp.Message = "User can not found"
        return
    }

    if !pass.Check(c.Params.Get("passwd"), rsu[0]["pass"].(string)) {
        rsp.Message = "Current Password can not match"
        return
    }

    itemlogin := map[string]interface{}{
        "email":   c.Params.Get("email"),
        "updated": rdc.TimeNow("datetime"),
    }
    ft := rdc.NewFilter()
    ft.And("uid", s.Uid)
    dcn.Update("ids_login", itemlogin, ft)

    rsp.Status = 200
    rsp.Message = "Successfully Updated"
}
