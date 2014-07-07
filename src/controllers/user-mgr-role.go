package controllers

import (
	"../../deps/lessgo/data/rdo"
	"../../deps/lessgo/data/rdo/base"
	"../../deps/lessgo/utils"
	"io"
	"net/http"
	"strings"
)

func (c UserMgr) RoleListAction() {

	if !c.Session.AccessAllowed("user.admin") {
		c.RenderError(200, "Access Denied")
		return
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		return
	}

	q := base.NewQuerySet().From("ids_role").Limit(1000)
	q.Where.And("status", 1)
	rsr, err := dcn.Base.Query(q)

	ls := []map[string]interface{}{}

	if err == nil && len(rsr) > 0 {
		for _, v := range rsr {
			ls = append(ls, map[string]interface{}{
				"rid":     v.Field("rid").String(),
				"name":    v.Field("name").String(),
				"desc":    v.Field("desc").String(),
				"created": v.Field("created").TimeParse("datetime"),
				"updated": v.Field("updated").TimeParse("datetime"),
			})
		}
		c.ViewData["list"] = ls
	}
}

type UserMgrInstance struct {
	InstanceId string
	AppTitle   string
	Version    string
	Privileges map[string]UserMgrPrivilege
}
type UserMgrPrivilege struct {
	Desc    string
	Checked bool
}

func (c UserMgr) RoleEditAction() {

	if !c.Session.AccessAllowed("user.admin") {
		c.RenderError(200, "Access Denied")
		return
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		c.RenderError(500, http.StatusText(500))
		return
	}

	rolePrivileges := []string{}

	if c.Params.Get("rid") != "" {

		q := base.NewQuerySet().From("ids_role").Limit(1)
		q.Where.And("rid", c.Params.Get("rid"))
		rsrole, err := dcn.Base.Query(q)
		if err != nil || len(rsrole) == 0 {
			c.RenderError(400, http.StatusText(400))
			return
		}

		pls := strings.Split(rsrole[0].Field("privileges").String(), ",")
		for _, v := range pls {
			if v == "" {
				continue
			}
			rolePrivileges = append(rolePrivileges, v)
		}

		c.ViewData["rid"] = c.Params.Get("rid")
		c.ViewData["name"] = rsrole[0].Field("name").String()
		c.ViewData["desc"] = rsrole[0].Field("desc").String()
		c.ViewData["status"] = rsrole[0].Field("status").String()

		c.ViewData["panel_title"] = "Edit Role"
		c.ViewData["rid"] = c.Params.Get("rid")
	} else {
		c.ViewData["panel_title"] = "New Role"
		c.ViewData["rid"] = ""
	}

	instances := map[string]UserMgrInstance{}
	q := base.NewQuerySet().From("ids_instance").Limit(1000)
	rsins, err := dcn.Base.Query(q)
	if err == nil && len(rsins) > 0 {
		for _, v := range rsins {
			instances[v.Field("id").String()] = UserMgrInstance{
				InstanceId: v.Field("id").String(),
				AppTitle:   v.Field("app_title").String(),
				Version:    v.Field("version").String(),
				Privileges: map[string]UserMgrPrivilege{},
			}
		}
	}

	//prePrivileges := map[string]string{}
	q = base.NewQuerySet().From("ids_privilege").Limit(1000)
	rspri, err := dcn.Base.Query(q)
	if err == nil && len(rspri) > 0 {
		for _, v := range rspri {

			if _, ok := instances[v.Field("instance").String()]; !ok {
				continue
			}

			pid := v.Field("pid").String()
			checked := false
			for _, rp := range rolePrivileges {
				if rp == pid {
					checked = true
				}
			}
			instances[v.Field("instance").String()].Privileges[pid] = UserMgrPrivilege{
				Desc:    v.Field("desc").String(),
				Checked: checked,
			}
		}
	}

	//c.ViewData["rolePrivileges"] = rolePrivileges
	c.ViewData["instances"] = instances
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

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Message = "Internal Server Error"
		return
	}

	q := base.NewQuerySet().From("ids_role").Limit(1)

	isNew := true
	roleset := map[string]interface{}{}

	if c.Params.Get("rid") != "" {

		q.Where.And("rid", c.Params.Get("rid"))

		rsrole, err := dcn.Base.Query(q)
		if err != nil || len(rsrole) == 0 {
			rsp.Status = 400
			rsp.Message = http.StatusText(400)
			return
		}

		isNew = false
	}

	roleset["updated"] = base.TimeNow("datetime")
	roleset["name"] = c.Params.Get("name")
	roleset["desc"] = c.Params.Get("desc")
	roleset["privileges"] = strings.Join(c.Params.Values["privileges"], ",")

	if isNew {

		si, err := c.Session.SessionFetch()
		if err != nil {
			return
		}

		roleset["created"] = base.TimeNow("datetime")
		roleset["uid"] = si.Uid
		roleset["status"] = 1

		_, err = dcn.Base.Insert("ids_role", roleset)
		if err != nil {
			rsp.Status = 500
			rsp.Message = "Can not write to database"
			return
		}

	} else {

		frupd := base.NewFilter()
		frupd.And("rid", c.Params.Get("rid"))
		if _, err := dcn.Base.Update("ids_role", roleset, frupd); err != nil {
			rsp.Status = 500
			rsp.Message = "Can not write to database"
			return
		}
	}

	rsp.Status = 200
	rsp.Message = ""
}
