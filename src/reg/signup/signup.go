package signup

import (
    "../../../deps/lessgo/pagelet"
    "errors"
    "regexp"
    "strings"
)

var emailPattern = regexp.MustCompile("^[_a-z0-9-]+(\\.[_a-z0-9-]+)*@[a-z0-9-]+(\\.[a-z0-9-]+)*(\\.[a-z]{2,10})$")

func Validate(params *pagelet.Params) error {

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
