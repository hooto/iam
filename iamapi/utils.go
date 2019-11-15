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
	"sync"
)

var (
	arrayStringMu sync.RWMutex
)

func ArrayStringHas(ls []string, a string) bool {

	for _, v := range ls {
		if v == a {
			return true
		}
	}

	return false
}

func ArrayStringSet(ls *[]string, a string) bool {

	arrayStringMu.Lock()
	defer arrayStringMu.Unlock()

	if ls == nil {
		*ls = []string{}
	}

	for _, v := range *ls {
		if v == a {
			return false
		}
	}
	*ls = append(*ls, a)

	return true
}

func ArrayStringDel(ls *[]string, a string) bool {

	arrayStringMu.Lock()
	defer arrayStringMu.Unlock()

	if ls != nil {
		for i, v := range *ls {
			if v == a {
				*ls = append((*ls)[:i], (*ls)[i+1:]...)
				return true
			}
		}
	}
	return false
}

func ArrayStringEqual(lsa, lsb []string) bool {

	if len(lsa) != len(lsb) {
		return false
	}

	for _, v := range lsa {

		hit := false

		for _, v2 := range lsb {

			if v == v2 {
				hit = true
				break
			}
		}

		if !hit {
			return false
		}
	}

	return true
}
