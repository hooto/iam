
function lfBoot()
{
    seajs.config({
        base: "/ids/",
    });

    var rqs = [
        "~/twitter-bootstrap/3.1.1/css/bootstrap.min.css",
        "~/jquery/1.10.2/jquery.min.js",
        "~/lessui/master/js/BrowserDetect.js",
    ];
    seajs.use(rqs, function() {

        var browser = BrowserDetect.browser;
        var version = BrowserDetect.version;
        var OS      = BrowserDetect.OS;

        if (!((browser == 'Chrome' && version >= 20)
            || (browser == 'Firefox' && version >= 3.6)
            || (browser == 'Safari' && version >= 5.0 && OS == 'Mac'))) {

            $('#body-content').load('/ids/error/browser');
            return;
        }

        rqs = [
            "~/lessui/master/css/lessui.min.css",
            "static/css/main.css?_="+ Math.random(),
            "~/twitter-bootstrap/3.1.1/js/bootstrap.min.js",
            "~/lessui/master/js/lessui.js",
        ];
        seajs.use(rqs, function() {
            lfPageWell();
        });        
    });
}

function lfPageWell()
{
    lfBodyLoader("index/well");
}

function lfAjax(obj, url)
{
    if (/\?/.test(url)) {
        url += "&_=";
    } else {
        url += "?_=";
    }
    url += Math.random();
    //console.log("req: ids/"+ url);
    $.ajax({
        url     : "/ids/"+ url,
        type    : "GET",
        timeout : 30000,
        success : function(rsp) {
            //console.log(rsp);
            $(obj).html(rsp);
        },
        error: function(xhr, textStatus, error) {

            if (xhr.status == 401) {
                console.log("access denied");
                //lfBodyLoader('user/login');
            } else {
                alert("Internal Server Error"); //+ xhr.responseText);
            }
        }
    });

    /*var uris = url.split("/");
    switch (uris[0]) {
    case "app":
    case "keeper":
    case "user":
        $(".p84ykc li.active").removeClass('active');
        $(".p84ykc a[href^='#"+uris[0]+"']").parent().addClass("active");
        break;
    }*/
}

function lfBodyLoader(uri)
{
    lfAjax("#body-content", uri);
}

function lfComLoader(uri)
{
    lfAjax("#com-content", uri);
}

function lfWorkLoader(uri)
{
    lfAjax("#work-content", uri);
}
