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
    <input type="hidden" name="userid" value="{[=it.meta.id]}">
    
    <label class="iam-form-group-title">Login Information (Required)</label>

    {[ if (it.meta.id == "") { ]}
    <div class="form-group">
      <label class="col-sm-2 control-label">Username</label>
      <div class="col-sm-10">
        <input type="text" class="form-control input-sm" name="username" value="{[=it.meta.name]}">
      </div>
    </div>
    {[ } ]}

    <div class="form-group">
      <label class="col-sm-2 control-label">Email</label>
      <div class="col-sm-10">
        <input type="text" class="form-control input-sm" name="email" value="{[=it.email]}">
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-2 control-label">Password</label>
      <div class="col-sm-10">
        <input type="text" class="form-control input-sm" name="auth" value="{[=it.auth]}">
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-2 control-label">Roles</label>
      <div class="col-sm-10">
        {[~it._roles.items :v]}
        <span class="iam-form-checkbox">
          {[ if (v.id == 100) { ]}
            <input type="checkbox" name="roles" value="{[=v.id]}" checked="checked" onclick="return false"> {[=v.meta.name]}
          {[ } else if (v.checked) { ]}
            <input type="checkbox" name="roles" value="{[=v.id]}" checked="checked"> {[=v.meta.name]}
          {[ } else { ]}
            <input type="checkbox" name="roles" value="{[=v.id]}"> {[=v.meta.name]}
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
            <input type="radio" name="status" value="{[=v.status]}" {[ if (v.status == it.status) { ]}checked="checked"{[ } ]}> {[=v.title]}
          </span>
        {[~]}
      </div>
    </div>

    <label class="iam-form-group-title">Profile Information (Optional)</label>

    <div class="form-group">
      <label class="col-sm-2 control-label">Nickname</label>
      <div class="col-sm-10">
        <input type="text" class="form-control input-sm" name="name" value="{[=it.name]}">
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-2 control-label">Birthday</label>
      <div class="col-sm-10">
        <input type="text" class="form-control input-sm" name="birthday" placeholder="Example : 1970-01-01" value="{[=it.profile.birthday]}">
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-2 control-label">About</label>
      <div class="col-sm-10">
        <textarea class="form-control input-sm" rows="3" name="about">{[=it.profile.about]}</textarea>
      </div>
    </div>

</div>
