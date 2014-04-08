package login

import (
    "../../../deps/lessgo/pagelet"
    "errors"
)

func PassSetValidate(params *pagelet.Params) error {

    if len(params.Get("passwd")) < 8 || len(params.Get("passwd")) > 30 {
        return errors.New("Password must be between 8 and 30 characters long")
    }

    if params.Get("passwd") != params.Get("passwd_confirm") {
        return errors.New("Passwords do not match")
    }

    return nil
}
