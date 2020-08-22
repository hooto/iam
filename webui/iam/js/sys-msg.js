var iamSysMsg = {
    data: null,
    actionls: [{
        action: 1 << 1,
        title: "OK",
    }, {
        action: 1 << 2,
        title: "Error",
    }, {
        action: 1 << 3,
        title: "Timeout",
    }],
}

iamSysMsg.Index = function() {
    l4iTemplate.Render({
        dstid: "com-content",
        tplurl: iam.TplPath("sys-msg/index"),
        data: {},
        callback: function() {
            iam.OpToolsClean();
            l4i.UrlEventClean("iam-module-navbar-menus");
            l4i.UrlEventRegister("sys-msg/list", iamSysMsg.List, "iam-module-navbar-menus");
            l4i.UrlEventHandler("sys-msg/list", true);
        },
    });
}

iamSysMsg.List = function() {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            if (!data) {
                return;
            }

            if (!data.items) {
                data.items = [];
            }

            // $("#work-content").html(tpl);


            for (var i in data.items) {

                if (!data.items[i].to_email) {
                    data.items[i].to_email = "";
                }
            }

            data._actionls = iamSysMsg.actionls;

            l4iTemplate.Render({
                dstid: "work-content",
                tplsrc: tpl,
                data: data,
                callback: function() {
                    iam.OpToolsRefresh("#iam-msglist-optools");
                },
            });
        });

        ep.fail(function(err) {
            alert(l4i.T("network error, please try again later"));
        });

        iam.ApiCmd("sys-msg/list", {
            callback: ep.done('data'),
        });

        iam.TplCmd("sys-msg/list", {
            callback: ep.done('tpl'),
        });
    });
}

iamSysMsg.Info = function(msg_id) {
    if (!msg_id) {
        return;
    }
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            if (!data || !data.kind) {
                return;
            }

            if (!data.to_email) {
                data.to_email = "";
            }

            data._actionls = iamSysMsg.actionls;

            l4iModal.Open({
                tplsrc: tpl,
                width: 1600,
                height: 800,
                data: data,
                title: l4i.T("Message"),
                buttons: [{
                    title: l4i.T("Cancel"),
                    onclick: "l4iModal.Close()",
                }],
            });
        });

        ep.fail(function(err) {
            alert(l4i.T("network error, please try again later"));
        });

        iam.ApiCmd("sys-msg/item?id=" + msg_id, {
            callback: ep.done('data'),
        });

        iam.TplCmd("sys-msg/info", {
            callback: ep.done('tpl'),
        });
    });
}

