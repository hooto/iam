<!DOCTYPE html>
<html lang="en">
{{template "Common/HtmlHeaderReg.tpl" .}}
<body>


<div id="iam-reg-passreset-frame" class="iam-reg-frame">
<div id="iam-reg-passreset-box" class="iam-reg-box">

  <div class="iam-reg-msg01">{{T . "Reset your Password"}}</div>

  {{if .pass_reset_id}}
  <form id="iam-reg-passreset-form"  class="iam-reg-form" action="#forgot-pass">

    <input type="hidden" name="id" value="{{.pass_reset_id}}">
    <div class="iam-key alert alert-dark">{{.pass_reset_id}}</div>

    <div id="iam-reg-passreset-form-alert" class="alert hide"></div>

    <div class="iam-group">
      <input type="text" class="iam-input" name="email" placeholder="{{T . "Confirm your Email"}}">
    </div>

    <div class="iam-group">
      <input type="password" class="iam-input" name="passwd" placeholder="{{T . "Your new password"}}">
    </div>

    <div class="iam-group">
      <input type="password" class="iam-input" name="passwd_confirm" placeholder="{{T . "Confirm your new password"}}">
    </div>    

    <div class="iam-group">
      <button type="submit" class="iam-btn">{{T . "Next"}}</button>
    </div>

  </form>
  {{else}}
    <div class="alert alert-danger">The Token is not valid or Expired</div>
  {{end}}

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

$("input[name=email]").focus();

//
$("#iam-reg-passreset-form").submit(function(event) {

    event.preventDefault();

    innerAlert("#iam-reg-passreset-form-alert", 'alert-info', "Pending");

    $.ajax({
        type    : "POST",
        url     : "/iam/reg/pass-reset-put",
        data    : $("#iam-reg-passreset-form").serialize(),
        timeout : 3000,
        success : function(data) {

            if (data.error) {
                return innerAlert("#iam-reg-passreset-form-alert", 'alert-danger', data.error.message);
            }

            if (data.kind != "UserAuth") {
                return innerAlert("#iam-reg-passreset-form-alert", 'alert-danger', "Unknown Error");
            }
                
            innerAlert("#iam-reg-passreset-form-alert", 'alert-success', "Successfully Updated. Page redirecting");
            $(".iam-group").hide(200);

            window.setTimeout(function(){
                window.location = "/iam/service/login?redirect_token={{.redirect_token}}";
            }, 2000);
        },
        error: function(xhr, textStatus, error) {
            innerAlert("#iam-reg-passreset-form-alert", 'alert-danger', '{{T . "Internal Server Error"}}');
        }
    });
});

</script>

</body>
</html>
