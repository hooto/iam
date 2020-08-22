var iamAccessKey = {
    aks: null,
    statuses: [{
        status: 1 << 1,
        title: "Active",
    }, {
        status: 0,
        title: "Disable",
    }],
    def: {
        kind: "AccessKey",
        id: "",
        secret: "",
        status: (1 << 1),
        scopes: [],
        description: "",
    },
}

iamAccessKey.Index = function() {
    l4iTemplate.Render({
        dstid: "com-content",
        tplurl: iam.TplPath("access-key/index"),
        // data: {},
        callback: function() {
            iam.OpToolsClean();
            l4i.UrlEventClean("iam-module-navbar-menus");
            l4i.UrlEventRegister("access-key/list", iamAccessKey.List, "iam-module-navbar-menus");
            l4i.UrlEventHandler("access-key/list", true);
        },
    });
}

iamAccessKey.List = function() {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'aks', function(tpl, aks) {

            if (!aks) {
                return;
            }

            if (!aks.items) {
                aks.items = [];
            }

            // $("#work-content").html(tpl);


            for (var i in aks.items) {

                if (!aks.items[i].description) {
                    aks.items[i].description = "";
                }

                if (!aks.items[i].status) {
                    aks.items[i].status = 0;
                }

                if (!aks.items[i].scopes) {
                    aks.items[i].scopes = [];
                }
            }

            aks._statuses = iamAccessKey.statuses;

            l4iTemplate.Render({
                dstid: "work-content",
                tplsrc: tpl, // "iam-aklist-tpl",
                data: aks,
                callback: function() {
                    iam.OpToolsRefresh("#iam-aklist-optools");
                },
            });
        });

        ep.fail(function(err) {
            alert(l4i.T("network error, please try again later"));
        });

        iam.ApiCmd("access-key/list?status=0", {
            callback: ep.done('aks'),
        });

        iam.TplCmd("access-key/list", {
            callback: ep.done('tpl'),
        });
    });
}

iamAccessKey.Info = function(akid) {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            if (!data || !data.kind) {
                return;
            }

            if (!data.item) {
                data.item = {};
            }

            if (!data.item.description) {
                data.item.description = "";
            }
            if (!data.item.id) {
                data.item.id = "";
            }
            if (!data.item.status) {
                data.item.status = 0;
            }
            if (!data.item.scopes) {
                data.item.scopes = [];
            }

            data.item._statuses = iamAccessKey.statuses;

            l4iModal.Open({
                tplsrc: tpl,
                width: 800,
                height: 350,
                data: data.item,
                title: "Access Key Info",
                buttons: [{
                    title: l4i.T("Cancel"),
                    onclick: "l4iModal.Close()",
                }],
            });
        });

        ep.fail(function(err) {
            alert(l4i.T("network error, please try again later"));
        });

        if (akid) {
            iam.ApiCmd("access-key/entry?access_key_id=" + akid, {
                callback: ep.done('data'),
            });
        } else {
            ep.emit("data", l4i.Clone(iamAccessKey.def));
        }

        iam.TplCmd("access-key/info", {
            callback: ep.done('tpl'),
        });
    });
}


iamAccessKey.Set = function(akid) {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            if (!data || !data.kind) {
                return;
            }

            if (data._isnew) {
                data._form_title = l4i.T("New %s", "Access Key");
            } else {
                data._form_title = l4i.T("%s Settings", "Access Key");
            }


            if (!data.item) {
                data.item = {};
            }

            if (!data.item.description) {
                data.item.description = "";
            }
            if (!data.item.id) {
                data.item.id = "";
            }
            if (!data.item.status) {
                data.item.status = 0;
            }
            if (!data.item.scopes) {
                data.item.scopes = [];
            }

            data.item._statuses = iamAccessKey.statuses;

            l4iModal.Open({
                tplsrc: tpl,
                width: 800,
                height: 300,
                data: data.item,
                title: data._form_title,
                buttons: [{
                    title: l4i.T("Cancel"),
                    onclick: "l4iModal.Close()",
                }, {
                    title: l4i.T("Save"),
                    onclick: "iamAccessKey.SetCommit()",
                    style: "btn btn-primary",
                }],
            });
        });

        ep.fail(function(err) {
            alert(l4i.T("network error, please try again later"));
        });

        if (akid) {
            iam.ApiCmd("access-key/entry?access_key_id=" + akid, {
                callback: ep.done('data'),
            });
        } else {
            ep.emit("data", l4i.Clone(iamAccessKey.def));
        }

        iam.TplCmd("access-key/set", {
            callback: ep.done('tpl'),
        });
    });
}

iamAccessKey.SetCommit = function() {
    var form = $("#iam-ak-set"),
        alert_id = "#iam-ak-set-alert";

    var req = {
        id: form.find("input[name=id]").val(),
        status: parseInt(form.find("input[name=status]:checked").val()),
        description: form.find("input[name=description]").val(),
        scopes: [],
    }

    form.find("input[name=bound_item]").each(function() {
        var val = $(this).val();
        var ar = val.split("=");
        if (ar.length == 2) {
            req.scopes.push({
                name: ar[0].trim(),
                value: ar[1].trim(),
            });
        }
    });


    iam.ApiCmd("access-key/set", {
        method: "PUT",
        data: JSON.stringify(req),
        callback: function(err, data) {

            if (err) {
                return l4i.InnerAlert(alert_id, 'alert-danger', err);
            }

            if (!data || data.error) {
                return l4i.InnerAlert(alert_id, 'alert-danger', data.error.message);
            }

            l4i.InnerAlert(alert_id, 'alert-success', l4i.T(l4i.T("Successfully %s", l4i.T("updated"))));

            window.setTimeout(function() {
                l4iModal.Close();
                iamAccessKey.List();
            }, 1000);
        },
    });
}

iamAccessKey.Del = function(akid) {
    l4iModal.Open({
        title: l4i.T("Delete"),
        tplsrc: '<div id="iam-ak-del" class="alert alert-danger">' + l4i.T("Are you sure to delete this") + '?</div>',
        height: "200px",
        buttons: [{
            title: l4i.T("Confirm and remove"),
            onclick: 'iamAccessKey.DelCommit("' + akid + '")',
            style: "btn-danger",
        }, {
            title: l4i.T("Cancel"),
            onclick: "l4iModal.Close()",
        }],
    });
}

iamAccessKey.DelCommit = function(akid) {
    var alertid = "#iam-ak-del";
    var uri = "access_key_id=" + akid;

    iam.ApiCmd("access-key/del?" + uri, {
        callback: function(err, data) {

            if (!data || data.kind != "AccessKey") {
                return l4i.InnerAlert(alertid, 'alert-danger', data.error.message);
            }

            l4i.InnerAlert(alertid, 'alert-success', l4i.T("Successfully %s", l4i.T("deleted")));
            setTimeout(function() {
                iamAccessKey.List();
                l4iModal.Close();
            }, 500);
        }
    });
}

iamAccessKey.Bind = function(akid) {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', function(tpl) {

            l4iModal.Open({
                tplsrc: tpl,
                width: 800,
                height: 260,
                title: l4i.T("Bind new instance to this %s", "AccessKey"),
                data: {
                    id: akid,
                },
                buttons: [{
                    title: l4i.T("Cancel"),
                    onclick: "l4iModal.Close()",
                }, {
                    title: l4i.T("Save"),
                    onclick: "iamAccessKey.BindCommit()",
                    style: "btn btn-primary",
                }],
            });
        });

        ep.fail(function(err) {
            alert(l4i.T("network error, please try again later"));
        });

        iam.TplCmd("access-key/bind", {
            callback: ep.done('tpl'),
        });
    });
}


iamAccessKey.BindCommit = function() {
    var form = $("#iam-ak-bind"),
        alert_id = "#iam-ak-bind-alert";

    var url = "?access_key_id=" + form.find("input[name=id]").val();
    url += "&scope_content=" + form.find("input[name=scope_content]").val();

    iam.ApiCmd("access-key/bind" + url, {
        method: "GET",
        callback: function(err, data) {

            if (err) {
                return l4i.InnerAlert(alert_id, 'alert-danger', err);
            }

            if (!data || data.error) {
                return l4i.InnerAlert(alert_id, 'alert-danger', data.error.message);
            }

            l4i.InnerAlert(alert_id, 'alert-success', l4i.T("Successfully %s", l4i.T("updated")));

            window.setTimeout(function() {
                l4iModal.Close();
                iamAccessKey.List();
            }, 1000);
        },
    });
}

iamAccessKey.UnBind = function(akid, name) {

    l4iModal.Open({
        title: l4i.T("Delete"),
        tplsrc: '<div id="iam-ak-unbind" class="alert alert-danger">' + l4i.T("Are you sure to delete this") + '?</div>',
        height: "200px",
        buttons: [{
            title: l4i.T("Confirm and remove"),
            onclick: 'iamAccessKey.UnBindCommit("' + akid + '", "' + name + '")',
            style: "btn-danger",
        }, {
            title: l4i.T("Cancel"),
            onclick: "l4iModal.Close()",
        }],
    });
}

iamAccessKey.UnBindCommit = function(akid, name) {

    var alertid = "#iam-ak-unbind";
    var url = "?access_key_id=" + akid + "&scope_content=" + name;

    iam.ApiCmd("access-key/unbind" + url, {
        method: "GET",
        callback: function(err, data) {

            if (!data || data.kind != "AccessKey") {
                return l4i.InnerAlert(alertid, 'alert-danger', data.error.message);
            }

            l4i.InnerAlert(alertid, 'alert-success', l4i.T("Successfully %s", l4i.T("deleted")));
            setTimeout(function() {
                iamAccessKey.List();
                l4iModal.Close();
            }, 500);
        },
    });
}

