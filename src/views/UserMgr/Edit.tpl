
<style>

.ids-um-formtable {
  width: 100%;
  border: 0;
  margin: 10px 0;
}

.ids-um-formtable td {
  padding: 5px 0;
}

</style>

<div class="panel panel-default">
  <div class="panel-heading">{{.panel_title}}</div>
  <div class="panel-body">
    <div id="o0jg5l" class="alert hide"></div>
    <form id="ids-usermgr-new-form" class="form-horizontal" action="#">
    <input type="hidden" name="uid" value="{{.uid}}">
    
    <label class="ids-form-group-title">Login Information (Required)</label>

    <div class="form-group">
      <label class="col-sm-2 control-label">Username</label>
      <div class="col-sm-10">
        <input type="text" class="form-control" name="uname" value="{{.uname}}">
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-2 control-label">Email</label>
      <div class="col-sm-10">
        <input type="text" class="form-control" name="email" value="{{.email}}">
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-2 control-label">Password</label>
      <div class="col-sm-10">
        <input type="text" class="form-control" name="passwd" value="{{.passwd}}">
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-2 control-label">Roles</label>
      <div class="col-sm-10">
        {{range .roles}}
          <span class="ids-form-checkbox">
            {{if eq .Rid "100"}}
              <input type="checkbox" name="roles" value="{{.Rid}}" checked="checked" onclick="return false"> {{.Name}}
            {{else if eq .Checked "1"}}
              <input type="checkbox" name="roles" value="{{.Rid}}" checked="checked"> {{.Name}}
            {{else}}
              <input type="checkbox" name="roles" value="{{.Rid}}"> {{.Name}}
            {{end}}
          </span>
        {{end}}
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-2 control-label">Status</label>
      <div class="col-sm-10">
        {{range $k, $v := .statuses}}
          <span class="ids-form-checkbox">
            <input type="radio" name="status" value="{{$k}}" {{if eq $k $.status}}checked="checked"{{end}}> {{$v}}
          </span>
        {{end}}
      </div>
    </div>

    <label class="ids-form-group-title">Profile Information (Optional)</label>

    <div class="form-group">
      <label class="col-sm-2 control-label">Nickname</label>
      <div class="col-sm-10">
        <input type="text" class="form-control" name="name" value="{{.name}}">
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-2 control-label">{{T . "Birthday"}}</label>
      <div class="col-sm-10">
        <input type="text" class="form-control" name="birthday" placeholder="{{T . "Example"}} : 1970-01-01" value="{{.birthday}}">
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-2 control-label">About</label>
      <div class="col-sm-10">
        <textarea class="form-control" rows="3" name="aboutme">{{.aboutme}}</textarea>
      </div>
    </div>

    <div class="form-group">
      <div class="col-sm-offset-2 col-sm-10">
        <button type="submit" class="btn btn-primary">{{T . "Submit"}}</button>
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
        url     : "/ids/user-mgr/save",
        data    : $("#ids-usermgr-new-form").serialize(),
        timeout : 3000,
        success : function(rsp) {

            var rsj = JSON.parse(rsp);

            if (rsj.status == 200) {
                
                lessAlert("#o0jg5l", 'alert-success', "Successfully saved");
                
                window.setTimeout(function(){
                    idsWorkLoader("user-mgr/list");
                }, 1500);

            } else {
                lessAlert("#o0jg5l", 'alert-danger', rsj.message);
            }
        },
        error: function(xhr, textStatus, error) {
            lessAlert("#o0jg5l", 'alert-danger', '{{T . "Internal Server Error"}}');
        }
    });
});

</script>
