package controllers

import (
    "../../deps/lessgo/data/rdc"
    "../../deps/lessgo/pagelet"
    "../conf"
    "fmt"
)

type SysMgr struct {
    *pagelet.Controller
}

func (c SysMgr) IndexAction() {

    if !c.Session.AccessAllowed("user.admin") {
        c.RenderError(401, "Access Denied")
        return
    }
}

func (c SysMgr) GenSetAction() {

    if !c.Session.AccessAllowed("user.admin") {
        c.RenderError(401, "Access Denied")
        return
    }

    dcn, err := rdc.InstancePull("def")
    if err != nil {
        c.RenderError(401, "Access Denied")
        return
    }

    q := rdc.NewQuerySet().From("ids_sysconfig").Limit(10)
    q.Where.And("key", "service_name").Or("key", "webui_banner_title")
    rs, err := dcn.Query(q)
    if err != nil || len(rs) < 1 {
        return
    }

    for _, v := range rs {
        key := fmt.Sprintf("%v", v["key"])
        val := fmt.Sprintf("%v", v["value"])

        if val == "" {

            switch key {
            case "service_name":
                val = conf.ConfigFetch().ServiceName
            case "webui_banner_title":
                val = conf.ConfigFetch().WebUiBannerTitle
            }
        }

        c.ViewData[key] = val
    }
}

func (c SysMgr) GenSetSaveAction() {

    if !c.Session.AccessAllowed("user.admin") {
        c.RenderError(401, "Access Denied")
        return
    }

    dcn, err := rdc.InstancePull("def")
    if err != nil {
        c.RenderError(401, "Access Denied")
        return
    }

    if conf.ConfigFetch().ServiceName != c.Params.Get("service_name") {
        itemset := map[string]interface{}{
            "value":   c.Params.Get("service_name"),
            "updated": rdc.TimeNow("datetime"),
        }
        ft := rdc.NewFilter()
        ft.And("key", "service_name")
        if _, err := dcn.Update("ids_sysconfig", itemset, ft); err == nil {
            conf.ConfigFetch().ServiceName = c.Params.Get("service_name")
        }
    }

    if conf.ConfigFetch().WebUiBannerTitle != c.Params.Get("webui_banner_title") {
        itemset := map[string]interface{}{
            "value":   c.Params.Get("webui_banner_title"),
            "updated": rdc.TimeNow("datetime"),
        }
        ft := rdc.NewFilter()
        ft.And("key", "webui_banner_title")
        if _, err := dcn.Update("ids_sysconfig", itemset, ft); err == nil {
            conf.ConfigFetch().WebUiBannerTitle = c.Params.Get("webui_banner_title")
        }
    }

    c.RenderError(200, "Successfully updated")
}
