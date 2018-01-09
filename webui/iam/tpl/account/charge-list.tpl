<div class="iam-div-light">
  <table class="table table-hover">
    <thead>
      <tr>
        <th>ID</th>
        <th>Product</th>
        <th>Prepay</th>
        <th>Pay Start</th>
        <th>Pay Close</th>
        <th>Updated</th>
      </tr>
    </thead>
    <tbody id="iam-acc-chargelist"></tbody>
  </table>
</div>

<script id="iam-acc-chargelist-tpl" type="text/html">
{[~it.items :v]}
<tr>
  <td class="iam-monofont">
    {[=v.id.substr(8)]}
  </td>
  <td class="iam-monofont">
    {[=v.product]}
  </td>
  <td>{[=v.prepay]}</td>
  <td>{[=l4i.UnixTimeFormat(v.time_start, "Y-m-d H:i")]}</td>
  <td>{[=l4i.UnixTimeFormat(v.time_close, "Y-m-d H:i")]}</td>
  <td>{[=l4i.MetaTimeParseFormat(v.updated, "Y-m-d H:i")]}</td>
</tr>
{[~]}
</script>
