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
  <img class="iup-photo" src="/ids/service/photo/{{.login_uid}}" />
  <ul class="iup-info">
    <li>{{T . "Name"}}: <strong>{{.login_name}}</strong></li>
    <li><a class="ids-useridx-click" href="#profile-set">{{T . "Change Profile"}}</a></li>
    <li><a class="ids-useridx-click" href="#photo-set">{{T . "Change Photo"}}</a></li>
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
          <li><a class="ids-useridx-click" href="#pass-set">{{T . "Change Password"}}</a></li>
        </ul>
      </td> 
    </tr>

    <tr> 
      <td class="iup-subtitle">{{T . "Email"}}</td> 
      <td>
        <ul> 
          <li>{{.login_email}}</li> 
          <li><a class="ids-useridx-click" href="#email-set">{{T . "Change"}}</a></li> 
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

$(".ids-useridx-click").click(function(){
    var uri = $(this).attr("href").substr(1);

    switch (uri) {
    case "pass-set":
        lessModalOpen("/ids/user/"+ uri, 1, 500, 350, "{{T . "Change Password"}}", null);
        break;
    case "email-set":
        lessModalOpen("/ids/user/"+ uri, 1, 500, 300, "{{T . "Change Email"}}", null);
        break;
    case "profile-set":
        lessModalOpen("/ids/user/"+ uri, 1, 700, 400, "{{T . "Change Profile"}}", null);
        break;
    case "photo-set":
        lessModalOpen("/ids/user/"+ uri, 1, 600, 400, "{{T . "Change Photo"}}", null);
        break;
    }
    
});

</script>

{{template "Common/HtmlFooter.tpl" .}}
