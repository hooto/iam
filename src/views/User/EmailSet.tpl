<style type="text/css">
#obzrht table td {
    padding: 5px 0;
}
</style>

<div id="kr82np" class="alert hide"></div>

<form id="obzrht" action="#email-set" method="post" >
  <table width="100%">
    <tr>
      <td width="160px">{{T . "Email"}}</td>
      <td >
        <input name="email" class="form-control" type="text" value="{{.login_email}}"></td>
    </tr>
    <tr>
      <td>{{T . "Password"}}</td>
      <td ><input name="passwd" class="form-control" type="password" value=""></td>
    </tr>

  </table>
</form>



<script>

$("input[name=email]").focus();


lessModalButtonAdd("dh23d9", "{{T . "Cancel"}}", "lessModalClose()", "");
lessModalButtonAdd("wtr5qw", "{{T . "Submit"}}", "_user_email_set_submit()", "btn-inverse");

$("#obzrht").submit(function(event) {

    event.preventDefault();

    _user_email_set_submit();
});

function _user_email_set_submit()
{
    $.ajax({
        type    : "POST",
        url     : "/ids/user/email-put?_="+Math.random(),
        data    : $("#obzrht").serialize(),
        timeout : 3000,
        success : function(rsp) {

            var rsj = JSON.parse(rsp);

            if (rsj.status == 200) {

                lessAlert("#kr82np", "alert-success", rsj.message);

                window.setTimeout(function(){
                    lessModalClose();
                    window.location = "/ids/user/index";
                }, 1000);

            } else {

                lessAlert("#kr82np", "alert-danger", rsj.message);
            }
        },
        error: function(xhr, textStatus, error) {
            lessAlert("#kr82np", "alert-danger", "{{T . "Internal Server Error"}}");
        }
    });
}
</script>

