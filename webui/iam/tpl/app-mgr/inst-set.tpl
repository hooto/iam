<style>
.form-horizontal {
  margin: 0 15px;
  padding: 2px;
}
.form-group {
  margin-bottom: 5px;
}
</style>

<div id="iam-appmgr-instset-alert"></div>
    
<div id="iam-appmgr-instset" class="form-horizontal" action="#">
    
    <input type="hidden" name="instid" value="{[=it.meta.id]}">
    
    <label class="iam-form-group-title">Application Information</label>

    <div class="form-group">
      <label class="col-sm-2 control-label">Name</label>
      <div class="col-sm-10">
        <input type="text" class="form-control" name="app_title" value="{[=it.app_title]}">
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

    <div class="form-group">
      <label class="col-sm-2 control-label">Access URL</label>
      <div class="col-sm-10">
        <input type="text" class="form-control" name="url" value="{[=it.url]}">
      </div>
    </div>

    {[ if (it.privileges.length > 0) { ]}
    <label class="iam-form-group-title">Privilege Information</label>

    <div class="form-group">
      <label class="col-sm-2 control-label">Privileges</label>
      <div class="col-sm-10">
        <table class="table">
        <thead>
          <tr>
            <th>Privilege</th>
            <th>Roles</th>
          </tr>
        </thead>
        <tbody>
          {[~it.privileges :v]}
          <tr>
            <td>
              <p><strong>{[=v.desc]}</strong></p>
              <p>{[=v.privilege]}</p>
            </td>
            <td>
            {[ if (v.roles) { ]}
            {[~v.roles :rv]}
              {[~it._roles.items :drv]}
              {[ if (rv == drv.idxid) { ]}
                {[=drv.meta.name]}
              {[ } ]}
              {[~]}
            {[~]}
            {[ } ]}
            </td>
          </tr>
          {[~]}
        </tbody>
        </table>
      </div>
    </div>
    {[ } ]}

</div>

