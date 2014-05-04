<table width="100%">
<tr>
  <!-- <td width="200px" valign="top">
    <div>

      <ul id="bgi4w4" class="nav nav-pills nav-stacked">
        <li class="active"><a href="#user-mgr/list">Browse</a></li>
        <li><a href="#user-mgr/new">New User</a></li>
      </ul>

    </div>
  </td>
  <td width="10px"></td> -->
  <td valign="top">
    <div id="work-content"></div>
  </td>
</tr>
</table>

<script type="text/javascript">
$("#bgi4w4 a").click(function(event) {

    $("#bgi4w4 li.active").removeClass('active');

    var uri = $(this).attr("href").substr(1);

    switch (uri) {
    case "user-mgr/list":
    case "user-mgr/new":
    case "user-mgr/edit":
        $(this).parent().addClass("active");
        idsWorkLoader(uri);
        break;
    }
});

idsWorkLoader("user-mgr/list");
</script>