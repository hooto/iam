<nav class="navbar navbar-default">

  <div class="container-fluid">
    
    <div class="navbar-header">
      <a class="navbar-brand" href="#">
        <i class="iam-nav-logo"></i>
        <span>lessFly Manager</span>
      </a>
    </div>

    <div id="etpntn" class="collapse navbar-collapse">
      <ul class="nav navbar-nav">
        <li class="active"><a href="#lpm/index">Package Management</a></li>
        <li><a href="#pkg-service/index">Package Service</a></li>
      </ul>

      <ul class="nav navbar-nav navbar-right">
        <li><a href="#">Sign Out</a></li>
      </ul>
    </div>
  </div>
</nav>

<script type="text/javascript">

$("#etpntn a").click(function(event) {

    $("#etpntn li.active").removeClass('active');

    var uri = $(this).attr("href").substr(1);

    switch (uri) {
    //case "index/index":
    case "lpm/index":
    case "pkg-service/index":
    case "container/index":
        $(this).parent().addClass("active");
        iam.ComLoader(uri);
        break;
    }
});

</script>
