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
	"encoding/base64"
	"strings"

	"github.com/lessos/lessgo/encoding/json"
)

func UserNameFilter(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	name2 := ""
	for _, v := range name {
		if (v >= 'a' && v <= 'z') || (v >= '0' || v <= '9') || (v == '-') || (v == '_') {
			name2 += string(v)
		}
	}
	return name2
}

type ServiceRedirectToken struct {
	RedirectUri string `json:"uri,omitempty" toml:"uri,omitempty"`
	State       string `json:"state,omitempty" toml:"state,omitempty"`
	ClientId    string `json:"cid,omitempty" toml:"cid,omitempty"`
	Persistent  int    `json:"pt,omitempty" toml:"pt,omitempty"`
}

func ServiceRedirectTokenValid(tokenstr string) bool {
	if _, err := base64.StdEncoding.DecodeString(tokenstr); err == nil {
		return true
	}
	return false
}

func (s *ServiceRedirectToken) Encode() string {

	if len(s.RedirectUri) > 200 {
		s.RedirectUri = s.RedirectUri[:200]
	}
	if len(s.State) > 100 {
		s.State = s.State[:100]
	}
	if len(s.ClientId) > 40 {
		s.ClientId = s.ClientId[:40]
	}

	js, _ := json.Encode(s, "")
	return base64.StdEncoding.EncodeToString(js)
}

func ServiceRedirectTokenDecode(tokenstr string) ServiceRedirectToken {
	var token ServiceRedirectToken
	if jsb, err := base64.StdEncoding.DecodeString(tokenstr); err == nil {
		json.Decode(jsb, &token)
	}
	return token
}
