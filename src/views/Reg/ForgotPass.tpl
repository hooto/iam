{{template "Common/HtmlHeader.tpl" .}}

<style>
body {
  margin: 0 auto;
  position: relative;
  font-size: 13px;
  font-family: Arial, sans-serif;
  background-color: #fff;
}

#ids-resetpass-box {
  width: 400px;
  position: absolute;
  left: 50%;
  top: 20%;
  margin-left: -200px;
  color: #555;
}

#ids-resetpass-form {
  background-color: #f7f7f7;
  border-radius: 4px;
  padding: 30px 30px 20px 30px;
  box-shadow: 0px 2px 2px 0px #999;
}

.ids-resetpass-msg01 {
  font-size: 20px;
  margin: 20px 0;
  text-align: center;
}

#ids-resetpass-form .ilf-group {
  padding: 0 0 10px 0; 
}

#ids-resetpass-form .ilf-input {
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

#ids-resetpass-form .ilf-input:focus {
  border-color: #66afe9;
  outline: 0;
  box-shadow: inset 0 1px 1px rgba(0,0,0,.075), 0 0 8px rgba(102, 175, 233, .6);
}

#ids-resetpass-form .ilf-btn {
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

#ids-resetpass-form .ilf-btn:hover {
  color: #fff;
  background-color: #3276b1;
  border-color: #285e8e;
}

#ids-resetpass-box .ilb-resetpass {
  margin: 10px 0;
  text-align: center;
  font-size: 15px;
}
#ids-resetpass-box .ilb-resetpass a {
  color: #427fed;
}

#ids-resetpass-box .ilb-footer {
  text-align: center;
  margin: 20px 0;
  font-size: 14px;
}
#ids-resetpass-box .ilb-footer a {
  color: #777;
}
#ids-resetpass-box .ilb-footer img {
  width: 16px;
  height: 16px;
}
</style>

<div id="ids-resetpass-box">

  <div class="ids-resetpass-msg01">{{T . "Reset your password"}}</div>

  <form id="ids-resetpass-form" action="#forgot-pass">

    <input type="hidden" name="continue" value="{{.continue}}">

    <div id="ids-resetpass-form-alert" class="alert alert-info ilf-groups">
    Enter the email address you use to sign in, the System will sent a URL to your email to reset the password.
    </div>

    <div class="ilf-group">
      <input type="text" class="ilf-input" name="email" placeholder="{{T . "Email"}}">
    </div>

    <div class="ilf-group">
      <button type="submit" class="ilf-btn">{{T . "Next"}}</button>
    </div>

  </form>

  <div class="ilb-resetpass">
    <a href="/ids/service/login?continue={{.continue}}">Sign in with your Account</a>
  </div>

  <div class="ilb-footer">
    <img src="/ids/~/ids/img/ids-s2-32.png"> 
    <a href="http://www.lesscompute.com" target="_blank">less Identity Server</a>
  </div>

</div>

<script>

//
$("input[name=name]").focus();

//
var ids_eh = $("#ids-resetpass-box").height();
$("#ids-resetpass-box").css({
    "top": "50%",
    "margin-top": - (ids_eh / 2) + "px" 
});

//
$("#ids-resetpass-form").submit(function(event) {

    event.preventDefault();

    lessAlert("#ids-resetpass-form-alert", 'alert-info', "Pending");

    $.ajax({
        type    : "POST",
        url     : "/ids/reg/forgot-pass-put",
        data    : $("#ids-resetpass-form").serialize(),
        timeout : 3000,
        //contentType: "application/json; charset=utf-8",
        success : function(rsp) {

            var rsj = JSON.parse(rsp);
            //console.log(rsp);

            if (rsj.status == 200) {
                
                $(".ilf-group").hide(600);

                lessAlert("#ids-resetpass-form-alert", 'alert-success', "The reset URL has been sent to your mailbox, please check your email");
                
                window.setTimeout(function(){
                    window.location = "/ids/service/login?continue={{.continue}}";
                }, 10000);

            } else {
                lessAlert("#ids-resetpass-form-alert", 'alert-danger', rsj.message);
            }
        },
        error: function(xhr, textStatus, error) {
            lessAlert("#ids-resetpass-form-alert", 'alert-danger', '{{T . "Internal Server Error"}}');
        }
    });
});

</script>

{{template "Common/HtmlFooter.tpl" .}}
