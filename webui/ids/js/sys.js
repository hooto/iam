var idssys = {

}

idssys.Init = function()
{
    l4i.UrlEventRegister("sys-mgr/index", idssys.Index);
    l4i.UrlEventRegister("sys-mgr/general-set", idssys.GeneralSet);
    l4i.UrlEventRegister("sys-mgr/mailer-set", idssys.MailerSet);
}

idssys.Index = function()
{
    ids.TplCmd("sys-mgr/index", {
        callback: function(err, data) {
            $("#com-content").html(data);
            idssys.GeneralSet();
        },
    });
}

idssys.GeneralSet = function()
{
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function (tpl, data) {

            if (!data || !data.items) {
                return;
            }

            data._items = {};
            for (var i in data.items) {
                data._items[data.items[i]["key"]] = data.items[i]["val"];
            }

            l4iTemplate.Render({
                dstid  : "work-content",
                tplsrc : tpl,
                data   : data,
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });
    
        ids.ApiCmd("sys-config/general", {
            callback: ep.done('data'),
        });

        ids.TplCmd("sys-mgr/general-set", {
            callback: ep.done('tpl'),           
        });
    });
}

idssys.GeneralSetCommit = function()
{
    var form = $("#ids-sysmgr-generalset");
    
    var req = {
        items : [{
            key: "service_name",
            val: form.find("input[name=service_name]").val(),
        }, {
            key: "webui_banner_title",
            val: form.find("input[name=webui_banner_title]").val(),
        }],
    };

    console.log(req);

    ids.ApiCmd("sys-config/general-set", {
        method : "PUT",
        data   : JSON.stringify(req),
        callback : function(err, data) {
            
            if (err) {
                return l4i.InnerAlert("#ids-sysmgr-generalset-alert", 'alert-danger', err);
            }

            if (!data || data.error) {
                return l4i.InnerAlert("#ids-sysmgr-generalset-alert", 'alert-danger', data.error.message);
            }

            l4i.InnerAlert("#ids-sysmgr-generalset-alert", 'alert-success', "Successfully updated");

            // window.setTimeout(function(){
            //     idssys.GeneralSet();
            // }, 1000);
        },
    });
}


idssys.MailerSet = function()
{
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function (tpl, data) {

            if (!data || !data.items) {
                return;
            }

            data._items = {};
            for (var i in data.items) {
                data._items[data.items[i]["key"]] = data.items[i]["val"];
            }

            l4iTemplate.Render({
                dstid  : "work-content",
                tplsrc : tpl,
                data   : data,
                success : function() {
                    // idsuser.Overview();
                },
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });
    
        ids.ApiCmd("sys-config/mailer", {
            callback: ep.done('data'),
        });

        ids.TplCmd("sys-mgr/mailer-set", {
            callback: ep.done('tpl'),           
        });
    });  
}


idssys.MailerSetCommit = function()
{
    var form = $("#ids-sysmgr-mailerset");
    
    var req = {
        items : [{
            key: "mailer_smtp_host",
            val: form.find("input[name=mailer_smtp_host]").val(),
        }, {
            key: "mailer_smtp_port",
            val: form.find("input[name=mailer_smtp_port]").val(),
        }, {
            key: "mailer_smtp_user",
            val: form.find("input[name=mailer_smtp_user]").val(),
        }, {
            key: "mailer_smtp_pass",
            val: form.find("input[name=mailer_smtp_pass]").val(),
        }],
    };

    ids.ApiCmd("sys-config/mailer-set", {
        method : "PUT",
        data   : JSON.stringify(req),
        callback : function(err, data) {
            
            if (err) {
                return l4i.InnerAlert("#ids-sysmgr-mailerset-alert", 'alert-danger', err);
            }

            if (!data || data.error) {
                return l4i.InnerAlert("#ids-sysmgr-mailerset-alert", 'alert-danger', data.error.message);
            }

            l4i.InnerAlert("#ids-sysmgr-mailerset-alert", 'alert-success', "Successfully updated");

            // window.setTimeout(function(){
            //     idssys.GeneralSet();
            // }, 1000);
        },
    });
}

