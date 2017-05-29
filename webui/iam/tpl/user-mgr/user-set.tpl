<style>
.form-horizontal {
  margin: 0 15px;
  padding: 2px;
}
.form-group {
  margin-bottom: 5px;
}
</style>

<div id="iam-usermgr-userset-alert"></div>

<div id="iam-usermgr-userset" class="form-horizontal">

    <label class="iam-form-group-title">Login Information (Required)</label>

    {[ if (it.login.name == "") { ]}
    <div class="form-group">
      <label class="col-sm-2 control-label">Username</label>
      <div class="col-sm-10">
        <input type="text" class="form-control input-sm" name="login_name" value="{[=it.login.name]}">
      </div>
    </div>
    {[ } else {]}
    <input type="hidden" name="login_name" value="{[=it.login.name]}">
    {[ } ]}

    <div class="form-group">
      <label class="col-sm-2 control-label">Email</label>
      <div class="col-sm-10">
        <input type="text" class="form-control input-sm" name="login_email" value="{[=it.login.email]}">
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-2 control-label">Password</label>
      <div class="col-sm-10">
        <input type="text" class="form-control input-sm" name="login_auth" value="{[=it.login._auth]}">
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-2 control-label">Roles</label>
      <div class="col-sm-10">
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
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-2 control-label">Status</label>
      <div class="col-sm-10">
        {[~it._statusls :v]}
          <span class="iam-form-checkbox">
            <input type="radio" name="login_status" value="{[=v.status]}" {[ if (v.status == it.login.status) { ]}checked="checked"{[ } ]}> {[=v.title]}
          </span>
        {[~]}
      </div>
    </div>

    <label class="iam-form-group-title">Profile Information (Optional)</label>

    <div class="form-group">
      <label class="col-sm-2 control-label">Display Name</label>
      <div class="col-sm-10">
        <input type="text" class="form-control input-sm" name="login_display_name" value="{[=it.login.display_name]}">
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-2 control-label">Birthday</label>
      <div class="col-sm-10">
        <input type="text" class="form-control input-sm" name="profile_birthday" placeholder="Example : 1970-01-01" value="{[=it.profile.birthday]}">
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-2 control-label">About</label>
      <div class="col-sm-10">
        <textarea class="form-control input-sm" rows="3" name="profile_about">{[=it.profile.about]}</textarea>
      </div>
    </div>

</div>
