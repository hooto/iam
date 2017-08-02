var iam = {
    base: "/iam/",
    baseui: "/iam/~/",
    basetpl: "/iam/~/iam/tpl/",
    api: "/iam/v1/",
    mgrapi: "/iam/",
    Session: null,
    debug: true,
    OpToolActive: null,
}

iam.debug_uri = function() {
    if (!iam.debug) {
        return "";
    }
    return "?_=" + Math.random();
}

iam.Boot = function() {
    seajs.config({
        base: iam.base,
        alias: {
            ep: "~/lessui/js/eventproxy.js"
        },
    });

    seajs.use([
        "~/twbs/css/bootstrap.min.css",
        "~/jquery/jquery.min.js",
        "~/lessui/js/browser-detect.js",
    ], function() {

        var browser = BrowserDetect.browser;
        var version = BrowserDetect.version;
        var OS = BrowserDetect.OS;

        if (!((browser == 'Chrome' && version >= 20)
            || (browser == 'Firefox' && version >= 3.6)
            || (browser == 'Safari' && version >= 5.0 && OS == 'Mac'))) {

            $('body').load(iam.basetpl + "error/browser.tpl");
            return;
        }

        seajs.use([
            "~/lessui/js/lessui.js",
            "~/lessui/css/lessui.css",
            "~/purecss/css/pure.css",
        ], function() {

            seajs.use([
                "~/iam/css/main.css" + iam.debug_uri(),
                "~/twbs/js/bootstrap.min.js",
                "~/iam/js/mgr.js" + iam.debug_uri(),
                "~/iam/js/user.js" + iam.debug_uri(),
                "~/iam/js/myapp.js" + iam.debug_uri(),
                "~/iam/js/sys.js" + iam.debug_uri(),
                "~/iam/js/usermgr.js" + iam.debug_uri(),
                "~/iam/js/appmgr.js" + iam.debug_uri(),
                "~/iam/js/access-key.js" + iam.debug_uri(),
            ], iam.load_index);
        });
    });
}

iam.load_index = function() {
    l4i.debug = iam.debug;

    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create("tpl", "session", "pinfo", function(tpl, session, pinfo) {

            if (!session || session.username == "") {
                return alert("Network is unreachable, Please try again later");
            }
            iam.Session = session;

            if (!pinfo.topnav) {
                pinfo.topnav = [];
            }
            if (!pinfo.webui_banner_title) {
                pinfo.webui_banner_title = "Account Center";
            }

            // console.log(pinfo);

            l4i.UrlEventRegister("index", iamUser.Overview, "iam-topbar-nav-menus");

            for (var i in pinfo.topnav) {
                switch (pinfo.topnav[i].path) {
                    case "my-app/index":
                        l4i.UrlEventRegister("my-app/index", iamMyApp.Index, "iam-topbar-nav-menus");
                        break;

                    case "access-key/index":
                        l4i.UrlEventRegister("access-key/index", iamAccessKey.Index, "iam-topbar-nav-menus");
                        break;

                    case "user-mgr/index":
                        l4i.UrlEventRegister("user-mgr/index", iamUserMgr.Index, "iam-topbar-nav-menus");
                        break;

                    case "app-mgr/index":
                        l4i.UrlEventRegister("app-mgr/index", iamAppMgr.Index, "iam-topbar-nav-menus");
                        break;

                    case "sys-mgr/index":
                        l4i.UrlEventRegister("sys-mgr/index", iamSys.Index, "iam-topbar-nav-menus");
                        break;
                }
            }

            l4iTemplate.Render({
                dstid: "body-content",
                tplsrc: tpl,
                data: {
                    pinfo: pinfo,
                    session: session,
                },
                success: function() {
                    l4i.UrlEventHandler("index", true);
                },
            });
        });

        ep.fail(function(err) {
            if (err && err == "AuthSession") {
                iam.AlertUserLogin();
            } else {
                alert("Network is unreachable, Please try again later");
            }
        });

        l4i.Ajax(iam.base + "auth/session", {
            nocache: true,
            callback: function(err, data) {
                if (!data || data.kind != "AuthSession") {
                    return ep.emit('error', "AuthSession");
                }
                ep.emit("session", data);
            },
        });

        iam.MgrApiCmd("user/panel-info", {
            callback: ep.done('pinfo'),
        });

        iam.TplCmd("index", {
            callback: ep.done("tpl"),
        });
    });
}

iam.AlertUserLogin = function() {
    l4iAlert.Open("warn", "You are not logged in, or your login session has expired. Please sign in again", {
        close: false,
        buttons: [{
            title: "SIGN IN",
            href: iam.base + "auth/login",
        }],
    });
}

iam.ApiCmd = function(url, options) {
    iam.api_cmd(iam.api + url, options);
}

iam.MgrApiCmd = function(url, options) {
    iam.api_cmd(iam.mgrapi + url, options);
}

iam.api_cmd = function(url, options) {
    var appcb = null;
    if (options.callback) {
        appcb = options.callback;
    }
    options.callback = function(err, data) {
        if (err == "Unauthorized") {
            return iam.AlertUserLogin();
        }
        if (appcb) {
            appcb(err, data);
        }
    }
    options.nocache = true;

    l4i.Ajax(url, options);
}

iam.TplCmd = function(url, options) {
    l4i.Ajax(iam.basetpl + url + ".tpl", options);
}

iam.Loader = function(target, uri) {
    l4i.Ajax(iam.basetpl + uri + ".tpl", {
        callback: function(err, data) {
            $(target).html(data);
        }
    });
}

iam.BodyLoader = function(uri) {
    iam.Loader("#body-content", uri);
}

iam.ComLoader = function(uri) {
    iam.Loader("#com-content", uri);
}

iam.WorkLoader = function(uri) {
    iam.Loader("#work-content", uri);
}

iam.OpToolsRefresh = function(div_target) {
    if (!div_target || typeof div_target == "string" && div_target == iam.OpToolActive) {
        return;
    }

    $("#iam-module-navbar-optools").empty();

    if (typeof div_target == "string") {

        var opt = $("#work-content").find(div_target);
        if (opt) {
            $("#iam-module-navbar-optools").html(opt.html());
            iam.OpToolActive = div_target;
        }
    }
}
