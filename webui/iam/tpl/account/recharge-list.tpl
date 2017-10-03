<div class="iam-div-light">
  <table class="table table-hover">
    <thead>
      <tr>
        <th>ID</th>
        <th>Type</th>
        <th>Amount</th>
        <th>Cash</th>
		<th>Product Limits</th>
		<th>Product Max</th>
		<th>Product Inpay</th>
        <th>Comment</th>
		<th>Operator</th>
        <th>Created</th>
        <th></th>
      </tr>
    </thead>
    <tbody id="iam-acc-rechargelist"></tbody>
  </table>
</div>

<script id="iam-acc-rechargelist-tpl" type="text/html">
{[~it.items :v]}
<tr>
  <td class="iam-monofont">
    {[=v.id.substr(8)]}
  </td>
  <td>
    {[~it._ecoin_types :sv]}
    {[ if (v.type == sv.value) { ]}{[=sv.name]}{[ } ]}
    {[~]}
  </td>
  <td>{[=v.amount]}</td>
  <td>{[=v.cash_amount]}</td>
  <td>{[=v._exp_product_limits]}</td>
  <td>{[=v.exp_product_max]}</td>
  <td>{[=v._exp_product_inpay]}</td>
  <td>{[=v.comment]}</td>
  <td>{[=v.user_opr]}</td>
  <td>{[=l4i.MetaTimeParseFormat(v.created, "Y-m-d")]}</td>
  <td align="right">
    <!--
	<button class="pure-button button-xsmall"
      onclick="iamAccessKey.Set('{[=v.access_key]}')">
      Setting
    </button>
	-->
  </td>
</tr>
{[~]}
</script>

<script type="text/html" id="iam-acc-rechargelist-optools">
<li class="iam-btn iam-btn-primary">
  <a href="#" onclick="iamAccessKey.Set()">
     ECoin Recharge
  </a>
</li>
</script>

