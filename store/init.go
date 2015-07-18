package store

import (
	"fmt"

	"github.com/lessos/bigtree/btapi"
	"github.com/lessos/lessgo/types"
	"github.com/lessos/lessgo/utils"
	"github.com/lessos/lessgo/utilx"

	"github.com/lessos/lessids/idsapi"
)

type InitNew struct {
}

func (i InitNew) Init() {

	//
	// privilege_sys_admin := idsapi.UserPrivilege{
	// 	ID:   "sys.admin",
	// 	Desc: "System Management",
	// }

	// privilege_user_admin := idsapi.UserPrivilege{
	// 	ID:   "user.admin",
	// 	Desc: "User Management",
	// }

	//
	role := idsapi.UserRole{
		Meta: types.ObjectMeta{
			ID:      "1",
			Name:    "Administrator",
			UserID:  utils.StringEncode16("sysadmin", 8),
			Created: utilx.TimeNow("atom"),
			Updated: utilx.TimeNow("atom"),
		},
		Desc:   "Root System Administrator",
		Status: 1,
	}
	rolejs, _ := utils.JsonEncode(role)
	BtAgent.ObjectSet(btapi.ObjectProposal{
		Meta: btapi.ObjectMeta{
			Path: fmt.Sprintf("/role/%s", role.Meta.ID),
		},
		Data: rolejs,
	})

	//
	role.Meta.ID = "100"
	role.Meta.Name = "Member"
	role.Desc = "Universal Member"
	rolejs, _ = utils.JsonEncode(role)
	BtAgent.ObjectSet(btapi.ObjectProposal{
		Meta: btapi.ObjectMeta{
			Path: fmt.Sprintf("/role/%s", role.Meta.ID),
		},
		Data: rolejs,
	})

	//
	role.Meta.ID = "101"
	role.Meta.Name = "Anonymous"
	role.Desc = "Anonymous Member"
	rolejs, _ = utils.JsonEncode(role)
	BtAgent.ObjectSet(btapi.ObjectProposal{
		Meta: btapi.ObjectMeta{
			Path: fmt.Sprintf("/role/%s", role.Meta.ID),
		},
		Data: rolejs,
	})
}
