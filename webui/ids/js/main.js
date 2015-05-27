var ids = {
    base    : "/ids/",
    baseui  : "/ids/~/",
    basetpl : "/ids/~/ids/tpl/",
    api     : "/ids/v1/",
    mgrapi  : "/ids/",
}

ids.Boot = function()
{
    seajs.config({
        alias: {
            ep: ids.baseui +"lessui/js/eventproxy.js"
        },
    });

    seajs.use([
        ids.baseui +"twitter-bootstrap/css/bootstrap.min.css",
        ids.baseui +"jquery/jquery.min.js",
        ids.baseui +"lessui/js/BrowserDetect.js",
    ], function() {

        var browser = BrowserDetect.browser;
        var version = BrowserDetect.version;
        var OS      = BrowserDetect.OS;

        if (!((browser == 'Chrome' && version >= 20)
            || (browser == 'Firefox' && version >= 3.6)
            || (browser == 'Safari' && version >= 5.0 && OS == 'Mac'))) {

            $('body').load(ids.basetpl +"error/browser.tpl");
            return;
        }

        seajs.use([
            ids.baseui +"lessui/js/lessui.js",
            ids.baseui +"lessui/css/lessui.min.css",
        ], function() {
            
            seajs.use([
                ids.baseui +"ids/css/main.css?_="+ Math.random(),
                ids.baseui +"twitter-bootstrap/js/bootstrap.min.js",
                ids.baseui +"ids/js/mgr.js?_="+ Math.random(),
                ids.baseui +"ids/js/user.js?_="+ Math.random(),
                ids.baseui +"ids/js/myapp.js?_="+ Math.random(),
                ids.baseui +"ids/js/sys.js?_="+ Math.random(),
                ids.baseui +"ids/js/usermgr.js?_="+ Math.random(),
                ids.baseui +"ids/js/appmgr.js?_="+ Math.random(),
            ], function() {
                idsmgr.Index();
                
                idsuser.Init();
                idsmyapp.Init();

                idssys.Init();
                idsusrmgr.Init();
                idsappmgr.Init();
            });
        });        
    });
}

ids.Ajax = function(url, options)
{
    options = options || {};

    //
    if (url.substr(0, 1) != "/" && url.substr(0, 4) != "http") {
        url = bt.base + url;
    }

    //
    if (/\?/.test(url)) {
        url += "&_=";
    } else {
        url += "?_=";
    }
    url += Math.random();

    //
    url += "&access_token="+ l4iCookie.Get("access_token");

    //
    if (options.method === undefined) {
        options.method = "GET";
    }

    //
    if (options.timeout === undefined) {
        options.timeout = 10000;
    }

    //
    $.ajax({
        url     : url,
        type    : options.method,
        data    : options.data,
        timeout : options.timeout,
        success : function(rsp) {
            if (typeof options.callback === "function") {
                options.callback(null, rsp);
            }
            if (typeof options.success === "function") {
                options.success(rsp);
            }
        },
        error: function(xhr, textStatus, error) {
            // console.log(xhr.responseText);
            if (typeof options.callback === "function") {
                options.callback(xhr.responseText, null);
            }
            if (typeof options.error === "function") {
                options.error(xhr, textStatus, error);
            }
        }
    });
}


ids.ApiCmd = function(url, options)
{
    ids.Ajax(ids.api + url, options);
}

ids.MgrApiCmd = function(url, options)
{
    ids.Ajax(ids.mgrapi + url, options);
}

ids.TplCmd = function(url, options)
{
    ids.Ajax(ids.basetpl + url +".tpl", options);
}

ids.Loader = function(target, uri)
{
    ids.Ajax(ids.basetpl + uri +".tpl", {
        callback: function(err, data) {
            $(target).html(data);
        }
    });
}

ids.BodyLoader = function(uri)
{
    ids.Loader("#body-content", uri);
}

ids.ComLoader = function(uri)
{
    ids.Loader("#com-content", uri);
}

ids.WorkLoader = function(uri)
{
    ids.Loader("#work-content", uri);
}
