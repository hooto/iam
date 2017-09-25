// Copyright 2014 Eryx <evorui аt gmаil dοt cοm>, All rights reserved.
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

package profile

import (
	"errors"
	"strings"
	"time"

	"github.com/hooto/iam/iamapi"
)

func PutValidate(set iamapi.UserProfile) (iamapi.UserProfile, error) {

	set.Login.DisplayName = strings.TrimSpace(set.Login.DisplayName)
	if len(set.Login.DisplayName) < 1 || len(set.Login.DisplayName) > 30 {
		return set, errors.New("DisplayName must be between 1 and 30 characters long")
	}

	if _, err := time.Parse("2006-01-02", set.Birthday); err != nil {
		return set, errors.New("Birthday is not valid")
	}

	if set.About == "" {
		return set, errors.New("About Me can not be null")
	}

	return set, nil
}
