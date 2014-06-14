package controllers

import (
    "../../deps/lessgo/data/rdo"
    "../../deps/lessgo/data/rdo/base"
    "../../deps/lessgo/net/email"
    "../../deps/lessgo/utils"
    "../conf"
)

func (c SysMgr) EmailSetAction() {

    if !c.Session.AccessAllowed("sys.admin") {
        c.RenderError(200, "Access Denied")
        return
    }

    dcn, err := rdo.ClientPull("def")
    if err != nil {
        c.RenderError(401, "Access Denied")
        return
    }

    q := base.NewQuerySet().From("ids_sysconfig").Limit(10)
    q.Where.And("key", "mailer")
    rs, err := dcn.Base.Query(q)
    if err != nil || len(rs) < 1 {
        return
    }

    var mailer conf.ConfigMailer
    err = utils.JsonDecode(rs[0].Field("value").String(), &mailer)
    if err != nil {
        return
    }

    c.ViewData["mailer"] = mailer
}

func (c SysMgr) EmailSetSaveAction() {

    if !c.Session.AccessAllowed("sys.admin") {
        c.RenderError(401, "Access Denied")
        return
    }

    dcn, err := rdo.ClientPull("def")
    if err != nil {
        c.RenderError(401, "Access Denied")
        return
    }

    isNew := true
    q := base.NewQuerySet().From("ids_sysconfig").Limit(1)
    q.Where.And("key", "mailer")
    rs, err := dcn.Base.Query(q)
    if err == nil && len(rs) == 1 {
        isNew = false
    }

    mailer := conf.ConfigMailer{
        SmtpHost: c.Params.Get("mailer_smtp_host"),
        SmtpPort: c.Params.Get("mailer_smtp_port"),
        SmtpUser: c.Params.Get("mailer_smtp_user"),
        SmtpPass: c.Params.Get("mailer_smtp_pass"),
    }

    preMailer := conf.ConfigFetch().Mailer
    if mailer.SmtpHost != preMailer.SmtpHost ||
        mailer.SmtpPort != preMailer.SmtpPort ||
        mailer.SmtpUser != preMailer.SmtpUser ||
        mailer.SmtpPass != preMailer.SmtpPass {

        val, err := utils.JsonEncode(mailer)
        if err != nil {
            c.RenderError(500, "InternalServerError")
            return
        }
        itemset := map[string]interface{}{
            "value":   val,
            "updated": base.TimeNow("datetime"),
        }
        if isNew {
            itemset["key"] = "mailer"
            _, err = dcn.Base.Insert("ids_sysconfig", itemset)
        } else {
            ft := base.NewFilter()
            ft.And("key", "mailer")
            _, err = dcn.Base.Update("ids_sysconfig", itemset, ft)
        }

        if err != nil {
            c.RenderError(500, "InternalServerError")
            return
        }

        conf.ConfigFetch().Mailer = mailer
        email.MailerRegister("def",
            mailer.SmtpHost,
            mailer.SmtpPort,
            mailer.SmtpUser,
            mailer.SmtpPass)
    }

    c.RenderError(200, "Successfully updated")
}
