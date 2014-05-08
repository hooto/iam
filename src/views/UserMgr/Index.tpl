<style type="text/css">
#bgi4w4 li a {
  padding: 5px 10px;
  margin-bottom: 10px;
}
</style>

<ul id="bgi4w4" class="nav nav-pills">
  <li class="active"><a href="#user-mgr/list">Users</a></li>
  <li><a href="#user-mgr/role-list">Role Settings</a></li>
</ul>

<div id="work-content" class="ids-user-panel">loading</div>

<script type="text/javascript">
$("#bgi4w4 a").click(function(event) {

    $("#bgi4w4 li.active").removeClass('active');

    var uri = $(this).attr("href").substr(1);

    switch (uri) {
    case "user-mgr/list":
    //case "user-mgr/new":
    //case "user-mgr/edit":
    case "user-mgr/role-list":
        $(this).parent().addClass("active");
        idsWorkLoader(uri);
        break;
    }
});

idsWorkLoader("user-mgr/list");
</script>
