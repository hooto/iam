// Copyright 2014 lessos Authors, All rights reserved.
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

package store

import (
	"fmt"
	"regexp"
	"strings"

	"code.hooto.com/lessos/iam/config"
	"code.hooto.com/lessos/iam/iamapi"
	"code.hooto.com/lynkdb/iomix/skv"
	"github.com/lessos/lessgo/crypto/idhash"
	"github.com/lessos/lessgo/crypto/phash"
	"github.com/lessos/lessgo/net/email"
	"github.com/lessos/lessgo/types"
)

var (
	path_prefix           = ""
	Data                  skv.Connector
	def_sysadmin          = "sysadmin"
	def_sysadmin_password = "changeme"
	app_inst_id_re        = regexp.MustCompile("^[0-9a-f]{16}$")
)

//
func PvGet(path string) *skv.Result {
	return Data.PvGet(PathPrefixAppend(path))
}

func PvPut(path string, v interface{}, opts *skv.PathWriteOptions) *skv.Result {
	return Data.PvPut(PathPrefixAppend(path), v, opts)
}

func PvDel(path string, opts *skv.PathWriteOptions) *skv.Result {
	return Data.PvDel(PathPrefixAppend(path), opts)
}

func PvScan(path, offset, cutset string, limit int) *skv.Result {
	return Data.PvScan(PathPrefixAppend(path), offset, cutset, limit)
}

func PvRevScan(path, offset, cutset string, limit int) *skv.Result {
	return Data.PvRevScan(PathPrefixAppend(path), offset, cutset, limit)
}

//
func PoGet(path string, key interface{}) *skv.Result {
	return Data.PoGet(PathPrefixAppend(path), key)
}

func PoPut(path string, key, value interface{}, opts *skv.PathWriteOptions) *skv.Result {
	return Data.PoPut(PathPrefixAppend(path), key, value, opts)
}

func PoDel(path string, key interface{}, opts *skv.PathWriteOptions) *skv.Result {
	return Data.PoDel(PathPrefixAppend(path), key, opts)
}

func PoScan(path string, offset, cutset interface{}, limit int) *skv.Result {
	return Data.PoScan(PathPrefixAppend(path), offset, cutset, limit)
}

func PoRevScan(path string, offset, cutset interface{}, limit int) *skv.Result {
	return Data.PoRevScan(PathPrefixAppend(path), offset, cutset, limit)
}

//
func PathPrefixSet(path string) {
	path_prefix = path
}

func PathPrefixAppend(path string) string {

	if path_prefix == "" {
		return path
	}

	return path_prefix + "/" + path
}

func Init() error {

	if Data == nil {
		return fmt.Errorf("iam.store connect not ready")
	}

	if rs := Data.PoPut(PathPrefixAppend("iam-test"), uint32(1234), "test", &skv.PathWriteOptions{
		Ttl: 3000,
	}); !rs.OK() {
		return fmt.Errorf("iam.store connect not ready")
	}

	if rs := Data.PoGet(PathPrefixAppend("iam-test"), uint32(1234)); !rs.OK() ||
		rs.Bytex().String() != "test" {
		return fmt.Errorf("iam.store connect not ready")
	}

	return nil
}

func InitData() (err error) {

	if err := Init(); err != nil {
		return err
	}

	if len(config.Config.InstanceID) < 16 {
		return fmt.Errorf("No InstanceID Setup")
	}

	//
	role := iamapi.UserRole{
		Id:      1,
		Name:    "Administrator",
		User:    def_sysadmin,
		Desc:    "Root System Administrator",
		Status:  1,
		Created: types.MetaTimeNow(),
		Updated: types.MetaTimeNow(),
	}
	PoPut("role", role.Id, role, nil)

	//
	role.Id = 100
	role.Name = "Member"
	role.Desc = "Universal Member"
	PoPut("role", role.Id, role, nil)

	//
	role.Id = 101
	role.Name = "Developer"
	role.Desc = "Universal Developer"
	PoPut("role", role.Id, role, nil)

	//
	role.Id = 1000
	role.Name = "Anonymous"
	role.Desc = "Anonymous Member"
	PoPut("role", role.Id, role, nil)

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
		Meta: types.InnerObjectMeta{
			ID:      config.Config.InstanceID,
			User:    def_sysadmin,
			Created: types.MetaTimeNow(),
			Updated: types.MetaTimeNow(),
		},
		Version:    config.Version,
		AppID:      "iam",
		AppTitle:   "lessOS IAM Service",
		Status:     1,
		Url:        "",
		Privileges: ps,
	}
	PoPut("app-instance", inst.Meta.ID, inst, nil)

	AppInstanceRegister(inst)

	/*
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
			PoPut(fmt.Sprintf("role-privilege/%d/%s", rid, inst.Meta.ID), strings.Join(v, ","), nil)
		}
	*/

	// Init Super SysAdmin Account
	uid := idhash.HashToHexString([]byte(def_sysadmin), 8)
	if rs := PoGet("user", uid); rs.NotFound() {

		sysadm := iamapi.User{
			Id:          uid,
			Name:        def_sysadmin,
			DisplayName: "System Administrator",
			Email:       "",
			Roles:       []uint32{1, 100},
			Status:      1,
			Created:     types.MetaTimeNow(),
			Updated:     types.MetaTimeNow(),
		}

		auth, err := phash.Generate(def_sysadmin_password)
		if err != nil {
			return err
		}

		sysadm.Keys.Set(iamapi.UserKeyDefault, auth)

		PoPut("user", uid, sysadm, nil)
	}

	return nil
}

func SysConfigRefresh() {

	//
	var mailer iamapi.SysConfigMailer
	if obj := PvGet("sys-config/mailer"); obj.OK() {

		obj.Decode(&mailer)

		if mailer.SmtpHost != "" {

			email.MailerRegister(
				"def",
				mailer.SmtpHost,
				mailer.SmtpPort,
				mailer.SmtpUser,
				mailer.SmtpPass,
			)
		}
	}

	if obj := PvGet("sys-config/service_name"); obj.OK() {

		if obj.Bytex().String() != "" {
			config.Config.ServiceName = obj.Bytex().String()
		}
	}

	if obj := PvGet("sys-config/webui_banner_title"); obj.OK() {

		if obj.Bytex().String() != "" {
			config.Config.WebUiBannerTitle = obj.Bytex().String()
		}
	}

	if obj := PvGet("sys-config/user_reg_disable"); obj.OK() {

		if obj.Bytex().String() == "1" {
			config.UserRegistrationDisabled = true
		} else {
			config.UserRegistrationDisabled = false
		}
	}
}

func AppInstanceRegister(inst iamapi.AppInstance) error {

	if !app_inst_id_re.MatchString(inst.Meta.ID) {
		return fmt.Errorf("Invalid meta.id (%s)", inst.Meta.ID)
	}

	if rs := PoGet("user", iamapi.UserId(def_sysadmin)); rs.NotFound() {
		inst.Meta.User = def_sysadmin
	}

	var prev iamapi.AppInstance
	if rs := PoGet("app-instance", inst.Meta.ID); rs.OK() {
		rs.Decode(&prev)
	}

	if prev.Meta.ID == "" {
		inst.Meta.Created = types.MetaTimeNow()
		inst.Meta.Updated = types.MetaTimeNow()
		PoPut("app-instance", inst.Meta.ID, inst, nil)
	} else {

		prev.Meta.Updated = types.MetaTimeNow()

		if inst.Version != "" && inst.Version != prev.Version {
			prev.Version = inst.Version
		}

		if inst.AppTitle != "" && inst.AppTitle != prev.AppTitle {
			prev.AppTitle = inst.AppTitle
		}

		if inst.Status != prev.Status {
			prev.Status = inst.Status
		}

		if inst.Url != "" && inst.Url != prev.Url {
			prev.Url = inst.Url
		}

		PoPut("app-instance", prev.Meta.ID, prev, nil)

		// TODO remove unused privileges
	}

	// privilege
	rps := map[uint32][]string{}
	for _, v := range inst.Privileges {

		for _, rid := range v.Roles {

			if _, ok := rps[rid]; !ok {
				rps[rid] = []string{}
			}

			rps[rid] = append(rps[rid], v.Privilege)
		}
	}

	for rid, v := range rps {
		PoPut(fmt.Sprintf("role-privilege/%d", rid), inst.Meta.ID, strings.Join(v, ","), nil)
	}

	return nil
}
