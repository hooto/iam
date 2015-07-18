package role

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/lessos/bigtree/btapi"

	"github.com/lessos/lessids/idsapi"
	"github.com/lessos/lessids/store"
)

var (
	locker      sync.Mutex
	nextRefresh = time.Now()
	roles       = map[string][]string{}
	privileges  = map[string]string{}
)

//func innerRefresh() {

func innerRefresh() {

	//fmt.Println("init once")

	if nextRefresh.After(time.Now()) {
		return
	}

	locker.Lock()
	defer locker.Unlock()

	if objs := store.BtAgent.ObjectList(btapi.ObjectProposal{
		Meta: btapi.ObjectMeta{
			Path: "/app-instance/",
		},
	}); objs.Error == nil {

		for _, obj := range objs.Items {

			var inst idsapi.AppInstance
			if err := obj.JsonDecode(&inst); err == nil {

				for _, priv := range inst.Privileges {

					pkey := inst.Meta.ID + "." + priv.Privilege

					if _, ok := privileges[pkey]; ok {
						continue
					}

					privileges[pkey] = fmt.Sprintf("%d", priv.ID)
				}
			}
		}
	}

	//
	if objs := store.BtAgent.ObjectList(btapi.ObjectProposal{
		Meta: btapi.ObjectMeta{
			Path: "/role/",
		},
	}); objs.Error == nil {

		for _, obj := range objs.Items {

			var role idsapi.UserRole
			if err := obj.JsonDecode(&role); err == nil {

				if _, ok := roles[role.Meta.ID]; ok {
					continue
				}

				roles[role.Meta.ID] = role.Privileges
			}
		}
	}

	nextRefresh = time.Now().Add(time.Second * 60)
}

func AccessAllowed(role, instance, privilege string) bool {

	innerRefresh()

	pkey := instance + "." + privilege
	pid, ok := privileges[pkey]
	if !ok {
		return false
	}

	rs := strings.Split(role, ",")
	for _, v := range rs {

		if v2, ok := roles[v]; ok {

			for _, pid2 := range v2 {
				if pid2 == pid {
					return true
				}
			}
		}
	}

	return false
}
