<style>
.form-horizontal {
  margin: 0 15px;
  padding: 2px;
}
.form-group {
  margin-bottom: 5px;
}
</style>


<div id="iam-accmgr-recharge-alert" class="alert hide"></div>


<div id="iam-accmgr-recharge-form" class="form-horizontal">
 
    <div class="form-group">
      <label class="col-sm-3 control-label">Username</label>
      <div class="col-sm-9">
        <input type="text" class="form-control input-sm" name="user" value="{[=it.user]}">
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-3 control-label">Recharge Type</label>
      <div class="col-sm-9">
	    <select name="type" class="form-control">
        {[~it._recharge_types :v]}
          <option value="{[=v.value]}" {[ if (v.default) { ]}selected{[ } ]}> {[=v.name]}
        {[~]}
		</select>
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-3 control-label">Ecoin Amount</label>
      <div class="col-sm-9">
        <input type="text" class="form-control input-sm" name="amount" value="">
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-3 control-label">Product Limits</label>
      <div class="col-sm-9">
        <input type="text" class="form-control input-sm" name="exp_product_limits" value="sys/pod">
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-3 control-label">Product Max</label>
      <div class="col-sm-9">
        <input type="text" class="form-control input-sm" name="exp_product_max" value="0">
      </div>
    </div>


    <div class="form-group">
      <label class="col-sm-3 control-label">Comment</label>
      <div class="col-sm-9">
        <input type="text" class="form-control input-sm" name="comment" value="">
      </div>
    </div>

</div>
