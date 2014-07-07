package controllers

import (
	"../../deps/lessgo/data/rdo"
	"../../deps/lessgo/data/rdo/base"
	"../../deps/lessgo/pagelet"
	"../../deps/lessgo/pass"
	"../../deps/lessgo/utils"
	"../reg/signup"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	userMgrPasswdHidden = "************"
	userMgrPageLimit    = 20
)

var (
	userMgrStatus = map[string]string{
		//0: "Deleted",
		"1": "Active",
		"2": "Banned",
	}
)

type RoleEntry struct {
	Rid, Name, Checked string
}

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

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		return
	}

	rdict := map[string]string{}
	q := base.NewQuerySet().From("ids_role").Limit(100)
	rsr, err := dcn.Base.Query(q)
	if err == nil && len(rsr) > 0 {
		for _, v := range rsr {
			rdict[v.Field("rid").String()] = v.Field("name").String()
		}
	}

	page, err := strconv.Atoi(c.Params.Get("page"))
	if err != nil {
		page = 0
	}

	// filter: query_text
	q = base.NewQuerySet().From("ids_login").Limit(userMgrPageLimit)
	if query_text := c.Params.Get("query_text"); query_text != "" {
		q.Where.And("name.like", "%"+query_text+"%").
			Or("uname.like", "%"+query_text+"%").
			Or("email.like", "%"+query_text+"%")
		c.ViewData["query_text"] = query_text
	}

	count, _ := dcn.Base.Count("ids_login", q.Where)
	pager := pagelet.Pager(page, int(count), userMgrPageLimit, 10)
	c.ViewData["pager"] = pager

	if pager.CurrentPageNumber > 1 {
		q.Offset(int64((pager.CurrentPageNumber - 1) * userMgrPageLimit))
	}
	rsl, err := dcn.Base.Query(q)

	if err == nil && len(rsl) > 0 {

		ls := []map[string]interface{}{}

		for _, v := range rsl {

			rids := strings.Split(v.Field("roles").String(), ",")
			for rk, rv := range rids {

				rname, ok := rdict[rv]
				if !ok {
					continue
				}

				rids[rk] = rname
			}

			status_display := ""
			if vd, ok := userMgrStatus[v.Field("status").String()]; ok {
				status_display = vd
			}

			ls = append(ls, map[string]interface{}{
				"uid":            v.Field("uid").String(),
				"uname":          v.Field("uname").String(),
				"name":           v.Field("name").String(),
				"email":          v.Field("email").String(),
				"timezone":       v.Field("timezone").String(),
				"status_display": status_display,
				"roles_display":  rids,
				"created":        v.Field("created").TimeParse("datetime"),
				"updated":        v.Field("updated").TimeParse("datetime"),
			})
		}

		c.ViewData["list"] = ls
	}

	c.ViewData["query_role"] = c.Params.Get("query_role")
}

func (c UserMgr) EditAction() {

	if !c.Session.AccessAllowed("user.admin") {
		c.RenderError(200, "Access Denied")
		return
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		c.RenderError(500, http.StatusText(500))
		return
	}

	//
	roles := []RoleEntry{}
	q := base.NewQuerySet().From("ids_role").Limit(100)
	q.Where.And("status", 1)
	rsr, err := dcn.Base.Query(q)
	if err == nil && len(rsr) > 0 {

		for _, v := range rsr {
			roles = append(roles, RoleEntry{
				v.Field("rid").String(),
				v.Field("name").String(),
				""})
		}
	}

	if c.Params.Get("uid") != "" {

		q := base.NewQuerySet().From("ids_login").Limit(1)
		q.Where.And("uid", c.Params.Get("uid"))
		rslogin, err := dcn.Base.Query(q)
		if err != nil || len(rslogin) == 0 {
			c.RenderError(400, http.StatusText(400))
			return
		}

		rls := strings.Split(rslogin[0].Field("roles").String(), ",")
		for _, v := range rls {
			for k2, v2 := range roles {
				if v2.Rid == v {
					roles[k2].Checked = "1"
					break
				}
			}
		}

		c.ViewData["uid"] = c.Params.Get("uid")
		c.ViewData["uname"] = rslogin[0].Field("uname").String()
		c.ViewData["email"] = rslogin[0].Field("email").String()
		c.ViewData["passwd"] = userMgrPasswdHidden
		c.ViewData["name"] = rslogin[0].Field("name").String()
		c.ViewData["status"] = rslogin[0].Field("status").String()

		q.From("ids_profile")
		rsprofile, err := dcn.Base.Query(q)
		if err == nil && len(rsprofile) == 1 {
			c.ViewData["birthday"] = rsprofile[0].Field("birthday").String()
			c.ViewData["aboutme"] = rsprofile[0].Field("aboutme").String()
		}

		c.ViewData["panel_title"] = "Edit Account"
		c.ViewData["uid"] = c.Params.Get("uid")
	} else {

		c.ViewData["panel_title"] = "New Account"
		c.ViewData["uid"] = ""
		c.ViewData["status"] = "1"
	}

	c.ViewData["roles"] = roles
	c.ViewData["statuses"] = userMgrStatus
}

func (c UserMgr) SaveAction() {

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

	if err := signup.Validate(c.Params); err != nil {
		rsp.Message = err.Error()
		return
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Message = "Internal Server Error"
		return
	}

	q := base.NewQuerySet().From("ids_login").Limit(1)

	isNew := true
	loginset := map[string]interface{}{}

	if c.Params.Get("uid") != "" {

		q.Where.And("uid", c.Params.Get("uid"))

		rslogin, err := dcn.Base.Query(q)
		if err != nil || len(rslogin) == 0 {
			c.RenderError(400, http.StatusText(400))
			return
		}

		isNew = false
	}

	//
	q = base.NewQuerySet().From("ids_login").Limit(1)
	q.Where.And("email", c.Params.Get("email"))
	rsu, err := dcn.Base.Query(q)
	if err == nil && len(rsu) == 1 {

		if isNew || rsu[0].Field("uid").String() != c.Params.Get("uid") {
			rsp.Message = "The `Email` already exists, please choose another one"
			return
		}

	} else {
		loginset["email"] = c.Params.Get("email")
	}

	//
	q = base.NewQuerySet().From("ids_login").Limit(1)
	q.Where.And("uname", c.Params.Get("uname"))
	rsu, err = dcn.Base.Query(q)
	if err == nil && len(rsu) == 1 {

		if isNew || rsu[0].Field("uid").String() != c.Params.Get("uid") {
			rsp.Message = "The `Username` already exists, please choose another one"
			return
		}

	} else {
		loginset["uname"] = c.Params.Get("uname")
	}

	if c.Params.Get("passwd") != userMgrPasswdHidden {

		pass, err := pass.HashDefault(c.Params.Get("passwd"))
		if err != nil {
			return
		}
		loginset["pass"] = pass
	}

	if isNew {
		loginset["created"] = base.TimeNow("datetime")
		loginset["timezone"] = "UTC"
	}
	loginset["status"] = c.Params.Get("status")
	loginset["updated"] = base.TimeNow("datetime")
	loginset["name"] = c.Params.Get("name")
	loginset["roles"] = strings.Join(c.Params.Values["roles"], ",")

	frupd := base.NewFilter()

	if isNew {
		rst, err := dcn.Base.Insert("ids_login", loginset)
		if err != nil {
			rsp.Status = 500
			rsp.Message = "Can not write to database"
			return
		}

		lastid, err := rst.LastInsertId()
		if err != nil || lastid == 0 {
			rsp.Status = 500
			rsp.Message = "Can not write to database"
			return
		}

		c.Params.Set("uid", fmt.Sprintf("%v", lastid))

	} else {

		frupd.And("uid", c.Params.Get("uid"))
		if _, err := dcn.Base.Update("ids_login", loginset, frupd); err != nil {
			rsp.Status = 500
			rsp.Message = "Can not write to database"
			return
		}
	}

	if _, err := time.Parse("2006-01-02", c.Params.Get("birthday")); err != nil {
		c.Params.Set("birthday", "0000-00-00")
	}

	profile := map[string]interface{}{
		"birthday": c.Params.Get("birthday"),
		"aboutme":  c.Params.Get("aboutme"),
		"updated":  base.TimeNow("datetime"),
	}
	if isNew {
		profile["uid"] = c.Params.Get("uid")
		profile["uuid"] = rsu[0].Field("uuid").String()
		profile["gender"] = 0
		profile["created"] = base.TimeNow("datetime")

		dcn.Base.Insert("ids_profile", profile)
	} else {
		dcn.Base.Update("ids_profile", profile, frupd)
	}

	rsp.Status = 200
	rsp.Message = ""
}
