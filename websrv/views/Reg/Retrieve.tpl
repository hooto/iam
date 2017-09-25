<!DOCTYPE html>
<html lang="en">
{{template "Common/HtmlHeader.tpl" .}}
<body>

<style type="text/css">
body {
  margin: 0 auto;
  position: relative;
  font-size: 13px;
  font-family: Arial, sans-serif;
  background-color: #222;
  color: #eee;
}

#iam-resetpass-frame {
  width: 100%;
  height: 100%;
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;
  min-width: 500px;
}

#iam-resetpass-box {
  width: 550px;
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  margin: 0 auto;
}

#iam-resetpass-form {
  background-color: #f7f7f7;
  border-radius: 4px;
  padding: 30px 30px 20px 30px;
  /*box-shadow: 0px 2px 2px 0px #999;*/
}

.iam-resetpass-msg01 {
  font-size: 28px;
  padding: 40px 0;
  text-align: center;
}

#iam-resetpass-form .ilf-group {
  padding: 0 0 10px 0; 
}

#iam-resetpass-form .ilf-input {
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

#iam-resetpass-form .ilf-input:focus {
  border-color: #66afe9;
  outline: 0;
  box-shadow: inset 0 1px 1px rgba(0,0,0,.075), 0 0 8px rgba(102, 175, 233, .6);
}

#iam-resetpass-form .ilf-btn {
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

#iam-resetpass-form .ilf-btn:hover {
  color: #fff;
  background-color: #3276b1;
  border-color: #285e8e;
}

#iam-resetpass-box .ilb-resetpass {
  margin: 20px 0;
  text-align: center;
  font-size: 15px;
}
#iam-resetpass-box .ilb-resetpass a {
  font-size: 16px;
  color: #fff;
}

#iam-resetpass-box .ilb-footer {
  text-align: center;
  margin: 20px 0;
  font-size: 14px;
}
#iam-resetpass-box .ilb-footer a {
  color: #777;
}
#iam-resetpass-box .ilb-footer img {
  width: 16px;
  height: 16px;
}
</style>

<div id="iam-resetpass-frame">
<div id="iam-resetpass-box">

  <div class="iam-resetpass-msg01">{{T . "Reset your password"}}</div>

  <form id="iam-resetpass-form" action="#retrieve">

    <input type="hidden" name="redirect_token" value="{{.redirect_token}}">

    <div id="iam-resetpass-form-alert" class="alert alert-info ilf-groups">
      <p>Enter the username and email address you use to sign in.</p>
      <p>The System will sent a URL to your email to reset the password.</p>
    </div>

    <div class="ilf-group">
      <input type="text" class="ilf-input" name="username" placeholder="{{T . "Username"}}">
    </div>

    <div class="ilf-group">
      <input type="text" class="ilf-input" name="email" placeholder="{{T . "Email"}}">
    </div>

    <div class="ilf-group">
      <button type="submit" class="ilf-btn">{{T . "Next"}}</button>
    </div>

  </form>

  <div class="ilb-resetpass">
    <a href="/iam/service/login?redirect_token={{.redirect_token}}">Sign in with your Account</a>
  </div>

  <div class="ilb-footer">
    <img src="/iam/~/iam/img/iam-s2-32.png"> 
    <a href="https://github.com/hooto/iam" target="_blank">hooto IAM</a>
  </div>

</div>
</div>

<script>

//
$("input[name=username]").focus();


//
$("#iam-resetpass-form").submit(function(event) {

    event.preventDefault();

    var alertid = "#iam-resetpass-form-alert";

    l4i.InnerAlert(alertid, 'alert-info', "Pending");

    $.ajax({
        type    : "POST",
        url     : "/iam/reg/retrieve-put",
        data    : $("#iam-resetpass-form").serialize(),
        timeout : 3000,
        //contentType: "application/json; charset=utf-8",
        success : function(data) {

            if (data.error) {
                return l4i.InnerAlert(alertid, 'alert-danger', data.error.message);
            }

            if (data.kind != "UserAuth") {
                return l4i.InnerAlert(alertid, 'alert-danger', "Unknown Error");
            }
                
            l4i.InnerAlert(alertid, 'alert-success', "The reset URL has been sent to your mailbox, please check your email.");
            $(".ilf-group").hide(200);

            // window.setTimeout(function(){
            //     window.location = "/iam/service/login";
            // }, 1500);
        },
        error: function(xhr, textStatus, error) {
            l4i.InnerAlert(alertid, 'alert-danger', '{{T . "Internal Server Error"}}');
        }
    });
});

</script>

</body>
</html>
