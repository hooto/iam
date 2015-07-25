var idsappmgr = {
    statusls : [{
        status: 1, title: "Active",
    }, {
        status: 2, title: "Banned",
    }],
    roles : null,
    appinstdef : {
        meta : {
            id : "",
        },
        app_id : "",
        app_title : "",
        version : "",
        status : 1,
        url : "",
    },
}

idsappmgr.Init = function()
{
    l4i.UrlEventRegister("app-mgr/index", idsappmgr.Index);
    l4i.UrlEventRegister("app-mgr/inst-list", idsappmgr.InstList);
}

idsappmgr.Index = function()
{
    ids.TplCmd("app-mgr/index", {
        callback: function(err, data) {
            $("#com-content").html(data);
            idsappmgr.InstList();
        },
    });
}

idsappmgr.InstList = function()
{
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function (tpl, data) {

            if (!data || !data.items) {
                return;
            }

            data._statusls = idsusrmgr.statusls;

            for (var i in data.items) {

                if (!data.items[i].privileges) {
                    data.items[i].privileges = [];
                }

                if (!data.items[i].version) {
                    data.items[i].version = "";
                }

                data.items[i]._privilegeNumber = data.items[i].privileges.length;
            }

            $("#work-content").html(tpl);

            l4iTemplate.Render({
                dstid  : "ids-appmgr-insts",
                tplid  : "ids-appmgr-insts-tpl",
                data   : data,
                success : function() {
                    l4iTemplate.Render({
                        dstid  : "ids-appmgr-insts-pager",
                        tplid  : "ids-appmgr-insts-pager-tpl",
                        data   : l4i.Pager(data.meta),
                    });
                },
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        ids.ApiCmd("app-mgr/inst-list", {
            callback: ep.done('data'),
        });

        ids.TplCmd("app-mgr/inst-list", {
            callback: ep.done('tpl'),           
        });
    });
}

idsappmgr.InstSetForm = function(instid)
{
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'roles', 'data', function (tpl, roles, data) {

            if (!data || !data.kind) {
                return;
            }

            if (!data.privileges) {
                data.privileges = [];
            }
            data._privilegeNumber = data.privileges.length;
            
            data._statusls = idsappmgr.statusls;
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
                tplsrc  : tpl,
                width   : 900,
                height  : 600,
                data    : data,
                title   : "App Instance Setting",
                buttons : [{
                    title : "Cancel",
                    onclick : "l4iModal.Close()",
                }, {
                    title : "Save",
                    onclick : "idsappmgr.InstSetCommit()",
                    style   : "btn btn-primary",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        if (idsappmgr.roles) {
            ep.emit("roles", idsappmgr.roles)
        } else {
            ids.ApiCmd("user-mgr/role-list?status=0", {
                callback: function(err, roles) {
                    idsappmgr.roles = roles;
                    ep.emit("roles", roles);
                },
            });
        }

        if (instid) {
        
            ids.ApiCmd("app-mgr/inst-entry?instid="+ instid, {
                callback: ep.done('data'),
            });
        
        } else {
            ep.emit("data", l4i.Clone(idsappmgr.appinstdef));
        }

        ids.TplCmd("app-mgr/inst-set", {
            callback: ep.done('tpl'),           
        });
    });
}

idsappmgr.InstSetCommit = function()
{
    var form = $("#ids-appmgr-instset");
    
    var req = l4i.Clone(idsappmgr.appinstdef)

    req.meta.id = form.find("input[name=instid]").val();
    req.app_title = form.find("input[name=app_title]").val();
    req.url = form.find("input[name=url]").val();

    req.status = parseInt(form.find("input[name=status]:checked").val());

    try {

        // form.find("input[name=roles]:checked").each(function() {
            
        //     var val = parseInt($(this).val());
        //     if (val > 0) {
        //         req.roles.push(val);
        //     }
        // });

    } catch (err) {
        return l4i.InnerAlert("#ids-appmgr-instset-alert", 'alert-danger', err);
    }

    ids.ApiCmd("app-mgr/inst-set", {
        method : "PUT",
        data   : JSON.stringify(req),
        callback : function(err, data) {
            
            if (err) {
                return l4i.InnerAlert("#ids-appmgr-instset-alert", 'alert-danger', err);
            }

            if (!data || data.error) {
                return l4i.InnerAlert("#ids-appmgr-instset-alert", 'alert-danger', data.error.message);
            }

            l4i.InnerAlert("#ids-appmgr-instset-alert", 'alert-success', "Successfully updated");

            window.setTimeout(function(){
                l4iModal.Close();
                idsappmgr.InstList();
            }, 1000);
        },
    });
}
