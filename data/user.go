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

package data

import (
	"sync"
	"time"

	"github.com/hooto/hlog4g/hlog"
	"github.com/hooto/iam/iamapi"
)

var (
	userMu          sync.RWMutex
	userCaches      = map[string]*iamapi.User{}
	userGroupCaches = map[string]*iamapi.User{}
	userGroupMaps   = map[string][]string{}
	userRefreshed   = int64(0)
)

func UserList() map[string]*iamapi.User {

	userCacheRefresh()

	userMu.RLock()
	defer userMu.RUnlock()

	return userCaches
}

func UserGet(name string) *iamapi.User {

	userCacheRefresh()

	userMu.RLock()
	defer userMu.RUnlock()

	if p, ok := userCaches[name]; ok {
		return p
	}

	return nil
}

func UserSet(user *iamapi.User) bool {

	userCacheRefresh()

	userMu.Lock()
	defer userMu.Unlock()

	p, ok := userCaches[user.Name]
	if !ok || p.Updated >= user.Updated {

		if rs := Data.NewWriter(iamapi.ObjKeyUser(user.Name), nil).SetJsonValue(user).
			SetIncr(0, "user").Exec(); !rs.OK() {
			return false
		}

		userCaches[user.Name] = user
		if user.Type == iamapi.UserTypeGroup {
			userGroupCaches[user.Name] = user
		}
	}
	return true
}

func UserGroupList() map[string]*iamapi.User {

	userCacheRefresh()

	userMu.RLock()
	defer userMu.RUnlock()

	return userGroupCaches
}

func UserGroups(name string) []string {

	userMu.RLock()
	defer userMu.RUnlock()

	if v, ok := userGroupMaps[name]; ok {
		return v
	}
	return nil
}

func userCacheRefresh() {

	userMu.Lock()
	defer userMu.Unlock()

	tn := time.Now().Unix()

	if (userRefreshed + 600) > tn {
		return
	}

	offset := iamapi.ObjKeyUser("")
	cutset := append(iamapi.ObjKeyUser(""), 0xff)

	for {

		rs := Data.NewRanger(offset, cutset).SetLimit(1000).Exec()
		if !rs.OK() {
			break
		}

		for _, obj := range rs.Items {

			var user iamapi.User
			if err := obj.JsonDecode(&user); err != nil {
				continue
			}

			offset = iamapi.ObjKeyUser(user.Name)

			p, ok := userCaches[user.Name]
			if ok && p.Updated >= user.Updated {
				continue
			}

			userCaches[user.Name] = &user

			if user.Type == iamapi.UserTypeGroup {
				userGroupCaches[user.Name] = &user

				for _, kv := range append(user.Members, user.Owners...) {
					if !iamapi.ArrayStringHas(userGroupMaps[kv], user.Name) {
						userGroupMaps[kv] = append(userGroupMaps[kv], user.Name)
					}
				}
			}
		}

		if !rs.NextResultSet {
			break
		}
	}

	hlog.Printf("debug", "data/user cache refreshed %d", len(userCaches))

	userRefreshed = tn
}
