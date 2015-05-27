<table width="100%" style="margin-top:10px">
<tr>

  <td width="40%">
    <div class="ids-user-panel ids-user-profile">
      <div class="iup-title">Profile</div>
      <img class="iup-photo" src="/ids/v1/service/photo/{[=it.login.meta.id]}" />
      <ul class="iup-info">
        <li><strong>{[=it.login.name]}</strong></li>
        <li><a class="ids-useridx-click" href="#profile-set" onclick="idsuser.ProfileSetForm()">Change Profile</a></li>
        <li><a class="ids-useridx-click" href="#photo-set" onclick="idsuser.PhotoSetForm('{[=it.login.meta.id]}')">Change Photo</a></li>
      </ul>
    </div>
  </td>

  <td width="20px"></td>

  <td>
    <div class="ids-user-panel ids-user-personal">
      <div class="iup-title">Personal Settings</div>
      <table> 
    
        <tr> 
          <td class="iup-subtitle">Security</td> 
          <td> 
            <ul> 
              <li><a class="ids-useridx-click" href="#pass-set" onclick="idsuser.PassSetForm()">Change Password</a></li>
            </ul>
          </td> 
        </tr>
    
        <tr> 
          <td class="iup-subtitle">Email</td> 
          <td>
            <ul> 
              <li>{[=it.login.email]}</li> 
              <li><a class="ids-useridx-click" href="#email-set" onclick="idsuser.EmailSetForm()">Change</a></li> 
            </ul>
          </td>
        </tr> 
    
      </table> 
    </div>
  </td>

</tr>
</table>
