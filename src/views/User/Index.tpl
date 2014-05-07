{{template "Common/HtmlHeader.tpl" .}}

<div class="ids-container" style="padding:10px 0;">
<table id="ids-user-header">
  <tr>
    <td align="left" width="220px">
      <div class="iuh-brand">
        <img class="iuh-brand-logo" src="/ids/~/ids/img/ids-s2-32.png"> 
        <span class="iuh-brand-title">{{.webui_banner_title}}</span>
      </div>
    </td>
    <td>
      <div class="iuh-menu">
        {{range .menus}}
        <a href="{{.path}}">{{.title}}</a>
        {{end}}
      </div>
    </td>
    <td align="right">
      <div id="ids-userbox">
        <span class="btn btn-default ids-userbox-signout">Sign Out</span>
        <img class="ids-userbox-ico" src="/ids/~/ids/img/user-default.png">
      </div>
    </td>
  </tr>
</table>
</div>


<div id="com-content" class="ids-container"></div>


{{template "Common/Footer.tpl" .}}
{{template "Common/HtmlFooter.tpl" .}}

<script type="text/javascript">

function _user_menugo(uri)
{
    switch (uri) {
    case "user/my":
    case "sys-mgr/index":
    case "user-mgr/index":
        $(".iuh-menu a.active").removeClass('active');
        idsComLoader(uri);
        $(".iuh-menu").find("a[href='#"+uri+"']").addClass("active");
        break;
    }
}

$(".iuh-menu a").click(function(event) {
    event.preventDefault();
    var uri = $(this).attr("href").substr(1);
    _user_menugo(uri);
});

_user_menugo("user/my");


$(".ids-userbox-signout").click(function(){
    lessCookie.Del("access_token");
    window.setTimeout(function(){    
        window.location = "/ids/service/login";
    }, 500);

});
</script>
