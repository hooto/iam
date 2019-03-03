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

  <div class="iam-form-group-title">Application Information</div>

  <table class="iam-formtable">
    <tbody>
    <tr>
      <td width="200px">Name</td>
      <td>
        <input type="text" class="form-control" name="app_title" value="{[=it.app_title]}">
      </td>
    </tr>

    <tr>
      <td>Status</td>
      <td>
        {[~it._statusls :v]}
          <span class="iam-form-checkbox">
            <input type="radio" name="status" value="{[=v.status]}" {[ if (v.status == it.status) { ]}checked="checked"{[ } ]}> {[=v.title]}
          </span>
        {[~]}
      </td>
    </tr>

    <tr>
      <td>Access URL</td>
      <td>
        <input type="text" class="form-control" name="url" value="{[=it.url]}">
      </td>
    </tr>
    </tbody>
  </table>


  {[ if (it.privileges.length > 0) { ]}
  <div class="iam-form-group-title">Privilege Information</div>
  <table class="iam-formtable">
    <tbody>
    <tr>
      <td width="200px">Privileges</td>
      <td>
        <table>
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
              {[ if (rv == drv.id) { ]}
                {[=drv.name]}
              {[ } ]}
              {[~]}
            {[~]}
            {[ } else { ]}
              Owner
            {[ } ]}
            </td>
          </tr>
          {[~]}
        </tbody>
        </table>
      </td>
    </tr>
    </tbody>
  </table>
  {[ } ]}

</div>

