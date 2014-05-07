<style type="text/css">
#n66g5e li a {
  padding: 5px 10px;
  margin-bottom: 10px;
}
</style>

<ul id="n66g5e" class="nav nav-pills">
  <li class="active"><a href="#sys-mgr/gen-set">General Settings</a></li>
  <li><a href="#sys-mgr/email-set">Email Server Settings</a></li>
</ul>

<div id="work-content" class="ids-user-panel">loading</div>

<script type="text/javascript">
$("#n66g5e a").click(function(event) {

    $("#n66g5e li.active").removeClass('active');

    var uri = $(this).attr("href").substr(1);

    switch (uri) {
    case "sys-mgr/gen-set":
    case "sys-mgr/email-set":
        $(this).parent().addClass("active");
        idsWorkLoader(uri);
        break;
    }
});

idsWorkLoader("sys-mgr/email-set");
</script>