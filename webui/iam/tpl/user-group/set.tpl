<div id="iam-user-group-set-alert"></div>

<div id="iam-user-group-set" class="form-horizontal">

  <table class="iam-formtable">
    <tbody>
 
    <tr>
      <td width="200px">{[=l4i.T("Username")]}</td>
      <td>
        <input type="text" class="form-control" name="name" value="{[=it.name]}" {[? it.name.length > 0]}readonly{[?]}>
      </td>
    </tr>

    <tr>
      <td>{[=l4i.T("Display Name")]}</td>
      <td>
        <input type="text" class="form-control" name="display_name" value="{[=it.display_name]}">
      </td>
    </tr>

    <tr>
      <td>{[=l4i.T("Owner")]}</td>
      <td>
        <input type="text" class="form-control" name="owners" value="{[=it._owners]}">
		<div class="form-text text-muted">example: username1,username2,...</div>
      </td>
    </tr>

    <tr>
      <td>{[=l4i.T("Member")]}</td>
      <td>
        <textarea class="form-control" name="members" rows="2">{[=it._members]}</textarea>
		<div class="form-text text-muted">example: username1,username2,...</div>
      </td>
    </tr>


    <tr>
      <td>{[=l4i.T("Status")]}</td>
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

