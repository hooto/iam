
<style>

.ids-um-formtable {
  width: 100%;
  border: 0;
  margin: 10px 0;
}

.ids-um-formtable td {
  padding: 5px 0;
}

.ids-form-checkbox {
  margin: 0 20px 0 0;
  padding: 10px 0;
}

</style>

<div class="panel panel-default">
  <div class="panel-heading">{{T . "New Account"}}</div>
  <div class="panel-body">
    <div id="ids-usermgr-new-form-alert" class="alert hide"></div>
    <form id="ids-usermgr-new-form" class="form-horizontal" action="#">
    
    <label class="ids-form-group-title">Login Information (Required)</label>

    <div class="form-group">
      <label class="col-sm-2 control-label">Username</label>
      <div class="col-sm-10">
        <input type="text" class="form-control" name="uname">
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-2 control-label">Email</label>
      <div class="col-sm-10">
        <input type="text" class="form-control" name="email">
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-2 control-label">Password</label>
      <div class="col-sm-10">
        <input type="text" class="form-control" name="passwd">
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-2 control-label">Roles</label>
      <div class="col-sm-10">
        {{if .roles}}
        {{range .roles}}
          <span class="ids-form-checkbox">
            {{if eq .rid 100}}
              <input type="checkbox" name="roles" value="{{.rid}}" sdisabled="sdisabled" checked="checked"> {{.name}}
            {{else}}
              <input type="checkbox" name="roles" value="{{.rid}}"> {{.name}}
            {{end}}
          </span>
        {{end}}
        {{end}}
      </div>
    </div>

    <label class="ids-form-group-title">Profile Information (Optional)</label>

    <div class="form-group">
      <label class="col-sm-2 control-label">Nickname</label>
      <div class="col-sm-10">
        <input type="text" class="form-control" name="name">
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-2 control-label">{{T . "Birthday"}}</label>
      <div class="col-sm-10">
        <input type="text" class="form-control" name="birthday" placeholder="{{T . "Example"}} : 1970-01-01">
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-2 control-label">About</label>
      <div class="col-sm-10">
        <textarea class="form-control" rows="3" name="aboutme"></textarea>
      </div>
    </div>

    <div class="form-group">
      <div class="col-sm-offset-2 col-sm-10">
        <button type="submit" class="btn btn-primary">{{T . "Create Account"}}</button>
      </div>
    </div>

    </form>
  </div>
</div>

<script>

//
$("input[name=uname]").focus();

//
$("#ids-usermgr-new-form").submit(function(event) {

    event.preventDefault();
    
    $.ajax({
        type    : "POST",
        url     : "/ids/user-mgr/new-save",
        data    : $("#ids-usermgr-new-form").serialize(),
        timeout : 3000,
        success : function(rsp) {

            var rsj = JSON.parse(rsp);

            if (rsj.status == 200) {
                
                lessAlert("#ids-usermgr-new-form-alert", 'alert-success', "Successfully created");
                
                window.setTimeout(function(){
                    idsWorkLoader("user-mgr/list");
                }, 1500);

            } else {
                lessAlert("#ids-usermgr-new-form-alert", 'alert-danger', rsj.message);
            }
        },
        error: function(xhr, textStatus, error) {
            lessAlert("#ids-usermgr-new-form-alert", 'alert-danger', '{{T . "Internal Server Error"}}');
        }
    });
});

</script>
