var iamAcc = {
    ecoin_types: [{
        name: "Cash",
        value: 1,
    }, {
        name: "Virtual",
        value: 32,
        default: true,
    }],
}

iamAcc.Index = function() {
    iam.OpToolActive = null;
    iam.TplCmd("account/index", {
        callback: function(err, data) {
            $("#com-content").html(data);
            l4i.UrlEventClean("iam-module-navbar-menus");
            l4i.UrlEventRegister("account/recharge-list", iamAcc.RechargeList, "iam-module-navbar-menus");
            l4i.UrlEventRegister("account/active-list", iamAcc.ActiveList, "iam-module-navbar-menus");
            l4i.UrlEventRegister("account/charge-list", iamAcc.ChargeList, "iam-module-navbar-menus");
            l4i.UrlEventRegister("account/payout-list", iamAcc.PayoutList, "iam-module-navbar-menus");
            l4i.UrlEventHandler("account/charge-list", true);
        },
    });
}

iamAcc.RechargeList = function() {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            if (!data) {
                return;
            }

            if (!data.items) {
                data.items = [];
            }

            $("#work-content").html(tpl);
            // iam.OpToolsRefresh("#iam-acc-activelist-optools");
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
            }

            data._ecoin_types = l4i.Clone(iamAcc.ecoin_types);

            l4iTemplate.Render({
                dstid: "iam-acc-rechargelist",
                tplid: "iam-acc-rechargelist-tpl",
                data: data,
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        iam.ApiCmd("account/recharge-list", {
            callback: ep.done('data'),
        });

        iam.TplCmd("account/recharge-list", {
            callback: ep.done('tpl'),
        });
    });
}

iamAcc.ActiveList = function() {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            if (!data) {
                return;
            }

            if (!data.items) {
                data.items = [];
            }

            for (var i in data.items) {

                if (!data.items[i].amount) {
                    data.items[i].amount = 0;
                }
                data.items[i].amount = data.items[i].amount.toFixed(4);
                if (!data.items[i].payout) {
                    data.items[i].payout = 0;
                }
                data.items[i].payout = data.items[i].payout.toFixed(4);
                if (!data.items[i].prepay) {
                    data.items[i].prepay = 0;
                }
                data.items[i].prepay = data.items[i].prepay.toFixed(4);

                if (!data.items[i].exp_product_limits) {
                    data.items[i]._exp_product_limits = "";
                } else {
                    data.items[i]._exp_product_limits = data.items[i].exp_product_limits.join(", ");
                }
            }

            $("#work-content").html(tpl);
            // iam.OpToolsRefresh("#iam-acc-activelist-optools");

            data._ecoin_types = l4i.Clone(iamAcc.ecoin_types);

            l4iTemplate.Render({
                dstid: "iam-acc-activelist",
                tplid: "iam-acc-activelist-tpl",
                data: data,
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        iam.ApiCmd("account/active-list", {
            callback: ep.done('data'),
        });

        iam.TplCmd("account/active-list", {
            callback: ep.done('tpl'),
        });
    });
}

iamAcc.ChargeList = function() {
    seajs.use(["ep"], function(EventProxy) {

        var ep = EventProxy.create('tpl', 'data', function(tpl, data) {

            if (!data) {
                return;
            }

            if (!data.items) {
                data.items = [];
            }

            for (var i in data.items) {
                if (!data.items[i].exp_product_limits) {
                    data.items[i]._exp_product_limits = "";
                } else {
                    data.items[i]._exp_product_limits = data.items[i].exp_product_limits.join(", ");
                }
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

            data._ecoin_types = l4i.Clone(iamAcc.ecoin_types);

            l4iTemplate.Render({
                dstid: "iam-acc-chargelist",
                tplid: "iam-acc-chargelist-tpl",
                data: data,
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        iam.ApiCmd("account/charge-list", {
            callback: ep.done('data'),
        });

        iam.TplCmd("account/charge-list", {
            callback: ep.done('tpl'),
        });
    });
}

iamAcc.PayoutList = function() {
    var alert_id = "#iam-acc-payoutlist-alert";
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
            // iam.OpToolsRefresh("#iam-acc-payoutlist-optools");
            //
            if (data.items.length < 1) {
                return l4i.InnerAlert(alert_id, 'alert-info', "No Payout Statement Found at this Moment");
            }

            data._ecoin_types = l4i.Clone(iamAcc.ecoin_types);

            l4iTemplate.Render({
                dstid: "iam-acc-payoutlist",
                tplid: "iam-acc-payoutlist-tpl",
                data: data,
            });
        });

        ep.fail(function(err) {
            alert("Error: Please try again later");
        });

        iam.ApiCmd("account/charge-payout-list", {
            callback: ep.done('data'),
        });

        iam.TplCmd("account/payout-list", {
            callback: ep.done('tpl'),
        });
    });
}

