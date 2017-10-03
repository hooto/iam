var iamUser = {

}

iamUser.Overview = function() {
    iam.OpToolsClean();
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', 'acc', function(tpl, data, acc) {

            if (!data.login.display_name) {
                data.login.display_name = data.login.name;
            }

            if (!acc.balance) {
                acc.balance = 0;
            }
            if (!acc.prepay) {
                acc.prepay = 0;
            }

            data.account = acc;

            l4iTemplate.Render({
                dstid: "com-content",
                tplsrc: tpl,
                data: data,
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        iam.ApiCmd("account/user", {
            callback: ep.done('acc'),
        });

        iam.ApiCmd("user/profile", {
            callback: ep.done('data'),
        });

        iam.TplCmd("user/overview", {
            callback: ep.done('tpl'),
        });
    });
}

iamUser.PassSetForm = function() {
    iam.TplCmd("user/pass-set", {

        callback: function(err, tpl) {

            l4iModal.Open({
                tplsrc: tpl,
                width: 500,
                height: 350,
                title: "Change Password",
                buttons: [{
                    title: "Cancel",
                    onclick: "l4iModal.Close()",
                }, {
                    title: "Save",
                    onclick: "iamUser.PassSetCommit()",
                    style: "btn btn-primary",
                }],
            });
        },
    });
}

iamUser.PassSetCommit = function() {
    var form = $("#iam-user-pass-set"),
        alert_id = "#iam-user-pass-set-alert";

    var req = {
        current_password: form.find("input[name=passwd_current]").val(),
        new_password: form.find("input[name=passwd_new]").val(),
    };

    if (req.new_password != form.find("input[name=passwd_confirm]").val()) {
        return l4i.InnerAlert(alert_id, 'alert-danger', "Passwords do not match");
    }

    iam.ApiCmd("user/pass-set", {
        method: "PUT",
        data: JSON.stringify(req),
        callback: function(err, data) {

            if (err) {
                return l4i.InnerAlert(alert_id, 'alert-danger', err);
            }

            if (!data || data.error) {
                return l4i.InnerAlert(alert_id, 'alert-danger', data.error.message);
            }

            l4i.InnerAlert(alert_id, 'alert-success', "Successfully updated");

            window.setTimeout(function() {
                l4iModal.Close();
            }, 1000);
        },
    });
}


iamUser.EmailSetForm = function() {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            if (!data.login.email) {
                data.login.email = "";
            }

            l4iModal.Open({
                tplsrc: tpl,
                data: {
                    email: data.login.email,
                },
                width: 500,
                height: 350,
                title: "Change Email",
                buttons: [{
                    title: "Cancel",
                    onclick: "l4iModal.Close()",
                }, {
                    title: "Save",
                    onclick: "iamUser.EmailSetCommit()",
                    style: "btn btn-primary",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        iam.ApiCmd("user/profile", {
            callback: ep.done('data'),
        });

        iam.TplCmd("user/email-set", {
            callback: ep.done('tpl'),
        });
    });
}

iamUser.EmailSetCommit = function() {
    var form = $("#iam-user-email-set"),
        alert_id = "#iam-user-email-set-alert";

    var req = {
        email: form.find("input[name=email]").val(),
        auth: form.find("input[name=auth]").val(),
    };

    iam.ApiCmd("user/email-set", {
        method: "PUT",
        data: JSON.stringify(req),
        callback: function(err, data) {

            if (err) {
                return l4i.InnerAlert(alert_id, 'alert-danger', err);
            }

            if (!data || data.error) {
                return l4i.InnerAlert(alert_id, 'alert-danger', data.error.message);
            }

            l4i.InnerAlert(alert_id, 'alert-success', "Successfully updated");

            window.setTimeout(function() {
                l4iModal.Close();
                iamUser.Overview();
            }, 1000);
        },
    });
}


iamUser.ProfileSetForm = function() {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            if (!data.login.display_name) {
                data.login.display_name = "No Name";
            }

            if (!data.birthday) {
                data.birthday = "";
            }

            if (!data.about) {
                data.about = "";
            }

            l4iModal.Open({
                tplsrc: tpl,
                data: data,
                width: 600,
                height: 400,
                title: "Change Profile",
                buttons: [{
                    title: "Cancel",
                    onclick: "l4iModal.Close()",
                }, {
                    title: "Save",
                    onclick: "iamUser.ProfileSetCommit()",
                    style: "btn btn-primary",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        iam.ApiCmd("user/profile", {
            callback: ep.done('data'),
        });

        iam.TplCmd("user/profile-set", {
            callback: ep.done('tpl'),
        });
    });
}

iamUser.ProfileSetCommit = function() {
    var form = $("#iam-user-profile-set"),
        alert_id = "#iam-user-profile-set-alert";

    var req = {
        login: {
            display_name: form.find("input[name=display_name]").val(),
        },
        birthday: form.find("input[name=birthday]").val(),
        about: form.find("textarea[name=about]").val(),
    };

    iam.ApiCmd("user/profile-set", {
        method: "PUT",
        data: JSON.stringify(req),
        callback: function(err, data) {

            if (err) {
                return l4i.InnerAlert(alert_id, 'alert-danger', err);
            }

            if (!data || data.error) {
                return l4i.InnerAlert(alert_id, 'alert-danger', data.error.message);
            }

            l4i.InnerAlert(alert_id, 'alert-success', "Successfully updated");

            window.setTimeout(function() {
                l4iModal.Close();
                iamUser.Overview();
            }, 1000);
        },
    });
}

iamUser.PhotoSetForm = function(username) {
    iam.TplCmd("user/photo-set", {

        callback: function(err, tpl) {

            l4iModal.Open({
                tplsrc: tpl,
                width: 600,
                height: 400,
                data: {
                    username: username
                },
                title: "Change Photo",
                buttons: [{
                    title: "Cancel",
                    onclick: "l4iModal.Close()",
                }, {
                    title: "Save",
                    onclick: "iamUser.PhotoSetCommit()",
                    style: "btn btn-primary",
                }],
            });
        },
    });
}

iamUser.PhotoSetCommit = function() {
    var files = document.getElementById('iam-user-photo-set-file').files,
        alert_id = "#iam-user-photo-set-alert";

    if (!files.length) {
        return l4i.InnerAlert(alert_id, "alert-danger", 'Please select a file!');
    }

    for (var i = 0, file; file = files[i]; ++i) {

        if (file.size > 2 * 1024 * 1024) {
            return l4i.InnerAlert("iam-user-photo-set-alert", 'alert-danger', 'The file is too large to upload');
        }

        var reader = new FileReader();
        reader.onload = (function(file) {
            return function(e) {

                if (e.target.readyState != FileReader.DONE) {
                    return;
                }

                var req = {
                    size: file.size,
                    name: file.name,
                    data: e.target.result,
                }

                iam.ApiCmd("user/photo-set", {
                    method: "PUT",
                    data: JSON.stringify(req),
                    callback: function(err, data) {

                        if (err) {
                            return l4i.InnerAlert(alert_id, 'alert-danger', err);
                        }

                        if (!data || data.error) {
                            return l4i.InnerAlert(alert_id, 'alert-danger', data.error.message);
                        }

                        l4i.InnerAlert(alert_id, 'alert-success', "Successfully updated");

                        window.setTimeout(function() {
                            l4iModal.Close();
                            iamUser.Overview();
                        }, 1000);
                    },
                });
            };

        })(file);

        reader.readAsDataURL(file);
    }
}

iamUser.SignOut = function() {
    // l4iCookie.Del("access_token");
    window.setTimeout(function() {
        window.location = "/iam/service/sign-out";
    }, 500);
}
