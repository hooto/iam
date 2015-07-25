// Copyright 2015 lessOS.com, All rights reserved.
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

package role

import (
	"sync"
	"time"

	"github.com/lessos/bigtree/btapi"
	"github.com/lessos/lessgo/utils"
	"github.com/lessos/lessids/idsapi"
	"github.com/lessos/lessids/store"
)

var (
	locker sync.Mutex
	// nextRefresh = time.Now().Add(time.Second * -61)
	// roles       = map[uint32]role_bound{}
	// privileges  = map[string]string{}
	inst2perms = map[string]*perm_map{}
)

type perm_map struct {
	refresh time.Time
	maps    map[string][]uint32
}

// func innerRefresh() {

// 	//fmt.Println("init once")

// 	if nextRefresh.After(time.Now()) {
// 		return
// 	}

// 	locker.Lock()
// 	defer locker.Unlock()

// 	//
// 	if objs := store.BtAgent.ObjectList(btapi.ObjectProposal{
// 		Meta: btapi.ObjectMeta{
// 			Path: "/role/",
// 		},
// 	}); objs.Error == nil {

// 		if len(objs.Items) < 3 {
// 			it := store.InitNew{}
// 			it.Init()
// 		}

// 		for _, obj := range objs.Items {

// 			var role idsapi.UserRole
// 			if err := obj.JsonDecode(&role); err == nil {

// 				if _, ok := roles[role.IdxID]; ok {
// 					continue
// 				}

// 				roles[role.IdxID] = role.Privileges
// 			}
// 		}
// 	}

// 	nextRefresh = time.Now().Add(time.Second * 60)
// }

func instPerms(instanceid string) *perm_map {

	locker.Lock()
	defer locker.Unlock()

	perm, ok := inst2perms[instanceid]
	if ok {

		if perm.refresh.After(time.Now()) {
			return perm
		}

		perm.refresh = time.Now().Add(time.Second * 60)

	} else {

		perm = &perm_map{
			refresh: time.Now().Add(time.Second * 60),
			maps:    map[string][]uint32{},
		}
	}

	//
	if obj := store.BtAgent.ObjectGet(btapi.ObjectProposal{
		Meta: btapi.ObjectMeta{
			Path: "/app-instance/" + instanceid,
		},
	}); obj.Error == nil {

		var inst idsapi.AppInstance

		if err := obj.JsonDecode(&inst); err == nil {

			for _, ip := range inst.Privileges {

				if len(ip.Roles) > 0 {
					perm.maps[ip.Privilege] = ip.Roles
				}
			}

			inst2perms[instanceid] = perm
		}
	}

	return perm
}

func AccessAllowed(roles []uint32, instanceid, privilege string) bool {

	if instanceid == "" {
		instanceid = utils.StringEncode16("lessids", 12)
	}

	p := instPerms(instanceid)

	if mroles, ok := p.maps[privilege]; ok {

		for _, rid := range mroles {

			for _, diffrid := range roles {

				if rid == diffrid {
					return true
				}
			}
		}
	}

	return false
}
