<div id="iam-user-profile-set-alert" class="alert hide"></div>

<table id="iam-user-profile-set" class="iam-formtable">
  <tr>
    <td width="200px">Display Name</td>
    <td>
      <input name="display_name" type="text" class="form-control" value="{[=it.login.display_name]}">
    </td>
  </tr>
  <tr>
    <td>Birthday</td>
    <td>
      <input name="birthday" type="text" class="form-control" value="{[=it.birthday]}">
      <small class="form-text text-muted">Example : 1970-01-01</small>
    </td>
  </tr>
  <tr>
    <td>About me</td>
    <td>
      <textarea name="about" class="form-control" rows="3">{[=it.about]}</textarea>
    </td>
  </tr>
</table>
