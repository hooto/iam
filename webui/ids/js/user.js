var idsuser = {

}

idsuser.Init = function()
{
    l4i.UrlEventRegister("user/overview", idsuser.Overview);
}

idsuser.Overview = function()
{
    // console.log("overview");
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function (tpl, data) {

            l4iTemplate.Render({
                dstid  : "com-content",
                tplsrc : tpl,
                data   : data,
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });
    
        ids.ApiCmd("user/profile", {
            callback: ep.done('data'),
        });

        ids.TplCmd("user/overview", {
            callback: ep.done('tpl'),           
        });
    });
}

idsuser.PassSetForm = function()
{
    ids.TplCmd("user/pass-set", {

        callback : function(err, tpl) {
            
            l4iModal.Open({
                tplsrc  : tpl,
                width   : 500,
                height  : 350,
                title   : "Change Password",
                buttons : [{
                    title : "Cancel",
                    onclick : "l4iModal.Close()",
                }, {
                    title : "Save",
                    onclick : "idsuser.PassSetCommit()",
                    style   : "btn btn-primary",
                }],
            });
        },
    });    
}

idsuser.PassSetCommit = function()
{
    var form = $("#ids-user-pass-set");
    
    var req = {
        currentPassword  : form.find("input[name=passwd_current]").val(),
        newPassword : form.find("input[name=passwd_new]").val(),
    };

    if (req.newPassword != form.find("input[name=passwd_confirm]").val()) {
        return l4i.InnerAlert("#ids-user-pass-set-alert", 'alert-danger', "Passwords do not match");
    }

    ids.ApiCmd("user/pass-set", {
        method : "PUT",
        data   : JSON.stringify(req),
        callback : function(err, data) {
            
            if (err) {
                return l4i.InnerAlert("#ids-user-pass-set-alert", 'alert-danger', err);
            }

            if (!data || data.error) {
                return l4i.InnerAlert("#ids-user-pass-set-alert", 'alert-danger', data.error.message);
            }

            l4i.InnerAlert("#ids-user-pass-set-alert", 'alert-success', "Successfully updated");

            window.setTimeout(function(){
                l4iModal.Close();
            }, 1000);
        },
    });
}


idsuser.EmailSetForm = function()
{
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function (tpl, data) {

            l4iModal.Open({
                tplsrc  : tpl,
                data    : {
                    email : data.login.email,
                },
                width   : 500,
                height  : 350,
                title   : "Change Email",
                buttons : [{
                    title : "Cancel",
                    onclick : "l4iModal.Close()",
                }, {
                    title : "Save",
                    onclick : "idsuser.EmailSetCommit()",
                    style   : "btn btn-primary",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });
    
        ids.ApiCmd("user/profile", {
            callback: ep.done('data'),
        });

        ids.TplCmd("user/email-set", {
            callback: ep.done('tpl'),           
        });
    });
}

idsuser.EmailSetCommit = function()
{
    var form = $("#ids-user-email-set");
    
    var req = {
        email : form.find("input[name=email]").val(),
        auth  : form.find("input[name=auth]").val(),
    };

    ids.ApiCmd("user/email-set", {
        method : "PUT",
        data   : JSON.stringify(req),
        callback : function(err, data) {
            
            if (err) {
                return l4i.InnerAlert("#ids-user-email-set-alert", 'alert-danger', err);
            }

            if (!data || data.error) {
                return l4i.InnerAlert("#ids-user-email-set-alert", 'alert-danger', data.error.message);
            }

            l4i.InnerAlert("#ids-user-email-set-alert", 'alert-success', "Successfully updated");

            window.setTimeout(function(){
                l4iModal.Close();
                idsuser.Overview();
            }, 1000);
        },
    });
}


idsuser.ProfileSetForm = function()
{
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function (tpl, data) {

            l4iModal.Open({
                tplsrc  : tpl,
                data    : data,
                width   : 600,
                height  : 400,
                title   : "Change Profile",
                buttons : [{
                    title : "Cancel",
                    onclick : "l4iModal.Close()",
                }, {
                    title : "Save",
                    onclick : "idsuser.ProfileSetCommit()",
                    style   : "btn btn-primary",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });
    
        ids.ApiCmd("user/profile", {
            callback: ep.done('data'),
        });

        ids.TplCmd("user/profile-set", {
            callback: ep.done('tpl'),           
        });
    });
}

idsuser.ProfileSetCommit = function()
{
    var form = $("#ids-user-profile-set");
    
    var req = {
        name : form.find("input[name=name]").val(),
        birthday : form.find("input[name=birthday]").val(),
        about  : form.find("textarea[name=about]").val(),
    };

    ids.ApiCmd("user/profile-set", {
        method : "PUT",
        data   : JSON.stringify(req),
        callback : function(err, data) {
            
            if (err) {
                return l4i.InnerAlert("#ids-user-profile-set-alert", 'alert-danger', err);
            }

            if (!data || data.error) {
                return l4i.InnerAlert("#ids-user-profile-set-alert", 'alert-danger', data.error.message);
            }

            l4i.InnerAlert("#ids-user-profile-set-alert", 'alert-success', "Successfully updated");

            window.setTimeout(function(){
                l4iModal.Close();
                idsuser.Overview();
            }, 1000);
        },
    });
}

idsuser.PhotoSetForm = function(uuid)
{
    ids.TplCmd("user/photo-set", {

        callback : function(err, tpl) {
            
            l4iModal.Open({
                tplsrc  : tpl,
                width   : 600,
                height  : 400,
                data    : {login_id: uuid},
                title   : "Change Photo",
                buttons : [{
                    title : "Cancel",
                    onclick : "l4iModal.Close()",
                }, {
                    title : "Save",
                    onclick : "idsuser.PhotoSetCommit()",
                    style   : "btn btn-primary",
                }],
            });
        },
    });
}

idsuser.PhotoSetCommit = function()
{
    var files = document.getElementById('ids-user-photo-set-file').files;
    
    if (!files.length) {
        return l4i.InnerAlert("#ids-user-photo-set-alert", "alert-danger", 'Please select a file!');
    }

    for (var i = 0, file; file = files[i]; ++i) {
        
        if (file.size > 2 * 1024 * 1024) {
            return l4i.InnerAlert("ids-user-photo-set-alert", 'alert-danger', 'The file is too large to upload');
        }
                
        var reader = new FileReader();
        reader.onload = (function(file) {  
            return function(e) {
                
                if (e.target.readyState != FileReader.DONE) {
                    return;
                }

                var req = {
                    size : file.size,
                    name : file.name,
                    data : e.target.result,
                }

                ids.ApiCmd("user/photo-set", {
                    method : "PUT",
                    data   : JSON.stringify(req),
                    callback : function(err, data) {
            
                        if (err) {
                            return l4i.InnerAlert("#ids-user-photo-set-alert", 'alert-danger', err);
                        }

                        if (!data || data.error) {
                            return l4i.InnerAlert("#ids-user-photo-set-alert", 'alert-danger', data.error.message);
                        }

                        l4i.InnerAlert("#ids-user-photo-set-alert", 'alert-success', "Successfully updated");

                        window.setTimeout(function(){
                            l4iModal.Close();
                            idsuser.Overview();
                        }, 1000);
                    },
                });
            };

        })(file);
        
        reader.readAsDataURL(file);
    }
}

idsuser.SignOut = function()
{
    // l4iCookie.Del("access_token");
    window.setTimeout(function(){    
        window.location = "/ids/service/sign-out";
    }, 500);
}