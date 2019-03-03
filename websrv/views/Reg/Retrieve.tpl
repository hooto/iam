<!DOCTYPE html>
<html lang="en">
{{template "Common/HtmlHeaderReg.tpl" .}}
<body>

<div id="iam-resetpass-frame" class="iam-reg-frame">
<div id="iam-resetpass-box" class="iam-reg-box">

  <div class="iam-reg-msg01">{{T .LANG "Reset your password"}}</div>

  <form id="iam-resetpass-form" class="iam-reg-form" action="#retrieve" onsubmit="return false;">

    <input type="hidden" name="redirect_token" value="{{.redirect_token}}">

    <div id="iam-resetpass-form-alert" class="alert alert-info iam-groups">
      <p>{{T .LANG "Enter the username and email address you use to sign in"}}.</p>
      <p>{{T .LANG "The System will sent a URL to your email to reset the password"}}.</p>
    </div>

    <div class="iam-group">
      <input type="text" class="iam-input" name="username" placeholder="{{T .LANG "Username"}}">
    </div>

    <div class="iam-group">
      <input type="text" class="iam-input" name="email" placeholder="{{T .LANG "Email"}}">
    </div>

    <div class="iam-group">
      <button type="submit" class="iam-btn" onclick="iamLogin.RetrieveCommit()">{{T .LANG "Next"}}</button>
    </div>

  </form>

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
    $("#iam-resetpass-form").find("input[name=username]").focus();
}, 200);
</script>

</body>
</html>
