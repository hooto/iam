<div class="iam-user-overview-box">

    <div class="iam-user-panel iam-user-profile">
      <div class="iup-title">{[=l4i.T("Personal Info")]}</div>
      <img class="iup-photo" src="/iam/v1/service/photo/{[=it.login.name]}" />
      <ul class="iup-info">
        <li><strong>{[=it.login.display_name]}</strong></li>
        <li><a class="iam-useridx-click" href="#profile-set" onclick="iamUser.ProfileSetForm()">{[=l4i.T("Change Profile")]}</a></li>
        <li><a class="iam-useridx-click" href="#photo-set" onclick="iamUser.PhotoSetForm('{[=it.login.name]}')">{[=l4i.T("Change Photo")]}</a></li>
      </ul>
    </div>


    <div class="iam-user-panel iam-user-personal">
      <div class="iup-title">{[=l4i.T("Security Settings")]}</div>
      <table>

        <tr>
          <td class="iup-subtitle">{[=l4i.T("Password")]}</td>
          <td>
            <ul>
              <li><a class="iam-useridx-click" href="#pass-set" onclick="iamUser.PassSetForm()">{[=l4i.T("Change %s", l4i.T("Password"))]}</a></li>
            </ul>
          </td>
        </tr>

        <tr>
          <td class="iup-subtitle">{[=l4i.T("Email")]}</td>
          <td>
            <ul>
              {[if (it.login.email) {]}
              <li>{[=it.login.email]}</li>
              {[}]}
              <li><a class="iam-useridx-click" href="#email-set" onclick="iamUser.EmailSetForm()">{[=l4i.T("Change")]}</a></li>
            </ul>
          </td>
        </tr>

      </table>
    </div>

    <div class="iam-user-panel iam-user-acc-balance">
      <div class="iup-title">{[=l4i.T("Account Balance")]}</div>
	  <div class="acc-balance-value">{[=it.account.balance]}</div>
	</div>
</div>
