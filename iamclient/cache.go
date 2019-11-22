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

package iamclient

import (
	"sync"
	"time"

	"github.com/hooto/iam/iamapi"
)

var (
	cacheRefreshed    = int64(0)
	userSessionMu     sync.RWMutex
	userSessionCaches = map[string]*iamapi.UserSession{}
	akSessionCaches   = map[string]*iamapi.AccessKeySession{}
)

func cacheRefresh() {

	userSessionMu.Lock()
	defer userSessionMu.Unlock()

	//
	tn := time.Now().Unix()
	if (cacheRefreshed + 60) > tn {
		return
	}

	//
	dels := []string{}
	for k, v := range userSessionCaches {
		if v.Expired <= tn {
			dels = append(dels, k)
		}
	}
	for _, k := range dels {
		delete(userSessionCaches, k)
	}

	//
	dels = []string{}
	for k, v := range akSessionCaches {
		if v.Expired <= tn {
			dels = append(dels, k)
		}
	}
	for _, k := range dels {
		delete(akSessionCaches, k)
	}

	cacheRefreshed = tn
}

func sessionCacheRefresh(v *iamapi.UserSession, tn int64) bool {
	if tn == 0 {
		tn = time.Now().Unix()
	}
	if (v.Cached + 60) < tn {
		return true
	}
	return false
}

func SessionCache(id string) *iamapi.UserSession {

	cacheRefresh()

	userSessionMu.RLock()
	defer userSessionMu.RUnlock()

	if v, ok := userSessionCaches[id]; ok {
		return v
	}
	return nil
}

func SessionSync(v *iamapi.UserSession, tn int64) {
	userSessionMu.Lock()
	defer userSessionMu.Unlock()
	if tn > 0 {
		v.Cached = tn
	} else {
		v.Cached = time.Now().Unix()
	}
	userSessionCaches[v.AccessToken] = v
}

func akSessionCache(id string) *iamapi.AccessKeySession {

	userSessionMu.RLock()
	defer userSessionMu.RUnlock()

	if v, ok := akSessionCaches[id]; ok {
		return v
	}
	return nil
}

func akSessionSync(v *iamapi.AccessKeySession) {
	userSessionMu.Lock()
	defer userSessionMu.Unlock()
	akSessionCaches[v.AccessKey] = v
}
