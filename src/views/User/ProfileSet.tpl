<style type="text/css">
#i1o9wr table td {
    padding: 5px 0;
    vertical-align: top;
}
</style>

<div id="h40i6o" class="alert hide"></div>

<form id="i1o9wr" action="#profile-set" method="post">
  <table class="" width="100%" style="">
    <tr>
      <td width="120px"><strong>{{T . "Nickname"}}</strong></td>
      <td>
        <input name="id" type="text" class="form-control" value="{{.login_name}}" readonly="readonly">
      </td>
    </tr>
    <tr>
      <td><strong>{{T . "Birthday"}}</strong></td>
      <td>
        <input name="birthday" type="text" class="form-control" value="{{.profile_birthday}}">
        <div>{{T . "Example"}} : 1970-01-01</div>
      </td>
    </tr>
    <tr>
      <td valign="top"><strong>{{T . "About me"}}</strong></td>
      <td>
        <textarea name="aboutme" class="form-control" rows="5">{{.profile_aboutme}}</textarea>
      </td>
    </tr>
  </table>
</form>

<script>

lessModalButtonAdd("fdp2ja", "{{T . "Cancel"}}", "lessModalClose()", "");
lessModalButtonAdd("hpvnk3", "{{T . "Save"}}", "_ids_user_proset()", "btn-inverse");

$("#i1o9wr").submit(function(event) {

    event.preventDefault();

    _ids_user_proset();
});

function _ids_user_proset()
{
    $.ajax({
        type    : "POST",
        url     : "/ids/user/profile-put?_="+Math.random(),
        data    : $("#i1o9wr").serialize(),
        timeout : 3000,
        success : function(rsp) {

            var rsj = JSON.parse(rsp);

            if (rsj.status == 200) {

                lessAlert("#h40i6o", "alert-success", rsj.message);

                window.setTimeout(function(){
                    lessModalClose();
                    window.location = "/ids/user/index";
                }, 1000);

            } else {

                lessAlert("#h40i6o", "alert-danger", rsj.message);
            }
        },
        error: function(xhr, textStatus, error) {
            lessAlert("#h40i6o", "alert-danger", "{{T . "Internal Server Error"}}");
        }
    });
}
</script>
