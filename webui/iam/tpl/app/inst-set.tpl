<style>
.form-horizontal {
  margin: 0 15px;
  padding: 2px;
}
.form-group {
  margin-bottom: 5px;
}
</style>

<div id="iam-app-instset-alert"></div>

<div id="iam-app-instset" class="form-horizontal" action="#">

  <input type="hidden" name="instid" value="{[=it.meta.id]}">

  <label class="iam-form-group-title">{[=l4i.T("%s Information", l4i.T("Application"))]}</label>

  <table class="iam-formtable">
    <tbody>
    <tr>
      <td width="200px">
        {[=l4i.T("Name")]}
      </td>
      <td>
        <input type="text" class="form-control" name="app_title" value="{[=it.app_title]}">
      </td>
    </tr>

    <tr>
      <td>
        {[=l4i.T("Status")]}
      </td>
      <td>
        {[~it._statusls :v]}
        {[ if (v.status == it.status) { ]}
          {[=v.title]}
        {[ } ]}
        {[~]}
      </td>
    </tr>

    <tr>
      <td>{[=l4i.T("%s URL", l4i.T("Access"))]}</td>
      <td>
        <input type="text" class="form-control" name="url" value="{[=it.url]}">
      </td>
    </tr>
    </tbody>
  </table>

  {[? it.privileges.length > 0]}
  <label class="iam-form-group-title">{[=l4i.T("%s Information", l4i.T("Privilege"))]}</label>

  <table class="iam-formtable">
    <tbody>
    <tr>
      <td width="200px">{[=l4i.T("Privileges")]}</td>
      <td>
        <table>
        <thead>
          <tr>
            <th>{[=l4i.T("Privilege")]}</th>
            <th>{[=l4i.T("Description")]}</th>
            <th>{[=l4i.T("Roles")]}</th>
          </tr>
        </thead>
        <tbody>
          {[~it.privileges :v]}
          <tr>
            <td>
              <strong>{[=v.privilege]}</strong>
            </td>
            <td>
              <strong>{[=v.desc]}</strong>
            </td>
            <td>
            {[ if (v.roles) { ]}
            {[~v.roles :rv]}
              {[~it._roles.items :drv]}
              {[ if (rv == drv.id) { ]}
                {[=drv.name]}&nbsp;
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
    </tbody
  </table>
  {[?]}

</div>

