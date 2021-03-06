var iam = {
    version: "0.2",
    base: "/iam/",
    baseui: "/iam/~/",
    basetpl: "/iam/~/iam/tpl/",
    api: "/iam/v1/",
    mgrapi: "/iam/",
    Session: null,
    OpToolActive: null,
}

iam.hash_uri = function() {
    return "?_=" + iam.version;
}

iam.Boot = function() {
    seajs.config({
        base: iam.base,
        alias: {
            ep: "~/lessui/js/eventproxy.js"
        },
    });

    seajs.use([
        "~/bs/4/css/bootstrap.css",
        "~/jquery/jquery.js",
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
            "~/fa/css/fa.css",
        ], function() {

            seajs.use([
                "~/iam/css/main.css" + iam.hash_uri(),
                // "~/bs/4/js/bootstrap.js",
                "~/iam/js/mgr.js" + iam.hash_uri(),
                "~/iam/js/user.js" + iam.hash_uri(),
                "~/iam/js/user-group.js" + iam.hash_uri(),
                "~/iam/js/app.js" + iam.hash_uri(),
                "~/iam/js/access-key.js" + iam.hash_uri(),
                "~/iam/js/account.js" + iam.hash_uri(),
                "~/iam/js/sys.js" + iam.hash_uri(),
                "~/iam/js/sys-msg.js" + iam.hash_uri(),
                "~/iam/js/user-mgr.js" + iam.hash_uri(),
                "~/iam/js/app-mgr.js" + iam.hash_uri(),
                "~/iam/js/account-mgr.js" + iam.hash_uri(),
            ], iam.load_index);
        });
    });
}

iam.load_index = function() {

    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create("tpl", "session", "pinfo", "lang", function(tpl, session, pinfo, lang) {

            if (lang && lang.items) {
                l4i.LangSync(lang.items, lang.locale);
            }

            if (!session || session.username == "") {
                return alert(l4i.T("network error, please try again later"));
            }
            iam.Session = session;

            if (!pinfo.topnav) {
                pinfo.topnav = [];
            }
            if (!pinfo.webui_banner_title) {
                pinfo.webui_banner_title = l4i.T("User Center");
            }

            l4i.UrlEventRegister("index", iamUser.Overview, "iam-topbar-nav-menus");

            for (var i in pinfo.topnav) {
                switch (pinfo.topnav[i].path) {
                    case "app/index":
                        l4i.UrlEventRegister("app/index", iamApp.Index, "iam-topbar-nav-menus");
                        break;

                    case "access-key/index":
                        l4i.UrlEventRegister("access-key/index", iamAccessKey.Index, "iam-topbar-nav-menus");
                        break;

                    case "account/index":
                        l4i.UrlEventRegister("account/index", iamAcc.Index, "iam-topbar-nav-menus");
                        break;

                    case "user-group/index":
                        l4i.UrlEventRegister("user-group/index", iamUserGroup.Index, "iam-topbar-nav-menus");
                        break;

                    case "user-mgr/index":
                        l4i.UrlEventRegister("user-mgr/index", iamUserMgr.Index, "iam-topbar-nav-menus");
                        break;

                    case "acc-mgr/index":
                        l4i.UrlEventRegister("acc-mgr/index", iamAccMgr.Index, "iam-topbar-nav-menus");
                        break;

                    case "app-mgr/index":
                        l4i.UrlEventRegister("app-mgr/index", iamAppMgr.Index, "iam-topbar-nav-menus");
                        break;

                    case "sys-mgr/index":
                        l4i.UrlEventRegister("sys-mgr/index", iamSys.Index, "iam-topbar-nav-menus");
                        break;

                    case "sys-msg/index":
                        l4i.UrlEventRegister("sys-msg/index", iamSysMsg.Index, "iam-topbar-nav-menus");
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
                callback: function() {
                    l4i.UrlEventHandler("index", true);
                },
            });
        });

        ep.fail(function(err) {
            if (err && err == "AuthSession") {
                iam.AlertUserLogin();
            } else {
                alert(l4i.T("network error, please try again later"));
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

        iam.ApiCmd("langsrv/locale", {
            callback: ep.done("lang"),
        });

        iam.TplCmd("index", {
            callback: ep.done("tpl"),
        });
    });
}

iam.AlertUserLogin = function() {
    l4iAlert.Open("warn",
        l4i.T("You are not logged in, or your login session has expired. Please sign in again"), {
            close: false,
            buttons: [{
                title: l4i.T("SIGN IN"),
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

iam.TplPath = function(url) {
    return iam.basetpl + url + ".tpl";
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
            l4iTemplate.Render({
                dstid: "iam-module-navbar-optools",
                tplsrc: opt.html(),
                data: {},
            });
            iam.OpToolActive = div_target;
        }
    }
}

iam.OpToolsClean = function() {
    $("#iam-module-navbar-optools").empty();
    iam.OpToolActive = null;
}
