var iamMyApp = {
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

iamMyApp.Init = function()
{
    l4i.UrlEventRegister("my-app/index", iamMyApp.Index);
    l4i.UrlEventRegister("my-app/inst-list", iamMyApp.InstList);
}

iamMyApp.Index = function()
{
    iam.TplCmd("my-app/index", {
        callback: function(err, data) {
            $("#com-content").html(data);
            iamMyApp.InstList();
        },
    });
}

iamMyApp.InstList = function()
{
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function (tpl, data) {

            if (!data || !data.items) {
                $("#work-content").text("There are currently no applications available");
                return;
            }

            data._statusls = iamUserMgr.statusls;

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
                dstid  : "iam-myapp-insts",
                tplid  : "iam-myapp-insts-tpl",
                data   : data,
                // success : function() {
                //     l4iTemplate.Render({
                //         dstid  : "iam-myapp-insts-pager",
                //         tplid  : "iam-myapp-insts-pager-tpl",
                //         data   : l4i.Pager(data.meta),
                //     });
                // },
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        iam.ApiCmd("my-app/inst-list", {
            callback: ep.done('data'),
        });

        iam.TplCmd("my-app/inst-list", {
            callback: ep.done('tpl'),           
        });
    });
}

iamMyApp.InstSetForm = function(instid)
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
            
            data._statusls = iamMyApp.statusls;
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
                    onclick : "iamMyApp.InstSetCommit()",
                    style   : "btn btn-primary",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        if (iamMyApp.roles) {
            ep.emit("roles", iamMyApp.roles)
        } else {
            iam.ApiCmd("user/role-list?status=0", {
                callback: function(err, roles) {
                    iamMyApp.roles = roles;
                    ep.emit("roles", roles);
                },
            });
        }

        if (instid) {
        
            iam.ApiCmd("my-app/inst-entry?instid="+ instid, {
                callback: ep.done('data'),
            });
        
        } else {
            ep.emit("data", l4i.Clone(iamMyApp.appinstdef));
        }

        iam.TplCmd("my-app/inst-set", {
            callback: ep.done('tpl'),           
        });
    });
}

iamMyApp.InstSetCommit = function()
{
    var form = $("#iam-myapp-instset");
    
    var req = l4i.Clone(iamMyApp.appinstdef)

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
        return l4i.InnerAlert("#iam-myapp-instset-alert", 'alert-danger', err);
    }

    iam.ApiCmd("my-app/inst-set", {
        method : "PUT",
        data   : JSON.stringify(req),
        callback : function(err, data) {
            
            if (err) {
                return l4i.InnerAlert("#iam-myapp-instset-alert", 'alert-danger', err);
            }

            if (!data || data.error) {
                return l4i.InnerAlert("#iam-myapp-instset-alert", 'alert-danger', data.error.message);
            }

            l4i.InnerAlert("#iam-myapp-instset-alert", 'alert-success', "Successfully updated");

            window.setTimeout(function(){
                l4iModal.Close();
                iamMyApp.InstList();
            }, 1000);
        },
    });
}
