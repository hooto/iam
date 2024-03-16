package signup

import (
	"errors"
	"strings"

	"github.com/hooto/httpsrv"

	"github.com/hooto/iam/iamapi"
)

func Validate(params *httpsrv.Params) error {

	params.SetValue("uname", strings.ToLower(strings.TrimSpace(params.Value("uname"))))
	if len(params.Value("uname")) < 4 || len(params.Value("uname")) > 30 {
		return errors.New("Username must be between 4 and 30 characters long")
	}
	if matched := iamapi.UsernameRE.MatchString(params.Value("uname")); !matched {
		return errors.New("Username must consist of letters or numbers, and begin with a letter")
	}

	params.SetValue("email", strings.TrimSpace(strings.ToLower(params.Value("email"))))
	if matched := iamapi.EmailRE.MatchString(params.Value("email")); !matched {
		return errors.New("Email is not valid")
	}

	// params.SetValue("name", strings.TrimSpace(params.Value("name")))
	// if len(params.Value("name")) == 0 || len(params.Value("name")) > 30 {
	// 	return errors.New("Name must be between 1 and 30 characters long")
	// }

	if len(params.Value("passwd")) < 8 || len(params.Value("passwd")) > 30 {
		return errors.New("Password must be between 8 and 30 characters long")
	}

	return nil
}

func ValidateEmail(user *iamapi.User) error {

	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	if matched := iamapi.EmailRE.MatchString(user.Email); !matched {
		return errors.New("Email is not valid")
	}

	return nil
}

func ValidateUsername(user *iamapi.User) error {

	user.Name = strings.ToLower(strings.TrimSpace(user.Name))
	if len(user.Name) < 4 || len(user.Name) > 30 {
		return errors.New("Username must be between 4 and 30 characters long")
	}
	if matched := iamapi.UsernameRE.MatchString(user.Name); !matched {
		return errors.New("Username must consist of letters or numbers, and begin with a letter")
	}

	return nil
}
