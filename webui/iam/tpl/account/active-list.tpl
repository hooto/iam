<div class="iam-div-light">
  <table class="table table-hover">
    <thead>
      <tr>
        <th>ID</th>
        <th>Type</th>
        <th>Amount</th>
        <th>Payout</th>
        <th>Pre Paying</th>
		<th>Product Limits</th>
        <th>Actived</th>
        <th></th>
      </tr>
    </thead>
    <tbody id="iam-acc-activelist"></tbody>
  </table>
</div>

<script id="iam-acc-activelist-tpl" type="text/html">
{[~it.items :v]}
<tr>
  <td class="iam-monofont">
    {[=v.id]}
  </td>
  <td>
    {[~it._ecoin_types :sv]}
    {[ if (v.type == sv.value) { ]}{[=sv.name]}{[ } ]}
    {[~]}
  </td>
  <td>{[=v.amount]}</td>
  <td>{[=v.payout]}</td>
  <td>{[=v.prepay]}</td>
  <td>{[=v._exp_product_limits]}</td>
  <td>{[=l4i.MetaTimeParseFormat(v.created, "Y-m-d H:i")]}</td>
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

<script type="text/html" id="iam-acc-activelist-optools">
<li class="iam-btn iam-btn-primary">
  <a href="#" onclick="iamAccessKey.Set()">
     ECoin Recharge
  </a>
</li>
</script>
