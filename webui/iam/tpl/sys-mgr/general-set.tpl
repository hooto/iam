<style>
.page-header {
	margin: 10px 0;
	font-height: 100%;
}
</style>

<div class="page-header">
  <h2>General Settings <small></small></h2>
</div>

<div id="iam-sysmgr-generalset-alert" class="alert hide"></div>

  <div id="iam-sysmgr-generalset" class="form-horizontal">
   
    <label class="iam-form-group-title">Service Information</label>

    <div class="form-group">
      <label class="col-sm-2 control-label">Service Name</label>
      <div class="col-sm-4">
        <input type="text" class="form-control" name="service_name" value="{[=it._items.service_name]}">
      </div>
      <div class="col-sm-6">
        <div class="iam-callout iam-callout-primary">
          <h4>IAM Service Name</h4>
          <p>
            Will be used in the website name, e-mail signature ...
          </p>
        </div>
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-2 control-label">Banner title</label>
      <div class="col-sm-4">
        <input type="text" class="form-control" name="webui_banner_title" value="{[=it._items.webui_banner_title]}">
      </div>
      <div class="col-sm-6">
        <div class="iam-callout iam-callout-primary">
          <h4>Banner Title</h4>
          <p>
            Will be used in the website banner
          </p>
        </div>
      </div>
    </div>

    <label class="iam-form-group-title">User Registration Settings</label>

    <div class="form-group">
      <label class="col-sm-2 control-label">Disable Registration</label>
      <div class="col-sm-4">
        <div class="checkbox">
          <label>
            <input type="checkbox" name="user_reg_disable" value="1" {[ if (it._items.user_reg_disable == "1") { ]} checked {[ } ]}>
          </label>
        </div>
      </div>
      <div class="col-sm-6">
        <div class="iam-callout iam-callout-primary">
          <h4>Disable Registration</h4>
          <p>
            Disable Visitors to register accounts
          </p>
        </div>
      </div>
    </div>

    <!--
    <div class="form-group">
      <label class="col-sm-3">ICON</label>
      <div class="col-sm-3">
        <input id="service_info_icon" name="service_info_icon" size="20" type="file" class="form-control">
      </div>
      <div class="col-sm-6">
        <div class="iam-callout iam-callout-primary">
          <h4>Custom a Icon</h4>
          <p>
            Will be used in the site banner, browser shortcut icon ...
          </p>
        </div>
      </div>
    </div>
    -->

    <div class="form-group">
      <div class="col-sm-offset-2 col-sm-3">
        <button type="submit" class="btn btn-primary" onclick="iamsys.GeneralSetCommit()">Save</button>
      </div>
    </div>

  </div>
