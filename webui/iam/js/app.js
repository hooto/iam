var iamApp = {
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

iamApp.Index = function() {
    iam.TplCmd("app/index", {
        callback: function(err, data) {
            iam.OpToolsClean();
            $("#com-content").html(data);
            iamApp.InstList();
        },
    });
}

iamApp.InstList = function() {
    var alert_id = "#iam-app-insts-alert";

    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            $("#work-content").html(tpl);

            if (!data || !data.items) {
                return l4i.InnerAlert(alert_id, 'alert-info', "<strong>No authorized applications</strong><br>You have no applications authorized to access your account.");
            }

            data._statusls = iamUserMgr.statusls;

            for (var i in data.items) {

                if (!data.items[i].privileges) {
                    data.items[i].privileges = []
                }

                if (!data.items[i].version) {
                    data.items[i].version = "";
                }

                if (!data.items[i].app_title) {
                    data.items[i].app_title = "";
                }

                data.items[i]._privilegeNumber = data.items[i].privileges.length;
            }

            l4iTemplate.Render({
                dstid: "iam-app-insts",
                tplid: "iam-app-insts-tpl",
                data: data,
            // success : function() {
            //     l4iTemplate.Render({
            //         dstid  : "iam-app-insts-pager",
            //         tplid  : "iam-app-insts-pager-tpl",
            //         data   : l4i.Pager(data.meta),
            //     });
            // },
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        iam.ApiCmd("app/inst-list", {
            callback: ep.done('data'),
        });

        iam.TplCmd("app/inst-list", {
            callback: ep.done('tpl'),
        });
    });
}

iamApp.InstSetForm = function(instid) {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'roles', 'data', function(tpl, roles, data) {

            if (!data || !data.kind) {
                return;
            }

            if (!data.privileges) {
                data.privileges = [];
            }
            data._privilegeNumber = data.privileges.length;

            data._statusls = iamApp.statusls;
            data._roles = l4i.Clone(roles);

            for (var i in data.privileges) {
                if (!data.privileges[i].extroles) {
                    data.privileges[i].extroles = [];
                }
            }

            if (!data.url) {
                data.url = "";
            }

            if (!data.app_title) {
                data.app_title = "";
            }

            l4iModal.Open({
                tplsrc: tpl,
                width: 900,
                height: 600,
                data: data,
                title: "App Instance Setting",
                buttons: [{
                    title: "Cancel",
                    onclick: "l4iModal.Close()",
                }, {
                    title: "Save",
                    onclick: "iamApp.InstSetCommit()",
                    style: "btn btn-primary",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        if (iamApp.roles) {
            ep.emit("roles", iamApp.roles)
        } else {
            iam.ApiCmd("user/role-list?status=0", {
                callback: function(err, roles) {
                    iamApp.roles = roles;
                    ep.emit("roles", roles);
                },
            });
        }

        if (instid) {

            iam.ApiCmd("app/inst-entry?instid=" + instid, {
                callback: ep.done('data'),
            });

        } else {
            ep.emit("data", l4i.Clone(iamApp.appinstdef));
        }

        iam.TplCmd("app/inst-set", {
            callback: ep.done('tpl'),
        });
    });
}

iamApp.InstSetCommit = function() {
    var form = $("#iam-app-instset"),
        alert_id = "#iam-app-instset-alert",
        req = l4i.Clone(iamApp.appinstdef);

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

    iam.ApiCmd("app/inst-set", {
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
                iamApp.InstList();
            }, 1000);
        },
    });
}
