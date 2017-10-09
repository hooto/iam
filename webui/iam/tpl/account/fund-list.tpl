<style>
.iam-div-light table th {
  vertical-align: middle !important;
}
</style>
<div class="iam-div-light">
  <table class="table table-hover">
    <thead>
      <tr>
        <th>ID</th>
        <th>Type</th>
        <th>Amount</th>
        <th>Cash</th>
        <th>Prepay</th>
        <th>Payout</th>
		<th>Product<br>Limits</th>
		<th>Product<br>Max</th>
		<th>Product<br>Inpay</th>
        <th>Actived</th>
      </tr>
    </thead>
    <tbody id="iam-acc-fundlist"></tbody>
  </table>
</div>

<script id="iam-acc-fundlist-tpl" type="text/html">
{[~it.items :v]}
<tr>
  <td class="iam-monofont">
    {[=v.id.substr(8)]}
  </td>
  <td>
    {[~it._fund_types :sv]}
    {[ if (v.type == sv.value) { ]}{[=sv.name]}{[ } ]}
    {[~]}
  </td>
  <td>{[=v.amount]}</td>
  <td>{[=v.cash_amount]}</td>
  <td>{[=v.prepay]}</td>
  <td>{[=v.payout]}</td>
  <td>{[=v._exp_product_limits]}</td>
  <td>{[=v.exp_product_max]}</td>
  <td class="iam-monofont">
    {[~v.exp_product_inpay :pv]}
	<div>{[=pv]}</div>
    {[~]}
  </td>
  <td>{[=l4i.MetaTimeParseFormat(v.created, "Y-m-d H:i")]}</td>
</tr>
{[~]}
</script>

