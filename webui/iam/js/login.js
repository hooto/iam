var iamLogin = {
    version: null,
};

iamLogin.InnerAlert = function(alertId, typeUI, msg, fadeTime) {

    if (!fadeTime) {
        fadeTime = 200;
    }
    if (!typeUI) {
        return $(alertId).fadeOut(fadeTime);
    }
    var elem = $(alertId);
    if (elem) {
        elem.removeClass().addClass("alert " + typeUI).html(msg);
        elem.fadeOut(fadeTime, function() {
            elem.fadeIn(fadeTime);
        });
    }
}

iamLogin.LoginCommit = function() {

    var alertId = "#iam-login-form-alert";
    var formData = $("#iam-login-form").serialize();

    $(".iam-form-item").attr("disabled", "disabled");
    iamLogin.InnerAlert(alertId, 'alert-info', "pending ...", 100);

    setTimeout(function() {
        $.ajax({
            type: "POST",
            url: "/iam/v1/service/login-auth",
            data: formData,
            timeout: 10000,
            success: function(data) {

                if (data.error) {
                    $(".iam-form-item").removeAttr("disabled");
                    document.getElementById("iam-login-form-pwd").focus();
                    return iamLogin.InnerAlert(alertId, 'alert-danger', data.error.message);
                }

                if (data.kind != "ServiceLoginAuth") {
                    $(".iam-form-item").removeAttr("disabled");
                    return iamLogin.InnerAlert(alertId, 'alert-danger', "Unknown Error");
                }

                iamLogin.InnerAlert(alertId, 'alert-success',
                    l4i.T("Successfully %s", l4i.T("Sign in")) + ", " + l4i.T("Page redirecting ..."));
                    // $("#iam-login-input-frame").hide(100);

                window.setTimeout(function() {
                    window.location = data.redirect_uri;;
                }, 1500);
            },
            error: function(xhr, textStatus, error) {
                $(".iam-form-item").removeAttr("disabled");
                iamLogin.InnerAlert(alertId, 'alert-danger', l4i.T("Internal Server Error"));
            }
        });
    }, 300);
}


//
iamLogin.RetrieveCommit = function() {
    var alertId = "#iam-resetpass-form-alert";

    iamLogin.InnerAlert(alertId, 'alert-info', "Pending");

    $.ajax({
        type: "POST",
        url: "/iam/reg/retrieve-put",
        data: $("#iam-resetpass-form").serialize(),
        timeout: 10000,
        success: function(data) {

            if (data.error) {
                return iamLogin.InnerAlert(alertId, 'alert-danger', data.error.message);
            }

            if (data.kind != "UserAuth") {
                return iamLogin.InnerAlert(alertId, 'alert-danger', "Unknown Error");
            }

            iamLogin.InnerAlert(alertId, 'alert-success',
                l4i.T("The reset URL has been sent to your mailbox, please check your email"));
            $(".iam-group").hide(200);
        },
        error: function(xhr, textStatus, error) {
            iamLogin.InnerAlert(alertId, 'alert-danger', l4i.T("Internal Server Error"));
        }
    });
};

iamLogin.PassResetCommit = function() {

    var alertId = "#iam-reg-passreset-form-alert",
        form = $("#iam-reg-passreset-form");

    iamLogin.InnerAlert(alertId, 'alert-info', l4i.T("Pending"));

    $.ajax({
        type: "POST",
        url: "/iam/reg/pass-reset-put",
        data: form.serialize(),
        timeout: 3000,
        success: function(data) {

            if (data.error) {
                return iamLogin.InnerAlert(alertId, 'alert-danger', data.error.message);
            }

            if (data.kind != "UserAuth") {
                return iamLogin.InnerAlert(alertId, 'alert-danger', "Unknown Error");
            }

            iamLogin.InnerAlert(alertId, 'alert-success', l4i.T("Successfully %s", l4i.T("Updated")) + ", " + l4i.T("Page redirecting ..."));
            $(".iam-group").hide(200);

            window.setTimeout(function() {
                window.location = "/iam/service/login?redirect_token=" + form.find("input[name=redirect_token]").val();
            }, 2000);
        },
        error: function(xhr, textStatus, error) {
            iamLogin.InnerAlert(alertId, 'alert-danger', l4i.T("Internal Server Error"));
        }
    });
}

iamLogin.SignupCommit = function() {
    var alertid = "#iam-signup-form-alert",
        form = $("#iam-signup-form");

    $.ajax({
        type: "POST",
        url: "/iam/reg/sign-up-reg",
        data: form.serialize(),
        timeout: 3000,
        success: function(data) {

            if (data.error) {
                return iamLogin.InnerAlert(alertid, 'alert-danger', data.error.message);
            }

            if (data.kind != "User") {
                return iamLogin.InnerAlert(alertid, 'alert-danger', l4i.T("Unknown Error"));
            }

            iamLogin.InnerAlert(alertid, 'alert-success',
                l4i.T("Successfully %s", l4i.T("registration")) + ", " + l4i.T("Page redirecting ..."));
            $(".iam-group").hide(600);

            window.setTimeout(function() {
                window.location = "/iam/service/login?redirect_uri=" + form.find("input[name=redirect_token]").val();
            }, 1500);
        },
        error: function(xhr, textStatus, error) {
            iamLogin.InnerAlert(alertid, 'alert-danger', l4i.T("Internal Server Error"));
        }
    });
}


