<div class="iam-div-light">
  <table class="table table-hover valign-middle">
    <thead>
      <tr>
        <th>ID</th>
        <th>User</th>
        <th>Product</th>
        <th>Prepay</th>
        <th>Payout</th>
        <th>Pay Start</th>
        <th>Pay Close</th>
        <th>Updated</th>
		<th></th>
      </tr>
    </thead>
    <tbody id="iam-accmgr-chargelist"></tbody>
  </table>
</div>

<script id="iam-accmgr-chargelist-tpl" type="text/html">
{[~it.items :v]}
<tr>
  <td class="iam-monofont">
    {[=v.id.substr(8)]}
  </td>
  <td class="iam-monofont">
    {[=v.user]}
  </td>
  <td class="iam-monofont">
    {[=v.product]}
  </td>
  <td>{[=v.prepay]}</td>
  <td>{[=v.payout]}</td>
  <td>{[=l4i.UnixTimeFormat(v.time_start, "Y-m-d H:i")]}</td>
  <td>{[=l4i.UnixTimeFormat(v.time_close, "Y-m-d H:i")]}</td>
  <td>{[=l4i.MetaTimeParseFormat(v.updated, "Y-m-d H:i")]}</td>
  <td align="right">
    {[if (v.prepay > 0 && v.payout == 0) {]}
	<button class="pure-button button-small"
      onclick="iamAccMgr.ChargeSetPayout('{[=v.user]}', '{[=v.id]}')">
      Close
    </button>
	{[}]}
  </td>
</tr>
{[~]}
</script>
