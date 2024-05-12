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

package data

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/hooto/hauth/go/hauth/v1"
	"github.com/hooto/hlog4g/hlog"
	"github.com/lessos/lessgo/crypto/idhash"
	"github.com/lessos/lessgo/crypto/phash"
	"github.com/lessos/lessgo/net/email"
	"github.com/lessos/lessgo/types"
	"github.com/lynkdb/kvgo/v2/pkg/kvapi"

	"github.com/hooto/iam/config"
	"github.com/hooto/iam/iamapi"
)

var (
	Data                    kvapi.Client
	def_sysadmin            = "sysadmin"
	DefaultSysadminPassword = "changeme"
	app_inst_id_re          = regexp.MustCompile("^[0-9a-f]{16}$")
	KeyMgr                  = hauth.NewAccessKeyManager()
	mu                      sync.RWMutex
	inited                  = false
)

func Setup() error {

	mu.Lock()
	defer mu.Unlock()

	if inited {
		return nil
	}

	if Data == nil {
		return fmt.Errorf("iam.data connect not ready #1")
	}

	if rs := Data.NewWriter([]byte("iam:test"), []byte("test")).
		SetTTL(1000).Exec(); !rs.OK() {
		return fmt.Errorf("iam.data connect not ready #2 " + rs.String())

	} else if rs = Data.NewReader([]byte("iam:test")).Exec(); !rs.OK() || rs.Item().StringValue() != "test" {
		return fmt.Errorf("iam.data connect not ready #3")
	} else {
		hlog.Printf("info", "iam/data connect ok")
	}

	// AccessKey
	var (
		ak0 = iamapi.NsAccessKey("", "")
		akz = append(iamapi.NsAccessKey("", ""), []byte{0xff}...)
		akn = 0
	)
	for {

		rs := Data.NewRanger(ak0, akz).SetLimit(1000).Exec()
		if !rs.OK() {
			break
		}

		for _, v := range rs.Items {

			ak0 = v.Key

			var ak hauth.AccessKey
			if err := v.JsonDecode(&ak); err == nil {
				KeyMgr.KeySet(&ak)
				hlog.Printf("debug", "iam/access_key load %s", ak.Id)
				continue
			}
		}

		akn += len(rs.Items)

		if len(rs.Items) < 1000 {
			break
		}
	}
	hlog.Printf("info", "iam/access_key load %d from database", akn)

	//
	for _, v := range config.Config.AccessKeys {
		KeyMgr.KeySet(v)
	}
	hlog.Printf("info", "iam/access_key load %d from config", len(config.Config.AccessKeys))

	inited = true

	return nil
}

func InitData() (err error) {

	if err := Setup(); err != nil {
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

		rw := Data.NewWriter(iamapi.ObjKeyRole(v.Name), nil).SetJsonValue(v).
			SetIncr(uint64(v.Id), "role").SetCreateOnly(true)
		if rs := rw.Exec(); !rs.OK() {
			return fmt.Errorf("db err init-role %s", rs.ErrorMessage())
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

	if rs := Data.NewWriter(iamapi.ObjKeyAppInstance(inst.Meta.ID), nil).SetJsonValue(inst).
		SetCreateOnly(true).Exec(); !rs.OK() {
		return fmt.Errorf("db err init-data %s", rs.ErrorMessage())
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
	if rs := Data.NewReader(iamapi.ObjKeyUser(def_sysadmin)).Exec(); rs.NotFound() {

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

		auth, err := phash.Generate(DefaultSysadminPassword)
		if err != nil {
			return err
		}

		sysadm.Keys.Set(iamapi.UserKeyDefault, auth)

		ow := Data.NewWriter(iamapi.ObjKeyUser(def_sysadmin), nil).SetJsonValue(sysadm).
			SetCreateOnly(true).SetIncr(1, "user")
		if rs := ow.Exec(); !rs.OK() {
			return fmt.Errorf("db err init-sysadmin %s", rs.ErrorMessage())
		}

		var (
			tn     = time.Now()
			set_id = append(iamapi.Uint32ToBytes(uint32(tn.Unix())), idhash.Rand(8)...) // 4 + 8
		)

		set := iamapi.AccountFund{
			Id:               iamapi.BytesToHexString(set_id),
			Amount:           100000,
			Type:             iamapi.AccountCurrencyTypeVirtual,
			User:             def_sysadmin,
			Operator:         def_sysadmin,
			ExpProductMax:    100,
			ExpProductLimits: []types.NameIdentifier{"sys/pod"},
			Created:          types.MetaTimeSet(tn),
			Updated:          types.MetaTimeSet(tn),
			Priority:         8,
		}

		acc_user := iamapi.AccountUser{
			User:    def_sysadmin,
			Balance: set.Amount,
			Updated: set.Updated,
		}

		type keyValue struct {
			Key   []byte
			Value interface{}
		}

		sets := []keyValue{
			{
				Key:   iamapi.ObjKeyAccFundMgr(set.Id),
				Value: set,
			},
			{
				Key:   iamapi.ObjKeyAccFundUser(set.User, set.Id),
				Value: set,
			},
			{
				Key:   iamapi.ObjKeyAccUser(set.User),
				Value: acc_user,
			},
		}

		for _, v := range sets {
			Data.NewWriter(v.Key, v.Value).Exec()
		}
	}

	return nil
}

func SysConfigRefresh() {

	//
	if rs := Data.NewReader(iamapi.ObjKeySysConfig("mailer")).Exec(); rs.OK() {

		var mailer iamapi.SysConfigMailer

		if err := rs.Item().JsonDecode(&mailer); err == nil && mailer.SmtpHost != "" {

			email.MailerRegister(
				"def",
				mailer.SmtpHost,
				mailer.SmtpPort,
				mailer.SmtpUser,
				mailer.SmtpPass,
			)
		}
	}

	if rs := Data.NewReader(iamapi.ObjKeySysConfig("service_name")).Exec(); rs.OK() &&
		rs.Item().StringValue() != "" {
		config.Config.ServiceName = rs.Item().StringValue()
	}

	if rs := Data.NewReader(iamapi.ObjKeySysConfig("webui_banner_title")).Exec(); rs.OK() &&
		rs.Item().StringValue() != "" {
		config.Config.WebUiBannerTitle = rs.Item().StringValue()
	}

	if rs := Data.NewReader(iamapi.ObjKeySysConfig("service_login_form_alert_msg")).Exec(); rs.OK() {
		config.Config.ServiceLoginFormAlertMsg = rs.Item().StringValue()
	}

	if rs := Data.NewReader(iamapi.ObjKeySysConfig("user_reg_disable")).Exec(); rs.OK() {

		if rs.Item().StringValue() == "1" {
			config.UserRegistrationDisabled = true
		} else {
			config.UserRegistrationDisabled = false
		}
	}
}

func AccessKeyInitData(ak *hauth.AccessKey) error {

	if ak.User == "" {
		ak.User = "sysadmin"
	}

	ak.Status = hauth.AccessKeyStatusActive

	if rs := Data.NewReader(iamapi.NsAccessKey(ak.User, ak.Id)).Exec(); rs.OK() {
		var prev hauth.AccessKey
		if err := rs.Item().JsonDecode(&prev); err != nil {
			return err
		}
		if !ak.Equal(&prev) {
			if rs := Data.NewWriter(iamapi.NsAccessKey(ak.User, ak.Id), nil).SetJsonValue(ak).
				Exec(); !rs.OK() {
				return rs.Error()
			}
			hlog.Printf("warn", "IAM/AK ID %s, SEC %s, Refreshed", ak.Id, ak.Secret[:8])
		}
		hlog.Printf("warn", "IAM/AK ID %s, SEC %s, Refreshed", ak.Id, ak.Secret[:8])
	} else {

		if rs := Data.NewWriter(iamapi.NsAccessKey(ak.User, ak.Id), nil).SetJsonValue(ak).
			SetCreateOnly(true).Exec(); !rs.OK() {
			return rs.Error()
		}
		hlog.Printf("warn", "IAM/AK ID %s, SEC %s, Refreshed", ak.Id, ak.Secret[:8])
	}

	KeyMgr.KeySet(ak)

	return nil
}

func AccessKeyReset(ak *hauth.AccessKey) error {

	if ak.User == "" {
		return errors.New("No User Set")
	}

	ak.Status = hauth.AccessKeyStatusActive

	if rs := Data.NewWriter(iamapi.NsAccessKey(ak.User, ak.Id), nil).SetJsonValue(ak).
		Exec(); !rs.OK() {
		return rs.Error()
	}

	KeyMgr.KeySet(ak)

	return nil
}

func AppInstanceRegister(inst iamapi.AppInstance) error {

	if !app_inst_id_re.MatchString(inst.Meta.ID) {
		return fmt.Errorf("Invalid meta.id (%s)", inst.Meta.ID)
	}

	if rs := Data.NewReader(iamapi.ObjKeyUser(inst.Meta.User)).Exec(); rs.NotFound() {
		inst.Meta.User = def_sysadmin
	}

	var prev iamapi.AppInstance
	if rs := Data.NewReader(iamapi.ObjKeyAppInstance(inst.Meta.ID)).Exec(); rs.OK() {
		rs.Item().JsonDecode(&prev)
	}

	if prev.Meta.ID == "" {
		inst.Meta.Created = types.MetaTimeNow()
		inst.Meta.Updated = types.MetaTimeNow()

		Data.NewWriter(iamapi.ObjKeyAppInstance(inst.Meta.ID), nil).SetJsonValue(inst).
			SetCreateOnly(true).Exec()

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

		Data.NewWriter(iamapi.ObjKeyAppInstance(prev.Meta.ID), nil).SetJsonValue(prev).
			SetCreateOnly(true).Exec()

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
			iamapi.ObjKeyRolePrivilege(rid, inst.Meta.ID), []byte(strings.Join(v, ","))).
			SetCreateOnly(true).Exec(); !rs.OK() {
			return fmt.Errorf("db err %s", rs.ErrorMessage())
		}
	}

	return nil
}
