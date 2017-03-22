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
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"strings"

	"code.hooto.com/lynkdb/iomix/skv"
	"github.com/eryx/imaging"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/pass"
	"github.com/lessos/lessgo/types"
	"github.com/lessos/lessgo/utilx"

	"code.hooto.com/lessos/iam/base/login"
	"code.hooto.com/lessos/iam/base/profile"
	"code.hooto.com/lessos/iam/iamapi"
	"code.hooto.com/lessos/iam/iamclient"
	"code.hooto.com/lessos/iam/store"
)

type User struct {
	*httpsrv.Controller
}

func (c User) ProfileAction() {

	rsp := iamapi.UserProfile{}

	defer c.RenderJson(&rsp)

	session, err := iamclient.SessionInstance(c.Session)

	if err != nil || !session.IsLogin() {
		rsp.Error = &types.ErrorMeta{"401", "Access Denied"}
		return
	}

	// profile
	if obj := store.PvGet("user-profile/" + session.UserID); obj.OK() {
		obj.Decode(&rsp)
	}

	if rsp.Name == "" {

		rsp.Name = session.Name

		store.PvPut("user-profile/"+session.UserID, rsp, nil)
	}

	// login
	var user iamapi.User
	if obj := store.PvGet("user/" + session.UserID); obj.OK() {
		obj.Decode(&user)
	}

	if user.Meta.ID != session.UserID {
		rsp.Error = &types.ErrorMeta{"401", "Access Denied"}
		return
	}

	rsp.Login.Meta = user.Meta
	rsp.Login.Name = user.Name
	rsp.Login.Email = user.Email
	rsp.Name = user.Name

	rsp.Kind = "UserProfile"
}

func (c User) ProfileSetAction() {

	var (
		rsp types.TypeMeta
		req iamapi.UserProfile
		err error
	)

	defer c.RenderJson(&rsp)

	if err := c.Request.JsonDecode(&req); err != nil {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "Bad Request"}
		return
	}

	if req, err = profile.PutValidate(req); err != nil {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, err.Error()}
		return
	}

	session, err := iamclient.SessionInstance(c.Session)

	if err != nil || !session.IsLogin() {
		rsp.Error = &types.ErrorMeta{"401", "Access Denied"}
		return
	}

	// login
	var user iamapi.User
	uobj := store.PvGet("user/" + session.UserID)
	if uobj.OK() {
		uobj.Decode(&user)
	}

	if user.Meta.ID != session.UserID {
		rsp.Error = &types.ErrorMeta{"404", "User Not Found"}
		return
	}
	user.Name = req.Name

	store.PvPut("user/"+session.UserID, user, &skv.PvWriteOptions{
		PrevVersion: uobj.Meta().Version,
	})

	// profile
	var profile iamapi.UserProfile
	pobj := store.PvGet("user-profile/" + session.UserID)
	if pobj.OK() {
		pobj.Decode(&profile)
	}

	profile.Name = req.Name
	profile.Birthday = req.Birthday
	profile.About = req.About

	store.PvPut("user-profile/"+session.UserID, profile, &skv.PvWriteOptions{
		PrevVersion: pobj.Meta().Version,
	})

	rsp.Kind = "UserProfile"
}

func (c User) PassSetAction() {

	var (
		rsp types.TypeMeta
		req iamapi.UserPasswordSet
	)

	defer c.RenderJson(&rsp)

	if err := c.Request.JsonDecode(&req); err != nil {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "Bad Request"}
		return
	}

	if err := login.PassSetValidate(req); err != nil {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, err.Error()}
		return
	}

	session, err := iamclient.SessionInstance(c.Session)

	if err != nil || !session.IsLogin() {
		rsp.Error = &types.ErrorMeta{"401", "Access Denied"}
		return
	}

	var user iamapi.User
	uobj := store.PvGet("user/" + session.UserID)
	if uobj.OK() {
		uobj.Decode(&user)
	}

	if user.Meta.ID != session.UserID {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "User Not Found"}
		return
	}

	if !pass.Check(req.CurrentPassword, user.Auth) {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "Current Password can not match"}
		return
	}

	user.Meta.Updated = utilx.TimeNow("atom")
	user.Auth, _ = pass.HashDefault(req.NewPassword)

	store.PvPut("user/"+user.Meta.ID, user, &skv.PvWriteOptions{
		PrevVersion: uobj.Meta().Version,
	})

	rsp.Kind = "UserPassword"
}

func (c User) EmailSetAction() {

	var (
		rsp types.TypeMeta
		req iamapi.UserEmailSet
	)

	defer c.RenderJson(&rsp)

	if err := c.Request.JsonDecode(&req); err != nil {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "Bad Request"}
		return
	}

	if email, err := login.EmailSetValidate(req.Email); err != nil {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, err.Error()}
		return
	} else {
		req.Email = email
	}

	session, err := iamclient.SessionInstance(c.Session)

	if err != nil || !session.IsLogin() {
		rsp.Error = &types.ErrorMeta{"401", "Access Denied"}
		return
	}

	var user iamapi.User
	uobj := store.PvGet("user/" + session.UserID)
	if uobj.OK() {
		uobj.Decode(&user)
	}

	if user.Meta.ID != session.UserID {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "User Not Found"}
		return
	}

	if !pass.Check(req.Auth, user.Auth) {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "Password can not match"}
		return
	}

	user.Email = req.Email
	user.Meta.Updated = utilx.TimeNow("atom")

	store.PvPut("user/"+user.Meta.ID, user, &skv.PvWriteOptions{
		PrevVersion: uobj.Meta().Version,
	})

	rsp.Kind = "UserEmail"
}

func (c User) PhotoSetAction() {

	var (
		rsp types.TypeMeta
		req iamapi.UserPhotoSet
	)

	defer c.RenderJson(&rsp)

	if err := c.Request.JsonDecode(&req); err != nil {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "Bad Request"}
		return
	}

	session, err := iamclient.SessionInstance(c.Session)

	if err != nil || !session.IsLogin() {
		rsp.Error = &types.ErrorMeta{"401", "Access Denied"}
		return
	}

	//
	img64 := strings.SplitAfter(req.Data, ";base64,")
	if len(img64) != 2 {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, "Bad Request"}
		return
	}
	imgreader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(img64[1]))
	imgsrc, _, err := image.Decode(imgreader)
	if err != nil {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, err.Error()}
		return
	}
	imgnew := imaging.Thumbnail(imgsrc, 96, 96, imaging.CatmullRom)

	var imgbuf bytes.Buffer
	err = png.Encode(&imgbuf, imgnew)
	if err != nil {
		rsp.Error = &types.ErrorMeta{iamapi.ErrCodeInvalidArgument, err.Error()}
		return
	}
	imgphoto := base64.StdEncoding.EncodeToString(imgbuf.Bytes())

	// profile
	var profile iamapi.UserProfile
	pobj := store.PvGet("user-profile/" + session.UserID)
	if pobj.OK() {
		pobj.Decode(&profile)
	}

	profile.Photo = "data:image/png;base64," + imgphoto
	profile.PhotoSource = req.Data

	store.PvPut("user-profile/"+session.UserID, profile, &skv.PvWriteOptions{
		PrevVersion: pobj.Meta().Version,
	})

	rsp.Kind = "UserPhoto"
}

func (c User) RoleListAction() {

	ls := iamapi.UserRoleList{}

	defer c.RenderJson(&ls)

	session, err := iamclient.SessionInstance(c.Session)

	if err != nil || !session.IsLogin() {
		ls.Error = &types.ErrorMeta{iamapi.ErrCodeUnauthorized, "Access Denied"}
		return
	}

	// TODO page
	if objs := store.PvScan("role/", "", "", 1000); objs.OK() {

		rss := objs.KvList()
		for _, obj := range rss {

			var role iamapi.UserRole
			if err := obj.Decode(&role); err == nil {

				if role.IdxID <= 1000 || role.Meta.UserID == session.UserID {
					ls.Items = append(ls.Items, role)
				}
			}
		}
	}

	ls.Kind = "UserRoleList"
}
