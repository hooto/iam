<style type="text/css">
#iam-user-profile-set td {
    padding: 5px 0;
    vertical-align: top;
}
</style>

<div id="iam-user-profile-set-alert" class="alert hide"></div>

  <table id="iam-user-profile-set" width="100%" style="">
    <tr>
      <td width="120px"><strong>Name</strong></td>
      <td>
        <input name="name" type="text" class="form-control" value="{[=it.login.name]}">
      </td>
    </tr>
    <tr>
      <td><strong>Birthday</strong></td>
      <td>
        <input name="birthday" type="text" class="form-control" value="{[=it.birthday]}">
        <div>Example : 1970-01-01</div>
      </td>
    </tr>
    <tr>
      <td valign="top"><strong>About me</strong></td>
      <td>
        <textarea name="about" class="form-control" rows="5">{[=it.about]}</textarea>
      </td>
    </tr>
  </table>
