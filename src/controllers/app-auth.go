package controllers

import (
	"../../deps/lessgo/data/rdo"
	"../../deps/lessgo/data/rdo/base"
	"../../deps/lessgo/pagelet"
	"../../deps/lessgo/utils"
	"io"
)

type AppAuth struct {
	*pagelet.Controller
}

func (c AppAuth) InfoAction() {

	c.AutoRender = false

	var rsp struct {
		ResponseJson
		Data struct {
			InstanceId string `json:"instance_id"`
			AppId      string `json:"app_id"`
			Version    string `json:"version"`
		} `json:"data"`
	}
	rsp.ApiVersion = apiVersion
	rsp.Status = 400
	rsp.Message = "Bad Request"

	defer func() {
		if rspj, err := utils.JsonEncode(rsp); err == nil {
			io.WriteString(c.Response.Out, rspj)
		}
	}()

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Status = 500
		rsp.Message = "Internal Server Error"
		return
	}

	instanceid := c.Params.Get("instanceid")

	q := base.NewQuerySet().From("ids_instance").Limit(1)
	q.Where.And("id", instanceid)
	rs, err := dcn.Base.Query(q)
	if err == nil && len(rs) == 0 {
		rsp.Status = 404
		rsp.Message = "Instance Not Found"
		return
	}

	rsp.Data.InstanceId = rs[0].Field("id").String()
	rsp.Data.AppId = rs[0].Field("app_id").String()
	rsp.Data.Version = rs[0].Field("version").String()

	rsp.Status = 200
	rsp.Message = ""
}

func (c AppAuth) RegisterAction() {

	c.AutoRender = false

	var rsp struct {
		ResponseJson
		Data struct {
			InstanceId string `json:"instance_id"`
			AppId      string `json:"app_id"`
			AppTitle   string `json:"app_title"`
			Version    string `json:"version"`
			Url        string `json:"url"`
			Privileges []struct {
				Key  string `json:"key"`
				Desc string `json:"desc"`
			} `json:"privileges"`
			Continue    string `json:"continue"`
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}
	rsp.ApiVersion = apiVersion
	rsp.Status = 400
	rsp.Message = "Bad Request"

	defer func() {
		if rspj, err := utils.JsonEncode(rsp); err == nil {
			io.WriteString(c.Response.Out, rspj)
		}
	}()

	body := c.Request.RawBodyString()
	if body == "" {
		return
	}

	var req struct {
		AccessToken string `json:"access_token"`
		Data        struct {
			InstanceId  string `json:"instance_id"`
			InstanceUrl string `json:"instance_url"`
			AppId       string `json:"app_id"`
			AppTitle    string `json:"app_title"`
			Version     string `json:"version"`
			Privileges  []struct {
				Key  string `json:"key"`
				Desc string `json:"desc"`
			} `json:"privileges"`
		} `json:"data"`
	}

	err := utils.JsonDecode(body, &req)
	if err != nil {
		rsp.Message = err.Error()
		return
	}
	if req.AccessToken == "" {
		return
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Message = "Internal Server Error"
		return
	}

	if !c.Session.AccessAllowed("sys.admin") {
		rsp.Status = 401
		rsp.Message = "Unauthorized"
		return
	}

	sess, err := c.Session.SessionFetch()

	q := base.NewQuerySet().From("ids_instance").Limit(1)
	q.Where.And("id", req.Data.InstanceId)
	rs, err := dcn.Base.Query(q)
	if err != nil {
		rsp.Message = "Internal Server Error"
		return
	}
	if len(rs) == 0 {

		item := map[string]interface{}{
			"id":        req.Data.InstanceId,
			"uid":       sess.Uid,
			"status":    1,
			"app_id":    req.Data.AppId,
			"app_title": req.Data.AppTitle,
			"version":   req.Data.Version,
			"url":       req.Data.InstanceUrl,
			"created":   base.TimeNow("datetime"),
			"updated":   base.TimeNow("datetime"),
		}

		if _, err := dcn.Base.Insert("ids_instance", item); err != nil {
			rsp.Status = 500
			rsp.Message = "Can not write to database" + err.Error()
			return
		}

	} else {

		item := map[string]interface{}{
			"app_title": req.Data.AppTitle,
			"version":   req.Data.Version,
			"url":       req.Data.InstanceUrl,
			"updated":   base.TimeNow("datetime"),
		}
		frupd := base.NewFilter()
		frupd.And("id", req.Data.InstanceId)
		if _, err := dcn.Base.Update("ids_instance", item, frupd); err != nil {
			rsp.Status = 500
			rsp.Message = "Can not write to database" + err.Error()
			return
		}
	}

	//
	q = base.NewQuerySet().From("ids_privilege").Limit(1000)
	q.Where.And("instance", req.Data.InstanceId)
	rs, err = dcn.Base.Query(q)
	if err != nil {
		rsp.Message = "Internal Server Error"
		return
	}

	for _, prePriv := range rs {

		isExist := false
		for _, curPrev := range req.Data.Privileges {

			if prePriv.Field("privilege").String() == curPrev.Key {
				isExist = true
				break
			}
		}

		if !isExist {
			frupd := base.NewFilter()
			frupd.And("instance", req.Data.InstanceId).And("privilege", prePriv.Field("privilege").String())
			dcn.Base.Delete("ids_privilege", frupd)
		}
	}

	for _, curPrev := range req.Data.Privileges {

		isExist := false

		for _, prePriv := range rs {

			if prePriv.Field("privilege").String() == curPrev.Key {
				isExist = true
				break
			}
		}

		if !isExist {
			item := map[string]interface{}{
				"instance":  req.Data.InstanceId,
				"uid":       sess.Uid,
				"privilege": curPrev.Key,
				"desc":      curPrev.Desc,
				"created":   base.TimeNow("datetime"),
			}

			if _, err := dcn.Base.Insert("ids_privilege", item); err != nil {
				rsp.Status = 500
				rsp.Message = "Can not write to database" + err.Error()
				return
			}
		}
	}

	rsp.Status = 200
	rsp.Message = ""
}
