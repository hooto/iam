<table class="iam-formtable valign-middle">
<tbody>
  <tr>
    <td width="200px">Access Key ID</td>
    <td class="iam-monofont">
      {[=it.id]}
    </td>
  </tr>
  <tr>
    <td>Access Key Secret</td>
    <td class="iam-monofont">
      {[=it.secret]}
    </td>
  </tr>
  <tr>
    <td>{[=l4i.T("Description")]}</td>
    <td>
      {[=it.description]}
    </td>
  </tr>
  <tr>
    <td>{[=l4i.T("Status")]}</td>
    <td>
      {[~it._statuses :v]}
      {[if (it.status == v.status) {]}{[=v.title]}{[}]}
      {[~]}
    </td>
  </tr>
  <tr>
    <td>{[=l4i.T("Scopes")]}</td>
    <td class="iam-monofont">
      {[~it.scopes :bv]}
      <div style="padding-bottom:5px;">{[=bv.name]} = {[=bv.value]}</div>
      {[~]}
    </td>
  </tr>
</tbody>
</table>

