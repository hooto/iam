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

package v1

import (
	"strings"
	"time"

	"github.com/hooto/hlog4g/hlog"
	"github.com/hooto/httpsrv"
	"github.com/lessos/lessgo/crypto/idhash"
	"github.com/lessos/lessgo/types"
	iox_utils "github.com/lynkdb/iomix/utils"

	"github.com/hooto/hauth/go/hauth/v1"
	"github.com/hooto/iam/data"
	"github.com/hooto/iam/iamapi"
	"github.com/hooto/iam/iamclient"
)

type AppAuth struct {
	*httpsrv.Controller
}

func (c AppAuth) InfoAction() {

	set := iamapi.AppAuthInfo{}

	defer c.RenderJson(&set)

	instid := c.Params.Value("instance_id")
	if instid == "" {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeNotFound, "App Instance Not Found")
		return
	}

	var inst iamapi.AppInstance
	if obj := data.Data.NewReader(iamapi.ObjKeyAppInstance(instid)).Exec(); obj.OK() {
		obj.Item().JsonDecode(&inst)
	}

	if inst.Meta.ID == instid {

		set.InstanceID = instid
		set.AppID = inst.AppID
		// set.Version = inst.Version

		set.Kind = "AppAuthInfo"

	} else {

		set.Error = types.NewErrorMeta(iamapi.ErrCodeNotFound, "App Instance Not Found")
	}
}

func (c AppAuth) RegisterAction() {

	set := iamapi.AppInstanceRegister{}
	defer c.RenderJson(&set)

	if err := c.Request.JsonDecode(&set); err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Bad Argument")
		return
	}

	tn := uint32(time.Now().Unix())
	if len(set.Instance.Meta.ID) > 0 {

		if len(set.Instance.Meta.ID) < 16 || !iamapi.AppIdRE.MatchString(set.Instance.Meta.ID) {
			set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Invalid Instance ID "+set.Instance.Meta.ID)
			return
		}

		if !strings.HasPrefix(set.Instance.Meta.ID, "00") {
			time_seq := iox_utils.BytesToUint32(iox_utils.HexStringToBytes(set.Instance.Meta.ID[:8]))
			if time_seq < (tn-31104000) || time_seq > (tn+864000) {
				set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Invalid Instance ID (Prefix Error)")
				return
			}
		}
	}

	ap, err := hauth.NewUserValidator(set.AccessToken, data.KeyMgr)
	if err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, err.Error())
		return
	}

	if err := ap.SignValid(); err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "Unauthorized")
		return
	}

	if set.Instance.Meta.ID == "" {
		set.Instance.Meta.ID = iox_utils.Uint32ToHexString(tn) + idhash.RandHexString(8)
	}

	// if !c.Session.AccessAllowed("sys.admin") {
	//        set.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Unauthorized")
	// 	return
	// }

	// sess, err := c.Session.SessionFetch()

	var (
		prev iamapi.AppInstance
	)

	if obj := data.Data.NewReader(iamapi.ObjKeyAppInstance(set.Instance.Meta.ID)).Exec(); obj.OK() {
		obj.Item().JsonDecode(&prev)
	}

	if prev.Meta.ID == "" {

		set.Instance.Meta.Created = types.MetaTimeNow()
		set.Instance.Meta.Updated = types.MetaTimeNow()
		set.Instance.Status = 1
		set.Instance.Meta.User = ap.Id

	} else {

		if prev.Meta.User != ap.Id {
			set.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "Unauthorized")
			return
		}

		set.Instance.Meta.Created = prev.Meta.Created
		set.Instance.Meta.User = prev.Meta.User
		set.Instance.Status = prev.Status
	}

	if obj := data.Data.NewWriter(iamapi.ObjKeyAppInstance(set.Instance.Meta.ID), nil).SetJsonValue(set.Instance).
		Exec(); !obj.OK() {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInternalError, obj.ErrorMessage())
		return
	}

	//
	// q = base.NewQuerySet().From("iam_privilege").Limit(1000)
	// q.Where.And("instance", req.Data.InstanceId)
	// rs, err = dcn.Base.Query(q)
	// if err != nil {
	// 	rsp.ErrorMessage() = "Internal Server Error"
	// 	return
	// }

	// for _, prePriv := range rs {

	// 	isExist := false
	// 	for _, curPrev := range req.Data.Privileges {

	// 		if prePriv.Field("privilege").String() == curPrev.Key {
	// 			isExist = true
	// 			break
	// 		}
	// 	}

	// 	if !isExist {
	// 		frupd := base.NewFilter()
	// 		frupd.And("instance", req.Data.InstanceId).And("privilege", prePriv.Field("privilege").String())
	// 		dcn.Base.Delete("iam_privilege", frupd)
	// 	}
	// }

	// for _, curPrev := range req.Data.Privileges {

	// 	isExist := false

	// 	for _, prePriv := range rs {

	// 		if prePriv.Field("privilege").String() == curPrev.Key {
	// 			isExist = true
	// 			break
	// 		}
	// 	}

	// 	if !isExist {
	// 		item := map[string]interface{}{
	// 			"instance":  req.Data.InstanceId,
	// 			"uid":       sess.UserID,
	// 			"privilege": curPrev.Key,
	// 			"desc":      curPrev.Desc,
	// 			"created":   base.TimeNow("datetime"),
	// 		}

	// 		if _, err := dcn.Base.Insert("iam_privilege", item); err != nil {
	// 			rsp.Status = 500
	// 			rsp.ErrorMessage() = "Can not write to database" + err.Error()
	// 			return
	// 		}
	// 	}
	// }

	set.Kind = "AppInstanceRegister"
}

func (c AppAuth) RoleListAction() {

	sets := iamapi.UserRoleList{}
	defer c.RenderJson(&sets)

	// TODO app<->role
	if rs := data.Data.NewRanger(
		iamapi.ObjKeyRole(""), iamapi.ObjKeyRole("")).SetLimit(100).Exec(); rs.OK() {

		for _, obj := range rs.Items {

			var role iamapi.UserRole
			if err := obj.JsonDecode(&role); err == nil {

				if role.Status == 0 || obj.Meta.IncrId == 1 {
					continue
				}

				sets.Items = append(sets.Items, iamapi.UserRole{
					Id:   uint32(obj.Meta.IncrId),
					Name: role.Name,
					Desc: role.Desc,
				})
			}
		}
	}

	sets.Kind = "UserRoleList"
}

func (c AppAuth) UserAccessKeyAction() {

	var set struct {
		types.TypeMeta
		iamapi.AccessKeySession
	}
	defer c.RenderJson(&set)

	var (
		username   = c.Params.Value("user")
		access_key = c.Params.Value("access_key")
	)
	if username == "" || access_key == "" {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Bad Argument")
		return
	}

	app_auth := c.Request.Header.Get("Auth")
	if app_auth == "" {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "Unauthorized")
		return
	}

	app_aka, err := iamapi.AccessKeyAuthDecode(app_auth)
	if err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "Unauthorized")
		return
	}

	if err := app_aka.Valid(); err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, err.Error())
		return
	}

	var app iamapi.AppInstance
	if rs := data.Data.NewReader(iamapi.ObjKeyAppInstance(app_aka.Key)).Exec(); rs.OK() {
		rs.Item().JsonDecode(&app)
	}

	if app.Meta.ID != app_aka.Key {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Bad Argument")
		return
	}

	if err := iamclient.AccessKeyAuthValid(app_aka, app.SecretKey); err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "Unauthorized")
		return
	}

	var user_ak hauth.AccessKey
	if rs := data.Data.NewReader(iamapi.NsAccessKey(username, access_key)).Exec(); rs.OK() {
		rs.Item().JsonDecode(&user_ak)
	}

	if user_ak.Id != access_key ||
		user_ak.Status != hauth.AccessKeyStatusActive {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Access Key Not Found")
		return
	}

	if err := user_ak.ScopeAllow(hauth.NewScopeFilter("app", app_aka.Key)); err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, err.Error())
		return
	}

	var user iamapi.User
	if obj := data.Data.NewReader(iamapi.ObjKeyUser(username)).Exec(); obj.OK() {
		obj.Item().JsonDecode(&user)
	}

	set.Kind = "AccessKeySession"
	set.AccessKeySession = iamapi.AccessKeySession{
		User:      username,
		AccessKey: user_ak.Id,
		SecretKey: user_ak.Secret,
		Roles:     user.Roles,
		Expired:   time.Now().Unix() + 864000,
	}
	hlog.Printf("info", "app-auth AccessKeySession")
}
