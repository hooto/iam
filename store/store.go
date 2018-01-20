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

package store

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hooto/iam/config"
	"github.com/hooto/iam/iamapi"
	"github.com/lessos/lessgo/crypto/idhash"
	"github.com/lessos/lessgo/crypto/phash"
	"github.com/lessos/lessgo/net/email"
	"github.com/lessos/lessgo/types"
	"github.com/lynkdb/iomix/skv"
)

var (
	Data                  skv.Connector
	def_sysadmin          = "sysadmin"
	def_sysadmin_password = "changeme"
	app_inst_id_re        = regexp.MustCompile("^[0-9a-f]{16}$")
)

func Init() error {

	if Data == nil {
		return fmt.Errorf("iam.store connect not ready")
	}

	if rs := Data.ProgPut(
		skv.NewProgKey("iam", "test"),
		skv.NewValueObject("test"),
		&skv.ProgWriteOptions{
			Expired: uint64(time.Now().Add(3e9).UnixNano()),
		},
	); !rs.OK() {
		return fmt.Errorf("iam.store connect not ready")
	}

	if rs := Data.ProgGet(
		skv.NewProgKey("iam", "test"),
	); !rs.OK() ||
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
	Data.ProgNew(iamapi.DataRoleKey(role.Id), skv.NewValueObject(role), nil)

	//
	role.Id = 100
	role.Name = "Member"
	role.Desc = "Universal Member"
	Data.ProgNew(iamapi.DataRoleKey(role.Id), skv.NewValueObject(role), nil)

	//
	role.Id = 101
	role.Name = "Developer"
	role.Desc = "Universal Developer"
	Data.ProgNew(iamapi.DataRoleKey(role.Id), skv.NewValueObject(role), nil)

	//
	role.Id = 1000
	role.Name = "Anonymous"
	role.Desc = "Anonymous Member"
	Data.ProgNew(iamapi.DataRoleKey(role.Id), skv.NewValueObject(role), nil)

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
		AppTitle:   "hooto IAM Service",
		Status:     1,
		Url:        "",
		Privileges: ps,
	}
	Data.ProgNew(iamapi.DataAppInstanceKey(inst.Meta.ID), skv.NewValueObject(inst), nil)

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
	if rs := Data.ProgGet(iamapi.DataUserKey(def_sysadmin)); rs.NotFound() {

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

		Data.ProgNew(iamapi.DataUserKey(def_sysadmin), skv.NewValueObject(sysadm), nil)
	}

	return nil
}

func SysConfigRefresh() {

	//
	var mailer iamapi.SysConfigMailer
	if obj := Data.ProgGet(iamapi.DataSysConfigKey("mailer")); obj.OK() {

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

	if obj := Data.ProgGet(iamapi.DataSysConfigKey("service_name")); obj.OK() {

		if obj.Bytex().String() != "" {
			config.Config.ServiceName = obj.Bytex().String()
		}
	}

	if obj := Data.ProgGet(iamapi.DataSysConfigKey("webui_banner_title")); obj.OK() {

		if obj.Bytex().String() != "" {
			config.Config.WebUiBannerTitle = obj.Bytex().String()
		}
	}

	if obj := Data.ProgGet(iamapi.DataSysConfigKey("service_login_form_alert_msg")); obj.OK() {
		config.Config.ServiceLoginFormAlertMsg = obj.Bytex().String()
	}

	if obj := Data.ProgGet(iamapi.DataSysConfigKey("user_reg_disable")); obj.OK() {

		if obj.Bytex().String() == "1" {
			config.UserRegistrationDisabled = true
		} else {
			config.UserRegistrationDisabled = false
		}
	}
}

func AccessKeyInitData(ak iamapi.AccessKey) error {
	if ak.User == "" {
		return errors.New("No User Set")
	}

	ak.Created = uint64(types.MetaTimeNow())
	ak.Action = 1
	if rs := Data.ProgNew(
		iamapi.DataAccessKeyKey(ak.User, ak.AccessKey),
		skv.NewValueObject(ak),
		nil,
	); !rs.OK() {
		return errors.New(rs.Bytex().String())
	}

	return nil
}

func AppInstanceRegister(inst iamapi.AppInstance) error {

	if !app_inst_id_re.MatchString(inst.Meta.ID) {
		return fmt.Errorf("Invalid meta.id (%s)", inst.Meta.ID)
	}

	if rs := Data.ProgGet(iamapi.DataUserKey(inst.Meta.User)); rs.NotFound() {
		inst.Meta.User = def_sysadmin
	}

	var prev iamapi.AppInstance
	if rs := Data.ProgGet(iamapi.DataAppInstanceKey(inst.Meta.ID)); rs.OK() {
		rs.Decode(&prev)
	}

	if prev.Meta.ID == "" {
		inst.Meta.Created = types.MetaTimeNow()
		inst.Meta.Updated = types.MetaTimeNow()
		Data.ProgNew(
			iamapi.DataAppInstanceKey(inst.Meta.ID),
			skv.NewValueObject(inst),
			nil,
		)
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

		if inst.SecretKey != "" && inst.SecretKey != prev.SecretKey {
			prev.SecretKey = inst.SecretKey
		}

		Data.ProgNew(
			iamapi.DataAppInstanceKey(prev.Meta.ID),
			skv.NewValueObject(prev),
			nil,
		)

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
		Data.ProgNew(
			iamapi.DataRolePrivilegeKey(rid, inst.Meta.ID),
			skv.NewValueObject(strings.Join(v, ",")),
			nil,
		)
	}

	return nil
}
