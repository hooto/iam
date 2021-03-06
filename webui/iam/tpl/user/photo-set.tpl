<div id="iam-user-photo-set-alert" class="alert alert-info">
    <p>{[=l4i.T("You must upload a %s file", "JPG, GIF, PNG")]}</p>
</div>

<div id="iam-user-photo-set">
<form enctype="multipart/form-data" action="#" method="post">
  <table class="iam-formtable" width="100%">
    <tr>
      <td width="200px">{[=l4i.T("Photo Overview")]}</td>
      <td>
	    <img src="/iam/v1/service/photo/{[=it.username]}" width="96" height="96" >
      </td>
    </tr>
    <tr>
      <td>{[=l4i.T("Upload %s", l4i.T("Photo"))]}</td>
      <td>
        <div class="custom-file">
          <input type="file" class="custom-file-input" id="iam-user-photo-set-file">
          <label class="custom-file-label" for="customFile">{[=l4i.T("Choose file to upload")]}</label>
        </div>
      </td>
    </tr>
  </table>
</form>
</div>
