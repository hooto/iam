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


#iam-reg-passreset-box {
  width: 550px;
  /*position: absolute;*/
  left: 50%;
  top: 20px;
  margin: 0 auto;
}

#iam-reg-passreset-form {
  background-color: #f7f7f7;
  border-radius: 4px;
  padding: 30px 30px 20px 30px;
  /*box-shadow: 0px 2px 2px 0px #999;*/
}

.iam-reg-passreset-msg01 {
  font-size: 20px;
  padding: 40px 0;
  text-align: center;
}

#iam-reg-passreset-form .ilf-group {
  padding: 0 0 10px 0; 
}

#iam-reg-passreset-form .ilf-input {
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

#iam-reg-passreset-form .ilf-input:focus {
  border-color: #66afe9;
  outline: 0;
  box-shadow: inset 0 1px 1px rgba(0,0,0,.075), 0 0 8px rgba(102, 175, 233, .6);
}

#iam-reg-passreset-form .ilf-btn {
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

#iam-reg-passreset-form .ilf-btn:hover {
  color: #fff;
  background-color: #3276b1;
  border-color: #285e8e;
}

#iam-reg-passreset-form .ilf-key {
  font-family: monospace;
  font-size: 16px;
  font-weight: bold;
  padding:10px 5px;
  background-color: #555;
  color: #fff;
  width: 100%;
  display: inline-block;
  border-radius: 3px;
  text-align: center;
}

#iam-reg-passreset-box .ilb-reg-passreset {
  margin: 20px 0;
  text-align: center;
  font-size: 15px;
}
#iam-reg-passreset-box .ilb-reg-passreset a {
  font-size: 16px;
  color: #fff;
}

#iam-reg-passreset-box .ilb-footer {
  text-align: center;
  margin: 20px 0;
  font-size: 14px;
}
#iam-reg-passreset-box .ilb-footer a {
  color: #777;
}
#iam-reg-passreset-box .ilb-footer img {
  width: 16px;
  height: 16px;
}
</style>


<div id="iam-reg-passreset-box">

  <div class="iam-reg-passreset-msg01">{{T . "Reset your Password"}}</div>

  {{if .pass_reset_id}}
  <form id="iam-reg-passreset-form" action="#forgot-pass">

    <div id="iam-reg-passreset-form-alert" class="alert hide"></div>

    <div class="ilf-group">
      <input type="hidden" name="id" value="{{.pass_reset_id}}">
      <span class="ilf-key">{{.pass_reset_id}}</span>
    </div>

    <div class="ilf-group">
      <input type="text" class="ilf-input" name="email" placeholder="{{T . "Confirm your Email"}}">
    </div>

    <div class="ilf-group">
      <input type="password" class="ilf-input" name="passwd" placeholder="{{T . "Your new password"}}">
    </div>

    <div class="ilf-group">
      <input type="password" class="ilf-input" name="passwd_confirm" placeholder="{{T . "Confirm your new password"}}">
    </div>    

    <div class="ilf-group">
      <button type="submit" class="ilf-btn">{{T . "Next"}}</button>
    </div>

  </form>
  {{else}}
    <div class="alert alert-danger">The Token is not valid or Expired</div>
  {{end}}

  <div class="ilb-reg-passreset">
    <a href="/iam/service/login?continue={{.continue}}&redirect_token={{.redirect_token}}">Sign in with your Account</a>
  </div>

  <div class="ilb-footer">
    <img src="/iam/~/iam/img/iam-s2-32.png"> 
    <a href="http://www.lessos.com/p/iam" target="_blank">lessOS IAM</a>
  </div>

</div>


<script>

$("input[name=email]").focus();

//
$("#iam-reg-passreset-form").submit(function(event) {

    event.preventDefault();

    l4i.InnerAlert("#iam-reg-passreset-form-alert", 'alert-info', "Pending");

    $.ajax({
        type    : "POST",
        url     : "/iam/reg/pass-reset-put",
        data    : $("#iam-reg-passreset-form").serialize(),
        timeout : 3000,
        //contentType: "application/json; charset=utf-8",
        success : function(data) {

            if (data.error) {
                return l4i.InnerAlert("#iam-reg-passreset-form-alert", 'alert-danger', data.error.message);
            }

            if (data.kind != "UserAuth") {
                return l4i.InnerAlert("#iam-reg-passreset-form-alert", 'alert-danger', "Unknown Error");
            }
                
            l4i.InnerAlert("#iam-reg-passreset-form-alert", 'alert-success', "Successfully Updated. Page redirecting");
            $(".ilf-group").hide(200);

            window.setTimeout(function(){
                window.location = "/iam/service/login?redirect_token={{.redirect_token}}";
            }, 2000);
        },
        error: function(xhr, textStatus, error) {
            l4i.InnerAlert("#iam-reg-passreset-form-alert", 'alert-danger', '{{T . "Internal Server Error"}}');
        }
    });
});

</script>

</body>
</html>
