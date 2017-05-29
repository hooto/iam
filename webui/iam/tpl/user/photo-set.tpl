<style>
#iam-user-photo-set-file {
    border: 1px solid #ccc;
}
</style>

<div id="iam-user-photo-set-alert" class="alert alert-info">
    <p>You must upload a JPG, GIF, or PNG file</p>
</div>

<div id="iam-user-photo-set">
<table>
  <tr>
    <td valign="bottom">
        <img src="/iam/v1/service/photo/{[=it.username]}" width="96" height="96" >
        <b>Normal size</b>
    </td>
    <td width="30px"></td>
    <td valign="bottom">
        <img src="/iam/v1/service/photo/{[=it.username]}" width="48" height="48" /> <b>Small size</b>
    </td>
  </tr>
</table>
<br />
<form class="oqmbg4" enctype="multipart/form-data" action="#" method="post">
  <table width="100%">

    <tr>
      <td width="160px"><b>Select a New Picture</b></td>
      <td><input id="iam-user-photo-set-file" name="attachment" size="20" type="file" class="btn" /></td>
    </tr>

  </table>
</form>
</div>
