var iamAppMgr = {
    statusls: [{
        status: 1,
        title: "Active",
    }, {
        status: 2,
        title: "Banned",
    }],
    roles: null,
    appinstdef: {
        meta: {
            id: "",
        },
        app_id: "",
        app_title: "",
        version: "",
        status: 1,
        url: "",
    },
}

iamAppMgr.Index = function() {
    iam.TplCmd("app-mgr/index", {
        callback: function(err, data) {
            iam.OpToolsClean();
            $("#com-content").html(data);
            l4i.UrlEventClean("iam-module-navbar-menus");
            l4i.UrlEventRegister("app-mgr/inst-list", iamAppMgr.InstList);
            l4i.UrlEventHandler("app-mgr/inst-list", true);
        },
    });
}

iamAppMgr.InstList = function() {
    var uri = "",
        alertId = "#iam-appmgr-instls-alert";
    if (document.getElementById("iam_appmgr_instls_qry_text")) {
        var qt = $("#iam_appmgr_instls_qry_text").val();
        if (qt && qt.length > 0) {
            uri = "?qry_text=" + qt;
        }
    }

    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            $("#work-content").html(tpl);

            if (!data || !data.items) {
                return l4i.InnerAlert(alertId, 'info', "No Items Yet ...");
            }

            data._statusls = iamUserMgr.statusls;

            for (var i in data.items) {

                if (!data.items[i].privileges) {
                    data.items[i].privileges = [];
                }

                if (!data.items[i].version) {
                    data.items[i].version = "";
                }

                if (!data.items[i].app_title) {
                    data.items[i].app_title = "";
                }

                data.items[i]._privilegeNumber = data.items[i].privileges.length;
            }

            iam.OpToolsRefresh("#iam-appmgr-instls-optools");

            l4iTemplate.Render({
                dstid: "iam-appmgr-instls",
                tplid: "iam-appmgr-instls-tpl",
                data: data,
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        iam.ApiCmd("app-mgr/inst-list" + uri, {
            callback: ep.done('data'),
        });

        iam.TplCmd("app-mgr/inst-list", {
            callback: ep.done('tpl'),
        });
    });
}

iamAppMgr.InstSetForm = function(instid) {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'roles', 'data', function(tpl, roles, data) {

            if (!data || !data.kind) {
                return;
            }

            if (!data.privileges) {
                data.privileges = [];
            }
            data._privilegeNumber = data.privileges.length;

            data._statusls = iamAppMgr.statusls;
            data._roles = l4i.Clone(roles);

            for (var i in data.privileges) {
                if (!data.privileges[i].extroles) {
                    data.privileges[i].extroles = [];
                }
            }

            if (!data.url) {
                data.url = "";
            }

            l4iModal.Open({
                tplsrc: tpl,
                width: 1000,
                height: 600,
                data: data,
                title: "App Instance Setting",
                buttons: [{
                    title: "Delete",
                    onclick: "iamAppMgr.InstDelCommit()",
                    style: "btn btn-outline-danger",
                }, {
                    title: "Cancel",
                    onclick: "l4iModal.Close()",
                }, {
                    title: "Save",
                    onclick: "iamAppMgr.InstSetCommit()",
                    style: "btn btn-primary",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        if (iamAppMgr.roles) {
            ep.emit("roles", iamAppMgr.roles)
        } else {
            iam.ApiCmd("user-mgr/role-list?status=0", {
                callback: function(err, roles) {
                    iamAppMgr.roles = roles;
                    ep.emit("roles", roles);
                },
            });
        }

        if (instid) {

            iam.ApiCmd("app-mgr/inst-entry?instid=" + instid, {
                callback: ep.done('data'),
            });

        } else {
            ep.emit("data", l4i.Clone(iamAppMgr.appinstdef));
        }

        iam.TplCmd("app-mgr/inst-set", {
            callback: ep.done('tpl'),
        });
    });
}

iamAppMgr.InstSetCommit = function() {
    var form = $("#iam-appmgr-instset"),
        alert_id = "#iam-appmgr-instset-alert",
        req = l4i.Clone(iamAppMgr.appinstdef);

    try {
        req.meta.id = form.find("input[name=instid]").val();
        req.app_title = form.find("input[name=app_title]").val();
        req.url = form.find("input[name=url]").val();

        req.status = parseInt(form.find("input[name=status]:checked").val());

        // form.find("input[name=roles]:checked").each(function() {

        //     var val = parseInt($(this).val());
        //     if (val > 0) {
        //         req.roles.push(val);
        //     }
        // });

    } catch (err) {
        return l4i.InnerAlert(alert_id, 'alert-danger', err);
    }

    iam.ApiCmd("app-mgr/inst-set", {
        method: "PUT",
        data: JSON.stringify(req),
        callback: function(err, data) {

            if (err) {
                return l4i.InnerAlert(alert_id, 'alert-danger', err);
            }

            if (!data || data.error) {
                return l4i.InnerAlert(alert_id, 'alert-danger', data.error.message);
            }

            l4i.InnerAlert(alert_id, 'alert-success', "Successfully updated");

            window.setTimeout(function() {
                l4iModal.Close();
                iamAppMgr.InstList();
            }, 1000);
        },
    });
}


iamAppMgr.InstDelCommit = function() {
    var form = $("#iam-appmgr-instset"),
        alert_id = "#iam-appmgr-instset-alert";

    if (!form) {
        return;
    }

    iam.ApiCmd("app-mgr/inst-del?inst_id=" + form.find("input[name=instid]").val(), {
        callback: function(err, data) {

            if (err) {
                return l4i.InnerAlert(alert_id, 'alert-danger', err);
            }

            if (!data || data.error) {
                return l4i.InnerAlert(alert_id, 'alert-danger', data.error.message);
            }

            l4i.InnerAlert(alert_id, 'alert-success', "Successfully updated");

            window.setTimeout(function() {
                l4iModal.Close();
                iamAppMgr.InstList();
            }, 500);
        },
    });
}


