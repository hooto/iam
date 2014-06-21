<style type="text/css">
#t06uvj li a {
  padding: 5px 10px;
  margin-bottom: 10px;
}
</style>

<ul id="t06uvj" class="nav nav-pills">
  <li class="active"><a href="#app-mgr/list">Instances</a></li>
</ul>

<div id="work-content" class="">loading</div>

<script type="text/javascript">
$("#t06uvj a").click(function(event) {

    $("#t06uvj li.active").removeClass('active');

    var uri = $(this).attr("href").substr(1);

    switch (uri) {
    case "app-mgr/list":
        $(this).parent().addClass("active");
        idsWorkLoader(uri);
        break;
    }
});

idsWorkLoader("app-mgr/list");
</script>
