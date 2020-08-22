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

package iamapi

import (
	"errors"
	"regexp"

	"github.com/lessos/lessgo/types"
	"github.com/lynkdb/iomix/utils"
)

var (
	MsgIdReg = regexp.MustCompile("^[a-f0-9]{16,32}$")
)

const (
	MsgActionPostOK      uint32 = 1 << 1
	MsgActionPostError   uint32 = 1 << 2
	MsgActionPostTimeout uint32 = 1 << 3
)

func (it *MsgItem) Valid() error {

	if !MsgIdReg.MatchString(it.Id) {
		return errors.New("invalid msg id")
	}

	if !UsernameRE.MatchString(it.ToUser) {
		return errors.New("user not found")
	}

	return nil
}

func (it *MsgItem) SentId() string {
	return utils.Uint32ToHexString(it.Created) + it.Id
}

type MsgList struct {
	types.TypeMeta `json:",inline" toml:",inline"`
	Items          []*MsgItem `json:"items" toml:"items"`
}
