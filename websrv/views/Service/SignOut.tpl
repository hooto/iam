<!DOCTYPE html>
<html lang="en">

<head>
  <meta charset="utf-8">
  <title>Sign Out</title>
</head>

<body>

<h2>Sign Out ...</h2>

<script type="text/javascript">

function cookie_set(key, val, sec)
{
    var expires = "";
    
    if (sec) {
        var date = new Date();
        date.setTime(date.getTime() + (sec * 1000));
        expires = "; expires=" + date.toGMTString();
    }
    
    document.cookie = key + "=" + val + expires + "; path=/";
}

cookie_set("{{.access_token_key}}", "", -1);

window.setTimeout(function() {
    window.location = "/iam";
}, 2000);

</script>

</body>

</html>
