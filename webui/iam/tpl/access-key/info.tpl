<style type="text/css">
.iam-akinfo td {
    padding: 5px 0;
    vertical-align: top;
}
</style>

<table width="100%" class="iam-akinfo">
  <tr>
    <td width="120px"><strong>Access Key</strong></td>
    <td class="iam-monofont">
      {[=it.access_key]}
    </td>
  </tr>
  <tr>
    <td><strong>Secret Key</strong></td>
    <td class="iam-monofont">
      {[=it.secret_key]}
    </td>
  </tr>
  <tr>
    <td><strong>Description</strong></td>
    <td>
      {[=it.desc]}
    </td>
  </tr>
  <tr>
    <td><strong>Action</strong></td>
    <td>
      {[~it._actionls :v]}
      {[if (it.action == v.action) {]}{[=v.title]}{[}]}
      {[~]}
    </td>
  </tr>
  <tr>
    <td><strong>Bounds</strong></td>
    <td class="iam-monofont">
      {[~it.bounds :bv]}
      <div style="padding-bottom:5px;">{[=bv.name]}</div>
      {[~]}
    </td>
  </tr>
</table>

