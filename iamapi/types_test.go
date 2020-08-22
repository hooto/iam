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

package iamapi

import (
	"testing"
)

func Test_Username_Re(t *testing.T) {

	for _, uname := range []string{
		"12345",
		"-abcd",
		"abcd-",
		"ab",
		"abcd-01234567890123456789",
		"ABCD",
	} {
		if UsernameRE.MatchString(uname) {
			t.Fatalf("Username Test Fail %s", uname)
		}
	}

	for _, uname := range []string{
		"a123",
		"abc-1234",
	} {
		if !UsernameRE.MatchString(uname) {
			t.Fatalf("Username Test Fail %s", uname)
		}
	}
}
