<div id="iam-accmgr-fund-alert" class="alert hide"></div>

<div id="iam-accmgr-fund-form" class="form-horizontal">
  <input type="hidden" name="id" value="{[=it.id]}">
  <table class="iam-formtable">
    <tbody>
    <tr>
      <td width="200px">{[=l4i.T("Type")]}</td>
      <td>
	    <select name="type" class="form-control">
        {[~it._fund_types :v]}
          <option value="{[=v.value]}" {[ if (it.type) { ]}selected{[ } ]}> {[=v.name]}
        {[~]}
		</select>
      </td>
    </tr>

    <tr>
      <td>{[=l4i.T("Product Type Limit")]}</td>
      <td>
        <input type="text" class="form-control input-sm" name="exp_product_limits" value="{[=it.exp_product_limits]}">
      </td>
    </tr>

    <tr>
      <td>{[=l4i.T("Product Number Limit")]}</td>
      <td>
        <input type="text" class="form-control input-sm" name="exp_product_max" value="{[=it.exp_product_max]}">
      </td>
    </tr>

    <tr>
      <td>{[=l4i.T("Comment")]}</td>
      <td>
        <input type="text" class="form-control input-sm" name="comment" value="{[=it.comment]}">
      </td>
    </tr>
    </tbody>
  </table>
</div>
