package signup

import (
	"errors"
	"strings"

	"github.com/hooto/httpsrv"

	"code.hooto.com/lessos/iam/iamapi"
)

func Validate(params *httpsrv.Params) error {

	params.Set("uname", strings.ToLower(strings.TrimSpace(params.Get("uname"))))
	if len(params.Get("uname")) < 4 || len(params.Get("uname")) > 30 {
		return errors.New("Username must be between 4 and 30 characters long")
	}
	if matched := iamapi.UserNameRe2.MatchString(params.Get("uname")); !matched {
		return errors.New("Username must consist of letters or numbers, and begin with a letter")
	}

	params.Set("email", strings.TrimSpace(strings.ToLower(params.Get("email"))))
	if matched := iamapi.UserEmailRe2.MatchString(params.Get("email")); !matched {
		return errors.New("Email is not valid")
	}

	// params.Set("name", strings.TrimSpace(params.Get("name")))
	// if len(params.Get("name")) == 0 || len(params.Get("name")) > 30 {
	// 	return errors.New("Name must be between 1 and 30 characters long")
	// }

	if len(params.Get("passwd")) < 8 || len(params.Get("passwd")) > 30 {
		return errors.New("Password must be between 8 and 30 characters long")
	}

	return nil
}

func ValidateEmail(user *iamapi.User) error {

	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	if matched := iamapi.UserEmailRe2.MatchString(user.Email); !matched {
		return errors.New("Email is not valid")
	}

	return nil
}

func ValidateUsername(user *iamapi.User) error {

	user.Name = strings.ToLower(strings.TrimSpace(user.Name))
	if len(user.Name) < 4 || len(user.Name) > 30 {
		return errors.New("Username must be between 4 and 30 characters long")
	}
	if matched := iamapi.UserNameRe2.MatchString(user.Name); !matched {
		return errors.New("Username must consist of letters or numbers, and begin with a letter")
	}

	return nil
}
