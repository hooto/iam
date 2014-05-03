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
    "strings"
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
