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

	"github.com/hooto/hlog4g/hlog"
	"github.com/lessos/lessgo/crypto/phash"
	"github.com/lessos/lessgo/net/email"
	"github.com/lessos/lessgo/types"
	kv2 "github.com/lynkdb/kvspec/go/kvspec/v2"

	"github.com/hooto/iam/config"
	"github.com/hooto/iam/iamapi"
)

var (
	Data                  kv2.Client
	def_sysadmin          = "sysadmin"
	def_sysadmin_password = "changeme"
	app_inst_id_re        = regexp.MustCompile("^[0-9a-f]{16}$")
)

func Init() error {

	if Data == nil {
		return fmt.Errorf("iam.store connect not ready #1")
	}

	if rs := Data.NewWriter([]byte("iam:test"), "test").
		ExpireSet(1000).Commit(); !rs.OK() {
		return fmt.Errorf("iam.store connect not ready #2 " + rs.String())
	}

	if rs := Data.NewReader([]byte("iam:test")).Query(); !rs.OK() || rs.DataValue().String() != "test" {
		return fmt.Errorf("iam.store connect not ready #3")
	} else {
		hlog.Printf("info", "iam/data connect ok")
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

	tnm := types.MetaTimeNow()

	for _, v := range []iamapi.UserRole{
		{
			Id:   1,
			Name: "sysadmin",
			Desc: "Root System Administrator",
		},
		{
			Id:   100,
			Name: "member",
			Desc: "Universal Member",
		},
		{
			Id:   101,
			Name: "developer",
			Desc: "Universal Developer",
		},
		{
			Id:   1000,
			Name: "guest",
			Desc: "Anonymous Guest",
		},
	} {

		v.User = def_sysadmin
		v.Status = 1
		v.Created = tnm
		v.Updated = tnm

		rw := Data.NewWriter(iamapi.ObjKeyRole(v.Name), v).
			IncrNamespaceSet("role").ModeCreateSet(true)
		rw.Meta.IncrId = uint64(v.Id)
		if rs := rw.Commit(); !rs.OK() {
			return fmt.Errorf("db err %s", rs.Message)
		}
	}

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

	if rs := Data.NewWriter(iamapi.ObjKeyAppInstance(inst.Meta.ID), inst).
		ModeCreateSet(true).Commit(); !rs.OK() {
		return fmt.Errorf("db err %s", rs.Message)
	}

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
	// uid := idhash.HashToHexString([]byte(def_sysadmin), 8)
	if rs := Data.NewReader(iamapi.ObjKeyUser(def_sysadmin)).Query(); rs.NotFound() {

		sysadm := iamapi.User{
			// Id:          uid,
			Name:        def_sysadmin,
			DisplayName: "System Administrator",
			Email:       "",
			Roles:       []uint32{1, 100},
			Status:      1,
			Created:     tnm,
			Updated:     tnm,
		}

		auth, err := phash.Generate(def_sysadmin_password)
		if err != nil {
			return err
		}

		sysadm.Keys.Set(iamapi.UserKeyDefault, auth)

		ow := Data.NewWriter(iamapi.ObjKeyUser(def_sysadmin), sysadm).
			ModeCreateSet(true).IncrNamespaceSet("user")
		ow.Meta.IncrId = 1
		if rs := ow.Commit(); !rs.OK() {
			return fmt.Errorf("db err %s", rs.Message)
		}
	}

	return nil
}

func SysConfigRefresh() {

	//
	if rs := Data.NewReader(iamapi.ObjKeySysConfig("mailer")).Query(); rs.OK() {

		var mailer iamapi.SysConfigMailer

		if err := rs.DataValue().Decode(&mailer, nil); err == nil && mailer.SmtpHost != "" {

			email.MailerRegister(
				"def",
				mailer.SmtpHost,
				mailer.SmtpPort,
				mailer.SmtpUser,
				mailer.SmtpPass,
			)
		}
	}

	if rs := Data.NewReader(iamapi.ObjKeySysConfig("service_name")).Query(); rs.OK() &&
		rs.DataValue().String() != "" {
		config.Config.ServiceName = rs.DataValue().String()
	}

	if rs := Data.NewReader(iamapi.ObjKeySysConfig("webui_banner_title")).Query(); rs.OK() &&
		rs.DataValue().String() != "" {
		config.Config.WebUiBannerTitle = rs.DataValue().String()
	}

	if rs := Data.NewReader(iamapi.ObjKeySysConfig("service_login_form_alert_msg")).Query(); rs.OK() {
		config.Config.ServiceLoginFormAlertMsg = rs.DataValue().String()
	}

	if rs := Data.NewReader(iamapi.ObjKeySysConfig("user_reg_disable")).Query(); rs.OK() {

		if rs.DataValue().String() == "1" {
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

	ak.Created = types.MetaTimeNow()
	ak.Action = 1
	if rs := Data.NewWriter(iamapi.ObjKeyAccessKey(ak.User, ak.AccessKey), ak).
		ModeCreateSet(true).Commit(); !rs.OK() {
		return rs.Error()
	}

	return nil
}

func AccessKeyReset(ak iamapi.AccessKey) error {

	if ak.User == "" {
		return errors.New("No User Set")
	}

	ak.Created = types.MetaTimeNow()
	ak.Action = 1

	if rs := Data.NewWriter(iamapi.ObjKeyAccessKey(ak.User, ak.AccessKey), ak).
		Commit(); !rs.OK() {
		return rs.Error()
	}

	return nil
}

func AppInstanceRegister(inst iamapi.AppInstance) error {

	if !app_inst_id_re.MatchString(inst.Meta.ID) {
		return fmt.Errorf("Invalid meta.id (%s)", inst.Meta.ID)
	}

	if rs := Data.NewReader(iamapi.ObjKeyUser(inst.Meta.User)).Query(); rs.NotFound() {
		inst.Meta.User = def_sysadmin
	}

	var prev iamapi.AppInstance
	if rs := Data.NewReader(iamapi.ObjKeyAppInstance(inst.Meta.ID)).Query(); rs.OK() {
		rs.DataValue().Decode(&prev, nil)
	}

	if prev.Meta.ID == "" {
		inst.Meta.Created = types.MetaTimeNow()
		inst.Meta.Updated = types.MetaTimeNow()

		Data.NewWriter(iamapi.ObjKeyAppInstance(inst.Meta.ID), inst).
			ModeCreateSet(true).Commit()

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

		Data.NewWriter(iamapi.ObjKeyAppInstance(prev.Meta.ID), prev).
			ModeCreateSet(true).Commit()

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

		if rs := Data.NewWriter(
			iamapi.ObjKeyRolePrivilege(rid, inst.Meta.ID), strings.Join(v, ",")).
			ModeCreateSet(true).Commit(); !rs.OK() {
			return fmt.Errorf("db err %s", rs.Message)
		}
	}

	return nil
}
