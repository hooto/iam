<div class="iam-div-light" style="padding:5px 10px">

  <div id="iam-sysmgr-mailerset-alert" class="alert hide"></div>

  <div class="iam-form-group-title">{[=l4i.T("%s Settings", l4i.T("SMTP server"))]}</div>

  <table id="iam-sysmgr-mailerset" class="iam-formtable">
    <tbody>
 
    <tr>
      <td width="200px">SMTP server address</label>
      <td width="40%">
        <input type="text" class="form-control" name="mailer_smtp_host" value="{[=it.smtp_host]}">
      </td>
	  <td></td>
    </tr>

    <tr>
      <td>SMTP Port</label>
      <td>
        <input type="text" class="form-control" name="mailer_smtp_port" value="{[=it.smtp_port]}">
      </td>
      <td>
        <div class="iam-callout iam-callout-primary">
          <p>
            Defaults to 25 for unencrypted and TLS SMTP,<br>and 465 for SSL SMTP.
          </p>
        </div>
      </td>
    </tr>

    <tr>
      <td>SMTP username</label>
      <td>
        <input type="text" class="form-control" name="mailer_smtp_user" value="{[=it.smtp_user]}">
      </td>
      <td>
        <div class="iam-callout iam-callout-primary">
          <p>
            Only enter a username if your SMTP server requires it.
          </p>
        </div>
      </td>
    </tr>

    <tr>
      <td>SMTP password</label>
      <td>
        <input type="password" class="form-control" name="mailer_smtp_pass" value="{[=it.smtp_pass]}">
      </td>
      <td>
        <div class="iam-callout iam-callout-primary">
          <p>
            Only enter a password if your SMTP server requires it.
          </p>
        </div>
      </td>
    </tr>

    <!-- TODO <tr>
      <td>SMTP encryption</label>
      <td>
        <input type="text" class="form-control" name="" value="TODO" disabled>
      </div>
      <td>
        <div class="iam-callout iam-callout-primary">
          <p>
            Enter the transport layer encryption required by your SMTP server.
          </p>
        </div>
      </div>
    </div> -->

    <tr>
	  <td></td>
      <td>
        <button type="submit" class="pure-button pure-button-primary" onclick="iamSys.MailerSetCommit()">{[=l4i.T("Save")]}</button>
      </td>
	  <td></td>
    </tr>
    </tbody>
  </table>

</div>
