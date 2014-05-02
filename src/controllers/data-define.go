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

type ResponseSession struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    Uid          uint32 `json:"uid"`
    Uname        string `json:"uname"`
    Name         string `json:"name"`
    Data         string `json:"data"`
    Roles        string `json:"roles"`
    Expired      string `json:"expired"`
    Timezone     string `json:"timezone"`
}
