<div id="iam-usermgr-userset-alert"></div>

<div id="iam-usermgr-userset" class="form-horizontal">

  <div class="iam-form-group-title">Login Information (Required)</div>

  <table class="iam-formtable">
    <tbody>
 
    <tr>
      <td width="200px">Username</td>
      <td>
        <input type="text" class="form-control" name="login_name" value="{[=it.login.name]}" {[? it.login.name.length > 0]}readonly{[?]}>
      </td>
    </tr>

    <tr>
      <td width="200px">Email</td>
      <td>
        <input type="text" class="form-control" name="login_email" value="{[=it.login.email]}">
      </td>
    </tr>

    <tr>
      <td>Password</td>
      <td>
        <input type="text" class="form-control" name="login_auth" value="{[=it.login._auth]}">
      </td>
    </tr>

    <tr>
      <td>Roles</td>
      <td>
        {[~it._roles.items :v]}
        <span class="iam-form-checkbox">
          {[ if (v.id == 100) { ]}
            <input type="checkbox" name="login_roles" value="{[=v.id]}" checked="checked" onclick="return false"> {[=v.name]}
          {[ } else if (v.checked) { ]}
            <input type="checkbox" name="login_roles" value="{[=v.id]}" checked="checked"> {[=v.name]}
          {[ } else { ]}
            <input type="checkbox" name="login_roles" value="{[=v.id]}"> {[=v.name]}
          {[ } ]}
        </span>
        {[~]}
      </td>
    </tr>

    <tr>
      <td>Status</td>
      <td>
        {[~it._statusls :v]}
          <span class="iam-form-checkbox">
            <input type="radio" name="login_status" value="{[=v.status]}" {[ if (v.status == it.login.status) { ]}checked="checked"{[ } ]}> {[=v.title]}
          </span>
        {[~]}
      </td>
    </tr>
    </tbody>
  </table>

  <div class="iam-form-group-title">Profile Information (Optional)</div>

  <table class="iam-formtable">
    <tbody>
    <tr>
      <td width="200px">Display Name</td>
      <td>
        <input type="text" class="form-control" name="login_display_name" value="{[=it.login.display_name]}">
      </td>
    </tr>

    <tr>
      <td>Birthday</td>
      <td>
        <input type="text" class="form-control" name="profile_birthday" placeholder="Example : 1970-01-01" value="{[=it.profile.birthday]}">
      </td>
    </tr>

    <tr>
      <td>About</td>
      <td>
        <textarea class="form-control" rows="3" name="profile_about">{[=it.profile.about]}</textarea>
      </td>
    </tr>
    </tbody>
  </table>

</div>

