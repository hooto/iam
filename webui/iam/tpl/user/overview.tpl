<table width="100%" style="margin-top:30px">
<tr>

  <td width="40%">
    <div class="iam-user-panel iam-user-profile">
      <div class="iup-title">Profile</div>
      <img class="iup-photo" src="/iam/v1/service/photo/{[=it.login.meta.id]}" />
      <ul class="iup-info">
        <li><strong>{[=it.login.name]}</strong></li>
        <li><a class="iam-useridx-click" href="#profile-set" onclick="iamUser.ProfileSetForm()">Change Profile</a></li>
        <li><a class="iam-useridx-click" href="#photo-set" onclick="iamUser.PhotoSetForm('{[=it.login.meta.id]}')">Change Photo</a></li>
      </ul>
    </div>
  </td>

  <td width="20px"></td>

  <td>
    <div class="iam-user-panel iam-user-personal">
      <div class="iup-title">Personal Settings</div>
      <table> 
    
        <tr> 
          <td class="iup-subtitle">Security</td> 
          <td> 
            <ul> 
              <li><a class="iam-useridx-click" href="#pass-set" onclick="iamUser.PassSetForm()">Change Password</a></li>
            </ul>
          </td> 
        </tr>
    
        <tr> 
          <td class="iup-subtitle">Email</td> 
          <td>
            <ul> 
              {[if (it.login.email) {]}
              <li>{[=it.login.email]}</li> 
              {[}]}
              <li><a class="iam-useridx-click" href="#email-set" onclick="iamUser.EmailSetForm()">Change</a></li> 
            </ul>
          </td>
        </tr> 
    
      </table> 
    </div>
  </td>

</tr>
</table>
