package store

import (
	"fmt"
	"strings"

	"github.com/lessos/lessgo/types"
	"github.com/lessos/lessgo/utils"
	"github.com/lessos/lessgo/utilx"

	"github.com/lessos/iam/iamapi"
)

type InitNew struct {
}

func (i InitNew) Init() {

	//
	role := iamapi.UserRole{
		Meta: types.ObjectMeta{
			ID:      utils.StringEncode16("1", 8),
			Name:    "Administrator",
			UserID:  utils.StringEncode16("sysadmin", 8),
			Created: utilx.TimeNow("atom"),
			Updated: utilx.TimeNow("atom"),
		},
		IdxID:  1,
		Desc:   "Root System Administrator",
		Status: 1,
	}
	BtAgent.ObjectSet(fmt.Sprintf("/global/iam/role/%s", role.Meta.ID), role, nil)

	//
	role.Meta.ID = utils.StringEncode16("100", 8)
	role.Meta.Name = "Member"
	role.IdxID = 100
	role.Desc = "Universal Member"
	BtAgent.ObjectSet(fmt.Sprintf("/global/iam/role/%s", role.Meta.ID), role, nil)

	//
	role.Meta.ID = utils.StringEncode16("1000", 8)
	role.Meta.Name = "Anonymous"
	role.IdxID = 1000
	role.Desc = "Anonymous Member"
	BtAgent.ObjectSet(fmt.Sprintf("/global/iam/role/%s", role.Meta.ID), role, nil)

	//
	ps := []iamapi.AppPrivilege{
		{
			Privilege: "sys.admin",
			Roles:     []uint32{1},
			Desc:      "System Management",
		},
		{
			Privilege: "user.admin",
			Roles:     []uint32{1},
			Desc:      "User Management",
		},
	}

	inst := iamapi.AppInstance{
		Meta: types.ObjectMeta{
			ID:      utils.StringEncode16("lessos-iam", 12),
			Created: utilx.TimeNow("atom"),
			Updated: utilx.TimeNow("atom"),
		},
		AppID:      "iam",
		AppTitle:   "lessOS IAM Service",
		Status:     1,
		Url:        "",
		Privileges: ps,
	}
	BtAgent.ObjectSet("/global/iam/app-instance/"+inst.Meta.ID, inst, nil)

	// privilege
	rps := map[uint32][]string{}
	for _, v := range ps {

		for _, rid := range v.Roles {

			if _, ok := rps[rid]; !ok {
				rps[rid] = []string{}
			}

			rps[rid] = append(rps[rid], v.Privilege)
		}
	}

	for rid, v := range rps {

		BtAgent.ObjectSet(fmt.Sprintf("/global/iam/role-privilege/%d/%s", rid, inst.Meta.ID), strings.Join(v, ","), nil)
	}
}
