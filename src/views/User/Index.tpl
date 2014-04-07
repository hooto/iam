{{template "Common/HtmlHeader.tpl" .}}
<div class="ids-container" style="padding:20px 0;">
<table id="ids-user-header">
  <tr>
    <td>
      <img class="ids-header-logo" src="/ids/static/img/ids-s2-32.png"> 
      <span class="ids-header-title">Account Settings</span>
    </td>
    
    <td align="right">
      <div id="ids-userbox">
        <span class="btn btn-default ids-userbox-signout">Sign Out</span>
        <img class="ids-userbox-ico" src="/ids/static/img/user-default.png">
      </div>
    </td>
  </tr>
</table>
</div>

<div class="ids-container">
<table width="100%">
<tr>
  <td width="40%">

<div class="ids-user-panel ids-user-profile">
  <div class="iup-title">{{T . "Profile"}}</div>
  <img class="iup-photo" src="/ids/static/img/user-default.png" />
  <ul class="iup-info">
    <li>{{T . "Name"}}: <strong>{{.name}}</strong></li>
    <li><a class="" href="#profile-set">{{T . "Change Profile"}}</a></li>
    <li><a class="" href="#photo-set">{{T . "Change Photo"}}</a></li>
  </ul>
</div>

  </td>
  <td width="20px"></td>
  <td>

<div class="ids-user-panel ids-user-personal">
  <div class="iup-title">{{T . "Personal Settings"}}</div>
  <table> 

    <tr> 
      <td class="iup-subtitle">{{T . "Security"}}</td> 
      <td> 
        <ul> 
          <li><a class="ids-user-pro_cli" href="#pass-set">{{T . "Change Password"}}</a></li>
        </ul>
      </td> 
    </tr> 

    <tr> 
      <td class="iup-subtitle">{{T . "Email"}}</td> 
      <td>
        <ul> 
          <li>{{.email}}</li> 
          <li><a class="ids-user-pro_cli" href="#email-set">{{T . "Change"}}</a></li> 
        </ul>
      </td>
    </tr> 

  </table> 
</div>

  </td>
</tr>
</table>

</div>


<script>

//
$("input[name=email]").focus();

//
var ids_eh = $("#ids-login-box").height();
$("#ids-login-box").css({
    "top": "50%",
    "margin-top": - (ids_eh / 2) + "px" 
});

//
$("#ids-login-form").submit(function(event) {

    event.preventDefault();

    /* var req = {
        data: {
            "email": $("input[name=email]").val(),
            "passwd": $("input[name=passwd]").val(),
            "continue": $("input[name=continue]").val(),
            "persistent": $("input[name=persistent]").val(),
        }
    } */

    $.ajax({
        type    : "POST",
        url     : "/ids/service/login-auth",
        data    : $(this).serialize(),//JSON.stringify(req),
        timeout : 3000,
        //contentType: "application/json; charset=utf-8",
        success : function(rsp) {

            var rsj = JSON.parse(rsp);
            //console.log(rsp);

            if (rsj.status != 200) {
                lessAlert("#ids-login-form-alert", 'alert-danger', rsj.message);
                return;
            }

            lessAlert("#ids-login-form-alert", 'alert-success', "Successfully Sign-on. Page redirecting");
            $("#ilf-grp-input").hide(200);
            
            lessCookie.Set("access_token", rsj.data.access_token, 7200);

            window.setTimeout(function(){    
                window.location = rsj.data.continue;
            }, 1500);
        },
        error: function(xhr, textStatus, error) {
            lessAlert("#ids-login-form-alert", 'alert-danger', '{{T . "Internal Server Error"}}');
        }
    });
});

</script>

{{template "Common/HtmlFooter.tpl" .}}
