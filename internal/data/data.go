// Copyright 2014 Eryx <evorui at gmail dot com>, All rights reserved.
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
	"log/slog"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/lessos/lessgo/crypto/phash"
	"github.com/lessos/lessgo/net/email"
	"github.com/lynkdb/kvgo/v2/pkg/kvapi"
	"github.com/lynkdb/kvgo/v2/pkg/kvrep"
	"github.com/lynkdb/kvgo/v2/pkg/storage"
	"github.com/sysinner/incore/v2/pkg/inauth"

	"github.com/hooto/iam/v2/internal/config"
	"github.com/hooto/iam/v2/pkg/iamapi"
)

var (
	Data           kvapi.Client
	app_inst_id_re = regexp.MustCompile("^[0-9a-f]{16}$")
	KeyMgr         = inauth.NewAccessKeyManager()
	mu             sync.RWMutex
	inited         = false
)

func Setup() error {

	mu.Lock()
	defer mu.Unlock()

	if inited {
		return nil
	}

	if Data == nil {
		if config.Config.Database != nil &&
			config.Config.Database.Database != "" {
			c, err := config.Config.Database.NewClient()
			if err != nil {
				return err
			}
			Data = c
			slog.Info("database connected", "type", "remote")
		} else {

			if c, err := kvrep.NewReplica(&storage.Options{
				DataDirectory: config.Prefix + "/var/iam_db",
			}); err != nil {
				return err
			} else {
				Data = c
				slog.Info("database connected", "type", "local-replica")
			}
		}
	}

	slog.Info("iam/data setup")

	if rs := Data.NewWriter([]byte("iam:test"), []byte("test")).
		SetTTL(1000).Exec(); !rs.OK() {
		return fmt.Errorf("iam.data connect not ready #2 " + rs.String())

	} else if rs = Data.NewReader([]byte("iam:test")).Exec(); !rs.OK() || rs.Item().StringValue() != "test" {
		return fmt.Errorf("iam.data connect not ready #3")
	} else {
		slog.Info("iam/data connect ok")
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

			var ak inauth.AccessKey
			if err := v.JsonDecode(&ak); err == nil {
				KeyMgr.Set(&ak)
				slog.Debug("iam/access_key load", "id", ak.Id)
				continue
			}
		}

		akn += len(rs.Items)

		if len(rs.Items) < 1000 {
			break
		}
	}
	slog.Info("iam/access_key load from database", "count", akn)

	//
	for _, v := range config.Config.AccessKeys {
		KeyMgr.Set(v)
	}

	slog.Info("iam/access_key load from config", "count", len(config.Config.AccessKeys))

	if err := loadAppKeys(); err != nil {
		return err
	}

	if err := initData(); err != nil {
		return err
	}

	inited = true

	return nil
}

func loadAppKeys() error {
	var (
		offset = iamapi.NsAppInstance("")
		cutset = append(offset, []byte{0xff}...)
		num    = 0
	)
	for {

		rs := Data.NewRanger(offset, cutset).SetLimit(1000).Exec()
		if !rs.OK() {
			break
		}

		for _, v := range rs.Items {

			var app iamapi.AppInstance
			if err := v.JsonDecode(&app); err != nil {
				slog.Debug("iam/app_instance load", "id", app.ID)
				continue
			}

			ak := &inauth.AccessKey{
				Id:     app.ID,
				User:   app.User,
				Secret: app.SecretKey,
				Type:   "App",
				State:  inauth.AccessKey_State_Active,
			}

			KeyMgr.Set(ak)
		}

		num += len(rs.Items)

		if len(rs.Items) < 1000 {
			break
		}
	}
	slog.Info("iam/app_instance load from database", "count", num)

	return nil
}

func initData() (err error) {

	if len(config.Config.InstanceID) < 16 {
		return fmt.Errorf("No InstanceID Setup")
	}

	tn := time.Now().Unix()

	for _, v := range []iamapi.UserRole{
		{
			Name: iamapi.Role_Sysadmin,
			Desc: "Root System Administrator",
		},
		{
			Name: iamapi.Role_User,
			Desc: "Universal Member",
		},
		{
			Name: iamapi.Role_Developer,
			Desc: "Universal Developer",
		},
		{
			Name: iamapi.Role_Guest,
			Desc: "Anonymous Guest",
		},
	} {

		v.User = iamapi.UserSysadmin
		v.Status = 1
		v.Created = tn
		v.Updated = tn

		rw := Data.NewWriter(iamapi.NsRole(v.Name), nil).SetJsonValue(v).
			SetCreateOnly(true)
		if rs := rw.Exec(); !rs.OK() {
			return fmt.Errorf("db err init-role %s", rs.ErrorMessage())
		}
	}

	//
	ps := []*iamapi.AppPermission{
		{
			Permission: "sys.admin",
			Roles:      []string{iamapi.Role_Sysadmin},
			Summary:    "System Management",
		},
		{
			Permission: "user.admin",
			Roles:      []string{iamapi.Role_Sysadmin},
			Summary:    "User Management",
		},
	}

	secretKey, _ := inauth.GenerateSecretKeyBase62(40)

	inst := iamapi.AppInstance{
		ID:          config.Config.InstanceID,
		Name:        "IAM Service",
		User:        iamapi.UserSysadmin,
		Created:     tn,
		Updated:     tn,
		Version:     config.Version,
		Status:      1,
		Url:         "",
		Permissions: ps,
		SecretKey:   secretKey,
	}

	if rs := Data.NewWriter(iamapi.NsAppInstance(inst.ID), nil).SetJsonValue(inst).
		SetCreateOnly(true).Exec(); !rs.OK() {
		return fmt.Errorf("db err init-data %s", rs.ErrorMessage())
	}

	AppInstanceRegister(inst)

	// Init Super SysAdmin
	// uid := idhash.HashToHexString([]byte(iamapi.UserSysadmin), 8)
	if rs := Data.NewReader(iamapi.NsUser(iamapi.UserSysadmin)).Exec(); rs.NotFound() {

		sysadm := iamapi.User{
			Name:        iamapi.UserSysadmin,
			DisplayName: "System Administrator",
			Email:       "",
			Roles:       []string{iamapi.Role_Sysadmin, iamapi.Role_User},
			Status:      1,
			Created:     tn,
			Updated:     tn,
		}

		auth, err := phash.Generate(iamapi.DefaultPassword)
		if err != nil {
			return err
		}

		sysadm.Keys.Set(iamapi.UserKeyDefault, auth)

		ow := Data.NewWriter(iamapi.NsUser(iamapi.UserSysadmin), nil).SetJsonValue(sysadm).
			SetCreateOnly(true).SetIncr(1, "user")
		if rs := ow.Exec(); !rs.OK() {
			return fmt.Errorf("db err init-sysadmin %s", rs.ErrorMessage())
		}
	}

	return nil
}

func SysConfigRefresh() {

	//
	if rs := Data.NewReader(iamapi.NsSysConfig("mailer")).Exec(); rs.OK() {

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

	if rs := Data.NewReader(iamapi.NsSysConfig("service_name")).Exec(); rs.OK() &&
		rs.Item().StringValue() != "" {
		config.Config.ServiceName = rs.Item().StringValue()
	}

	if rs := Data.NewReader(iamapi.NsSysConfig("webui_banner_title")).Exec(); rs.OK() &&
		rs.Item().StringValue() != "" {
		config.Config.WebUiBannerTitle = rs.Item().StringValue()
	}

	if rs := Data.NewReader(iamapi.NsSysConfig("service_login_form_alert_msg")).Exec(); rs.OK() {
		config.Config.ServiceLoginFormAlertMsg = rs.Item().StringValue()
	}

	if rs := Data.NewReader(iamapi.NsSysConfig("allow_user_sign_up")).Exec(); rs.OK() {

		if rs.Item().StringValue() == "1" {
			config.AllowUserSignUp = true
		} else {
			config.AllowUserSignUp = false
		}
	}
}

func AccessKeyInitData(ak *inauth.AccessKey) error {

	if ak.User == "" {
		ak.User = "sysadmin"
	}

	ak.State = inauth.AccessKey_State_Active

	if rs := Data.NewReader(iamapi.NsAccessKey(ak.User, ak.Id)).Exec(); rs.OK() {
		var prev inauth.AccessKey
		if err := rs.Item().JsonDecode(&prev); err != nil {
			return err
		}
		if !inauth.ProtoEqual(ak, &prev) {
			if rs := Data.NewWriter(iamapi.NsAccessKey(ak.User, ak.Id), nil).SetJsonValue(ak).
				Exec(); !rs.OK() {
				return rs.Error()
			}
			slog.Warn("IAM/AK Refreshed", "id", ak.Id, "secret", ak.Secret[:8])
		}
		slog.Warn("IAM/AK Refreshed", "id", ak.Id, "secret", ak.Secret[:8])
	} else {

		if rs := Data.NewWriter(iamapi.NsAccessKey(ak.User, ak.Id), nil).SetJsonValue(ak).
			SetCreateOnly(true).Exec(); !rs.OK() {
			return rs.Error()
		}
		slog.Warn("IAM/AK Refreshed", "id", ak.Id, "secret", ak.Secret[:8])
	}

	KeyMgr.Set(ak)

	return nil
}

func AccessKeyReset(ak *inauth.AccessKey) error {

	if ak.User == "" {
		return errors.New("No User Set")
	}

	ak.State = inauth.AccessKey_State_Active

	if rs := Data.NewWriter(iamapi.NsAccessKey(ak.User, ak.Id), nil).SetJsonValue(ak).
		Exec(); !rs.OK() {
		return rs.Error()
	}

	KeyMgr.Set(ak)

	return nil
}

func AppInstanceRegister(inst iamapi.AppInstance) error {

	if !app_inst_id_re.MatchString(inst.ID) {
		return fmt.Errorf("Invalid meta.id (%s)", inst.ID)
	}

	if rs := Data.NewReader(iamapi.NsUser(inst.User)).Exec(); rs.NotFound() {
		inst.User = iamapi.UserSysadmin
	}

	var prev iamapi.AppInstance
	if rs := Data.NewReader(iamapi.NsAppInstance(inst.ID)).Exec(); rs.OK() {
		rs.Item().JsonDecode(&prev)
	}

	tn := time.Now().Unix()

	if prev.ID == "" {
		inst.Created = tn
		inst.Updated = tn

		Data.NewWriter(iamapi.NsAppInstance(inst.ID), nil).SetJsonValue(inst).
			SetCreateOnly(true).Exec()

	} else {

		prev.Updated = tn

		if inst.Version != "" && inst.Version != prev.Version {
			prev.Version = inst.Version
		}

		// if inst.AppTitle != "" && inst.AppTitle != prev.AppTitle {
		// 	prev.AppTitle = inst.AppTitle
		// }

		if inst.Status != prev.Status {
			prev.Status = inst.Status
		}

		if inst.Url != "" && inst.Url != prev.Url {
			prev.Url = inst.Url
		}

		if inst.SecretKey != "" && inst.SecretKey != prev.SecretKey {
			prev.SecretKey = inst.SecretKey
		}

		Data.NewWriter(iamapi.NsAppInstance(prev.ID), nil).SetJsonValue(prev).
			SetCreateOnly(true).Exec()

		// TODO remove unused privileges
	}

	// privilege
	rps := map[string][]string{}
	for _, v := range inst.Permissions {

		for _, rid := range v.Roles {

			if _, ok := rps[rid]; !ok {
				rps[rid] = []string{}
			}

			rps[rid] = append(rps[rid], v.Permission)
		}
	}

	for rid, v := range rps {

		if rs := Data.NewWriter(
			iamapi.NsRolePrivilege(rid, inst.ID), []byte(strings.Join(v, ","))).
			SetCreateOnly(true).Exec(); !rs.OK() {
			return fmt.Errorf("db err %s", rs.ErrorMessage())
		}
	}

	return nil
}
