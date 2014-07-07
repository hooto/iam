package controllers

import (
	"../../deps/lessgo/data/rdo"
	"../../deps/lessgo/data/rdo/base"
	"../../deps/lessgo/utils"
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

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		return
	}

	users := []interface{}{}

	q := base.NewQuerySet().From("ids_instance").Limit(1000)
	rsinst, err := dcn.Base.Query(q)
	if err != nil || len(rsinst) == 0 {
		return
	}

	ls := []map[string]interface{}{}

	for _, v := range rsinst {

		item := map[string]interface{}{
			"id":         v.Field("id").String(),
			"app_title":  v.Field("app_title").String(),
			"app_id":     v.Field("app_id").String(),
			"version":    v.Field("version").String(),
			"uid":        v.Field("uid").String(),
			"uid_name":   v.Field("uid").String(),
			"privileges": v.Field("privileges").String(),
			"created":    v.Field("created").TimeParse("datetime"),
			"updated":    v.Field("updated").TimeParse("datetime"),
		}

		if vd, ok := userMgrAuthStatus[v.Field("status").String()]; ok {
			item["status_display"] = vd
		}

		uid := v.Field("uid").String()

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

		ls = append(ls, item)
	}

	//
	q = base.NewQuerySet().From("ids_login").Limit(1000)
	q.Where.And("uid.in", users...)
	rslogin, err := dcn.Base.Query(q)
	if err == nil && len(rslogin) > 0 {
		for _, v := range rslogin {

			for k2, v2 := range ls {
				if v2["uid"] == v.Field("uid").String() {
					ls[k2]["uid_name"] = v.Field("name").String()
				}
			}
		}
	}

	c.ViewData["list"] = ls
}

func (c UserMgr) AuthEditAction() {

	if !c.Session.AccessAllowed("user.admin") {
		c.RenderError(200, "Access Denied")
		return
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		c.RenderError(500, http.StatusText(500))
		return
	}

	if c.Params.Get("instid") != "" {

		q := base.NewQuerySet().From("ids_instance").Limit(1)
		q.Where.And("id", c.Params.Get("instid"))
		rsinst, err := dcn.Base.Query(q)
		if err != nil || len(rsinst) == 0 {
			c.RenderError(400, http.StatusText(400))
			return
		}

		if vd, ok := userMgrAuthStatus[rsinst[0].Field("status").String()]; ok {
			c.ViewData["status_display"] = vd
		}

		c.ViewData["version"] = rsinst[0].Field("version").String()
		c.ViewData["app_id"] = rsinst[0].Field("app_id").String()
		c.ViewData["app_title"] = rsinst[0].Field("app_title").String()
		c.ViewData["status"] = rsinst[0].Field("status").String()

		c.ViewData["panel_title"] = "Edit Authorization"
		c.ViewData["id"] = c.Params.Get("instid")

	} else {
		c.ViewData["panel_title"] = "New Authorization"
		c.ViewData["id"] = ""
		c.ViewData["status"] = "1"
	}

	q := base.NewQuerySet().From("ids_privilege").Limit(1000)
	q.Where.And("instance", c.Params.Get("instid"))
	rspri, err := dcn.Base.Query(q)
	if err == nil && len(rspri) > 0 {
		ls := []map[string]interface{}{}
		for _, v := range rspri {
			ls = append(ls, map[string]interface{}{
				"pid":       v.Field("pid").String(),
				"privilege": v.Field("privilege").String(),
				"desc":      v.Field("desc").String(),
			})
		}
		c.ViewData["privileges"] = ls
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

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Message = "Internal Server Error"
		return
	}

	q := base.NewQuerySet().From("ids_instance").Limit(1)

	isNew := true
	instset := map[string]interface{}{}

	if c.Params.Get("instid") != "" {

		q.Where.And("id", c.Params.Get("instid"))

		rsinst, err := dcn.Base.Query(q)
		if err != nil || len(rsinst) == 0 {
			rsp.Status = 400
			rsp.Message = http.StatusText(400)
			return
		}

		isNew = false
	}

	instset["updated"] = base.TimeNow("datetime")
	instset["app_title"] = c.Params.Get("app_title")

	if isNew {

		// TODO

	} else {

		instset["status"] = c.Params.Get("status")

		frupd := base.NewFilter()
		frupd.And("id", c.Params.Get("instid"))
		if _, err := dcn.Base.Update("ids_instance", instset, frupd); err != nil {
			rsp.Status = 500
			rsp.Message = "Can not write to database"
			return
		}
	}

	rsp.Status = 200
	rsp.Message = ""
}
