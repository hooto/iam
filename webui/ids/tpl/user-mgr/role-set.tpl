
<style>
.form-horizontal {
  margin: 0 15px;
  padding: 2px;
}
.form-group {
  margin-bottom: 5px;
}

.ids-um-formtable {
  width: 100%;
  border: 0;
  margin: 10px 0;
}

.ids-um-formtable td {
  padding: 5px 0;
}

.ids-um-role-inst {
  text-align: right;
}
.ids-umri-title {
  font-weight: bold;
  text-align: right;
}
.ids-umri-attr {
  float: right;
  color: #555;
}
.ids-umri-attr td {
  padding-left: 20px;
}
.ids-umri-attr-ctn {
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

<div id="ids-usermgr-roleset-alert" class="alert hide"></div>
    
<div id="ids-usermgr-roleset" class="form-horizontal">
    <input type="hidden" name="roleid" value="{[=it.meta.id]}">
    
    <label class="ids-form-group-title">Role Information</label>

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
          <span class="ids-form-checkbox">
            <input type="radio" name="status" value="{[=v.status]}" {[ if (v.status == it.status) { ]}checked="checked"{[ } ]}> {[=v.title]}
          </span>
        {[~]}
      </div>
    </div>

    <!-- <label class="ids-form-group-title">Privileges</label>

    {{range $insid, $inst := .instances}}
    <div class="form-group">
      <div class="col-sm-3 ids-um-role-inst">
        <div class="ids-umri-title">{{$inst.AppTitle}}</div>
        <table class="ids-umri-attr">
          <tr><td>Instance ID:</td><td class="ids-umri-attr-ctn">{{$inst.InstanceId}}</td></tr>
          <tr><td>Version:</td><td class="ids-umri-attr-ctn">{{$inst.Version}}</td></tr>
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
//         url     : "/ids/user-mgr/role-save",
//         data    : $("#vukihr").serialize(),
//         timeout : 3000,
//         success : function(rsp) {

//             var rsj = JSON.parse(rsp);

//             if (rsj.status == 200) {
                
//                 lessAlert("#ids-usermgr-roleset-alert", 'alert-success', "Successfully saved");
                
//                 window.setTimeout(function(){
//                     idsWorkLoader("user-mgr/role-list");
//                 }, 1500);

//             } else {
//                 lessAlert("#ids-usermgr-roleset-alert", 'alert-danger', rsj.message);
//                 $("button[type=submit]", this).removeAttr('disabled');
//             }
//         },
//         error: function(xhr, textStatus, error) {
//             lessAlert("#ids-usermgr-roleset-alert", 'alert-danger', '{{T . "Internal Server Error"}}');
//             $("button[type=submit]", this).removeAttr('disabled');
//         }
//     });
// });

</script>
