package controllers

import (
    //"../../deps/lessgo/data/rdc"
    "../../deps/lessgo/pagelet"
    //"../../deps/lessgo/pass"
    //"../../deps/lessgo/utils"
    //"../models/login"
    //"../models/profile"
    //"../models/session"
    //"encoding/base64"
    "fmt"
    //"html"
    //"io"
    //"strings"
)

type UserAdmin struct {
    *pagelet.Controller
}

func (c UserAdmin) IndexAction() {

    if !c.Session.AccessAllowed("user.admin") {
        fmt.Println("user.admin denied")
    } else {
        fmt.Println("user.admin AccessAllowed")
    }

    if !c.Session.AccessAllowed("user.admin2") {
        fmt.Println("user.admin2 denied")
    } else {
        fmt.Println("user.admin2 AccessAllowed")
    }
}
