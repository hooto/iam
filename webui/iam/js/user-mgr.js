var iamUserMgr = {
    roles: null,
    statusls: [{
        status: 1,
        title: "Active",
    }, {
        status: 2,
        title: "Banned",
    }],
    userEntryDef: {
        kind: "User",
        _isnew: true,
        login: {
            id: "",
            name: "",
            email: "",
            display_name: "",
            roles: [],
            status: 1,
            keys: [{
                key: "std",
                value: "",
            }],
        },
        profile: {
            birthday: "",
            about: "",
        },
    },
    userProfileDef: {
        birthday: "",
        about: "",
    },
    roledef: {
        kind: "UserRole",
        id: 0,
        name: "",
        status: 1,
        desc: "",
    },
}

iamUserMgr.Index = function() {
    iam.TplCmd("user-mgr/index", {
        callback: function(err, data) {
            iam.OpToolsClean();
            $("#com-content").html(data);
            l4i.UrlEventClean("iam-module-navbar-menus");
            l4i.UrlEventRegister("user-mgr/user-list", iamUserMgr.UserList, "iam-module-navbar-menus");
            l4i.UrlEventRegister("user-mgr/role-list", iamUserMgr.RoleList, "iam-module-navbar-menus");
            l4i.UrlEventHandler("user-mgr/user-list", true);
        },
    });
}

iamUserMgr.UserList = function() {
    var uri = "";
    if (document.getElementById("iam_usermgr_list_qry_text")) {
        var qt = $("#iam_usermgr_list_qry_text").val();
        if (qt && qt.length > 0) {
            uri = "?qry_text=" + qt;
        }
    }

    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'roles', 'data', function(tpl, roles, data) {

            if (!data || !data.items) {
                return;
            }

            iamUserMgr.roles = l4i.Clone(roles);

            data._roles = roles;
            data._statusls = iamUserMgr.statusls;

            for (var i in data.items) {
                if (!data.items[i].email) {
                    data.items[i].email = "";
                }
                if (!data.items[i].display_name) {
                    data.items[i].display_name = "";
                }
            }

            $("#work-content").html(tpl);
            iam.OpToolsRefresh("#iam-usermgr-list-optools");

            l4iTemplate.Render({
                dstid: "iam-usermgr-list",
                tplid: "iam-usermgr-list-tpl",
                data: data,
                success: function() {
                    l4iTemplate.Render({
                        dstid: "iam-usermgr-list-pager",
                        tplid: "iam-usermgr-list-pager-tpl",
                        data: l4i.Pager(data.meta),
                    });
                },
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        iam.ApiCmd("user-mgr/user-list" + uri, {
            callback: ep.done('data'),
        });

        if (iamUserMgr.roles) {
            ep.emit("roles", iamUserMgr.roles)
        } else {
            iam.ApiCmd("user-mgr/role-list", {
                callback: ep.done('roles'),
            });
        }

        iam.TplCmd("user-mgr/user-list", {
            callback: ep.done('tpl'),
        });
    });
}

iamUserMgr.UserSetForm = function(username) {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'roles', 'data', function(tpl, roles, data) {

            if (!data || !data.kind) {
                return;
            }

            if (data._isnew) {
                data._form_title = "New User";
            } else {
                data._form_title = "User Setting";
            }

            if (!data.login.display_name) {
                data.login.display_name = "";
            }

            if (!data.login.email) {
                data.login.email = "";
            }

            for (var i in data.login.keys) {
                if (data.login.keys[i].key == "std") {
                    data.login._auth = data.login.keys[i].value;
                    break;
                }
            }
            if (!data.login._auth) {
                data.login._auth = "";
            }

            if (!data.profile) {
                data.profile = l4i.Clone(iamUserMgr.userProfileDef);
            }
            if (!data.profile.birthday) {
                data.profile.birthday = "";
            }
            if (!data.profile.about) {
                data.profile.about = "";
            }

            data._roles = l4i.Clone(roles);
            for (var i in data._roles.items) {

                data._roles.items[i].checked = false;

                for (var j in data.roles) {
                    if (data._roles.items[i].id == data.roles[j]) {
                        data._roles.items[i].checked = true;
                        break;
                    }
                }
            }

            data._statusls = iamUserMgr.statusls;

            l4iModal.Open({
                tplsrc: tpl,
                width: 800,
                height: 600,
                data: data,
                title: data._form_title,
                buttons: [{
                    title: "Cancel",
                    onclick: "l4iModal.Close()",
                }, {
                    title: "Save",
                    onclick: "iamUserMgr.UserSetCommit()",
                    style: "btn btn-primary",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        if (iamUserMgr.roles) {
            ep.emit("roles", iamUserMgr.roles)
        } else {
            iam.ApiCmd("user-mgr/role-list", {
                callback: function(err, roles) {
                    iamUserMgr.roles = roles;
                    ep.emit("roles", roles);
                },
            });
        }

        if (username && username.length > 0) {
            iam.ApiCmd("user-mgr/user-entry?username=" + username, {
                callback: ep.done('data'),
            });
        } else {
            ep.emit("data", l4i.Clone(iamUserMgr.userEntryDef));
        }

        iam.TplCmd("user-mgr/user-set", {
            callback: ep.done('tpl'),
        });
    });
}

iamUserMgr.UserSetCommit = function() {
    var form = $("#iam-usermgr-userset"),
        alert_id = "#iam-usermgr-userset-alert",
        req = l4i.Clone(iamUserMgr.userEntryDef);

    try {

        req.login.name = form.find("input[name=login_name]").val();
        req.login.email = form.find("input[name=login_email]").val();
        req.login.keys = [{
            key: "std",
            value: form.find("input[name=login_auth]").val(),
        }];
        req.login.display_name = form.find("input[name=login_display_name]").val();

        req.profile.birthday = form.find("input[name=profile_birthday]").val();
        req.profile.about = form.find("textarea[name=profile_about]").val();

        form.find("input[name=login_roles]:checked").each(function() {
            var val = parseInt($(this).val());
            if (val > 0) {
                req.login.roles.push(val);
            }
        });

    } catch (err) {
        return l4i.InnerAlert(alert_id, 'alert-danger', err);
    }

    iam.ApiCmd("user-mgr/user-set", {
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
                iamUserMgr.UserList();
            }, 1000);
        },
    });
}

// Role
iamUserMgr.RoleList = function() {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'roles', function(tpl, roles) {

            if (!roles || !roles.items) {
                return;
            }

            $("#work-content").html(tpl);
            iam.OpToolsRefresh("#iam-usermgr-rolelist-optools");

            for (var i in roles.items) {

                if (!roles.items[i].desc) {
                    roles.items[i].desc = "";
                }
            }

            roles._statusls = iamUserMgr.statusls;

            l4iTemplate.Render({
                dstid: "iam-usermgr-rolelist",
                tplid: "iam-usermgr-rolelist-tpl",
                data: roles,
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        iam.ApiCmd("user-mgr/role-list?status=0", {
            callback: ep.done('roles'),
        });

        iam.TplCmd("user-mgr/role-list", {
            callback: ep.done('tpl'),
        });
    });
}

iamUserMgr.RoleSet = function(roleid) {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            if (!data || !data.kind) {
                return;
            }

            if (data._isnew) {
                data._form_title = "New Role";
            } else {
                data._form_title = "Role Setting";
            }

            if (!data.desc) {
                data.desc = "";
            }

            data._statusls = iamUserMgr.statusls;

            l4iModal.Open({
                tplsrc: tpl,
                width: 600,
                height: 400,
                data: data,
                title: data._form_title,
                buttons: [{
                    title: "Cancel",
                    onclick: "l4iModal.Close()",
                }, {
                    title: "Save",
                    onclick: "iamUserMgr.RoleSetCommit()",
                    style: "btn btn-primary",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        if (roleid) {
            iam.ApiCmd("user-mgr/role-entry?roleid=" + roleid, {
                callback: ep.done('data'),
            });
        } else {
            ep.emit("data", l4i.Clone(iamUserMgr.roledef));
        }

        iam.TplCmd("user-mgr/role-set", {
            callback: ep.done('tpl'),
        });
    });
}

iamUserMgr.RoleSetCommit = function() {
    var form = $("#iam-usermgr-roleset"),
        alert_id = "#iam-usermgr-roleset-alert";

    var req = {
        id: parseInt(form.find("input[name=roleid]").val()),
        name: form.find("input[name=name]").val(),
        status: parseInt(form.find("input[name=status]:checked").val()),
        desc: form.find("input[name=desc]").val(),
    }

    iam.ApiCmd("user-mgr/role-set", {
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
                iamUserMgr.RoleList();
            }, 1000);
        },
    });
}

