package controllers

import (
    "../../deps/lessgo/data/rdc"
    "../../deps/lessgo/utils"
    "fmt"
    "io"
    "net/http"
)

var (
    userMgrAuthStatus = map[string]string{
        //0: "Deleted",
        "1": "Active",
        "2": "Banned",
    }
)

func (c UserMgr) AuthListAction() {

    if !c.Session.AccessAllowed("user.admin") {
        c.RenderError(200, "Access Denied")
        return
    }

    dcn, err := rdc.InstancePull("def")
    if err != nil {
        return
    }

    users := []interface{}{}

    q := rdc.NewQuerySet().From("ids_instance").Limit(1000)
    rsinst, err := dcn.Query(q)
    if err == nil && len(rsinst) > 0 {

        for k, v := range rsinst {

            if vd, ok := userMgrAuthStatus[fmt.Sprintf("%v", v["status"])]; ok {
                rsinst[k]["status_display"] = vd
            }

            uid := fmt.Sprintf("%v", v["uid"])
            rsinst[k]["uid"] = uid
            rsinst[k]["uid_name"] = uid

            inArray := false
            for _, vuid := range users {
                if vuid == uid {
                    inArray = true
                    break
                }
            }

            if !inArray {
                users = append(users, uid)
            }
        }

    } else {
        return
    }

    //
    q = rdc.NewQuerySet().From("ids_login").Limit(1000)
    q.Where.And("uid.in", users...)
    rslogin, err := dcn.Query(q)
    if err == nil && len(rslogin) > 0 {
        for _, v := range rslogin {

            for k2, v2 := range rsinst {
                if v2["uid"] == fmt.Sprintf("%v", v["uid"]) {
                    rsinst[k2]["uid_name"] = fmt.Sprintf("%v", v["name"])
                }
            }
        }
    }

    c.ViewData["list"] = rsinst
}

func (c UserMgr) AuthEditAction() {

    if !c.Session.AccessAllowed("user.admin") {
        c.RenderError(200, "Access Denied")
        return
    }

    dcn, err := rdc.InstancePull("def")
    if err != nil {
        c.RenderError(500, http.StatusText(500))
        return
    }

    if c.Params.Get("instid") != "" {

        q := rdc.NewQuerySet().From("ids_instance").Limit(1)
        q.Where.And("id", c.Params.Get("instid"))
        rsinst, err := dcn.Query(q)
        if err != nil || len(rsinst) == 0 {
            c.RenderError(400, http.StatusText(400))
            return
        }

        if vd, ok := userMgrAuthStatus[fmt.Sprintf("%v", rsinst[0]["status"])]; ok {
            c.ViewData["status_display"] = vd
        }

        c.ViewData["version"] = rsinst[0]["version"]
        c.ViewData["app_id"] = rsinst[0]["app_id"]
        c.ViewData["app_title"] = rsinst[0]["app_title"]
        c.ViewData["status"] = fmt.Sprintf("%v", rsinst[0]["status"])

        c.ViewData["panel_title"] = "Edit Authorization"
        c.ViewData["id"] = c.Params.Get("instid")

    } else {
        c.ViewData["panel_title"] = "New Authorization"
        c.ViewData["id"] = ""
        c.ViewData["status"] = "1"
    }

    q := rdc.NewQuerySet().From("ids_privilege").Limit(1000)
    q.Where.And("instance", c.Params.Get("instid"))
    rspri, err := dcn.Query(q)
    if err == nil && len(rspri) > 0 {
        c.ViewData["privileges"] = rspri
    }

    c.ViewData["statuses"] = userMgrAuthStatus
}

func (c UserMgr) AuthSaveAction() {

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

    q := rdc.NewQuerySet().From("ids_instance").Limit(1)

    isNew := true
    instset := map[string]interface{}{}

    if c.Params.Get("instid") != "" {

        q.Where.And("id", c.Params.Get("instid"))

        rsinst, err := dcn.Query(q)
        if err != nil || len(rsinst) == 0 {
            rsp.Status = 400
            rsp.Message = http.StatusText(400)
            return
        }

        isNew = false
    }

    instset["updated"] = rdc.TimeNow("datetime")
    instset["app_title"] = c.Params.Get("app_title")

    if isNew {

        // TODO

    } else {

        instset["status"] = c.Params.Get("status")

        frupd := rdc.NewFilter()
        frupd.And("id", c.Params.Get("instid"))
        if _, err := dcn.Update("ids_instance", instset, frupd); err != nil {
            rsp.Status = 500
            rsp.Message = "Can not write to database"
            return
        }
    }

    rsp.Status = 200
    rsp.Message = ""
}
