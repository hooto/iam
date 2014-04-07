package controllers

import (
    "../../deps/lessgo/pagelet"
    "../models/session"
    "fmt"
)

type User struct {
    *pagelet.Controller
}

func (c User) IndexAction() {

    s := session.GetSession(c.Request)

    fmt.Println("Login", s)

}
