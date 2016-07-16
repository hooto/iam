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

	"github.com/lessos/iam/iamapi"
)

func PassSetValidate(set iamapi.UserPasswordSet) error {

	if len(set.NewPassword) < 8 || len(set.NewPassword) > 30 {
		return errors.New("Password must be between 8 and 30 characters long")
	}

	return nil
}
