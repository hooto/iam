<style>
.form-horizontal {
  margin: 0 15px;
  padding: 2px;
}
.form-group {
  margin-bottom: 5px;
}
</style>

<div id="ids-myapp-instset-alert"></div>
    
<div id="ids-myapp-instset" class="form-horizontal" action="#">
    
    <input type="hidden" name="instid" value="{[=it.meta.id]}">
    
    <label class="ids-form-group-title">Application Information</label>

    <div class="form-group">
      <label class="col-sm-2 control-label">Name</label>
      <div class="col-sm-10">
        <input type="text" class="form-control" name="name" value="{[=it.meta.name]}">
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-2 control-label">Status</label>
      <div class="col-sm-10">
        {[~it._statusls :v]}
          <span class="ids-form-checkbox">
            <input type="radio" name="status" value="{[=v.status]}" {[ if (v.status == it.status) { ]}checked="checked"{[ } ]}> {[=v.title]}
          </span>
        {[~]}
      </div>
    </div>

    {[ if (it.privileges.length > 0) { ]}
    <label class="ids-form-group-title">Privilege Information</label>

    <div class="form-group">
      <label class="col-sm-2 control-label">Privileges</label>
      <div class="col-sm-10">
        <table class="table">
        <thead>
          <tr>
            <th>Privilege</th>
            <th>Description</th>
            <th>General Roles</th>
          </tr>
        </thead>
        <tbody>
          {[~it.privileges :v]}
          <tr>
            <td>{[=v.privilege]}</td>
            <td>{[=v.desc]}</td>
            <td>
            {[ if (v.roles) { ]}
            {[~v.roles :rv]}
              {[~it._roles.items :drv]}
              {[ if (rv == drv.meta.id) { ]}
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

