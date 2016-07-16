// Copyright 2014-2016 iam Author, All rights reserved.
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

package login

import (
	"errors"
	"regexp"
	"strings"
)

var emailPattern = regexp.MustCompile("^[_a-z0-9-]+(\\.[_a-z0-9-]+)*@[a-z0-9-]+(\\.[a-z0-9-]+)*(\\.[a-z]{2,10})$")

func EmailSetValidate(email string) (string, error) {

	email = strings.TrimSpace(strings.ToLower(email))
	if matched := emailPattern.MatchString(email); !matched {
		return email, errors.New("Email is not valid")
	}

	return email, nil
}
