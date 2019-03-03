<!DOCTYPE html>
<html lang="en">
{{template "Common/HtmlHeaderReg.tpl" .}}
<body>

<div id="iam-signup-frame" class="iam-reg-frame">
<div id="iam-signup-box" class="iam-reg-box">

  <div class="iam-reg-msg01">{{T . "Create your Account"}}</div>

  <form id="iam-signup-form" class="iam-reg-form" action="#" onsubmit="return false;">

    <div id="iam-signup-form-alert" class="alert hide iam-group"></div>

    <input type="hidden" name="redirect_token" value="{{.redirect_token}}">
    {{if eq .user_reg_disable false }}
    <div class="iam-group">
      <input type="text" class="iam-input" name="uname" value="{{.uname}}" placeholder="{{T . "Unique Username"}}">
    </div>

    <div class="iam-group">
      <input type="text" class="iam-input" name="email" placeholder="{{T . "Email"}}">
    </div>

    <div class="iam-group">
      <input type="password" class="iam-input" name="passwd" placeholder="{{T . "Password"}}">
    </div>

    <div class="iam-group">
      <button type="submit" class="iam-btn" onclick="iamLogin.SignupCommit()">{{T . "Create Account"}}</button>
    </div>
    {{else}}
    <div class="alert alert-danger">User registration was closed!<br>Please contact the administrator to manually register accounts</div>
    {{end}}

  </form>

  <div class="ref-action">
    <a href="/iam/service/login?redirect_token={{.redirect_token}}">Sign in with your Account</a>
  </div>

  <div class="footer">
    <img src="/iam/~/iam/img/iam-s2-32.png"> 
    <a href="https://github.com/hooto/iam" target="_blank">hooto IAM</a>
  </div>

</div>
</div>

<script>
setTimeout(function() {
    $("#iam-signup-form").find("input[name=uname]").focus();
}, 200);
</script>

</body>
</html>
