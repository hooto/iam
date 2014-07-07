package controllers

import (
	"../../deps/lessgo/data/rdo"
	"../../deps/lessgo/data/rdo/base"
	"../../deps/lessgo/pagelet"
)

type AppMgr struct {
	*pagelet.Controller
}

func (c AppMgr) IndexAction() {

	if !c.Session.AccessAllowed("sys.admin") {
		c.RenderError(401, "Access Denied")
		return
	}
}

func (c AppMgr) ListAction() {

	if !c.Session.AccessAllowed("sys.admin") {
		c.RenderError(200, "Access Denied")
		return
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		return
	}

	q := base.NewQuerySet().From("ids_instance").Limit(1000)
	rsl, err := dcn.Base.Query(q)

	if err == nil && len(rsl) > 0 {

		ls := []map[string]interface{}{}

		for _, v := range rsl {

			ls = append(ls, map[string]interface{}{
				"id":        v.Field("id").String(),
				"uid":       v.Field("uid").Uint(),
				"status":    v.Field("status").Uint(),
				"app_id":    v.Field("app_id").String(),
				"app_title": v.Field("app_title").String(),
				"version":   v.Field("version").String(),
				"created":   v.Field("created").TimeParse("datetime"),
				"updated":   v.Field("updated").TimeParse("datetime"),
			})
		}

		c.ViewData["list"] = ls
	}

}
