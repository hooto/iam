
<style>

.ids-um-formtable {
  width: 100%;
  border: 0;
  margin: 10px 0;
}

.ids-um-formtable td {
  padding: 5px 0;
}

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
    <input type="hidden" name="rid" value="{{.rid}}">
    
    <label class="ids-form-group-title">Role Information</label>

    <div class="form-group">
      <label class="col-sm-2 control-label">Name</label>
      <div class="col-sm-10">
        <input type="text" class="form-control" name="name" value="{{.name}}">
      </div>
    </div>

    <div class="form-group">
      <label class="col-sm-2 control-label">Description</label>
      <div class="col-sm-10">
        <input type="text" class="form-control" name="desc" value="{{.desc}}">
      </div>
    </div>

    <div class="form-group">
      <div class="col-sm-offset-2 col-sm-10">
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
        url     : "/ids/user-mgr/role-save",
        data    : $("#vukihr").serialize(),
        timeout : 3000,
        success : function(rsp) {

            var rsj = JSON.parse(rsp);

            if (rsj.status == 200) {
                
                lessAlert("#q6p2le", 'alert-success', "Successfully saved");
                
                window.setTimeout(function(){
                    idsWorkLoader("user-mgr/role-list");
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
