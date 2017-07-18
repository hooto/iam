<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>IAM Service</title>
  <script src="/iam/~/jquery/jquery.min.js"></script>
  <script src="/iam/~/lessui/js/lessui.js"></script>
  <link rel="stylesheet" href="/iam/~/twbs/css/bootstrap.min.css" type="text/css">
  <link rel="shortcut icon" href="/iam/~/iam/img/favicon.ico" type="image/x-icon">
</head>
<body>

<style type="text/css">
* {
  margin: 0;
  padding: 0;
}

html,
body {
  height: 100%;
}

body {
  margin: 0 auto !important;
  padding: 0;
  font-size: 13px;
  font-family: Arial, sans-serif;
  background-color: #222;
  color: #eee;
  min-width: 500px;
  display: block;
}

#iam-login-frame {
  width: 100%;
  height: 100%;
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;
  min-width: 500px;
  min-height: 400px;
}

#iam-login-box {
  width: 360px;
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  color: #555;
  margin: 0 auto;
}

#iam-login-form {
  background-color: #f7f7f7;
  border-radius: 4px;
  padding: 20px 30px 20px 30px;
}

.iam-login-msg01 {
  font-size: 24px;
  padding: 40px 0;
  text-align: center;
  color: #eee;
}

.iam-user-ico-default {
  width: 96px;
  height: 96px;
  padding: 2px;
  position: relative;
  left: 50%;
  margin: 10px 0 30px -48px;
  border-radius: 50%;
  background-color: #dfdfdf;
}

#iam-login-form .iam-input-row {
  padding: 0 0 10px 0; 
}

#iam-login-form .ilf-input {
  display: block;
  width: 100%;
  height: 40px;
  padding: 5px 10px;
  font-size: 14px;
  line-height: 1.42857143;
  color: #555;
  background-color: #fff;
  background-image: none;
  border: 1px solid #ccc;
  border-radius: 2px;
  box-shadow: inset 0 1px 1px rgba(0, 0, 0, .075);
  transition: border-color ease-in-out .15s, box-shadow ease-in-out .15s;
}

#iam-login-form .ilf-input:focus {
  border-color: #66afe9;
  outline: 0;
  box-shadow: inset 0 1px 1px rgba(0,0,0,.075), 0 0 8px rgba(102, 175, 233, .6);
}

#iam-login-form .ilf-btn {
  width: 100%;
  height: 40px;
  display: inline-block;
  padding: 5px 10px;
  margin-bottom: 0;
  font-size: 14px;
  font-weight: normal;
  line-height: 1.42857143;
  text-align: center;
  white-space: nowrap;
  vertical-align: middle;
  cursor: pointer;
  -webkit-user-select: none;
     -moz-user-select: none;
      -ms-user-select: none;
          user-select: none;
  background-image: none;
  border: 1px solid transparent;
  border-radius: 3px;
  color: #fff;
  background-color: #427fed;
  border-color: #357ebd;
}

#iam-login-form .ilf-btn:hover {
  color: #fff;
  background-color: #3276b1;
  border-color: #285e8e;
}

#iam-login-form .iam-input-row-checkbox {
  display: inline-block;
  min-height: 20px;
  font-weight: normal;
  color: #555;
}
#iam-login-form .iam-input-row-checkbox input[type="checkbox"] {
  float: left;
  margin: 3px 5px 0 0;
  padding: 0;
  border: 1px solid #c6c6c6;
  cursor: pointer;
  width: 14px;
  height: 14px;
  box-sizing: border-box;
  line-height: normal;
}

#iam-login-form .iam-input-row-help {
  display: inline-block;
  float: right;
}
#iam-login-form .iam-input-row-help a {
  color: #427fed;
}

#iam-login-box > .signup {
  margin: 20px 0;
  text-align: center;
  font-size: 15px;
}
#iam-login-box > .signup a {
  font-size: 16px;
  color: #fff;
}

#iam-login-box > .footer {
  text-align: center;
  margin: 20px 0;
  font-size: 14px;
}
#iam-login-box > .footer a {
  color: #ccc;
}
#iam-login-box > .footer img {
  width: 16px;
  height: 16px;
}
</style>

<div id="iam-login-frame">
<div id="iam-login-box">

  <div class="iam-login-msg01">{{T . "Sign in with your Account"}}</div>

  <form id="iam-login-form" onsubmit="return false;">

    <input type="hidden" name="redirect_token" value="{{.redirect_token}}">

    <img class="iam-user-ico-default"  src="/iam/~/iam/img/user-default.svg">

    <div id="iam-login-form-alert" class="alert iam-input-row {{if eq .alert_msg nil}}hide{{end}}">{{.alert_msg}}</div>

    <div id="iam-login-input-frame">
      <div class="iam-input-row">
        <input type="text" class="ilf-input" name="uname" value="{{.uname}}" placeholder="{{T . "Username"}}">
      </div>

      <div class="iam-input-row">
        <input type="password" class="ilf-input" name="passwd" placeholder="{{T . "Password"}}">
      </div>

      <div class="iam-input-row">
        <button type="submit" class="ilf-btn">{{T . "Sign in"}}</button>
      </div>

      <div>
        <div class="iam-input-row-checkbox">
          <input name="persistent" type="checkbox" value="1" checked="{{.persistent_checked}}"> Stay signed in
        </div>
        <div class="iam-input-row-help">
        <a href="/iam/reg/retrieve?redirect_token={{.redirect_token}}">Forgot Password?</a>
        </div>
      </div>
    </div>
  </form>

  <div class="signup">
    <a href="/iam/reg/sign-up?redirect_token={{.redirect_token}}">Don't have an account? Create Account</a>
  </div>

  <div class="footer">
    <img src="/iam/~/iam/img/iam-s2-32.png"> 
    <a href="http://www.lessos.com/p/iam" target="_blank">lessOS IAM</a>
  </div>

</div>
</div>

<script type="text/javascript">

window.onload = function()
{
    //
    $("#iam-login-form").find("input[name=uname]").focus();
    
    //
    $("#iam-login-form").submit(function(event) {

        event.preventDefault();

        var alertid = "#iam-login-form-alert";

        l4i.Ajax("/iam/v1/service/login-auth", {
            type    : "POST",
            data    : $(this).serialize(),
            timeout : 3000,
            callback : function(err, rsj) {

                if (err || !rsj || !rsj.kind || rsj.kind != "ServiceLoginAuth") {

                    if (err) {
                        return l4i.InnerAlert(alertid, 'alert-danger', err);
                    }

                    if (rsj && rsj.error) {
                        return l4i.InnerAlert(alertid, 'alert-danger', rsj.error.message);
                    }

                    return l4i.InnerAlert(alertid, 'alert-danger', "Network Connection Exception");
                }
    
                l4i.InnerAlert(alertid, 'alert-success', "Successfully Sign-on. Page redirecting ...");
                $("#iam-login-input-frame").hide(100);
    
                window.setTimeout(function() {
                    window.location = rsj.redirect_uri;
                }, 1500);
            },
        });
    });
}
</script>
<body>
</html>
