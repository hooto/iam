
{{if .list}}
<table class="table table-hover">
  <thead>
    <tr>
      <th>Instance ID</th>
      <th>Application</th>
      <th>Version</th>
      <th>Status</th>
      <th>Owner</th>
      <th>Updated</th>
      <th></th>
    </tr>
  </thead>
  <tbody>
    {{range .list}}
    <tr>
      <td>{{.id}}</td>
      <td>{{.app_title}}</td>
      <td>{{.version}}</td>
      <td>{{.status}}</td>
      <td>{{.uid}}</td>
      <td>{{date .updated}}</td>
      <td>
        <a class="sepv31" href="#{{.id}}">Edit</a>
      </td>
    </tr>
    {{end}}
  </tbody>
</table>

{{else}}
<div class="alert alert-info" style="margin:20px 0;">Data not found</div>
{{end}}



<script type="text/javascript">

$(".sepv31").click(function() {
    var uid = $(this).attr("href").substr(1);
    idsWorkLoader("user-mgr/edit?uid="+ uid);
});

function _appmgr_list_refresh(page)
{
    var uri = "query_text="+ $("#query_text").val();
    uri += "&page="+ page;

    $.ajax({
        type    : "POST",
        url     : "/ids/user-mgr/list",
        data    : uri,
        timeout : 3000,
        success : function(rsp) {
            $("#work-content").html(rsp);
        },
        error   : function(xhr, textStatus, error) {
            //lessAlert("#azt02e", 'alert-danger', textStatus+' '+xhr.responseText);
        }
    });
}

</script>
