var iam = {
    base    : "/iam/",
    baseui  : "/iam/~/",
    basetpl : "/iam/~/iam/tpl/",
    api     : "/iam/v1/",
    mgrapi  : "/iam/",
}

iam.Boot = function()
{
    seajs.config({
        alias: {
            ep: iam.baseui +"lessui/js/eventproxy.js"
        },
    });

    seajs.use([
        iam.baseui +"twitter-bootstrap/css/bootstrap.min.css",
        iam.baseui +"jquery/jquery.min.js",
        iam.baseui +"lessui/js/BrowserDetect.js",
    ], function() {

        var browser = BrowserDetect.browser;
        var version = BrowserDetect.version;
        var OS      = BrowserDetect.OS;

        if (!((browser == 'Chrome' && version >= 20)
            || (browser == 'Firefox' && version >= 3.6)
            || (browser == 'Safari' && version >= 5.0 && OS == 'Mac'))) {

            $('body').load(iam.basetpl +"error/browser.tpl");
            return;
        }

        seajs.use([
            iam.baseui +"lessui/js/lessui.js",
            iam.baseui +"lessui/css/lessui.min.css",
        ], function() {
            
            seajs.use([
                iam.baseui +"iam/css/main.css?_="+ Math.random(),
                iam.baseui +"twitter-bootstrap/js/bootstrap.min.js",
                iam.baseui +"iam/js/mgr.js?_="+ Math.random(),
                iam.baseui +"iam/js/user.js?_="+ Math.random(),
                iam.baseui +"iam/js/myapp.js?_="+ Math.random(),
                iam.baseui +"iam/js/sys.js?_="+ Math.random(),
                iam.baseui +"iam/js/usermgr.js?_="+ Math.random(),
                iam.baseui +"iam/js/appmgr.js?_="+ Math.random(),
            ], function() {
                iammgr.Index();
                
                iamuser.Init();
                iammyapp.Init();

                iamsys.Init();
                iamusrmgr.Init();
                iamappmgr.Init();
            });
        });        
    });
}

iam.Ajax = function(url, options)
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
    if (l4iCookie.Get("access_token")) {
        url += "&access_token="+ l4iCookie.Get("access_token");
    }
    
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


iam.ApiCmd = function(url, options)
{
    iam.Ajax(iam.api + url, options);
}

iam.MgrApiCmd = function(url, options)
{
    iam.Ajax(iam.mgrapi + url, options);
}

iam.TplCmd = function(url, options)
{
    iam.Ajax(iam.basetpl + url +".tpl", options);
}

iam.Loader = function(target, uri)
{
    iam.Ajax(iam.basetpl + uri +".tpl", {
        callback: function(err, data) {
            $(target).html(data);
        }
    });
}

iam.BodyLoader = function(uri)
{
    iam.Loader("#body-content", uri);
}

iam.ComLoader = function(uri)
{
    iam.Loader("#com-content", uri);
}

iam.WorkLoader = function(uri)
{
    iam.Loader("#work-content", uri);
}
