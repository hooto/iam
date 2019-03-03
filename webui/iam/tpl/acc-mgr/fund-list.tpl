<style>
.iam-div-light table th {
  vertical-align: middle !important;
}
</style>
<div class="iam-div-light">
  <table class="table table-hover valign-middle">
    <thead>
      <tr>
        <th>ID</th>
        <th>Type</th>
        <th>Amount</th>
        <th>Cash</th>
		<th>Product<br>Limits</th>
		<th>Product<br>Max</th>
        <th>User</th>
		<th>Operator</th>
        <th>Comment</th>
        <th>Created</th>
        <th></th>
      </tr>
    </thead>
    <tbody id="iam-accm-fundlist"></tbody>
  </table>
</div>

<script id="iam-accm-fundlist-tpl" type="text/html">
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
  <td>{[=v._exp_product_limits]}</td>
  <td>{[=v.exp_product_max]}</td>
  <td class="iam-monofont">
    {[=v.user]}
  </td>
  <td>{[=v.operator]}</td>
  <td>{[=v.comment]}</td>
  <td>{[=l4i.MetaTimeParseFormat(v.created, "Y-m-d")]}</td>
  <td align="right">
	<button class="pure-button button-small"
      onclick="iamAccMgr.FundSet('{[=v.id]}')">
	  <span class="fa fa-cog"></span>
      Setting
    </button>
  </td>
</tr>
{[~]}
</script>

<script type="text/html" id="iam-accm-fundlist-optools">
<li class="iam-btn iam-btn-primary">
  <a href="#" onclick="iamAccMgr.FundNew()">
     Recharge
  </a>
</li>
</script>

