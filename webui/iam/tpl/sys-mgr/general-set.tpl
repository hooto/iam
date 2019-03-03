<div class="iam-div-light" style="padding:5px 10px">

  <div id="iam-sysmgr-generalset-alert" class="alert hide"></div>

  <div id="iam-sysmgr-generalset" class="form-horizontal">

  <div class="iam-form-group-title">Service Information</div>
  <table class="iam-formtable">
    <tbody>
    <tr>
      <td width="200px">Service Name</td>
      <td width="40%">
        <input type="text" class="form-control" name="service_name" value="{[=it._items.service_name]}">
      </td>
      <td>
        <div class="iam-callout iam-callout-primary">
          <h4>IAM Service Name</h4>
          <p>
            Will be used in the website name, e-mail signature ...
          </p>
        </div>
      </td>
    </tr>

    <tr>
      <td>Banner title</td>
      <td>
        <input type="text" class="form-control" name="webui_banner_title" value="{[=it._items.webui_banner_title]}">
      </td>
      <td>
        <div class="iam-callout iam-callout-primary">
          <h4>Banner Title</h4>
          <p>
            Will be used in the website banner
          </p>
        </div>
      </td>
    </tr>
    </tbody>
  </table>

  <div class="iam-form-group-title">User Registration Settings</div>

  <table class="iam-formtable">
    <tbody>
    <tr>
      <td width="200px">Disable Registration</td>
      <td width="40%">
        <div class="checkbox">
          <label>
            <input type="checkbox" name="user_reg_disable" value="1" {[ if (it._items.user_reg_disable == "1") { ]} checked {[ } ]}>
          </label>
        </div>
      </div>
      <td>
        <div class="iam-callout iam-callout-primary">
          <h4>Disable Registration</h4>
          <p>
            Disable Visitors to register accounts
          </p>
        </div>
      </div>
    </div>
    </tbody>
  </table>


    <!--
    <tr>
      <label class="col-sm-3">ICON</label>
      <div class="col-sm-3">
        <input id="service_info_icon" name="service_info_icon" size="20" type="file" class="form-control">
      </div>
      <td>
        <div class="iam-callout iam-callout-primary">
          <h4>Custom a Icon</h4>
          <p>
            Will be used in the site banner, browser shortcut icon ...
          </p>
        </div>
      </div>
    </div>
    -->

  <div class="iam-form-group-title">Messages</div>

  <table class="iam-formtable">
    <tbody>
    <tr>
      <td width="200px">Login form alert</td>
      <td width="40%">
        <input type="text" class="form-control" name="service_login_form_alert_msg" value="{[=it._items.service_login_form_alert_msg]}">
      </td>
	  <td></td>
    </tr>

    <tr>
	  <td></td>
      <td>
        <button type="submit" class="pure-button pure-button-primary" onclick="iamSys.GeneralSetCommit()">Save</button>
      </td>
	  <td></td>
    </tr>
    </tbody>
  </table>

</div>
