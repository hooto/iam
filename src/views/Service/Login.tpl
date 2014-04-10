{{template "Common/HtmlHeader.tpl" .}}

<style>
body {
  margin: 0 auto;
  position: relative;
  font-size: 13px;
  font-family: Arial, sans-serif;
  background-color: #fff;
}

#ids-login-box {
  width: 360px;
  position: absolute;
  left: 50%;
  top: 20%;
  margin-left: -180px;
  color: #555;
}

#ids-login-form {
  background-color: #f7f7f7;
  border-radius: 4px;
  padding: 20px 30px 20px 30px;
  box-shadow: 0px 2px 2px 0px #999;
}

.ids-login-msg01 {
  font-size: 20px;
  margin: 20px 0;
  text-align: center;
}

.ids-user-ico-default {
  width: 80px;
  height: 80px;
  padding: 2px;
  position: relative;
  left: 50%;
  margin: 0 0 30px -40px;
  border-radius: 50%;
  background-color: #dfdfdf;
}

#ids-login-form .ilf-group {
  padding: 0 0 10px 0; 
}

#ids-login-form .ilf-input {
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

#ids-login-form .ilf-input:focus {
  border-color: #66afe9;
  outline: 0;
  box-shadow: inset 0 1px 1px rgba(0,0,0,.075), 0 0 8px rgba(102, 175, 233, .6);
}

#ids-login-form .ilf-btn {
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

#ids-login-form .ilf-btn:hover {
  color: #fff;
  background-color: #3276b1;
  border-color: #285e8e;
}

#ids-login-form .ilf-checkbox {
  display: inline-block;
  min-height: 20px;
  font-weight: normal;
  color: #555;
}
#ids-login-form .ilf-checkbox input[type="checkbox"] {
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

#ids-login-form .ilf-help {
  display: inline-block;
  float: right;
}
#ids-login-form .ilf-help a {
  color: #427fed;
}

#ids-login-box .ilb-signup {
  margin: 10px 0;
  text-align: center;
  font-size: 15px;
}
#ids-login-box .ilb-signup a {
  color: #427fed;
}

#ids-login-box .ilb-footer {
  text-align: center;
  margin: 20px 0;
  font-size: 14px;
}
#ids-login-box .ilb-footer a {
  color: #777;
}
#ids-login-box .ilb-footer img {
  width: 16px;
  height: 16px;
}
</style>

<div id="ids-login-box">

  <div class="ids-login-msg01">{{T . "Sign in with your Account"}}</div>

  <form id="ids-login-form" action="#">

    <input type="hidden" name="continue" value="{{.continue}}">

    <img class="ids-user-ico-default"  src="/ids/~/ids/img/user-default.png">

    <div id="ids-login-form-alert" class="alert hide ilf-group"></div>

    <div id="ilf-grp-input">
    <div class="ilf-group">
      <input type="text" class="ilf-input" name="email" value="{{.email}}" placeholder="{{T . "Email"}}">
    </div>

    <div class="ilf-group">
      <input type="password" class="ilf-input" name="passwd" placeholder="{{T . "Password"}}">
    </div>

    <div class="ilf-group">
      <button type="submit" class="ilf-btn">{{T . "Sign in"}}</button>
    </div>

    <div>
      <div class="ilf-checkbox">
        <input name="persistent" type="checkbox" value="1" checked="{{.persistentChecked}}"> Stay signed in
      </div>
      <div class="ilf-help">
        <a href="/ids/reg/forgot-pass">Forgot Password?</a>
      </div>
    </div>
    </div>
  </form>

  <div class="ilb-signup">
    <a href="/ids/reg/sign-up?continue={{.continue}}">Create an account</a>
  </div>

  <div class="ilb-footer">
    <img src="/ids/~/ids/img/ids-s2-32.png"> 
    <a href="http://www.lesscompute.com" target="_blank">less Identity Server</a>
  </div>

</div>

<script>

//
$("input[name=email]").focus();

//
var ids_eh = $("#ids-login-box").height();
$("#ids-login-box").css({
    "top": "50%",
    "margin-top": - (ids_eh / 2) + "px" 
});

//
$("#ids-login-form").submit(function(event) {

    event.preventDefault();

    /* var req = {
        data: {
            "email": $("input[name=email]").val(),
            "passwd": $("input[name=passwd]").val(),
            "continue": $("input[name=continue]").val(),
            "persistent": $("input[name=persistent]").val(),
        }
    } */

    $.ajax({
        type    : "POST",
        url     : "/ids/service/login-auth",
        data    : $(this).serialize(),//JSON.stringify(req),
        timeout : 3000,
        //contentType: "application/json; charset=utf-8",
        success : function(rsp) {

            var rsj = JSON.parse(rsp);
            //console.log(rsp);

            if (rsj.status != 200) {
                lessAlert("#ids-login-form-alert", 'alert-danger', rsj.message);
                return;
            }

            lessAlert("#ids-login-form-alert", 'alert-success', "Successfully Sign-on. Page redirecting");
            $("#ilf-grp-input").hide(200);
            
            lessCookie.Set("access_token", rsj.data.access_token, 864000);

            window.setTimeout(function(){    
                window.location = rsj.data.continue;
            }, 1500);
        },
        error: function(xhr, textStatus, error) {
            lessAlert("#ids-login-form-alert", 'alert-danger', '{{T . "Internal Server Error"}}');
        }
    });
});

</script>

{{template "Common/HtmlFooter.tpl" .}}
