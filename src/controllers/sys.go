package controllers

import (
    "../../deps/lessgo/pagelet"
    "../../deps/lessgo/utils"
    "io"
)

type Status struct {
    *pagelet.Controller
}

func (c Status) InfoAction() {

    c.AutoRender = false

    var rsp struct {
        ResponseJson
        Data struct {
            ServiceStatus string `json:"serviceStatus"`
        }   `json:"data"`
    }
    rsp.ApiVersion = apiVersion
    rsp.Status = 200
    rsp.Data.ServiceStatus = "ok"

    defer func() {
        if rspj, err := utils.JsonEncode(rsp); err == nil {
            io.WriteString(c.Response.Out, rspj)
        }
    }()
}
