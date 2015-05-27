package ctrl

import (
	"net/http"

	"github.com/lessos/lessgo/data/rdo"
	"github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/httpsrv"
)

type Service struct {
	*httpsrv.Controller
}

func (c Service) IndexAction() {

}

func (c Service) LoginAction() {

	c.Data["continue"] = c.Params.Get("continue")
	if c.Params.Get("persistent") == "1" {
		c.Data["persistentChecked"] = "checked"
	}
}

func (c Service) SignOutAction() {

	c.Data["continue"] = "/ids"

	token := c.Session.AccessToken
	if c.Params.Get("access_token") != "" {
		token = c.Params.Get("access_token")
	}

	if len(c.Params.Get("continue")) > 0 {
		c.Data["continue"] = c.Params.Get("continue")
	}

	dcn, err := rdo.ClientPull("def")
	if err == nil {
		ft := base.NewFilter()
		ft.And("token", token)
		if _, err := dcn.Base.Delete("ids_sessions", ft); err != nil {
			//
		}
	}

	ck := &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	}
	http.SetCookie(c.Response.Out, ck)
}
