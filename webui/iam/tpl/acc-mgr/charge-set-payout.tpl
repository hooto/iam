<style>
.form-horizontal {
  margin: 0 15px;
  padding: 2px;
}
.form-group {
  margin-bottom: 5px;
}
</style>


<div id="iam-accmgr-chargeset-payout-alert" class="alert hide"></div>

<div id="iam-accmgr-chargeset-payout-form" class="form-horizontal">
    <input type="hidden" name="id" value="{[=it.id]}">
    <input type="hidden" name="user" value="{[=it.user]}">

    <div class="form-group">
      <label class="col-sm-3 control-label">Charge</label>
      <div class="col-sm-9">
	     {[=it.id]}
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-3 control-label">Product</label>
      <div class="col-sm-9">
	     {[=it.product]}
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-3 control-label">Payout</label>
      <div class="col-sm-9">
        <input type="text" class="form-control input-sm" name="payout" value="{[=it.prepay]}">
      </div>
    </div>
</div>
