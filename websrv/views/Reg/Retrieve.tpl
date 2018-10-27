<!DOCTYPE html>
<html lang="en">
{{template "Common/HtmlHeaderReg.tpl" .}}
<body>

<div id="iam-resetpass-frame" class="iam-reg-frame">
<div id="iam-resetpass-box" class="iam-reg-box">

  <div class="iam-reg-msg01">{{T . "Reset your password"}}</div>

  <form id="iam-resetpass-form" class="iam-reg-form" action="#retrieve" onsubmit="retrieveCommit();return false;">

    <input type="hidden" name="redirect_token" value="{{.redirect_token}}">

    <div id="iam-resetpass-form-alert" class="alert alert-info iam-groups">
      <p>Enter the username and email address you use to sign in.</p>
      <p>The System will sent a URL to your email to reset the password.</p>
    </div>

    <div class="iam-group">
      <input type="text" class="iam-input" name="username" placeholder="{{T . "Username"}}">
    </div>

    <div class="iam-group">
      <input type="text" class="iam-input" name="email" placeholder="{{T . "Email"}}">
    </div>

    <div class="iam-group">
      <button type="submit" class="iam-btn" onclick="retrieveCommit()">{{T . "Next"}}</button>
    </div>

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

//
$("input[name=username]").focus();


//
function retrieveCommit() {
    event.preventDefault();

    var alertid = "#iam-resetpass-form-alert";

    innerAlert(alertid, 'alert-info', "Pending");

    $.ajax({
        type    : "POST",
        url     : "/iam/reg/retrieve-put",
        data    : $("#iam-resetpass-form").serialize(),
        timeout : 10000,
        //contentType: "application/json; charset=utf-8",
        success : function(data) {

            if (data.error) {
                return innerAlert(alertid, 'alert-danger', data.error.message);
            }

            if (data.kind != "UserAuth") {
                return innerAlert(alertid, 'alert-danger', "Unknown Error");
            }
                
            innerAlert(alertid, 'alert-success', "The reset URL has been sent to your mailbox, please check your email.");
            $(".iam-group").hide(200);

            // window.setTimeout(function(){
            //     window.location = "/iam/service/login";
            // }, 1500);
        },
        error: function(xhr, textStatus, error) {
            innerAlert(alertid, 'alert-danger', '{{T . "Internal Server Error"}}');
        }
    });
};

window.onload = function()
{
    $("#iam-resetpass-form").find("input[name=username]").focus();
}
</script>

</body>
</html>
