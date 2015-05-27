<style>
.page-header {
	margin: 10px 0;
	font-height: 100%;
}
</style>

<div class="page-header">
  <h2>Email Settings <small></small></h2>
</div>

<div id="ids-sysmgr-mailerset-alert" class="alert hide"></div>

  <div id="ids-sysmgr-mailerset" class="form-horizontal">
   
    <label class="ids-form-group-title">Email SMTP server settings</label>

    <div class="form-group">
      <label class="col-sm-3">SMTP server address</label>
      <div class="col-sm-4">
        <input type="text" class="form-control" name="mailer_smtp_host" value="{[=it._items.mailer_smtp_host]}">
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-3">SMTP Port</label>
      <div class="col-sm-4">
        <input type="text" class="form-control" name="mailer_smtp_port" value="{[=it._items.mailer_smtp_port]}">
      </div>
      <div class="col-sm-5">
        <div class="ids-callout ids-callout-primary">
          <p>
            Defaults to 25 for unencrypted and TLS SMTP,<br>and 465 for SSL SMTP.
          </p>
        </div>
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-3">SMTP username</label>
      <div class="col-sm-4">
        <input type="text" class="form-control" name="mailer_smtp_user" value="{[=it._items.mailer_smtp_user]}">
      </div>
      <div class="col-sm-5">
        <div class="ids-callout ids-callout-primary">
          <p>
            Only enter a username if your SMTP server requires it.
          </p>
        </div>
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-3">SMTP password</label>
      <div class="col-sm-4">
        <input type="password" class="form-control" name="mailer_smtp_pass" value="{[=it._items.mailer_smtp_pass]}">
      </div>
      <div class="col-sm-5">
        <div class="ids-callout ids-callout-primary">
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
        <div class="ids-callout ids-callout-primary">
          <p>
            Enter the transport layer encryption required by your SMTP server.
          </p>
        </div>
      </div>
    </div> -->

    <div class="form-group">
      <div class="col-sm-offset-3 col-sm-4">
        <button type="submit" class="btn btn-primary" onclick="idssys.MailerSetCommit()">Save</button>
      </div>
    </div>

  </div>
