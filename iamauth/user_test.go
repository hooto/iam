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
	"time"
)

var (
	tKeys = []*AuthKey{
		{AccessKey: "be2c1fcf532baaa9", SecretKey: "c9a1a8ca13740018f1dd840a073ffc2e"},
		{AccessKey: "d4d7d973aa8d3c70", SecretKey: "ec1f6f37c8d81b7bdb855b651523367e"},
	}
	tKeyErrs = []*AuthKey{
		{AccessKey: "be2c1fcf532baaa9", SecretKey: "c9a1a8ca13740018f"},
	}
	tKeyNull     = []*AuthKey{}
	tPayloadItem = &UserPayload{
		Id:      "guest",
		Roles:   []uint32{100, 200},
		Groups:  []string{"staff"},
		Expired: 2012345678,
	}
	tToken = tPayloadItem.SignToken(tKeys)
)

func Test_UserMain(t *testing.T) {

	pl := NewUserPayload(
		"guest",
		"Guest",
		[]uint32{100, 200},
		[]string{"guest"},
		86400)

	pl.Expired = time.Now().Unix() + 1

	token := pl.SignToken(tKeys)
	t.Logf("SignToken %s", token)

	rs, err := NewUserValidator(token)
	if rs == nil || err != nil {
		t.Fatal("Failed on UserValid")
	}
	if rs.Id != tPayloadItem.Id {
		t.Fatal("Failed on Token Decode")
	}

	if err := rs.SignValid(tKeyErrs); err == nil {
		t.Fatal("Failed on UserValid")
	}

	if err := rs.SignValid(tKeyNull); err == nil {
		t.Fatal("Failed on UserValid")
	}

	time.Sleep(2e9) // expired

	if err := rs.SignValid(tKeys); err == nil {
		t.Fatal("Failed on UserValid")
	}
}

func Benchmark_UserPayload_SignToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tPayloadItem.SignToken(tKeys)
	}
}

func Benchmark_UserValidator_SignValid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rs, _ := NewUserValidator(tToken)
		rs.SignValid(tKeys)
	}
}
