<div class="iam-div-light" style="padding:5px 10px">

  <div id="iam-sysmgr-mailerset-alert" class="alert hide"></div>

  <div id="iam-sysmgr-mailerset" class="form-horizontal">
   
    <label class="iam-form-group-title">Email SMTP server settings</label>

    <div class="form-group">
      <label class="col-sm-3">SMTP server address</label>
      <div class="col-sm-4">
        <input type="text" class="form-control" name="mailer_smtp_host" value="{[=it.smtp_host]}">
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-3">SMTP Port</label>
      <div class="col-sm-4">
        <input type="text" class="form-control" name="mailer_smtp_port" value="{[=it.smtp_port]}">
      </div>
      <div class="col-sm-5">
        <div class="iam-callout iam-callout-primary">
          <p>
            Defaults to 25 for unencrypted and TLS SMTP,<br>and 465 for SSL SMTP.
          </p>
        </div>
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-3">SMTP username</label>
      <div class="col-sm-4">
        <input type="text" class="form-control" name="mailer_smtp_user" value="{[=it.smtp_user]}">
      </div>
      <div class="col-sm-5">
        <div class="iam-callout iam-callout-primary">
          <p>
            Only enter a username if your SMTP server requires it.
          </p>
        </div>
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-3">SMTP password</label>
      <div class="col-sm-4">
        <input type="password" class="form-control" name="mailer_smtp_pass" value="{[=it.smtp_pass]}">
      </div>
      <div class="col-sm-5">
        <div class="iam-callout iam-callout-primary">
          <p>
            Only enter a password if your SMTP server requires it.
          </p>
        </div>
      </div>
    </div>

    <!-- TODO <div class="form-group">
      <label class="col-sm-3">SMTP encryption</label>
      <div class="col-sm-4">        
        <input type="text" class="form-control" name="" value="TODO" disabled>
      </div>
      <div class="col-sm-5">
        <div class="iam-callout iam-callout-primary">
          <p>
            Enter the transport layer encryption required by your SMTP server.
          </p>
        </div>
      </div>
    </div> -->

    <div class="form-group">
      <div class="col-sm-offset-3 col-sm-4">
        <button type="submit" class="pure-button pure-button-primary" onclick="iamSys.MailerSetCommit()">Save</button>
      </div>
    </div>

  </div>

</div>