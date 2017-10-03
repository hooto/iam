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

iamAccMgr.Index = function() {
    iam.TplCmd("acc-mgr/index", {
        callback: function(err, data) {
            iam.OpToolsClean();
            $("#com-content").html(data);
            l4i.UrlEventClean("iam-module-navbar-menus");
            l4i.UrlEventRegister("acc-mgr/recharge-list", iamAccMgr.RechargeList, "iam-module-navbar-menus");
            l4i.UrlEventHandler("acc-mgr/recharge-list", true);
        },
    });
}

iamAccMgr.RechargeList = function() {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            if (!data) {
                return;
            }

            if (!data.items) {
                data.items = [];
            }

            $("#work-content").html(tpl);
            iam.OpToolsRefresh("#iam-accm-rechargelist-optools");

            //
            for (var i in data.items) {
                if (!data.items[i].comment) {
                    data.items[i].comment = "";
                }
                if (!data.items[i].cash_type) {
                    data.items[i].cash_type = 0;
                }
                if (!data.items[i].cash_amount) {
                    data.items[i].cash_amount = 0;
                }
                if (!data.items[i].exp_product_limits) {
                    data.items[i]._exp_product_limits = "";
                } else {
                    data.items[i]._exp_product_limits = data.items[i].exp_product_limits.join(", ");
                }
                if (!data.items[i].exp_product_max) {
                    data.items[i].exp_product_max = 0;
                }
            }

            data._ecoin_types = l4i.Clone(iamAccMgr.ecoin_types);

            l4iTemplate.Render({
                dstid: "iam-accm-rechargelist",
                tplid: "iam-accm-rechargelist-tpl",
                data: data,
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        iam.ApiCmd("account-mgr/recharge-list", {
            callback: ep.done('data'),
        });

        iam.TplCmd("acc-mgr/recharge-list", {
            callback: ep.done('tpl'),
        });
    });
}


iamAccMgr.RechargeNew = function(username) {

    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', function(tpl) {

            var data = {
                user: username,
            };

            if (!data.user) {
                data.user = "";
            }

            data._recharge_types = l4i.Clone(iamAccMgr.ecoin_types);

            l4iModal.Open({
                tplsrc: tpl,
                width: 800,
                height: 400,
                data: data,
                title: "Account Recharge",
                buttons: [{
                    title: "Cancel",
                    onclick: "l4iModal.Close()",
                }, {
                    title: "Commit",
                    onclick: "iamAccMgr.RechargeNewCommit()",
                    style: "btn btn-primary",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        iam.TplCmd("acc-mgr/recharge-new", {
            callback: ep.done('tpl'),
        });
    });
}

iamAccMgr.RechargeNewCommit = function() {
    var form = $("#iam-accmgr-recharge-form"),
        alert_id = "#iam-accmgr-recharge-alert",
        req = l4i.Clone(iamAccMgr.recharge_def);

    try {
        if (!form) {
            throw "Can Not Found FORM";
        }

        req.user = form.find("input[name=user]").val();
        req.comment = form.find("input[name=comment]").val();
        req.type = parseInt(form.find("select[name=type]").val());
        req.amount = parseInt(form.find("input[name=amount]").val());
        req.exp_product_max = parseInt(form.find("input[name=exp_product_max]").val());
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
            }, 1000);
        },
    });
}

iamAccMgr.RechargeSet = function(id) {

    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            if (!data || !data.kind) {
                return;
            }

            data._recharge_types = l4i.Clone(iamAccMgr.ecoin_types);
            if (!data.exp_product_limits) {
                data.exp_product_limits = "";
            }
            if (!data.exp_product_max) {
                data.exp_product_max = 0;
            }
            if (!data.comment) {
                data.comment = "";
            }

            l4iModal.Open({
                tplsrc: tpl,
                width: 700,
                height: 400,
                data: data,
                title: "Account Recharge Edit",
                buttons: [{
                    title: "Cancel",
                    onclick: "l4iModal.Close()",
                }, {
                    title: "Commit",
                    onclick: "iamAccMgr.RechargeSetCommit()",
                    style: "btn btn-primary",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        iam.ApiCmd("account-mgr/recharge-entry?id=" + id, {
            callback: ep.done('data'),
        });

        iam.TplCmd("acc-mgr/recharge-set", {
            callback: ep.done('tpl'),
        });
    });
}

iamAccMgr.RechargeSetCommit = function() {

    var form = $("#iam-accmgr-recharge-form"),
        alert_id = "#iam-accmgr-recharge-alert",
        req = l4i.Clone(iamAccMgr.recharge_def);

    try {
        if (!form) {
            throw "Can Not Found FORM";
        }

        req.id = form.find("input[name=id]").val();
        req.comment = form.find("input[name=comment]").val();
        req.type = parseInt(form.find("select[name=type]").val());
        req.exp_product_max = parseInt(form.find("select[name=exp_product_max]").val());
        // req.payout = parseFloat(form.find("input[name=payout]").val());
        var plimits = form.find("input[name=exp_product_limits]").val();
        if (plimits && plimits.length > 0) {
            req.exp_product_limits = plimits.split(",");
        }

    } catch (err) {
        return l4i.InnerAlert(alert_id, 'alert-danger', err);
    }

    iam.ApiCmd("account-mgr/recharge-set", {
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
                iamAccMgr.RechargeList();
            }, 1000);
        },
    });
}

