package profile

import (
    "../../../deps/lessgo/pagelet"
    "errors"
    "strings"
    "time"
)

func PutValidate(params *pagelet.Params) error {

    params.Set("name", strings.TrimSpace(params.Get("name")))
    if len(params.Get("name")) == 0 || len(params.Get("name")) > 30 {
        return errors.New("Name must be between 1 and 30 characters long")
    }

    if _, err := time.Parse("2006-01-02", params.Get("birthday")); err != nil {
        return errors.New("Birthday is not valid")
    }

    if len(params.Get("aboutme")) == 0 {
        return errors.New("About Me can not be null")
    }

    return nil
}
