// Copyright 2014-2016 iam Author, All rights reserved.
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

	"github.com/lessos/iam/iamapi"
	"github.com/lessos/iam/store"
	"github.com/lessos/lessgo/utils"
)

var (
	locker     sync.Mutex
	inst2perms = map[string]*perm_map{}
)

type perm_map struct {
	refresh time.Time
	owner   string
	maps    map[string][]uint32
}

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
			owner:   "",
			maps:    map[string][]uint32{},
		}
	}

	//
	if obj := store.BtAgent.ObjectGet("/global/iam/app-instance/" + instanceid); obj.Error == nil {

		var inst iamapi.AppInstance

		if err := obj.JsonDecode(&inst); err == nil {

			for _, ip := range inst.Privileges {

				if len(ip.Roles) > 0 {
					perm.owner = inst.Meta.UserID
					perm.maps[ip.Privilege] = ip.Roles
				}
			}

			inst2perms[instanceid] = perm
		}
	}

	return perm
}

func AccessAllowed(owner string, roles []uint32, instanceid, privilege string) bool {

	if instanceid == "" {
		instanceid = utils.StringEncode16("iam", 12)
	}

	p := instPerms(instanceid)
	if p.owner == owner {
		// logger.Printf("info", "acl ok owner %s %s", instanceid, privilege)
		return true
	}

	if mroles, ok := p.maps[privilege]; ok {

		for _, rid := range mroles {

			for _, diffrid := range roles {

				if rid == diffrid {
					// logger.Printf("info", "acl ok role %s %s", instanceid, privilege)
					return true
				}
			}
		}
	}

	// logger.Printf("info", "acl !! %s %s", instanceid, privilege)

	return false
}
