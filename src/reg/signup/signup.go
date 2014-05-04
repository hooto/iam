package signup

import (
    "../../../deps/lessgo/pagelet"
    "errors"
    "regexp"
    "strings"
)

var (
    emailPattern = regexp.MustCompile("^[_a-z0-9-]+(\\.[_a-z0-9-]+)*@[a-z0-9-]+(\\.[a-z0-9-]+)*(\\.[a-z]{2,10})$")

    unamePattern = regexp.MustCompile("^[a-z]{1}[a-z0-9]{3,29}$")
)

func Validate(params *pagelet.Params) error {

    params.Set("uname", strings.TrimSpace(params.Get("uname")))
    if len(params.Get("uname")) < 4 || len(params.Get("uname")) > 30 {
        return errors.New("Username must be between 4 and 30 characters long")
    }
    if matched := unamePattern.MatchString(params.Get("uname")); !matched {
        return errors.New("Username must consist of letters or numbers, and begin with a letter")
    }

    params.Set("email", strings.TrimSpace(strings.ToLower(params.Get("email"))))
    if matched := emailPattern.MatchString(params.Get("email")); !matched {
        return errors.New("Email is not valid")
    }

    params.Set("name", strings.TrimSpace(params.Get("name")))
    if len(params.Get("name")) == 0 || len(params.Get("name")) > 30 {
        return errors.New("Name must be between 1 and 30 characters long")
    }

    if len(params.Get("passwd")) < 8 || len(params.Get("passwd")) > 30 {
        return errors.New("Password must be between 8 and 30 characters long")
    }

    return nil
}
