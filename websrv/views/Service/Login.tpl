{{template "Common/HtmlHeader.tpl" .}}

<style type="text/css">
body {
  margin: 0 !important;
  padding: 0;
  position: relative;
  font-size: 13px;
  font-family: Arial, sans-serif;
  background-color: #222;
  /*background-color: #07c;*/
  color: #eee;
}

#iam-login-box {
  width: 360px;
  /*position: absolute;*/
  left: 50%;
  top: 20px;
  color: #555;
  margin: 0 auto;
}

#iam-login-form {
  background-color: #f7f7f7;
  border-radius: 4px;
  padding: 20px 30px 20px 30px;
  /*box-shadow: 0px 2px 2px 0px #999;*/
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

#iam-login-form .ilf-group {
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

#iam-login-form .ilf-checkbox {
  display: inline-block;
  min-height: 20px;
  font-weight: normal;
  color: #555;
}
#iam-login-form .ilf-checkbox input[type="checkbox"] {
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

#iam-login-form .ilf-help {
  display: inline-block;
  float: right;
}
#iam-login-form .ilf-help a {
  color: #427fed;
}

#iam-login-box .ilb-signup {
  margin: 20px 0;
  text-align: center;
  font-size: 15px;
}
#iam-login-box .ilb-signup a {
  font-size: 16px;
  color: #fff;
}

#iam-login-box .ilb-footer {
  text-align: center;
  margin: 20px 0;
  font-size: 14px;
}
#iam-login-box .ilb-footer a {
  color: #ccc;
}
#iam-login-box .ilb-footer img {
  width: 16px;
  height: 16px;
}
</style>

<div id="iam-login-box">

  <div class="iam-login-msg01">{{T . "Sign in with your Account"}}</div>

  <form id="iam-login-form" onsubmit="return false;">

    <input type="hidden" name="redirect_uri" value="{{.redirect_uri}}">
    <input type="hidden" name="state" value="{{.state}}">

    <img class="iam-user-ico-default"  src="/iam/~/iam/img/user-default.png">

    <div id="iam-login-form-alert" class="alert ilf-group {{if eq .alert_msg nil}}hide{{end}}">{{.alert_msg}}</div>

    <div id="ilf-grp-input">
      <div class="ilf-group">
        <input type="text" class="ilf-input" name="uname" value="{{.uname}}" placeholder="{{T . "Username"}}">
      </div>

      <div class="ilf-group">
        <input type="password" class="ilf-input" name="passwd" placeholder="{{T . "Password"}}">
      </div>

      <div class="ilf-group">
        <button type="submit" class="ilf-btn">{{T . "Sign in"}}</button>
      </div>

      <div>
        <div class="ilf-checkbox">
          <input name="persistent" type="checkbox" value="1" checked="{{.persistent_checked}}"> Stay signed in
        </div>
        <div class="ilf-help">
        <a href="/iam/reg/forgot-pass">Forgot Password?</a>
        </div>
      </div>
    </div>
  </form>

  <div class="ilb-signup">
    <a href="/iam/reg/sign-up?redirect_uri={{.redirect_uri}}">Don't have an account? Create Account</a>
  </div>

  <div class="ilb-footer">
    <img src="/iam/~/iam/img/iam-s2-32.png"> 
    <a href="http://www.lessos.com/p/iam" target="_blank">lessOS IAM</a>
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

        iam.Ajax("/iam/v1/service/login-auth", {
            type    : "POST",
            data    : $(this).serialize(),
            timeout : 3000,
            success : function(rsj) {

                if (!rsj || !rsj.kind || rsj.kind != "ServiceLoginAuth") {

                    if (rsj.error) {
                        return l4i.InnerAlert(alertid, 'alert-danger', rsj.error.message);
                    }

                    return l4i.InnerAlert(alertid, 'alert-danger', "Network Connection Exception");
                }
    
                l4i.InnerAlert(alertid, 'alert-success', "Successfully Sign-on. Page redirecting ...");
                $("#ilf-grp-input").hide(100);
    
                window.setTimeout(function() {
                    window.location = rsj.redirect_uri;
                }, 1500);
            },
            error: function(xhr, textStatus, error) {
                l4i.InnerAlert(alertid, 'alert-danger', '{{T . "Internal Server Error"}}');
            }
        });
    });
}
</script>

{{template "Common/HtmlFooter.tpl" .}}
