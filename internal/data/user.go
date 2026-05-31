// Copyright 2014 Eryx <evorui at gmail dot com>, All rights reserved.
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
	"log/slog"
	"slices"
	"sync"
	"time"

	"github.com/hooto/iam/v2/pkg/iamapi"
)

var (
	userMu          sync.RWMutex
	userCaches      = map[string]*iamapi.User{}
	userGroupCaches = map[string]*iamapi.User{}
	userGroupMaps   = map[string][]string{}
	userRefreshed   = int64(0)

	Users Userset = &userSet{
		index:          map[string]*iamapi.User{},
		groupItems:     map[string]*iamapi.User{},
		userGroupIndex: map[string][]string{},
	}
)

type Userset interface {
	Len() int
	Iter(func(item *iamapi.User) bool)
	User(name string) (*iamapi.User, error)
}

type userSet struct {
	mu    sync.RWMutex
	items []*iamapi.User
	index map[string]*iamapi.User

	groupItems map[string]*iamapi.User
	// user -> groups
	userGroupIndex map[string][]string
}

func (it *userSet) Len() int {
	it.mu.RLock()
	defer it.mu.RUnlock()
	return len(it.items)
}

func (it *userSet) Iter(fn func(item *iamapi.User) bool) {
	it.mu.RLock()
	defer it.mu.RUnlock()
	for _, v := range it.items {
		if !fn(v) {
			break
		}
	}
}

func (it *userSet) User(username string) (*iamapi.User, error) {

	it.mu.Lock()
	defer it.mu.Unlock()

	if item, ok := it.index[username]; ok {
		return item, nil
	}

	var user iamapi.User
	if rs := Data.NewReader(iamapi.NsUser(username)).Exec(); rs.NotFound() {
		return nil, nil
	} else if !rs.OK() {
		return nil, rs.Error()
	} else if err := rs.Item().JsonDecode(&user); err != nil {
		return nil, err
	}

	it.index[username] = &user

	if user.Type == iamapi.UserTypeGroup {
		it.groupItems[username] = &user
		for _, kv := range append(user.Members, user.Owners...) {
			if !slices.Contains(it.userGroupIndex[kv], user.Name) {
				it.userGroupIndex[kv] = append(it.userGroupIndex[kv], user.Name)
			}
		}
	}

	return &user, nil
}

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

	var user iamapi.User
	if rs := Data.NewReader(iamapi.NsUser(name)).Exec(); !rs.OK() {
		return nil
	} else if err := rs.Item().JsonDecode(&user); err != nil {
		return nil
	}

	userCaches[name] = &user
	if user.Type == iamapi.UserTypeGroup {
		userGroupCaches[name] = &user
		for _, kv := range append(user.Members, user.Owners...) {
			if !slices.Contains(userGroupMaps[kv], user.Name) {
				userGroupMaps[kv] = append(userGroupMaps[kv], user.Name)
			}
		}
	}

	return &user
}

func UserSet(user *iamapi.User) bool {

	userCacheRefresh()

	userMu.Lock()
	defer userMu.Unlock()

	p, ok := userCaches[user.Name]
	if !ok || p.Updated >= user.Updated {

		if rs := Data.NewWriter(iamapi.NsUser(user.Name), nil).SetJsonValue(user).
			Exec(); !rs.OK() {
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

	offset := iamapi.NsUser("")
	cutset := append(iamapi.NsUser(""), 0xff)

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

			offset = iamapi.NsUser(user.Name)

			p, ok := userCaches[user.Name]
			if ok && p.Updated >= user.Updated {
				continue
			}

			userCaches[user.Name] = &user

			if user.Type == iamapi.UserTypeGroup {
				userGroupCaches[user.Name] = &user

				for _, kv := range append(user.Members, user.Owners...) {
					if !slices.Contains(userGroupMaps[kv], user.Name) {
						userGroupMaps[kv] = append(userGroupMaps[kv], user.Name)
					}
				}
			}
		}

		if !rs.NextResultSet {
			break
		}
	}

	slog.Debug("data/user cache refreshed", "count", len(userCaches))

	userRefreshed = tn
}
