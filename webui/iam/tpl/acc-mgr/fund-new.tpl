<div id="iam-accmgr-fund-alert" class="alert hide"></div>

<div id="iam-accmgr-fund-form" class="form-horizontal">

  <table class="iam-formtable">
    <tbody>
    <tr>
      <td width="200px">Username</td>
      <td>
        <input type="text" class="form-control" name="user" value="{[=it.user]}">
      </td>
    </tr>

    <tr>
      <td>{[=l4i.T("Type")]}</td>
      <td>
	    <select name="type" class="form-control">
        {[~it._fund_types :v]}
          <option value="{[=v.value]}" {[ if (v.default) { ]}selected{[ } ]}> {[=v.name]}
        {[~]}
		</select>
      </td>
    </tr>

    <tr>
      <td>{[=l4i.T("Amount")]}</td>
      <td>
        <input type="text" class="form-control" name="amount" value="">
      </td>
    </tr>

    <tr>
      <td>{[=l4i.T("Product Type Limit")]}</td>
      <td>
        <input type="text" class="form-control" name="exp_product_limits" value="sys/pod">
      </td>
    </tr>

    <tr>
      <td>{[=l4i.T("Product Number Limit")]}</td>
      <td>
        <input type="text" class="form-control" name="exp_product_max" value="1">
      </td>
    </tr>

    <tr>
      <td>{[=l4i.T("Comment")]}</td>
      <td>
        <input type="text" class="form-control" name="comment" value="">
      </td>
    </tr>
    </tbody>
  </table>

</div>
