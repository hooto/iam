<style type="text/css">
#qj3mr0 table td {
    padding: 5px 0;
}
</style>

<div id="in3v8z" class="alert hide"></div>

<form id="qj3mr0" action="#pass-set" method="post" >
  <table class="box" width="100%">
    <tr>
      <td width="200px" >{{T . "Current password"}}</td>
      <td ><input name="passwd_current" class="form-control" type="password"></td>
    </tr>
    <tr>
      <td>{{T . "New password"}}</td>
      <td ><input name="passwd" class="form-control" type="password"></td>
    </tr>
    <tr>
      <td>{{T . "Confirm new password"}}</td>
      <td><input name="passwd_confirm" class="form-control" type="password"></td>
    </tr>
  </table>
</form>

<script>

lessModalButtonAdd("huli6u", "{{T . "Cancel"}}", "lessModalClose()", "");
lessModalButtonAdd("of42v7", "{{T . "Save"}}", "_user_pass_set_submit()", "btn-inverse");

$("#qj3mr0").submit(function(event) {

    event.preventDefault();

    _user_pass_set_submit();
});

function _user_pass_set_submit()
{
    event.preventDefault();

    $.ajax({
        type    : "POST",
        url     : "/ids/user/pass-put?_="+Math.random(),
        data    : $("#qj3mr0").serialize(),
        timeout : 3000,
        success : function(rsp) {

            var rsj = JSON.parse(rsp);

            if (rsj.status == 200) {

                lessAlert("#in3v8z", "alert-success", rsj.message);

                window.setTimeout(function(){
                    lessModalClose();
                }, 1000);

            } else {

                lessAlert("#in3v8z", "alert-danger", rsj.message);
            }
        },
        error: function(xhr, textStatus, error) {
            lessAlert("#in3v8z", "alert-danger", "{{T . "Internal Server Error"}}");
        }
    });
}
</script>


