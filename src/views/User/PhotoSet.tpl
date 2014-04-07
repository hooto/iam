
<style>
#attachment {
    border: 1px solid #ccc;
}
</style>

<div id="ps40yi" class="alert alert-info">
    <p>{{T . "You must upload a JPG, GIF, or PNG file"}}</p>
</div>

<table> 
  <tr> 
    <td valign="bottom">
        <img src="/ids/service/photo/{{.login_uid}}" width="96" height="96" > 
        <b>{{T . "Normal size"}}</b>
    </td> 
    <td width="30px"></td> 
    <td valign="bottom">
        <img src="/ids/service/photo/{{.login_uid}}" width="48" height="48" /> <b>{{T . "Small size"}}</b>
    </td> 
  </tr> 
</table>
<br />
<form class="oqmbg4" enctype="multipart/form-data" action="#" method="post">
  <table width="100%">
    
    <tr>
      <td width="160px"><b>{{T . "Select a New Picture"}}</b></td>
      <td><input id="attachment" name="attachment" size="20" type="file" class="btn" /></td>
    </tr>

  </table>
</form>


<script type="text/javascript">

lessModalButtonAdd("ftd5ci", "{{T . "Cancel"}}", "lessModalClose()", "");
lessModalButtonAdd("v7lq98", "{{T . "Upload"}}", "_user_photo_set_fileupl()", "btn-inverse pull-right");


$(".oqmbg4").submit(function(event) {

    event.preventDefault(); 

    _user_photo_set_fileupl();
});

function _user_photo_set_fileupl()
{
    var files = document.getElementById('attachment').files;
    
    if (!files.length) {
        lessAlert("#ps40yi", "alert-danger", '{{T . "Please select a file"}}!');
        return;
    }

    console.log("AA");

    for (var i = 0, file; file = files[i]; ++i) {
        
        if (file.size > 2 * 1024 * 1024) {
            lessAlert("ps40yi", 'alert-danger', '{{T . "The file is too large to upload"}}');
            return;
        }
                
        var reader = new FileReader();
        reader.onload = (function(file) {  
            return function(e) {
                if (e.target.readyState != FileReader.DONE) {
                    return;
                }

                var req = {
                    data: {
                        size : file.size,
                        name : file.name,
                        data : e.target.result,
                    }
                }

                console.log(JSON.stringify(req));

                $.ajax({
                    type    : "POST",
                    url     : "/ids/user/photo-put",
                    data    : JSON.stringify(req),
                    timeout : 3000,
                    contentType: "application/json; charset=utf-8",
                    success : function(rsp) {

                        var obj = JSON.parse(rsp);

                        if (obj.status == 200) {
                            lessAlert("#ps40yi", 'alert-success', obj.message);

                            window.setTimeout(function(){
                                lessModalClose();
                                window.location = "/ids/user/index";
                            }, 2000);

                        } else {
                            lessAlert("#ps40yi", 'alert-danger', obj.message);
                        }

                    },
                    error   : function(xhr, textStatus, error) {
                        lessAlert("#ps40yi", 'alert-danger', textStatus+' '+xhr.responseText);
                    }
                });    

            };  
        })(file); 
        
        reader.readAsDataURL(file);
    }
}

</script>
