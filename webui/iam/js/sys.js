var iamSys = {

}

iamSys.Index = function() {

    l4iTemplate.Render({
        dstid: "com-content",
        tplurl: iam.TplPath("sys-mgr/index"),
        data: {},
        callback: function() {
            iam.OpToolsClean();
            l4i.UrlEventClean("iam-module-navbar-menus");
            l4i.UrlEventRegister("sys-mgr/general-set", iamSys.GeneralSet, "iam-module-navbar-menus");
            l4i.UrlEventRegister("sys-mgr/mailer-set", iamSys.MailerSet, "iam-module-navbar-menus");
            l4i.UrlEventHandler("sys-mgr/general-set", true);
        },
    });
}

iamSys.GeneralSet = function() {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            if (!data || !data.items) {
                return;
            }

            data._items = {};
            for (var i in data.items) {
                if (!data.items[i].value) {
                    data._items[data.items[i]["name"]] = "";
                } else {
                    data._items[data.items[i]["name"]] = data.items[i]["value"];
                }
            }

            l4iTemplate.Render({
                dstid: "work-content",
                tplsrc: tpl,
                data: data,
            });
        });

        ep.fail(function(err) {
            alert(l4i.T("network error, please try again later"));
        });

        iam.ApiCmd("sys-config/general", {
            callback: ep.done('data'),
        });

        iam.TplCmd("sys-mgr/general-set", {
            callback: ep.done('tpl'),
        });
    });
}

iamSys.GeneralSetCommit = function() {
    var form = $("#iam-sysmgr-generalset"),
        alert_id = "#iam-sysmgr-generalset-alert";

    var user_reg_disable = "0";
    if (form.find("input[name=user_reg_disable]").is(":checked")) {
        user_reg_disable = "1";
    }

    var req = {
        items: [{
            name: "service_name",
            value: form.find("input[name=service_name]").val(),
        }, {
            name: "webui_banner_title",
            value: form.find("input[name=webui_banner_title]").val(),
        }, {
            name: "user_reg_disable",
            value: user_reg_disable,
        }, {
            name: "service_login_form_alert_msg",
            value: form.find("input[name=service_login_form_alert_msg]").val(),
        }],
    };

    // console.log(req);

    iam.ApiCmd("sys-config/general-set", {
        method: "PUT",
        data: JSON.stringify(req),
        callback: function(err, data) {

            if (err) {
                return l4i.InnerAlert(alert_id, 'alert-danger', err);
            }

            if (!data || data.error) {
                return l4i.InnerAlert(alert_id, 'alert-danger', data.error.message);
            }

            l4i.InnerAlert(alert_id, 'alert-success', l4i.T("Successfully %s", l4i.T("updated")));
        },
    });
}


iamSys.MailerSet = function(name) {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            if (!data || !data.items) {
                return;
            }

            var mailer = {}
            if (data.items[0].value && data.items[0].value.length > 2) {
                mailer = JSON.parse(data.items[0].value);
                if (!mailer) {
                    mailer = {};
                }
            }

            if (!mailer.smtp_host) {
                mailer.smtp_host = "";
            }

            if (!mailer.smtp_port) {
                mailer.smtp_port = "";
            }

            if (!mailer.smtp_user) {
                mailer.smtp_user = "";
            }

            if (!mailer.smtp_pass) {
                mailer.smtp_pass = "";
            }

            l4iTemplate.Render({
                dstid: "work-content",
                tplsrc: tpl,
                data: mailer,
            });
        });

        ep.fail(function(err) {
            alert(l4i.T("network error, please try again later"));
        });

        iam.ApiCmd("sys-config/mailer", {
            callback: ep.done('data'),
        });

        iam.TplCmd("sys-mgr/mailer-set", {
            callback: ep.done('tpl'),
        });
    });
}


iamSys.MailerSetCommit = function() {

    var form = $("#iam-sysmgr-mailerset"),
        alert_id = "#iam-sysmgr-mailerset-alert";
    if (!form) {
        return;
    }

    var mailer = {
        "smtp_host": form.find("input[name=mailer_smtp_host]").val(),
        "smtp_port": form.find("input[name=mailer_smtp_port]").val(),
        "smtp_user": form.find("input[name=mailer_smtp_user]").val(),
        "smtp_pass": form.find("input[name=mailer_smtp_pass]").val(),
    };

    var req = {
        items: [{
            name: "mailer",
            value: JSON.stringify(mailer),
        }],
    };

    iam.ApiCmd("sys-config/general-set", {
        method: "PUT",
        data: JSON.stringify(req),
        callback: function(err, data) {

            if (err) {
                return l4i.InnerAlert(alert_id, 'alert-danger', err);
            }

            if (!data || data.error) {
                return l4i.InnerAlert(alert_id, 'alert-danger', data.error.message);
            }

            l4i.InnerAlert(alert_id, 'alert-success', l4i.T("Successfully %s", l4i.T("updated")));
        },
    });
}

