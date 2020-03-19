<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>IAM Service</title>
  <script src="/iam/~/zepto/zepto.js"></script>
  <script src="/iam/~/lessui/js/lessui.js"></script>
  <script src="/iam/~/iam/js/login.js?v={{.sys_version_hash}}"></script>
  <link rel="stylesheet" href="/iam/~/iam/css/reg.css?v={{.sys_version_hash}}" type="text/css">
  <link rel="shortcut icon" href="/iam/~/iam/img/iam-s2-32.png" type="image/x-icon">
</head>
<body>

<div id="iam-login-frame" class="iam-reg-frame">
<div id="iam-login-box" class="iam-reg-box">

  <div class="iam-reg-msg01">{{T .LANG "Sign in with your Account"}}</div>

  <form id="iam-login-form" class="iam-reg-form" onsubmit="return false;">

    <input type="hidden" name="redirect_token" value="{{.redirect_token}}">

    <img class="iam-user-ico-default"  src="/iam/~/iam/img/user-default.svg">

    <div id="iam-login-form-alert" class="alert alert-info {{if eq .alert_msg nil}}hide{{end}}">{{.alert_msg}}</div>

    <div id="iam-login-input-frame">
      <div class="iam-input-row">
        <input type="text" class="iam-input iam-form-item" name="uname" id="iam-login-form-username" value="{{.uname}}" placeholder="{{T .LANG "Username"}}">
      </div>

      <div class="iam-input-row">
        <input type="password" class="iam-input iam-form-item" name="passwd" id="iam-login-form-pwd" placeholder="{{T .LANG "Password"}}">
      </div>

      <div class="iam-input-row">
        <button type="submit" class="iam-btn iam-form-item" onclick="iamLogin.LoginCommit()">{{T .LANG "Sign in"}}</button>
      </div>

      <div>
        <div class="iam-input-row-checkbox">
          <input name="persistent" type="checkbox" value="1" checked="{{.persistent_checked}}"> {{T .LANG "Stay signed in"}}
        </div>
        <div class="iam-input-row-help">
        <a href="/iam/reg/retrieve?redirect_token={{.redirect_token}}">{{T .LANG "Forgot Password"}} ?</a>
        </div>
      </div>
    </div>
  </form>

  <div class="ref-action">
    <a href="/iam/reg/sign-up?redirect_token={{.redirect_token}}">{{T .LANG "Create your Account"}}</a>
  </div>

  <div class="footer">
    <img src="/iam/~/iam/img/iam-s2-32.png"> 
    <a href="https://github.com/hooto/iam" target="_blank">hooto IAM</a>
  </div>

</div>
</div>

<script type="text/javascript">
setTimeout(function() {
    document.getElementById("iam-login-form-username").focus();
}, 200);
</script>
</body>
</html>
