
<style>
.form-horizontal {
  margin: 0 15px;
  padding: 2px;
}
.form-group {
  margin-bottom: 5px;
}

.iam-um-formtable {
  width: 100%;
  border: 0;
  margin: 10px 0;
}

.iam-um-formtable td {
  padding: 5px 0;
}

.iam-um-role-inst {
  text-align: right;
}
.iam-umri-title {
  font-weight: bold;
  text-align: right;
}
.iam-umri-attr {
  float: right;
  color: #555;
}
.iam-umri-attr td {
  padding-left: 20px;
}
.iam-umri-attr-ctn {
  stext-align: left;
}

.r0330s .item {
    position: relative;
    width: 200px;
    font-size: 12px;
    float: left;
    margin: 3px 10px 3px 0;
}

.r0330s .item input {
    margin-bottom: 0;
}

</style>

<div id="iam-usermgr-roleset-alert" class="alert hide"></div>
    
<div id="iam-usermgr-roleset" class="form-horizontal">
    <input type="hidden" name="roleid" value="{[=it.meta.id]}">
    
    <label class="iam-form-group-title">Role Information</label>

    <div class="form-group">
      <label class="col-sm-3 control-label">Name</label>
      <div class="col-sm-9">
        <input type="text" class="form-control" name="name" value="{[=it.meta.name]}">
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-3 control-label">Description</label>
      <div class="col-sm-9">
        <input type="text" class="form-control" name="desc" value="{[=it.desc]}">
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-3 control-label">Status</label>
      <div class="col-sm-9">
        {[~it._statusls :v]}
          <span class="iam-form-checkbox">
            <input type="radio" name="status" value="{[=v.status]}" {[ if (v.status == it.status) { ]}checked="checked"{[ } ]}> {[=v.title]}
          </span>
        {[~]}
      </div>
    </div>

    <!-- <label class="iam-form-group-title">Privileges</label>

    {{range $insid, $inst := .instances}}
    <div class="form-group">
      <div class="col-sm-3 iam-um-role-inst">
        <div class="iam-umri-title">{{$inst.AppTitle}}</div>
        <table class="iam-umri-attr">
          <tr><td>Instance ID:</td><td class="iam-umri-attr-ctn">{{$inst.InstanceId}}</td></tr>
          <tr><td>Version:</td><td class="iam-umri-attr-ctn">{{$inst.Version}}</td></tr>
        </table>
      </div>
      <div class="col-sm-9 r0330s">
        {{range $pid, $priv  := $inst.Privileges}}
        <label class="item">
          <input type="checkbox" name="privileges" value="{{$pid}}" {{if $priv.Checked}}checked="checked"{{end}}> {{$priv.Desc}}
        </label>
        {{end}}
      </div>
    </div>
    {{end}} -->
</div>

<script>

// //
// $("#vukihr").submit(function(event) {

//     event.preventDefault();
    
//     $("button[type=submit]", this).attr('disabled', 'disabled');

//     $.ajax({
//         type    : "POST",
//         url     : "/iam/user-mgr/role-save",
//         data    : $("#vukihr").serialize(),
//         timeout : 3000,
//         success : function(rsp) {

//             var rsj = JSON.parse(rsp);

//             if (rsj.status == 200) {
                
//                 lessAlert("#iam-usermgr-roleset-alert", 'alert-success', "Successfully saved");
                
//                 window.setTimeout(function(){
//                     iamWorkLoader("user-mgr/role-list");
//                 }, 1500);

//             } else {
//                 lessAlert("#iam-usermgr-roleset-alert", 'alert-danger', rsj.message);
//                 $("button[type=submit]", this).removeAttr('disabled');
//             }
//         },
//         error: function(xhr, textStatus, error) {
//             lessAlert("#iam-usermgr-roleset-alert", 'alert-danger', '{{T . "Internal Server Error"}}');
//             $("button[type=submit]", this).removeAttr('disabled');
//         }
//     });
// });

</script>
