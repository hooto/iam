package role

import (
	"strings"
	"sync"
	"time"

	"github.com/lessos/lessgo/data/rdo"
	"github.com/lessos/lessgo/data/rdo/base"
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

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		return
	}

	q := base.NewQuerySet().From("ids_privilege").Limit(1000)
	rspri, err := dcn.Base.Query(q)
	if err != nil || len(rspri) == 0 {
		return
	}
	for _, v := range rspri {

		pkey := v.Field("instance").String() + "." + v.Field("privilege").String()
		if _, ok := privileges[pkey]; ok {
			continue
		}

		privileges[pkey] = v.Field("pid").String()
	}

	q = base.NewQuerySet().From("ids_role").Limit(1000)
	q.Where.And("status", 1)
	rsrole, err := dcn.Base.Query(q)
	if err != nil || len(rsrole) == 0 {
		return
	}

	for _, v := range rsrole {

		pid := v.Field("rid").String()

		if _, ok := roles[pid]; ok {
			continue
		}

		roles[pid] = strings.Split(v.Field("privileges").String(), ",")
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
