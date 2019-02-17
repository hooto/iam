<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>IAM Service</title>
  <script src="/iam/~/zepto/zepto.js"></script>
  <link rel="stylesheet" href="/iam/~/iam/css/reg.css" type="text/css">
  <link rel="shortcut icon" href="/iam/~/iam/img/iam-s2-32.png" type="image/x-icon">
</head>
<body>

<div id="iam-login-frame" class="iam-reg-frame">
<div id="iam-login-box" class="iam-reg-box">

  <div class="iam-reg-msg01">{{T . "Sign in with your Account"}}</div>

  <form id="iam-login-form" class="iam-reg-form" onsubmit="iamLoginCommit();return false;">

    <input type="hidden" name="redirect_token" value="{{.redirect_token}}">

    <img class="iam-user-ico-default"  src="/iam/~/iam/img/user-default.svg">

    <div id="iam-login-form-alert" class="alert alert-info {{if eq .alert_msg nil}}hide{{end}}">{{.alert_msg}}</div>

    <div id="iam-login-input-frame">
      <div class="iam-input-row">
        <input type="text" class="iam-input iam-form-item" name="uname" id="iam-login-form-username" value="{{.uname}}" placeholder="{{T . "Username"}}">
      </div>

      <div class="iam-input-row">
        <input type="password" class="iam-input iam-form-item" name="passwd" id="iam-login-form-pwd" placeholder="{{T . "Password"}}">
      </div>

      <div class="iam-input-row">
        <button type="submit" class="iam-btn iam-form-item" onclick="iamLoginCommit()">{{T . "Sign in"}}</button>
      </div>

      <div>
        <div class="iam-input-row-checkbox">
          <input name="persistent" type="checkbox" value="1" checked="{{.persistent_checked}}"> Stay signed in
        </div>
        <div class="iam-input-row-help">
        <a href="/iam/reg/retrieve?redirect_token={{.redirect_token}}">Forgot Password ?</a>
        </div>
      </div>
    </div>
  </form>

  <div class="ref-action">
    <a href="/iam/reg/sign-up?redirect_token={{.redirect_token}}">Don't have an account ? Create Account</a>
  </div>

  <div class="footer">
    <img src="/iam/~/iam/img/iam-s2-32.png"> 
    <a href="https://github.com/hooto/iam" target="_blank">hooto IAM</a>
  </div>

</div>
</div>

<script type="text/javascript">
function innerAlert (alertid, typeUI, msg, fadeTime) {
	if (!fadeTime) {
		fadeTime = 200;
	}
    if (!typeUI) {
        return $(alertid).fadeOut(fadeTime);
    }
    var elem = $(alertid);
    if (elem) {
        elem.removeClass().addClass("alert " + typeUI).html(msg);
        elem.fadeOut(fadeTime, function() {
             elem.fadeIn(fadeTime);
        });
    }
}

function iamLoginCommit() {
    event.preventDefault();
    var alertid = "#iam-login-form-alert";
    var formData = $("#iam-login-form").serialize();

	$(".iam-form-item").attr("disabled", "disabled");
    innerAlert(alertid, 'alert-info', "pending ...", 100);
    setTimeout(function() {
        $.ajax({
            type    : "POST",
            url     : "/iam/v1/service/login-auth",
            data    : formData,
            timeout : 10000,
            success : function(data) {
    
                if (data.error) {
    	            $(".iam-form-item").removeAttr("disabled");
                    document.getElementById("iam-login-form-pwd").focus();
                    return innerAlert(alertid, 'alert-danger', data.error.message);
                }
    
                if (data.kind != "ServiceLoginAuth") {
    	            $(".iam-form-item").removeAttr("disabled");
                    return innerAlert(alertid, 'alert-danger', "Unknown Error");
                }
    
                innerAlert(alertid, 'alert-success', "Successfully Sign-on. Page redirecting ...");
                // $("#iam-login-input-frame").hide(100);
    
                window.setTimeout(function(){
                    window.location = data.redirect_uri;;
                }, 1500);
            },
            error: function(xhr, textStatus, error) {
    	        $(".iam-form-item").removeAttr("disabled");
                innerAlert(alertid, 'alert-danger', '{{T . "Internal Server Error"}}');
            }
        });
    }, 300);
}

setTimeout(function() {
    document.getElementById("iam-login-form-username").focus();
}, 200);

</script>
</body>
</html>
