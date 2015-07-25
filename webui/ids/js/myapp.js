var idsmyapp = {
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

idsmyapp.Init = function()
{
    l4i.UrlEventRegister("my-app/index", idsmyapp.Index);
    l4i.UrlEventRegister("my-app/inst-list", idsmyapp.InstList);
}

idsmyapp.Index = function()
{
    ids.TplCmd("my-app/index", {
        callback: function(err, data) {
            $("#com-content").html(data);
            idsmyapp.InstList();
        },
    });
}

idsmyapp.InstList = function()
{
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function (tpl, data) {

            if (!data || !data.items) {
                return;
            }

            data._statusls = idsusrmgr.statusls;

            for (var i in data.items) {

                if (!data.items[i].privileges) {
                    data.items[i].privileges = []
                }

                if (!data.items[i].version) {
                    data.items[i].version = "";
                }

                data.items[i]._privilegeNumber = data.items[i].privileges.length;
            }

            $("#work-content").html(tpl);

            l4iTemplate.Render({
                dstid  : "ids-myapp-insts",
                tplid  : "ids-myapp-insts-tpl",
                data   : data,
                // success : function() {
                //     l4iTemplate.Render({
                //         dstid  : "ids-myapp-insts-pager",
                //         tplid  : "ids-myapp-insts-pager-tpl",
                //         data   : l4i.Pager(data.meta),
                //     });
                // },
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        ids.ApiCmd("my-app/inst-list", {
            callback: ep.done('data'),
        });

        ids.TplCmd("my-app/inst-list", {
            callback: ep.done('tpl'),           
        });
    });
}

idsmyapp.InstSetForm = function(instid)
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
            
            data._statusls = idsmyapp.statusls;
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
                    onclick : "idsmyapp.InstSetCommit()",
                    style   : "btn btn-primary",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        if (idsmyapp.roles) {
            ep.emit("roles", idsmyapp.roles)
        } else {
            ids.ApiCmd("user-mgr/role-list?status=0", {
                callback: function(err, roles) {
                    idsmyapp.roles = roles;
                    ep.emit("roles", roles);
                },
            });
        }

        if (instid) {
        
            ids.ApiCmd("my-app/inst-entry?instid="+ instid, {
                callback: ep.done('data'),
            });
        
        } else {
            ep.emit("data", l4i.Clone(idsmyapp.appinstdef));
        }

        ids.TplCmd("my-app/inst-set", {
            callback: ep.done('tpl'),           
        });
    });
}

idsmyapp.InstSetCommit = function()
{
    var form = $("#ids-myapp-instset");
    
    var req = l4i.Clone(idsmyapp.appinstdef)

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
        return l4i.InnerAlert("#ids-myapp-instset-alert", 'alert-danger', err);
    }

    ids.ApiCmd("my-app/user-set", {
        method : "PUT",
        data   : JSON.stringify(req),
        callback : function(err, data) {
            
            if (err) {
                return l4i.InnerAlert("#ids-myapp-instset-alert", 'alert-danger', err);
            }

            if (!data || data.error) {
                return l4i.InnerAlert("#ids-myapp-instset-alert", 'alert-danger', data.error.message);
            }

            l4i.InnerAlert("#ids-myapp-instset-alert", 'alert-success', "Successfully updated");

            window.setTimeout(function(){
                l4iModal.Close();
                idsmyapp.InstList();
            }, 1000);
        },
    });
}
