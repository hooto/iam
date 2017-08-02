<style type="text/css">
#iam-akset td {
    padding: 5px 0;
    vertical-align: top;
}
</style>


<div id="iam-akset-alert" class="alert hide"></div>
    
<table id="iam-akset" width="100%" style="">
  <input type="hidden" name="access_key" value="{[=it.access_key]}">
  <tr>
    <td width="120px"><strong>Description</strong></td>
    <td>
      <input name="desc" type="text" class="form-control" value="{[=it.desc]}">
    </td>
  </tr>
  <tr>
    <td><strong>Action</strong></td>
    <td>
      {[~it._actionls :v]}
        <span class="iam-form-checkbox">
          <input type="radio" name="action" value="{[=v.action]}" {[ if (v.action == it.action) { ]}checked="checked"{[ } ]}> {[=v.title]}
        </span>
      {[~]}
    </td>
  </tr>
</table>

