<style type="text/css">
#iam-ak-set td {
    padding: 5px 0;
    vertical-align: top;
}
</style>


<div id="iam-ak-set-alert" class="alert hide"></div>
    
<table id="iam-ak-set" class="iam-formtable valign-middle">
  <input type="hidden" name="id" value="{[=it.id]}">
  <tr>
    <td width="200px">{[=l4i.T("Description")]}</td>
    <td>
      <input name="desc" type="text" class="form-control" value="{[=it.description]}">
    </td>
  </tr>
  <tr>
    <td>{[=l4i.T("Status")]}</td>
    <td>
      {[~it._statuses :v]}
        <span class="iam-form-checkbox">
          <input type="radio" name="status" value="{[=v.status]}" {[ if (v.status == it.status) { ]}checked="checked"{[ } ]}> {[=v.title]}
        </span>
      {[~]}
    </td>
  </tr>
</table>

