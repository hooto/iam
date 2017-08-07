<!DOCTYPE html>
<html lang="en">
{{template "Common/HtmlHeader.tpl" .}}
<body>

<style type="text/css">
body {
  margin: 0 auto !important;
  position: relative;
  font-size: 13px;
  font-family: Arial, sans-serif;
  background-color: #222;
  color: #eee;
  min-width: 500px;
  display: block;
}

#iam-signup-frame {
  width: 100%;
  height: 100%;
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;
  min-width: 500px;
  min-height: 400px;
}

#iam-signup-box {
  width: 500px;
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  margin: 0 auto;
}

#iam-signup-form {
  background-color: #f7f7f7;
  border-radius: 4px;
  padding: 30px 30px 20px 30px;
  /*box-shadow: 0px 2px 2px 0px #999;*/
}

.iam-signup-msg01 {
  font-size: 28px;
  padding: 40px 0;
  text-align: center;
}

#iam-signup-form .ilf-group {
  padding: 0 0 10px 0; 
}

#iam-signup-form .ilf-input {
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

#iam-signup-form .ilf-input:focus {
  border-color: #66afe9;
  outline: 0;
  box-shadow: inset 0 1px 1px rgba(0,0,0,.075), 0 0 8px rgba(102, 175, 233, .6);
}

#iam-signup-form .ilf-btn {
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

#iam-signup-form .ilf-btn:hover {
  color: #fff;
  background-color: #3276b1;
  border-color: #285e8e;
}

#iam-signup-box .ilb-signup {
  margin: 20px 0;
  text-align: center;
  font-size: 15px;
}
#iam-signup-box .ilb-signup a {
  font-size: 16px;
  color: #fff;
}

#iam-signup-box .ilb-footer {
  text-align: center;
  margin: 20px 0;
  font-size: 14px;
}
#iam-signup-box .ilb-footer a {
  color: #777;
}
#iam-signup-box .ilb-footer img {
  width: 16px;
  height: 16px;
}
</style>

<div id="iam-signup-frame">
<div id="iam-signup-box">

  <div class="iam-signup-msg01">{{T . "Create your Account"}}</div>

  <form id="iam-signup-form" action="#">

    <input type="hidden" name="continue" value="{{.continue}}">

    <div id="iam-signup-form-alert" class="alert hide ilf-group"></div>

    {{if eq .user_reg_disable false }}
    <div class="ilf-group">
      <input type="text" class="ilf-input" name="uname" value="{{.uname}}" placeholder="{{T . "Unique Username"}}">
    </div>

    <div class="ilf-group">
      <input type="text" class="ilf-input" name="email" placeholder="{{T . "Email"}}">
    </div>

    <div class="ilf-group">
      <input type="password" class="ilf-input" name="passwd" placeholder="{{T . "Password"}}">
    </div>

    <div class="ilf-group">
      <button type="submit" class="ilf-btn">{{T . "Create Account"}}</button>
    </div>
    {{else}}
    <div class="alert alert-danger">User registration was closed!<br>Please contact the administrator to manually register accounts</div>
    {{end}}

  </form>

  <div class="ilb-signup">
    <a href="/iam/service/login?continue={{.continue}}&redirect_token={{.redirect_token}}">Sign in with your Account</a>
  </div>

  <div class="ilb-footer">
    <img src="/iam/~/iam/img/iam-s2-32.png"> 
    <a href="https://github.com/lessos/iam" target="_blank">lessOS IAM</a>
  </div>

</div>
</div>

<script>
$("input[name=name]").focus();
$("#iam-signup-form").submit(function(event) {
    event.preventDefault();
    var alertid = "#iam-signup-form-alert";
    $.ajax({
        type    : "POST",
        url     : "/iam/reg/sign-up-reg",
        data    : $("#iam-signup-form").serialize(),//JSON.stringify(req),
        timeout : 3000,
        success : function(data) {

            if (data.error) {
                return l4i.InnerAlert(alertid, 'alert-danger', data.error.message);
            }

            if (data.kind != "User") {
                return l4i.InnerAlert(alertid, 'alert-danger', "Unknown Error");
            }
                
            l4i.InnerAlert(alertid, 'alert-success', "Successfully registration. Page redirecting");
            $(".ilf-group").hide(600);

            window.setTimeout(function(){
                window.location = "/iam/service/login?continue={{.continue}}&redirect_uri={{.redirect_token}}";
            }, 1500);
        },
        error: function(xhr, textStatus, error) {
            l4i.InnerAlert(alertid, 'alert-danger', '{{T . "Internal Server Error"}}');
        }
    });
});
</script>

</body>
</html>
