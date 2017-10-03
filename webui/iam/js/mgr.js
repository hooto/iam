var iamMgr = {

}

iamMgr.NavAction = function(uri) {
    switch (uri) {
        case "user/my":
        case "user-mgr/index":
        case "auth-mgr/index":
        case "sys-mgr/index":
        case "app-mgr/index":
        case "app/my":
            $(".iuh-menu a.active").removeClass('active');
            iam.ComLoader(uri);
            $(".iuh-menu").find("a[href='#" + uri + "']").addClass("active");
            break;
    }

// $(".iuh-menu a").click(function(event) {
//     event.preventDefault();
//     var uri = $(this).attr("href").substr(1);
//     _user_menugo(uri);
// });
}

iamMgr.Index = function() {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            l4iTemplate.Render({
                dstid: "body-content",
                tplsrc: tpl,
                data: data,
                success: function() {
                    iamUser.Overview();
                },
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        iam.MgrApiCmd("/user/panel-info", {
            callback: ep.done('data'),
        });

        iam.TplCmd("user/well", {
            callback: ep.done('tpl'),
        });
    });
}
