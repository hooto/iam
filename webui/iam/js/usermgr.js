var iamUserMgr = {
    roles : null,
    statusls : [{
        status: 1, title: "Active",
    }, {
        status: 2, title: "Banned",
    }],
    userdef : {
        kind : "User",
        _isnew : true,
        meta : {
            id : "",
            name : "",
        },
        email : "",
        name : "",
        roles : [],
        status : 1,
        auth : "",
        profile : {
            name : "",
            birthday : "",
            about : "",
        },
    },
    roledef : {
        kind : "UserRole",
        meta : {
            id : "",
            name : "",
        },
        status : 1,
        desc : "",
    },
}

iamUserMgr.Index = function()
{
    iam.OpToolActive = null;
    iam.TplCmd("user-mgr/index", {
        callback: function(err, data) {
            $("#com-content").html(data);
            l4i.UrlEventClean("iam-module-navbar-menus");
            l4i.UrlEventRegister("user-mgr/user-list", iamUserMgr.UserList, "iam-module-navbar-menus");
            l4i.UrlEventRegister("user-mgr/role-list", iamUserMgr.RoleList, "iam-module-navbar-menus");
            l4i.UrlEventHandler("user-mgr/user-list", true);
        },
    });
}

iamUserMgr.UserList = function()
{
    var uri = "";
    if (document.getElementById("iam_usermgr_list_qry_text")) {
        var qt = $("#iam_usermgr_list_qry_text").val();
        if (qt && qt.length > 0) {
            uri = "?qry_text="+ qt;
        }
    }

    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'roles', 'data', function (tpl, roles, data) {

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
            }

            $("#work-content").html(tpl);
            iam.OpToolsRefresh("#iam-usermgr-list-optools");

            l4iTemplate.Render({
                dstid  : "iam-usermgr-list",
                tplid  : "iam-usermgr-list-tpl",
                data   : data,
                success : function() {
                    l4iTemplate.Render({
                        dstid  : "iam-usermgr-list-pager",
                        tplid  : "iam-usermgr-list-pager-tpl",
                        data   : l4i.Pager(data.meta),
                    });
                },
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        iam.ApiCmd("user-mgr/user-list"+ uri, {
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

iamUserMgr.UserSetForm = function(userid)
{
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'roles', 'data', function (tpl, roles, data) {

            if (!data || !data.kind) {
                return;
            }

            if (data._isnew) {
                data._form_title = "New User";
            } else {
                data._form_title = "User Setting";
            }

            if (!data.name) {
                data.name = "";
            }

            if (!data.email) {
                data.email = "";
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
                tplsrc  : tpl,
                width   : 800,
                height  : 600,
                data    : data,
                title   : data._form_title,
                buttons : [{
                    title : "Cancel",
                    onclick : "l4iModal.Close()",
                }, {
                    title : "Save",
                    onclick : "iamUserMgr.UserSetCommit()",
                    style   : "btn btn-primary",
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

        if (userid) {
        
            iam.ApiCmd("user-mgr/user-entry?userid="+ userid, {
                callback: ep.done('data'),
            });
        
        } else {
            ep.emit("data", l4i.Clone(iamUserMgr.userdef));
        }

        iam.TplCmd("user-mgr/user-set", {
            callback: ep.done('tpl'),           
        });
    });
}

iamUserMgr.UserSetCommit = function()
{
    var form = $("#iam-usermgr-userset");
    
    var req = l4i.Clone(iamUserMgr.userdef)

    req.meta.id = form.find("input[name=userid]").val();
    if (req.meta.id == "") {
        req.meta.name = form.find("input[name=username]").val();
    }
    req.email = form.find("input[name=email]").val();
    req.auth = form.find("input[name=auth]").val();
    req.name = form.find("input[name=name]").val();

    req.profile.birthday = form.find("input[name=birthday]").val();
    req.profile.about = form.find("textarea[name=about]").val();

    try {

        form.find("input[name=roles]:checked").each(function() {
            
            var val = parseInt($(this).val());
            if (val > 0) {
                req.roles.push(val);
            }            
        });

    } catch (err) {
        return l4i.InnerAlert("#iam-usermgr-userset-alert", 'alert-danger', err);
    }

    iam.ApiCmd("user-mgr/user-set", {
        method : "PUT",
        data   : JSON.stringify(req),
        callback : function(err, data) {
            
            if (err) {
                return l4i.InnerAlert("#iam-usermgr-userset-alert", 'alert-danger', err);
            }

            if (!data || data.error) {
                return l4i.InnerAlert("#iam-usermgr-userset-alert", 'alert-danger', data.error.message);
            }

            l4i.InnerAlert("#iam-usermgr-userset-alert", 'alert-success', "Successfully updated");

            window.setTimeout(function(){
                l4iModal.Close();
                iamUserMgr.UserList();
            }, 1000);
        },
    });
}

// Role
iamUserMgr.RoleList = function()
{
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'roles', function (tpl, roles) {

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
                dstid  : "iam-usermgr-rolelist",
                tplid  : "iam-usermgr-rolelist-tpl",
                data   : roles,
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

iamUserMgr.RoleSet = function(roleid)
{
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function (tpl, data) {

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
                tplsrc  : tpl,
                width   : 600,
                height  : 400,
                data    : data,
                title   : data._form_title,
                buttons : [{
                    title : "Cancel",
                    onclick : "l4iModal.Close()",
                }, {
                    title : "Save",
                    onclick : "iamUserMgr.RoleSetCommit()",
                    style   : "btn btn-primary",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        if (roleid) {
        
            iam.ApiCmd("user-mgr/role-entry?roleid="+ roleid, {
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

iamUserMgr.RoleSetCommit = function()
{
    var form = $("#iam-usermgr-roleset");
    
    var req = {
        meta : {
            id : form.find("input[name=roleid]").val(),
            name : form.find("input[name=name]").val(),
        },
        status : parseInt(form.find("input[name=status]:checked").val()),
        desc : form.find("input[name=desc]").val(),
    }

    iam.ApiCmd("user-mgr/role-set", {
        method : "PUT",
        data   : JSON.stringify(req),
        callback : function(err, data) {
            
            if (err) {
                return l4i.InnerAlert("#iam-usermgr-roleset-alert", 'alert-danger', err);
            }

            if (!data || data.error) {
                return l4i.InnerAlert("#iam-usermgr-roleset-alert", 'alert-danger', data.error.message);
            }

            l4i.InnerAlert("#iam-usermgr-roleset-alert", 'alert-success', "Successfully updated");

            window.setTimeout(function(){
                l4iModal.Close();
                iamUserMgr.RoleList();
            }, 1000);
        },
    });
}

