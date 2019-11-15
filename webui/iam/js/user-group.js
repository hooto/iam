var iamUserGroup = {
    itemDef: {
        kind: "UserGroupItem",
        _isnew: true,
        name: "",
        owners: [],
        members: [],
        status: 1,
    },
    statusls: [{
        status: 1,
        title: "Active",
    }, {
        status: 2,
        title: "Banned",
    }],
}

iamUserGroup.Index = function() {
    l4iTemplate.Render({
        dstid: "com-content",
        tplurl: iam.TplPath("user-group/index"),
        data: {},
        callback: function() {
            iam.OpToolsClean();
            l4i.UrlEventClean("iam-module-navbar-menus");
            l4i.UrlEventRegister("user-group/list", iamUserGroup.List, "iam-module-navbar-menus");
            l4i.UrlEventHandler("user-group/list", true);
        },
    });
}

iamUserGroup.List = function() {

    var alertId = "#iam-user-group-list-alert";

    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            if (!data || !data.kind || data.kind != "UserGroupList") {
                return l4i.InnerAlert(alertId, 'ok', l4i.T("Successfully %s", l4i.T("updated")));
            }

            if (!data.items) {
                data.items = [];
            }

            for (var i in data.items) {

                if (!data.items[i].owners) {
                    data.items[i].owners = [];
                }

                if (!data.items[i].members) {
                    data.items[i].members = [];
                }


                data.items[i]._owners = data.items[i].owners.join(",");
                data.items[i]._members = data.items[i].members.join(",");


                if (!data.items[i].display_name) {
                    data.items[i].display_name = "";
                }

                if (!data.items[i].status) {
                    data.items[i].status = 1;
                }
            }

            data._statusls = iamUserGroup.statusls;

            l4iTemplate.Render({
                dstid: "work-content",
                tplsrc: tpl,
                callback: function() {

                    iam.OpToolsRefresh("#iam-user-group-list-optools");

                    l4iTemplate.Render({
                        dstid: "iam-user-group-list",
                        tplid: "iam-user-group-list-tpl",
                        data: data,
                    });
                },
            });
        });

        ep.fail(function(err) {
            alert(l4i.T("network error, please try again later"));
        });

        iam.ApiCmd("user-group/list", {
            callback: ep.done('data'),
        });

        iam.TplCmd("user-group/list", {
            callback: ep.done('tpl'),
        });
    });
}

iamUserGroup.setOptions = null;
iamUserGroup.SetForm = function(username, options) {

    iamUserGroup.setOptions = options || {
        callback: iamUserGroup.List,
    };

    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            if (!data || !data.kind) {
                return;
            }

            if (data._isnew) {
                data._form_title = l4i.T("New %s", l4i.T("Group"));
            } else {
                data._form_title = l4i.T("%s Settings", l4i.T("Group"));
            }

            if (!data.name) {
                data.name = "";
            }

            if (!data.display_name) {
                data.display_name = "";
            }

            if (!data.owners || data.owners.length < 1) {
                data.owners = [iam.Session.username];
            }
            data._owners = data.owners.join(",");


            if (!data.members || data.members.length < 1) {
                data.members = [iam.Session.username];
            }
            data._members = data.members.join(",");


            if (!data.status) {
                data.status = 1;
            }

            data._statusls = iamUserGroup.statusls;

            l4iModal.Open({
                tplsrc: tpl,
                width: 900,
                height: 550,
                data: data,
                title: data._form_title,
                buttons: [{
                    title: l4i.T("Cancel"),
                    onclick: "l4iModal.Close()",
                }, {
                    title: l4i.T("Save"),
                    onclick: "iamUserGroup.SetCommit()",
                    style: "btn btn-primary",
                }],
            });
        });

        ep.fail(function(err) {
            alert(l4i.T("network error, please try again later"));
        });

        if (username && username.length > 0) {
            iam.ApiCmd("user-group/item?name=" + username, {
                callback: ep.done('data'),
            });
        } else {
            ep.emit("data", l4i.Clone(iamUserGroup.itemDef));
        }

        iam.TplCmd("user-group/set", {
            callback: ep.done('tpl'),
        });
    });
}

iamUserGroup.SetCommit = function() {

    var form = $("#iam-user-group-set"),
        req = l4i.Clone(iamUserGroup.itemDef);

    try {

        req.name = form.find("input[name=name]").val();
        req.display_name = form.find("input[name=display_name]").val();
        req.status = parseInt(form.find("input[name=status]:checked").val());

        //
        var arr = form.find("input[name=owners]").val();
        if (!arr) {
            throw "No Owner Found";
        }
        arr = arr.replace(/(?:\r\n|\r|\n| |)/g, '');
        req.owners = arr.split(",");

        //
        arr = form.find("textarea[name=members]").val();
        if (!arr) {
            throw "No Member Found";
        }
        arr = arr.replace(/(?:\r\n|\r|\n| |)/g, '');
        req.members = arr.split(",");


    } catch (err) {
        return l4iModal.FootAlert('error', err, 1000);
    }


    iam.ApiCmd("user-group/set", {
        method: "PUT",
        data: JSON.stringify(req),
        callback: function(err, data) {

            if (err) {
                return l4iModal.FootAlert('error', err, 1000);
            }

            if (!data || data.error) {
                return l4iModal.FootAlert('error', data.error.message, 1000);
            }

            l4iModal.FootAlert('ok', l4i.T("Successfully %s", l4i.T("updated")), 1000);

            window.setTimeout(function() {
                l4iModal.Close();
                if (iamUserGroup.setOptions.callback) {
                    iamUserGroup.setOptions.callback();
                }
            }, 1000);
        },
    });
}

