package signup

import (
	"errors"
	"regexp"
	"strings"

	"github.com/lessos/lessgo/httpsrv"

	"github.com/lessos/lessids/idsapi"
)

var (
	emailPattern = regexp.MustCompile("^[_a-z0-9-]+(\\.[_a-z0-9-]+)*@[a-z0-9-]+(\\.[a-z0-9-]+)*(\\.[a-z]{2,10})$")
	unamePattern = regexp.MustCompile("^[a-z]{1}[a-z0-9]{3,29}$")
	uidPattern   = regexp.MustCompile("^[0-9a-f]{8,8}$")
)

func Validate(params *httpsrv.Params) error {

	params.Set("uname", strings.ToLower(strings.TrimSpace(params.Get("uname"))))
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

func ValidateEmail(user *idsapi.User) error {

	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	if matched := emailPattern.MatchString(user.Email); !matched {
		return errors.New("Email is not valid")
	}

	return nil
}

func ValidateUserID(user *idsapi.User) error {

	user.Meta.ID = strings.ToLower(user.Meta.ID)
	if matched := uidPattern.MatchString(user.Meta.ID); !matched {
		return errors.New("UserID is not valid")
	}

	return nil
}

func ValidateUsername(user *idsapi.User) error {

	user.Meta.Name = strings.ToLower(strings.TrimSpace(user.Meta.Name))
	if len(user.Meta.Name) < 4 || len(user.Meta.Name) > 30 {
		return errors.New("Username must be between 4 and 30 characters long")
	}
	if matched := unamePattern.MatchString(user.Meta.Name); !matched {
		return errors.New("Username must consist of letters or numbers, and begin with a letter")
	}

	return nil
}
