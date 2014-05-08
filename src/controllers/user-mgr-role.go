package controllers

import (
    "../../deps/lessgo/data/rdc"
    "../../deps/lessgo/utils"
    "fmt"
    "io"
    "net/http"
)

func (c UserMgr) RoleListAction() {

    if !c.Session.AccessAllowed("user.admin") {
        c.RenderError(200, "Access Denied")
        return
    }

    dcn, err := rdc.InstancePull("def")
    if err != nil {
        return
    }

    q := rdc.NewQuerySet().From("ids_role").Limit(1000)
    q.Where.And("status", 1)
    rsr, err := dcn.Query(q)
    if err == nil && len(rsr) > 0 {
        c.ViewData["list"] = rsr
    }
}

func (c UserMgr) RoleEditAction() {

    if !c.Session.AccessAllowed("user.admin") {
        c.RenderError(200, "Access Denied")
        return
    }

    dcn, err := rdc.InstancePull("def")
    if err != nil {
        c.RenderError(500, http.StatusText(500))
        return
    }

    if c.Params.Get("rid") != "" {

        q := rdc.NewQuerySet().From("ids_role").Limit(1)
        q.Where.And("rid", c.Params.Get("rid"))
        rsrole, err := dcn.Query(q)
        if err != nil || len(rsrole) == 0 {
            c.RenderError(400, http.StatusText(400))
            return
        }

        /* pls := strings.Split(rsrole[0]["privileges"].(string), ",")
           for _, v := range pls {
               for k2, v2 := range roles {
                   if v2.Rid == v {
                       roles[k2].Checked = "1"
                       break
                   }
               }
           } */

        c.ViewData["rid"] = c.Params.Get("rid")
        c.ViewData["name"] = rsrole[0]["name"]
        c.ViewData["desc"] = rsrole[0]["desc"]
        c.ViewData["status"] = fmt.Sprintf("%v", rsrole[0]["status"])

        c.ViewData["panel_title"] = "Edit Role"
        c.ViewData["rid"] = c.Params.Get("rid")
    } else {
        c.ViewData["panel_title"] = "New Role"
        c.ViewData["rid"] = ""
    }
}

func (c UserMgr) RoleSaveAction() {

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

    if !c.Session.AccessAllowed("user.admin") {
        return
    }

    dcn, err := rdc.InstancePull("def")
    if err != nil {
        rsp.Message = "Internal Server Error"
        return
    }

    q := rdc.NewQuerySet().From("ids_role").Limit(1)

    isNew := true
    roleset := map[string]interface{}{}

    if c.Params.Get("rid") != "" {

        q.Where.And("rid", c.Params.Get("rid"))

        rsrole, err := dcn.Query(q)
        if err != nil || len(rsrole) == 0 {
            rsp.Status = 400
            rsp.Message = http.StatusText(400)
            return
        }

        isNew = false
    }

    roleset["updated"] = rdc.TimeNow("datetime")
    roleset["name"] = c.Params.Get("name")
    roleset["desc"] = c.Params.Get("desc")
    roleset["privileges"] = ""
    //roleset["roles"] = strings.Join(c.Params.Values["roles"], ",")

    if isNew {

        si, err := c.Session.SessionFetch()
        if err != nil {
            return
        }

        roleset["created"] = rdc.TimeNow("datetime")
        roleset["uid"] = si.Uid
        roleset["status"] = 1

        _, err = dcn.Insert("ids_role", roleset)
        if err != nil {
            rsp.Status = 500
            rsp.Message = "Can not write to database"
            return
        }

    } else {

        frupd := rdc.NewFilter()
        frupd.And("rid", c.Params.Get("rid"))
        if _, err := dcn.Update("ids_role", roleset, frupd); err != nil {
            rsp.Status = 500
            rsp.Message = "Can not write to database"
            return
        }
    }

    rsp.Status = 200
    rsp.Message = ""
}
