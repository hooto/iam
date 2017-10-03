var iamAccMgr = {
    fund_types: [{
        name: "Cash",
        value: 1,
    }, {
        name: "Virtual",
        value: 32,
        default: true,
    }],
    fund_def: {
        user: "",
        amount: 0,
        comment: "",
        exp_product_max: 0,
    },
    fund_set_user: null,
}

iamAccMgr.Index = function() {
    iam.TplCmd("acc-mgr/index", {
        callback: function(err, data) {
            iam.OpToolsClean();
            $("#com-content").html(data);
            l4i.UrlEventClean("iam-module-navbar-menus");
            l4i.UrlEventRegister("acc-mgr/fund-list", iamAccMgr.FundList, "iam-module-navbar-menus");
            l4i.UrlEventRegister("acc-mgr/charge-list", iamAccMgr.ChargeList, "iam-module-navbar-menus");
            l4i.UrlEventHandler("acc-mgr/fund-list", true);
        },
    });
}

iamAccMgr.FundList = function() {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            if (!data) {
                return;
            }

            if (!data.items) {
                data.items = [];
            }

            $("#work-content").html(tpl);
            iam.OpToolsRefresh("#iam-accm-fundlist-optools");

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
                if (!data.items[i].operator) {
                    data.items[i].operator = "";
                }
            }

            data._fund_types = l4i.Clone(iamAccMgr.fund_types);

            l4iTemplate.Render({
                dstid: "iam-accm-fundlist",
                tplid: "iam-accm-fundlist-tpl",
                data: data,
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        iam.ApiCmd("account-mgr/fund-list", {
            callback: ep.done('data'),
        });

        iam.TplCmd("acc-mgr/fund-list", {
            callback: ep.done('tpl'),
        });
    });
}


iamAccMgr.FundNew = function(username) {

    if (username) {
        iamAccMgr.fund_set_user = username;
    } else {
        iamAccMgr.fund_set_user = null;
    }

    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', function(tpl) {

            var data = {
                user: username,
            };

            if (!data.user) {
                data.user = "";
            }

            data._fund_types = l4i.Clone(iamAccMgr.fund_types);

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
                    onclick: "iamAccMgr.FundNewCommit()",
                    style: "btn btn-primary",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        iam.TplCmd("acc-mgr/fund-new", {
            callback: ep.done('tpl'),
        });
    });
}

iamAccMgr.FundNewCommit = function() {
    var form = $("#iam-accmgr-fund-form"),
        alert_id = "#iam-accmgr-fund-alert",
        req = l4i.Clone(iamAccMgr.fund_def);

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

    iam.ApiCmd("account-mgr/fund-new", {
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
                if (!iamAccMgr.fund_set_user) {
                    iamAccMgr.FundList();
                }
            }, 1000);
        },
    });
}

iamAccMgr.FundSet = function(id) {

    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            if (!data || !data.kind) {
                return;
            }

            data._fund_types = l4i.Clone(iamAccMgr.fund_types);
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
                title: "Account Fund Edit",
                buttons: [{
                    title: "Cancel",
                    onclick: "l4iModal.Close()",
                }, {
                    title: "Commit",
                    onclick: "iamAccMgr.FundSetCommit()",
                    style: "btn btn-primary",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        iam.ApiCmd("account-mgr/fund-entry?id=" + id, {
            callback: ep.done('data'),
        });

        iam.TplCmd("acc-mgr/fund-set", {
            callback: ep.done('tpl'),
        });
    });
}

iamAccMgr.FundSetCommit = function() {

    var form = $("#iam-accmgr-fund-form"),
        alert_id = "#iam-accmgr-fund-alert",
        req = l4i.Clone(iamAccMgr.fund_def);

    try {
        if (!form) {
            throw "Can Not Found FORM";
        }

        req.id = form.find("input[name=id]").val();
        req.comment = form.find("input[name=comment]").val();
        req.type = parseInt(form.find("select[name=type]").val());
        req.exp_product_max = parseInt(form.find("input[name=exp_product_max]").val());
        // req.payout = parseFloat(form.find("input[name=payout]").val());
        var plimits = form.find("input[name=exp_product_limits]").val();
        if (plimits && plimits.length > 0) {
            req.exp_product_limits = plimits.split(",");
        }

    } catch (err) {
        return l4i.InnerAlert(alert_id, 'alert-danger', err);
    }

    iam.ApiCmd("account-mgr/fund-set", {
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
                iamAccMgr.FundList();
            }, 1000);
        },
    });
}

iamAccMgr.ChargeList = function() {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            if (!data) {
                return;
            }

            if (!data.items) {
                data.items = [];
            }

            for (var i in data.items) {
                if (!data.items[i].payout) {
                    data.items[i].payout = 0;
                }
                data.items[i].payout = data.items[i].payout.toFixed(4);
                if (!data.items[i].prepay) {
                    data.items[i].prepay = 0;
                }
                data.items[i].prepay = data.items[i].prepay.toFixed(4);
            }

            $("#work-content").html(tpl);

            data._fund_types = l4i.Clone(iamAcc.fund_types);

            l4iTemplate.Render({
                dstid: "iam-accmgr-chargelist",
                tplid: "iam-accmgr-chargelist-tpl",
                data: data,
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        iam.ApiCmd("account-mgr/charge-list", {
            callback: ep.done('data'),
        });

        iam.TplCmd("acc-mgr/charge-list", {
            callback: ep.done('tpl'),
        });
    });
}

iamAccMgr.ChargeSetPayout = function(user, id) {

    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            if (!data || !data.kind) {
                return;
            }

            l4iModal.Open({
                tplsrc: tpl,
                width: 700,
                height: 400,
                data: data,
                title: "Account Charge Close",
                buttons: [{
                    title: "Cancel",
                    onclick: "l4iModal.Close()",
                }, {
                    title: "Commit",
                    onclick: "iamAccMgr.ChargeSetPayoutCommit()",
                    style: "btn btn-primary",
                }],
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        iam.ApiCmd("account-mgr/charge-entry?id=" + id + "&user=" + user, {
            callback: ep.done('data'),
        });

        iam.TplCmd("acc-mgr/charge-set-payout", {
            callback: ep.done('tpl'),
        });
    });
}

iamAccMgr.ChargeSetPayoutCommit = function() {

    var form = $("#iam-accmgr-chargeset-payout-form"),
        alert_id = "#iam-accmgr-chargeset-payout-alert",
        req = {};

    try {
        if (!form) {
            throw "Can Not Found FORM";
        }

        req.id = form.find("input[name=id]").val();
        req.user = form.find("input[name=user]").val();
        req.payout = parseFloat(form.find("input[name=payout]").val());

    } catch (err) {
        return l4i.InnerAlert(alert_id, 'alert-danger', err);
    }

    iam.ApiCmd("account-mgr/charge-set-payout", {
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
                iamAccMgr.ChargeList();
            }, 1000);
        },
    });
}


