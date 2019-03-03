var iamAccessKey = {
    aks: null,
    actionls: [{
        action: 1,
        title: "Active",
    }, {
        action: 2,
        title: "Stop",
    }],
    def: {
        kind: "AccessKey",
        access_key: "",
        secret_key: "",
        action: 1,
        desc: "",
        bounds: [],
    },
}

iamAccessKey.Index = function() {
    iam.TplCmd("access-key/index", {
        callback: function(err, data) {
            iam.OpToolsClean();
            $("#com-content").html(data);
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

            $("#work-content").html(tpl);
            iam.OpToolsRefresh("#iam-aklist-optools");

            for (var i in aks.items) {

                if (!aks.items[i].desc) {
                    aks.items[i].desc = "";
                }

                if (!aks.items[i].bounds) {
                    aks.items[i].bounds = [];
                }

                aks.items[i]._bounds = [];
                for (var j in aks.items[i].bounds) {
                    aks.items[i]._bounds.push(aks.items[i].bounds[j].name);
                }
            }

            aks._actionls = iamAccessKey.actionls;

            l4iTemplate.Render({
                dstid: "iam-aklist",
                tplid: "iam-aklist-tpl",
                data: aks,
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        iam.ApiCmd("access-key/list?action=0", {
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

            if (!data.desc) {
                data.desc = "";
            }
            if (!data.access_key) {
                data.access_key = "";
            }
            if (!data.action) {
                data.action = 1;
            }
            if (!data.bounds) {
                data.bounds = [];
            }

            data._actionls = iamAccessKey.actionls;

            l4iModal.Open({
                tplsrc: tpl,
                width: 800,
                height: 350,
                data: data,
                title: "Access Key Info",
                buttons: [{
                    title: "Cancel",
                    onclick: "l4iModal.Close()",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        if (akid) {
            iam.ApiCmd("access-key/entry?access_key=" + akid, {
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
                data._form_title = "New Access Key";
            } else {
                data._form_title = "Access Key Setting";
            }

            if (!data.desc) {
                data.desc = "";
            }
            if (!data.access_key) {
                data.access_key = "";
            }
            if (!data.action) {
                data.action = 1;
            }
            if (!data.bounds) {
                data.bounds = [];
            }

            data._actionls = iamAccessKey.actionls;

            l4iModal.Open({
                tplsrc: tpl,
                width: 800,
                height: 300,
                data: data,
                title: data._form_title,
                buttons: [{
                    title: "Cancel",
                    onclick: "l4iModal.Close()",
                }, {
                    title: "Save",
                    onclick: "iamAccessKey.SetCommit()",
                    style: "btn btn-primary",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        if (akid) {
            iam.ApiCmd("access-key/entry?access_key=" + akid, {
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
        access_key: form.find("input[name=access_key]").val(),
        action: parseInt(form.find("input[name=action]:checked").val()),
        desc: form.find("input[name=desc]").val(),
        bounds: [],
    }

    form.find("input[name=bound_item]").each(function() {
        var val = $(this).val();
        if (val.length > 0) {
            req.bounds.push({
                name: val,
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

            l4i.InnerAlert(alert_id, 'alert-success', "Successfully updated");

            window.setTimeout(function() {
                l4iModal.Close();
                iamAccessKey.List();
            }, 1000);
        },
    });
}

iamAccessKey.Del = function(akid) {
    l4iModal.Open({
        title: "Delete",
        tplsrc: '<div id="iam-ak-del" class="alert alert-danger">Are you sure to delete this?</div>',
        height: "200px",
        buttons: [{
            title: "Confirm to delete",
            onclick: 'iamAccessKey.DelCommit("' + akid + '")',
            style: "btn-danger",
        }, {
            title: "Cancel",
            onclick: "l4iModal.Close()",
        }],
    });
}

iamAccessKey.DelCommit = function(akid) {
    var alertid = "#iam-ak-del";
    var uri = "access_key=" + akid;

    iam.ApiCmd("access-key/del?" + uri, {
        callback: function(err, data) {

            if (!data || data.kind != "AccessKey") {
                return l4i.InnerAlert(alertid, 'alert-danger', data.error.message);
            }

            l4i.InnerAlert(alertid, 'alert-success', "Successful deleted");
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
                title: "Bind Instance to this AccessKey",
                data: {
                    access_key: akid,
                },
                buttons: [{
                    title: "Cancel",
                    onclick: "l4iModal.Close()",
                }, {
                    title: "Save",
                    onclick: "iamAccessKey.BindCommit()",
                    style: "btn btn-primary",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        iam.TplCmd("access-key/bind", {
            callback: ep.done('tpl'),
        });
    });
}


iamAccessKey.BindCommit = function() {
    var form = $("#iam-ak-bind"),
        alert_id = "#iam-ak-bind-alert";

    var url = "?access_key=" + form.find("input[name=access_key]").val();
    url += "&bound_name=" + form.find("input[name=bound_name]").val();

    iam.ApiCmd("access-key/bind" + url, {
        method: "GET",
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
                iamAccessKey.List();
            }, 1000);
        },
    });
}

iamAccessKey.UnBind = function(akid, name) {

    l4iModal.Open({
        title: "Delete",
        tplsrc: '<div id="iam-ak-unbind" class="alert alert-danger">Are you sure to delete this?</div>',
        height: "200px",
        buttons: [{
            title: "Confirm to delete",
            onclick: 'iamAccessKey.UnBindCommit("' + akid + '", "' + name + '")',
            style: "btn-danger",
        }, {
            title: "Cancel",
            onclick: "l4iModal.Close()",
        }],
    });
}

iamAccessKey.UnBindCommit = function(akid, name) {

    var alertid = "#iam-ak-unbind";
    var url = "?access_key=" + akid + "&bound_name=" + name;

    iam.ApiCmd("access-key/unbind" + url, {
        method: "GET",
        callback: function(err, data) {

            if (!data || data.kind != "AccessKey") {
                return l4i.InnerAlert(alertid, 'alert-danger', data.error.message);
            }

            l4i.InnerAlert(alertid, 'alert-success', "Successful deleted");
            setTimeout(function() {
                iamAccessKey.List();
                l4iModal.Close();
            }, 500);
        },
    });
}

