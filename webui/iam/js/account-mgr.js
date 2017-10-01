var iamAccMgr = {
    ecoin_types: [{
        name: "Cash",
        value: 1,
    }, {
        name: "Virtual",
        value: 32,
        default: true,
    }],
    recharge_def: {
        user: "",
        amount: 0,
        comment: "",
    },
}

iamAccMgr.Index = function() {}

iamAccMgr.Recharge = function(username) {

    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            if (!data || !data.kind) {
                return;
            }
            if (!data.ecoin_balance) {
                data.ecoin_balance = 0;
            }

            data._recharge_types = l4i.Clone(iamAccMgr.ecoin_types);

            l4iModal.Open({
                tplsrc: tpl,
                width: 600,
                height: 350,
                data: data,
                title: "Ecoin Recharge",
                buttons: [{
                    title: "Cancel",
                    onclick: "l4iModal.Close()",
                }, {
                    title: "Commit",
                    onclick: "iamAccMgr.RechargeCommit()",
                    style: "btn btn-primary",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        iam.ApiCmd("user-mgr/user-entry?username=" + username, {
            callback: ep.done('data'),
        });
        iam.TplCmd("account-mgr/recharge", {
            callback: ep.done('tpl'),
        });
    });
}

iamAccMgr.RechargeCommit = function() {
    var form = $("#iam-accmgr-recharge-form"),
        alert_id = "#iam-accmgr-recharge-alert",
        req = l4i.Clone(iamAccMgr.recharge_def);

    try {
        if (!form) {
            throw "Can Not Found FORM";
        }

        req.user = form.find("input[name=username]").val();
        req.comment = form.find("input[name=comment]").val();
        req.type = parseInt(form.find("select[name=recharge_type]").val());
        req.amount = parseInt(form.find("input[name=ecoin_amount]").val());
        var plimits = form.find("input[name=exp_product_limits]").val();
        if (plimits && plimits.length > 0) {
            req.exp_product_limits = plimits.split(",");
        }

    } catch (err) {
        return l4i.InnerAlert(alert_id, 'alert-danger', err);
    }

    iam.ApiCmd("account-mgr/recharge-new", {
        method: "PUT",
        data: JSON.stringify(req),
        callback: function(err, data) {

            if (err) {
                return l4i.InnerAlert(alert_id, 'alert-danger', err);
            }

            if (!data || data.error) {
                if (data.error) {
                    return l4i.InnerAlert(alert_id, 'alert-danger', data.error.message);
                }
                return l4i.InnerAlert(alert_id, 'alert-danger', "network error");
            }

            l4i.InnerAlert(alert_id, 'alert-success', "Successfully updated");

            window.setTimeout(function() {
                l4iModal.Close();
                iamUserMgr.UserList();
            }, 1000);
        },
    });
}

