<div id="iam-acc-payoutlist-alert" class="hide" style="margin:20px 0;"></div>

<div class="iam-div-light">
  <table class="table table-hover">
    <thead>
      <tr>
        <th>ID</th>
        <th>Product</th>
        <th>Payout</th>
        <th>Pay Start</th>
        <th>Pay Close</th>
        <th>Updated</th>
      </tr>
    </thead>
    <tbody id="iam-acc-payoutlist"></tbody>
  </table>
</div>

<script id="iam-acc-payoutlist-tpl" type="text/html">
{[~it.items :v]}
<tr>
  <td class="iam-monofont">
    {[=v.id.substr(8)]}
  </td>
  <td class="iam-monofont">
    {[=v.product]}
  </td>
  <td>{[=v.payout]}</td>
  <td>{[=l4i.UnixTimeFormat(v.time_start, "Y-m-d H:i")]}</td>
  <td>{[=l4i.UnixTimeFormat(v.time_close, "Y-m-d H:i")]}</td>
  <td>{[=l4i.MetaTimeParseFormat(v.updated, "Y-m-d H:i")]}</td>
</tr>
{[~]}
</script>
