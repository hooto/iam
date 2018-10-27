<!DOCTYPE html>
<html lang="en">
{{template "Common/HtmlHeaderReg.tpl" .}}
<body>

<div id="iam-signup-frame" class="iam-reg-frame">
<div id="iam-signup-box" class="iam-reg-box">

  <div class="iam-reg-msg01">{{T . "Create your Account"}}</div>

  <form id="iam-signup-form" class="iam-reg-form" action="#" onsubmit="signupCommit();return false;">

    <div id="iam-signup-form-alert" class="alert hide iam-group"></div>

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
      <button type="submit" class="iam-btn" onclick="signupCommit()">{{T . "Create Account"}}</button>
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
function innerAlert (alertid, type_ui, msg) {
    if (!type_ui) {
        return $(alertid).fadeOut(200);
    }
    var elem = $(alertid);
    if (elem) {
        elem.removeClass().addClass("alert " + type_ui).html(msg);
        elem.fadeOut(200, function() {
             elem.fadeIn(200);
        });
    }
}

function signupCommit() {
    event.preventDefault();
    var alertid = "#iam-signup-form-alert";
    $.ajax({
        type    : "POST",
        url     : "/iam/reg/sign-up-reg",
        data    : $("#iam-signup-form").serialize(),//JSON.stringify(req),
        timeout : 3000,
        success : function(data) {

            if (data.error) {
                return innerAlert(alertid, 'alert-danger', data.error.message);
            }

            if (data.kind != "User") {
                return innerAlert(alertid, 'alert-danger', "Unknown Error");
            }
                
            innerAlert(alertid, 'alert-success', "Successfully registration. Page redirecting");
            $(".iam-group").hide(600);

            window.setTimeout(function(){
                window.location = "/iam/service/login?redirect_uri={{.redirect_token}}";
            }, 1500);
        },
        error: function(xhr, textStatus, error) {
            innerAlert(alertid, 'alert-danger', '{{T . "Internal Server Error"}}');
        }
    });
}

window.onload = function()
{
    $("#iam-signup-form").find("input[name=uname]").focus();
}

</script>

</body>
</html>
