<div id="iam-accmgr-chargeset-payout-alert" class="alert hide"></div>

<div id="iam-accmgr-chargeset-payout-form">
  <input type="hidden" name="id" value="{[=it.id]}">
  <input type="hidden" name="user" value="{[=it.user]}">

  <table class="iam-formtable">
    <tbody>
    <tr>
      <td>ID</td>
      <td>
	     {[=it.id]}
      </td>
    </tr>

    <tr>
      <td>{[=l4i.T("Product")]}</td>
      <td>
	     {[=it.product]}
      </td>
    </tr>

    <tr>
      <td>{[=l4i.T("Amount")]}</td>
      <td>
        <input type="text" class="form-control input-sm" name="payout" value="{[=it.prepay]}">
      </td>
    </tr>
    </tbody>
  </table>
</div>
