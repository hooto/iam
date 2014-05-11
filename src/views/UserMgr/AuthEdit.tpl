<style>
.page-header {
  margin: 10px 0;
  font-height: 100%;
}
</style>

<div class="page-header">
  <h2>{{.panel_title}} <small></small></h2>
</div>

    <div id="q6p2le" class="alert hide"></div>
    
    <form id="vukihr" class="form-horizontal" action="#">
    
    <input type="hidden" name="instid" value="{{.id}}">
    
    <label class="ids-form-group-title">Application Information</label>

    <div class="form-group">
      <label class="col-sm-3 control-label">Instance ID</label>
      <div class="col-sm-9">
        <p class="form-control-static">{{.id}}</p>
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-3 control-label">Application ID</label>
      <div class="col-sm-9">
        <p class="form-control-static">{{.app_id}}</p>
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-3 control-label">Application Name</label>
      <div class="col-sm-9">
        <input type="text" class="form-control" name="app_title" value="{{.app_title}}">
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-3 control-label">Version</label>
      <div class="col-sm-9">
        <p class="form-control-static">{{.version}}</p>
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-3 control-label">Status</label>
      <div class="col-sm-9">
        {{range $k, $v := .statuses}}
          <span class="ids-form-checkbox">
            <input type="radio" name="status" value="{{$k}}" {{if eq $k $.status}}checked="checked"{{end}}> {{$v}}
          </span>
        {{end}}
      </div>
    </div>

    {{if .privileges}}
    <label class="ids-form-group-title">Privilege Information</label>

    <div class="form-group">
      <label class="col-sm-3 control-label">Privileges</label>
      <div class="col-sm-9">
        <table class="table">
        <thead>
          <tr>
            <th>#</th>
            <th>Privilege</th>
            <th>Description</th>
          </tr>
        </thead>
        <tbody>
          {{range $pid, $priv := .privileges}}
          <tr>
            <td>{{$priv.pid}}</td>
            <td>{{$priv.privilege}}</td>
            <td>{{$priv.desc}}</td>
          </tr>
          {{end}}
        </tbody>
        </table>
      </div>
    </div>
    {{end}}

    <div class="form-group">
      <div class="col-sm-offset-3">
        <button type="submit" class="btn btn-primary">{{T . "Submit"}}</button>
      </div>
    </div>

    </form>

<script>

//
$("#vukihr").submit(function(event) {

    event.preventDefault();
    
    $("button[type=submit]", this).attr('disabled', 'disabled');

    $.ajax({
        type    : "POST",
        url     : "/ids/user-mgr/auth-save",
        data    : $("#vukihr").serialize(),
        timeout : 3000,
        success : function(rsp) {

            var rsj = JSON.parse(rsp);

            if (rsj.status == 200) {
                
                lessAlert("#q6p2le", 'alert-success', "Successfully saved");
                
                window.setTimeout(function(){
                    idsWorkLoader("user-mgr/auth-list");
                }, 1500);

            } else {
                lessAlert("#q6p2le", 'alert-danger', rsj.message);
                $("button[type=submit]", this).removeAttr('disabled');
            }
        },
        error: function(xhr, textStatus, error) {
            lessAlert("#q6p2le", 'alert-danger', '{{T . "Internal Server Error"}}');
            $("button[type=submit]", this).removeAttr('disabled');
        }
    });
});

</script>
