// Copyright 2014 Eryx <evorui at gmail dot com>, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package iamapi

import (
	"errors"
	"regexp"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Validator func(string) error

// dnsLabelRegexp matches a single RFC 1123 DNS label:
// - lowercase letters, digits, and hyphens only
// - must not start or end with a hyphen
var dnsLabelRegexp = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$`)

var (
	UsernameValid Validator

	RoleValid Validator

	AppIdValid Validator

	ObjectIdValid = regexp.MustCompile("^[0-9a-f]{12,16}$")

	EmailValid Validator

	// RFC 1123 DNS label
	DNSLabelValid Validator

	validate = validator.New()

	trans ut.Translator
)

func init() {

	var (
		en  = en.New()
		uni = ut.New(en, en)
	)

	trans, _ = uni.GetTranslator("en")

	en_translations.RegisterDefaultTranslations(validate, trans)

	//
	UsernameValid = newValidator("required,min=3,max=30,alphanum")

	RoleValid = newValidator("required,min=2,max=20,alphanum")

	EmailValid = newValidator("required,email")

	AppIdValid = newValidator("required,min=8,max=30,alphanum")

	// Register custom RFC 1123 label validator
	validate.RegisterValidation("dns_label", func(fl validator.FieldLevel) bool {
		return dnsLabelRegexp.MatchString(fl.Field().String())
	})

	DNSLabelValid = newValidator("required,dns_label,min=3,max=63")
}

func newValidator(rule string) Validator {
	return func(str string) error {
		if err := validate.Var(str, rule); err != nil {
			errs := err.(validator.ValidationErrors)
			return errors.New(errs[0].Translate(trans))
		}
		return nil
	}
}
