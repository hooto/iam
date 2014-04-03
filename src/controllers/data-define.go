package controllers

import ()

const (
    apiVersion = "1.0.0"
)

type ResponseJson struct {
    Status     int    `json:"status"`
    Message    string `json:"message"`
    ApiVersion string `json:"apiVersion"`
}
