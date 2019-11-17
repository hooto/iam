<div id="iam-usermgr-roleset-alert" class="alert hide"></div>
    
<div id="iam-usermgr-roleset" class="form-horizontal">
    
  <div class="iam-form-group-title">{[=l4i.T("%s Information", l4i.T("Role"))]}</div>

  <table class="iam-formtable">
    <tbody>
    <tr>
      <td width="200px">{[=l4i.T("Name")]}</label>
      <td>
        <input type="text" class="form-control" name="name" value="{[=it.name]}" {[? it.name && it.name.length > 0]}readonly{[?]}>
      </td>
    </tr>

    <tr>
      <td>{[=l4i.T("Description")]}</label>
      <td>
        <input type="text" class="form-control" name="desc" value="{[=it.desc]}">
      </td>
    </tr>

    <tr>
      <td>{[=l4i.T("Status")]}</label>
      <td>
        {[~it._statusls :v]}
          <span class="iam-form-checkbox">
            <input type="radio" name="status" value="{[=v.status]}" {[ if (v.status == it.status) { ]}checked="checked"{[ } ]}> {[=v.title]}
          </span>
        {[~]}
      </td>
    </tr>

    </tbody>
  </table>
</div>

