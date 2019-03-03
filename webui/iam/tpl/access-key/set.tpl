<style type="text/css">
#iam-ak-set td {
    padding: 5px 0;
    vertical-align: top;
}
</style>


<div id="iam-ak-set-alert" class="alert hide"></div>
    
<table id="iam-ak-set" class="iam-formtable valign-middle">
  <input type="hidden" name="access_key" value="{[=it.access_key]}">
  <tr>
    <td width="200px">{[=l4i.T("Description")]}</td>
    <td>
      <input name="desc" type="text" class="form-control" value="{[=it.desc]}">
    </td>
  </tr>
  <tr>
    <td>{[=l4i.T("Action")]}</td>
    <td>
      {[~it._actionls :v]}
        <span class="iam-form-checkbox">
          <input type="radio" name="action" value="{[=v.action]}" {[ if (v.action == it.action) { ]}checked="checked"{[ } ]}> {[=v.title]}
        </span>
      {[~]}
    </td>
  </tr>
</table>

