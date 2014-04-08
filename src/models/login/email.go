package login

import (
    "../../../deps/lessgo/pagelet"
    "errors"
    "regexp"
    "strings"
)

var emailPattern = regexp.MustCompile("^[_a-z0-9-]+(\\.[_a-z0-9-]+)*@[a-z0-9-]+(\\.[a-z0-9-]+)*(\\.[a-z]{2,10})$")

func EmailSetValidate(params *pagelet.Params) error {

    params.Set("email", strings.TrimSpace(strings.ToLower(params.Get("email"))))
    if matched := emailPattern.MatchString(params.Get("email")); !matched {
        return errors.New("Email is not valid")
    }

    return nil
}
