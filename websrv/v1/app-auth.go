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

package v1

import (
	"code.hooto.com/lessos/iam/iamapi"
	"code.hooto.com/lessos/iam/store"
	"code.hooto.com/lynkdb/iomix/skv"
	"github.com/lessos/lessgo/crypto/idhash"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/types"
	"github.com/lessos/lessgo/utilx"
)

type AppAuth struct {
	*httpsrv.Controller
}

func (c AppAuth) InfoAction() {

	set := iamapi.AppAuthInfo{}

	defer c.RenderJson(&set)

	instid := c.Params.Get("instance_id")
	if instid == "" {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeNotFound, "App Instance Not Found")
		return
	}

	var inst iamapi.AppInstance
	if obj := store.PvGet("app-instance/" + instid); obj.OK() {
		obj.Decode(&inst)
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

	// if set.Instance.Meta.ID == "" {
	// 	set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Bad Argument")
	// 	return
	// }
	if len(set.AccessToken) < 30 {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "Unauthorized")
		return
	}
	set.AccessToken = set.AccessToken[:8] + "/" + set.AccessToken[9:]

	var session iamapi.UserSession
	if obj := store.PvGet("session/" + set.AccessToken); obj.OK() {
		obj.Decode(&session)
	}

	if !session.IsLogin() {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "Unauthorized")
		return
	}

	if set.Instance.Meta.ID == "" {
		set.Instance.Meta.ID = idhash.RandHexString(16)
	}

	// if !c.Session.AccessAllowed("sys.admin") {
	//        set.Error = types.NewErrorMeta(iamapi.ErrCodeAccessDenied, "Unauthorized")
	// 	return
	// }

	// sess, err := c.Session.SessionFetch()

	var (
		prevVersion uint64
		prev        iamapi.AppInstance
	)

	if obj := store.PvGet("app-instance/" + set.Instance.Meta.ID); obj.OK() {
		obj.Decode(&prev)
		prevVersion = obj.Meta().Version
	}

	if prev.Meta.ID == "" {

		set.Instance.Meta.Created = utilx.TimeNow("datetime")
		set.Instance.Meta.Updated = utilx.TimeNow("datetime")
		set.Instance.Status = 1
		set.Instance.Meta.UserID = session.UserID

	} else {

		if prev.Meta.UserID != session.UserID {
			set.Error = types.NewErrorMeta(iamapi.ErrCodeUnauthorized, "Unauthorized")
			return
		}

		set.Instance.Meta.Created = prev.Meta.Created
		set.Instance.Meta.UserID = prev.Meta.UserID
		set.Instance.Status = prev.Status
	}

	if obj := store.PvPut("app-instance/"+set.Instance.Meta.ID, set.Instance, &skv.PvWriteOptions{
		PrevVersion: prevVersion,
	}); !obj.OK() {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInternalError, obj.Bytex().String())
		return
	}

	//
	// q = base.NewQuerySet().From("iam_privilege").Limit(1000)
	// q.Where.And("instance", req.Data.InstanceId)
	// rs, err = dcn.Base.Query(q)
	// if err != nil {
	// 	rsp.Message = "Internal Server Error"
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
	// 			rsp.Message = "Can not write to database" + err.Error()
	// 			return
	// 		}
	// 	}
	// }

	set.Kind = "AppInstanceRegister"
}
