

<style type="text/css">
._list_query_input {
    padding: 5px 5px 5px 30px;
    background: url(/ids/~/ids/img/search-16.png) no-repeat 8px; 
    width: 220px;
}
.pagination {
  margin: 10px 0;
}
</style>

<table width="100%">
  <tr>
    <td>
      <form id="gir0c7" action="#" class="form-inlines">
        <input id="query_text" type="text"
          class="form-control _list_query_input" 
          placeholder="Enter to search" 
          value="{{.query_text}}">
      </form>
    </td>
    <td align="right">
      <!-- <button type="button" 
        class="btn btn-primary btn-sm" 
        onclick="idsWorkLoader('user-mgr/auth-edit')">
        New Instance
      </button> -->
    </td>
  </tr>
</table>
{{if .list}}
<table class="table table-hover">
  <thead>
    <tr>
      <th>Instance ID</th>
      <th>Application Name</th>
      <th>Application ID</th>
      <th>Version</th>
      <th>Status</th>
      <th>Owner</th>
      <th>Privileges</th>
      <th>Created</th>
      <th>Updated</th>
      <th></th>
    </tr>
  </thead>
  <tbody>
    {{range .list}}
    <tr>
      <td>{{.id}}</td>
      <td>{{.app_title}}</td>
      <td>{{.app_id}}</td>
      <td>{{.version}}</td>
      <td>{{.status_display}}</td>
      <td>{{.uid_name}}</td>
      <td>{{.privileges}}
      <td>{{date .created}}</td>
      <td>{{date .updated}}</td>
      <td>
        <a class="zaae6b" href="#{{.id}}">Edit</a>
      </td>
    </tr>
    {{end}}
  </tbody>
</table>

<div id="mxifon">
<ul class="pagination pagination-sm">
  {{if .pager.FirstPageNumber}}
  <li><a href="#{{.pager.FirstPageNumber}}">First</a></li>
  {{end}}
  {{range $index, $page := .pager.RangePages}}
  <li {{if eq $page $.pager.CurrentPageNumber}}class="active"{{end}}><a href="#{{$page}}">{{$page}}</a></li>
  {{end}}
  {{if .pager.LastPageNumber}}
  <li><a href="#{{.pager.LastPageNumber}}">Last</a></li>
  {{end}}
</ul>
</div>

{{else}}
<div class="alert alert-info" style="margin:20px 0;">Data not found</div>
{{end}}



<script type="text/javascript">

$(".zaae6b").click(function() {
    var id = $(this).attr("href").substr(1);
    idsWorkLoader("user-mgr/auth-edit?instid="+ id);
});

function _usermgr_authlist_refresh(page)
{
    var uri = "query_text="+ $("#query_text").val();
    uri += "&page="+ page;

    $.ajax({
        type    : "POST",
        url     : "/ids/user-mgr/auth-list",
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

$("#gir0c7").submit(function(event) {
    event.preventDefault();
    _usermgr_authlist_refresh(0);
});

$("#mxifon a").click(function() {
    var page = $(this).attr("href").substr(1);
    _usermgr_authlist_refresh(page);
});

</script>
