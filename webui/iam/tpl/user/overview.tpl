<div class="iam-user-overview-box">

    <div class="iam-user-panel iam-user-profile">
      <div class="iup-title">Profile</div>
      <img class="iup-photo" src="/iam/v1/service/photo/{[=it.login.name]}" />
      <ul class="iup-info">
        <li><strong>{[=it.login.display_name]}</strong></li>
        <li><a class="iam-useridx-click" href="#profile-set" onclick="iamUser.ProfileSetForm()">Change Profile</a></li>
        <li><a class="iam-useridx-click" href="#photo-set" onclick="iamUser.PhotoSetForm('{[=it.login.name]}')">Change Photo</a></li>
      </ul>
    </div>


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

    <div class="iam-user-panel iam-user-ecoin">
      <div class="iup-title">Account Balance</div>
	  <div class="ecoin-value">{[=it.login.ecoin_amount]}</div>
	</div>
</div>
