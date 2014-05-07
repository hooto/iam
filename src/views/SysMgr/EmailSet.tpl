<style>
.page-header {
	margin: 10px 0;
	font-height: 100%;
}
</style>

<div class="page-header">
  <h2>Email Settings <small></small></h2>
</div>

<div id="s61n1y" class="alert hide"></div>

  <form id="znslct" class="form-horizontal" action="#">
   
    <label class="ids-form-group-title">Email SMTP server settings</label>

    <div class="form-group">
      <label class="col-sm-3">SMTP server address</label>
      <div class="col-sm-4">
        <input type="text" class="form-control" name="mailer_smtp_host" value="{{.mailer.SmtpHost}}">
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-3">SMTP Port</label>
      <div class="col-sm-4">
        <input type="text" class="form-control" name="mailer_smtp_port" value="{{.mailer.SmtpPort}}">
      </div>
      <div class="col-sm-5">
        <div class="ids-callout ids-callout-primary">
          <p>
            Optional. Defaults to 25 for unencrypted and TLS SMTP, and 465 for SSL SMTP.
          </p>
        </div>
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-3">SMTP username</label>
      <div class="col-sm-4">
        <input type="text" class="form-control" name="mailer_smtp_user" value="{{.mailer.SmtpUser}}">
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
        <input type="password" class="form-control" name="mailer_smtp_pass" value="{{.mailer.SmtpPass}}">
      </div>
      <div class="col-sm-5">
        <div class="ids-callout ids-callout-primary">
          <p>
            Only enter a password if your SMTP server requires it.
          </p>
        </div>
      </div>
    </div>

    <div class="form-group">
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
    </div>

    <div class="form-group">
      <div class="col-sm-offset-3">
        <button type="submit" class="btn btn-primary">{{T . "Save"}}</button>
      </div>
    </div>

  </form>

<script type="text/javascript">
$("#znslct").submit(function(event) {

    event.preventDefault(); 

    _ids_sysmgr_genset();
});

function _ids_sysmgr_genset()
{
    $.ajax({
        type    : "POST",
        url     : "/ids/sys-mgr/email-set-save?_="+Math.random(),
        data    : $("#znslct").serialize(),
        timeout : 3000,
        success : function(rsp) {
          	lessAlert("#s61n1y", "alert-success", rsp);
        },
        error: function(xhr, textStatus, error) {
            lessAlert("#s61n1y", "alert-danger", "Error: "+ xhr.responseText);
        }
    });
}

</script>