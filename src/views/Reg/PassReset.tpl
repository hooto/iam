{{template "Common/HtmlHeader.tpl" .}}


<style type="text/css">
body {
  margin: 0 auto;
  position: relative;
  font-size: 13px;
  font-family: Arial, sans-serif;
  background-color: #fff;
}


#ids-reg-passreset-box {
  width: 480px;
  position: absolute;
  left: 50%;
  top: 20%;
  margin-left: -240px;
  color: #555;
}

#ids-reg-passreset-form {
  background-color: #f7f7f7;
  border-radius: 4px;
  padding: 30px 30px 20px 30px;
  box-shadow: 0px 2px 2px 0px #999;
}

.ids-reg-passreset-msg01 {
  font-size: 20px;
  margin: 20px 0;
  text-align: center;
}

#ids-reg-passreset-form .ilf-group {
  padding: 0 0 10px 0; 
}

#ids-reg-passreset-form .ilf-input {
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

#ids-reg-passreset-form .ilf-input:focus {
  border-color: #66afe9;
  outline: 0;
  box-shadow: inset 0 1px 1px rgba(0,0,0,.075), 0 0 8px rgba(102, 175, 233, .6);
}

#ids-reg-passreset-form .ilf-btn {
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

#ids-reg-passreset-form .ilf-btn:hover {
  color: #fff;
  background-color: #3276b1;
  border-color: #285e8e;
}

#ids-reg-passreset-form .ilf-key {
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

#ids-reg-passreset-box .ilb-reg-passreset {
  margin: 10px 0;
  text-align: center;
  font-size: 15px;
}
#ids-reg-passreset-box .ilb-reg-passreset a {
  color: #427fed;
}

#ids-reg-passreset-box .ilb-footer {
  text-align: center;
  margin: 20px 0;
  font-size: 14px;
}
#ids-reg-passreset-box .ilb-footer a {
  color: #777;
}
#ids-reg-passreset-box .ilb-footer img {
  width: 16px;
  height: 16px;
}
</style>


<div id="ids-reg-passreset-box">

  <div class="ids-reg-passreset-msg01">{{T . "Reset your password"}}</div>

  {{if .pass_reset_id}}
  <form id="ids-reg-passreset-form" action="#forgot-pass">

    <div id="ids-reg-passreset-form-alert" class="alert hide"></div>

    <div class="ilf-group">
      <input type="hidden" name="id" value="{{.pass_reset_id}}">
      <span class="ilf-key">{{.pass_reset_id}}</span>
    </div>

    <div class="ilf-group">
      <input type="text" class="ilf-input" name="email" placeholder="{{T . "Confirm your Email"}}">
    </div>

    <div class="ilf-group">
      <input type="text" class="ilf-input" name="birthday" placeholder="{{T . "Confirm your birthday. Example: 1970-01-01"}}">
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
    <a href="/ids/service/login?continue={{.continue}}">Sign in with your Account</a>
  </div>

  <div class="ilb-footer">
    <img src="/ids/~/ids/img/ids-s2-32.png"> 
    <a href="http://www.lesscompute.com" target="_blank">less Identity Server</a>
  </div>

</div>


<script>

$("input[name=name]").focus();

//
var ids_eh = $("#ids-reg-passreset-box").height();
$("#ids-reg-passreset-box").css({
    "top": "50%",
    "margin-top": - (ids_eh / 2) + "px" 
});

//
$("#ids-reg-passreset-form").submit(function(event) {

    event.preventDefault();

    lessAlert("#ids-reg-passreset-form-alert", 'alert-info', "Pending");

    $.ajax({
        type    : "POST",
        url     : "/ids/reg/pass-reset-put",
        data    : $("#ids-reg-passreset-form").serialize(),
        timeout : 3000,
        //contentType: "application/json; charset=utf-8",
        success : function(rsp) {

            var rsj = JSON.parse(rsp);
            console.log(rsp);

            if (rsj.status == 200) {
                
                $(".ilf-group").hide(600);

                lessAlert("#ids-reg-passreset-form-alert", 'alert-success', "Successfully Updated");
                
                window.setTimeout(function(){
                    window.location = "/ids/service/login";
                }, 5000);

            } else {
                lessAlert("#ids-reg-passreset-form-alert", 'alert-danger', rsj.message);
            }
        },
        error: function(xhr, textStatus, error) {
            lessAlert("#ids-reg-passreset-form-alert", 'alert-danger', '{{T . "Internal Server Error"}}');
        }
    });
});

</script>


{{template "Common/HtmlFooter.tpl" .}}
