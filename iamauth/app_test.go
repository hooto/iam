// Copyright 2019 Eryx <evorui аt gmail dοt com>, All rights reserved.
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

package iamauth

import (
	"testing"
)

var (
	tAppAuthKey = &AuthKey{
		User:      "guest",
		AccessKey: "be2c1fcf532baaa9",
		SecretKey: "c9a1a8ca13740018f1dd840a073ffc2e",
	}
	tAppAuthKeyErr = &AuthKey{
		User:      "guest",
		AccessKey: "be2c1fcf532baaa9",
		SecretKey: "c9a1a8ca13740018",
	}
	tAppData = []byte(`{"id": "1234", "data": "hello world"}`)
)

func Test_AppMain(t *testing.T) {

	pl := NewAppCredential(tAppAuthKey)

	token := pl.SignToken(tAppData)

	t.Logf("AppSignToken %s", token)

	rs, err := NewAppValidator(token)
	if rs == nil || err != nil {
		t.Fatal("Failed on AppValid")
	}

	if rs.User != tAppAuthKey.User {
		t.Fatal("Failed on Token Decode")
	}

	if err := rs.SignValid(tAppData, tAppAuthKey); err != nil {
		t.Fatal("Failed on AppValid")
	}

	if err := rs.SignValid(tAppData, tAppAuthKeyErr); err == nil {
		t.Fatal("Failed on AppValid")
	}
}

func Benchmark_AppCredential_SignToken(b *testing.B) {
	ac := NewAppCredential(tAppAuthKey)
	for i := 0; i < b.N; i++ {
		ac.SignToken(tAppData)
	}
}

func Benchmark_AppValidator_SignValid(b *testing.B) {
	token := NewAppCredential(tAppAuthKey).SignToken(tAppData)
	for i := 0; i < b.N; i++ {
		rs, _ := NewAppValidator(token)
		rs.SignValid(tAppData, tAppAuthKey)
	}
}
