<!DOCTYPE html>
<html lang="en">
{{template "Common/HtmlHeaderReg.tpl" .}}
<body>


<div id="iam-reg-passreset-frame" class="iam-reg-frame">
<div id="iam-reg-passreset-box" class="iam-reg-box">

  <div class="iam-reg-msg01">{{T .LANG "Reset your Password"}}</div>

  {{if .pass_reset_id}}
  <form id="iam-reg-passreset-form"  class="iam-reg-form" onsubmit="return false;">

    <input type="hidden" name="id" value="{{.pass_reset_id}}">
    <input type="hidden" name="redirect_token" value="{{.redirect_token}}">
    <div class="iam-key alert alert-dark">{{.pass_reset_id}}</div>

    <div id="iam-reg-passreset-form-alert" class="alert hide"></div>

    <div class="iam-group">
      <input type="text" class="iam-input" name="email" placeholder="{{T .LANG "Confirm your Email"}}">
    </div>

    <div class="iam-group">
      <input type="password" class="iam-input" name="passwd" placeholder="{{T .LANG "Your new password"}}">
    </div>

    <div class="iam-group">
      <input type="password" class="iam-input" name="passwd_confirm" placeholder="{{T .LANG "Confirm your new password"}}">
    </div>    

    <div class="iam-group">
      <button type="submit" class="iam-btn" onclick="iamLogin.PassResetCommit()">{{T .LANG "Next"}}</button>
    </div>

  </form>
  {{else}}
    <div class="alert alert-danger">{{T .LANG "The Token is not valid or Expired"}}</div>
  {{end}}

  <div class="ref-action">
    <a href="/iam/service/login?redirect_token={{.redirect_token}}">{{T .LANG "Sign in with your Account"}}</a>
  </div>

  <div class="footer">
    <img src="/iam/~/iam/img/iam-s2-32.png"> 
    <a href="https://github.com/hooto/iam" target="_blank">hooto IAM</a>
  </div>

</div>
</div>


<script>
setTimeout(function() {
    $("#iam-reg-passreset-form").find("input[name=email]").focus();
}, 200);
</script>

</body>
</html>
