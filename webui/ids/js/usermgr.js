var idsusrmgr = {
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

idsusrmgr.Init = function()
{
    l4i.UrlEventRegister("user-mgr/index", idsusrmgr.Index);
    l4i.UrlEventRegister("user-mgr/user-list", idsusrmgr.UserList);
    l4i.UrlEventRegister("user-mgr/role-list", idsusrmgr.RoleList);
}

idsusrmgr.Index = function()
{
    ids.TplCmd("user-mgr/index", {
        callback: function(err, data) {
            $("#com-content").html(data);
            idsusrmgr.UserList();
        },
    });
}

idsusrmgr.UserList = function()
{
    //     var uri = "query_text="+ $("#query_text").val();
    // uri += "&page="+ page;

    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'roles', 'data', function (tpl, roles, data) {

            if (!data || !data.items) {
                return;
            }

            data._roles = roles;
            data._statusls = idsusrmgr.statusls;

            $("#work-content").html(tpl);

            l4iTemplate.Render({
                dstid  : "ids-usermgr-list",
                tplid  : "ids-usermgr-list-tpl",
                data   : data,
                success : function() {
                    l4iTemplate.Render({
                        dstid  : "ids-usermgr-list-pager",
                        tplid  : "ids-usermgr-list-pager-tpl",
                        data   : l4i.Pager(data.meta),
                    });
                },
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        ids.ApiCmd("user-mgr/user-list", {
            callback: ep.done('data'),
        });

        if (idsusrmgr.roles) {
            ep.emit("roles", idsusrmgr.roles)
        } else {
            ids.ApiCmd("user-mgr/role-list", {
                callback: function(err, data) {
                    idsusrmgr.roles = data;
                    ep.emit("roles", data);
                },
            });
        }        

        ids.TplCmd("user-mgr/user-list", {
            callback: ep.done('tpl'),           
        });
    });
}

idsusrmgr.UserSetForm = function(userid)
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
                    if (data._roles.items[i].idxid == data.roles[j]) {
                        data._roles.items[i].checked = true;
                        break;
                    }
                }
            }

            data._statusls = idsusrmgr.statusls; 

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
                    onclick : "idsusrmgr.UserSetCommit()",
                    style   : "btn btn-primary",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        if (idsusrmgr.roles) {
            ep.emit("roles", idsusrmgr.roles)
        } else {
            ids.ApiCmd("user-mgr/role-list", {
                callback: function(err, roles) {
                    idsusrmgr.roles = roles;
                    ep.emit("roles", roles);
                },
            });
        }

        if (userid) {
        
            ids.ApiCmd("user-mgr/user-entry?userid="+ userid, {
                callback: ep.done('data'),
            });
        
        } else {
            ep.emit("data", l4i.Clone(idsusrmgr.userdef));
        }

        ids.TplCmd("user-mgr/user-set", {
            callback: ep.done('tpl'),           
        });
    });
}

idsusrmgr.UserSetCommit = function()
{
    var form = $("#ids-usermgr-userset");
    
    var req = l4i.Clone(idsusrmgr.userdef)

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
        return l4i.InnerAlert("#ids-usermgr-userset-alert", 'alert-danger', err);
    }

    ids.ApiCmd("user-mgr/user-set", {
        method : "PUT",
        data   : JSON.stringify(req),
        callback : function(err, data) {
            
            if (err) {
                return l4i.InnerAlert("#ids-usermgr-userset-alert", 'alert-danger', err);
            }

            if (!data || data.error) {
                return l4i.InnerAlert("#ids-usermgr-userset-alert", 'alert-danger', data.error.message);
            }

            l4i.InnerAlert("#ids-usermgr-userset-alert", 'alert-success', "Successfully updated");

            window.setTimeout(function(){
                l4iModal.Close();
                idsusrmgr.UserList();
            }, 1000);
        },
    });
}

// Role
idsusrmgr.RoleList = function()
{
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'roles', function (tpl, roles) {

            if (!roles || !roles.items) {
                return;
            }

            $("#work-content").html(tpl);

            roles._statusls = idsusrmgr.statusls;

            l4iTemplate.Render({
                dstid  : "ids-usermgr-rolelist",
                tplid  : "ids-usermgr-rolelist-tpl",
                data   : roles,
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        ids.ApiCmd("user-mgr/role-list?status=0", {
            callback: ep.done('roles'),
        });

        ids.TplCmd("user-mgr/role-list", {
            callback: ep.done('tpl'),           
        });
    });
}

idsusrmgr.RoleSetForm = function(roleid)
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

            data._statusls = idsusrmgr.statusls;

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
                    onclick : "idsusrmgr.RoleSetCommit()",
                    style   : "btn btn-primary",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        if (roleid) {
        
            ids.ApiCmd("user-mgr/role-entry?roleid="+ roleid, {
                callback: ep.done('data'),
            });
        
        } else {
            ep.emit("data", l4i.Clone(idsusrmgr.roledef));
        }

        ids.TplCmd("user-mgr/role-set", {
            callback: ep.done('tpl'),           
        });
    });
}

idsusrmgr.RoleSetCommit = function()
{
    var form = $("#ids-usermgr-roleset");
    
    var req = {
        meta : {
            id : form.find("input[name=roleid]").val(),
            name : form.find("input[name=name]").val(),
        },
        status : parseInt(form.find("input[name=status]:checked").val()),
        desc : form.find("input[name=desc]").val(),
    }

    ids.ApiCmd("user-mgr/role-set", {
        method : "PUT",
        data   : JSON.stringify(req),
        callback : function(err, data) {
            
            if (err) {
                return l4i.InnerAlert("#ids-usermgr-roleset-alert", 'alert-danger', err);
            }

            if (!data || data.error) {
                return l4i.InnerAlert("#ids-usermgr-roleset-alert", 'alert-danger', data.error.message);
            }

            l4i.InnerAlert("#ids-usermgr-roleset-alert", 'alert-success', "Successfully updated");

            window.setTimeout(function(){
                l4iModal.Close();
                idsusrmgr.RoleList();
            }, 1000);
        },
    });
}

