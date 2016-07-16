var iamsys = {

}

iamsys.Init = function()
{
    l4i.UrlEventRegister("sys-mgr/index", iamsys.Index);
    l4i.UrlEventRegister("sys-mgr/general-set", iamsys.GeneralSet);
    l4i.UrlEventRegister("sys-mgr/mailer-set", iamsys.MailerSet);
}

iamsys.Index = function()
{
    iam.TplCmd("sys-mgr/index", {
        callback: function(err, data) {
            $("#com-content").html(data);
            iamsys.GeneralSet();
        },
    });
}

iamsys.GeneralSet = function()
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
    
        iam.ApiCmd("sys-config/general", {
            callback: ep.done('data'),
        });

        iam.TplCmd("sys-mgr/general-set", {
            callback: ep.done('tpl'),           
        });
    });
}

iamsys.GeneralSetCommit = function()
{
    var form = $("#iam-sysmgr-generalset");

    var user_reg_disable = "0";
    if (form.find("input[name=user_reg_disable]").is(":checked")) {
        user_reg_disable = "1";
    }
    
    var req = {
        items : [{
            key: "service_name",
            val: form.find("input[name=service_name]").val(),
        }, {
            key: "webui_banner_title",
            val: form.find("input[name=webui_banner_title]").val(),
        }, {
            key: "user_reg_disable",
            val: user_reg_disable,
        }],
    };

    // console.log(req);

    iam.ApiCmd("sys-config/general-set", {
        method : "PUT",
        data   : JSON.stringify(req),
        callback : function(err, data) {
            
            if (err) {
                return l4i.InnerAlert("#iam-sysmgr-generalset-alert", 'alert-danger', err);
            }

            if (!data || data.error) {
                return l4i.InnerAlert("#iam-sysmgr-generalset-alert", 'alert-danger', data.error.message);
            }

            l4i.InnerAlert("#iam-sysmgr-generalset-alert", 'alert-success', "Successfully updated");

            // window.setTimeout(function(){
            //     iamsys.GeneralSet();
            // }, 1000);
        },
    });
}


iamsys.MailerSet = function(name)
{
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function (tpl, data) {

            if (!data || !data.items) {
                return;
            }

            var mailer = JSON.parse(data.items[0].val);
            if (!mailer) {
                mailer = {}
            }

            if (!mailer.smtp_host) {
                mailer.smtp_host = "";
            }

            if (!mailer.smtp_port) {
                mailer.smtp_port = "";
            }

            if (!mailer.smtp_user) {
                mailer.smtp_user = "";
            }

            if (!mailer.smtp_pass) {
                mailer.smtp_pass = "";
            }

            l4iTemplate.Render({
                dstid  : "work-content",
                tplsrc : tpl,
                data   : mailer,
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });
    
        iam.ApiCmd("sys-config/mailer", {
            callback: ep.done('data'),
        });

        iam.TplCmd("sys-mgr/mailer-set", {
            callback: ep.done('tpl'),           
        });
    });  
}


iamsys.MailerSetCommit = function()
{
    var form = $("#iam-sysmgr-mailerset");
    
    var mailer = {
        "smtp_host": form.find("input[name=mailer_smtp_host]").val(),
        "smtp_port": form.find("input[name=mailer_smtp_port]").val(),
        "smtp_user": form.find("input[name=mailer_smtp_user]").val(),
        "smtp_pass": form.find("input[name=mailer_smtp_pass]").val(),
    };

    var req = {
        items : [{
            key: "mailer",
            val: JSON.stringify(mailer),
        }],
    };

    iam.ApiCmd("sys-config/general-set", {
        method : "PUT",
        data   : JSON.stringify(req),
        callback : function(err, data) {
            
            if (err) {
                return l4i.InnerAlert("#iam-sysmgr-mailerset-alert", 'alert-danger', err);
            }

            if (!data || data.error) {
                return l4i.InnerAlert("#iam-sysmgr-mailerset-alert", 'alert-danger', data.error.message);
            }

            l4i.InnerAlert("#iam-sysmgr-mailerset-alert", 'alert-success', "Successfully updated");

            // window.setTimeout(function(){
            //     iamsys.GeneralSet();
            // }, 1000);
        },
    });
}

