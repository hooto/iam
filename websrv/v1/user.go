// Copyright 2015 lessOS.com, All rights reserved.
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

	"github.com/eryx/imaging"
	"github.com/lessos/bigtree/btapi"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/pass"
	"github.com/lessos/lessgo/types"
	"github.com/lessos/lessgo/utilx"

	"github.com/lessos/lessids/base/login"
	"github.com/lessos/lessids/base/profile"
	"github.com/lessos/lessids/idclient"
	"github.com/lessos/lessids/idsapi"
	"github.com/lessos/lessids/store"
)

type User struct {
	*httpsrv.Controller
}

func (c User) ProfileAction() {

	rsp := idsapi.UserProfile{}

	defer c.RenderJson(&rsp)

	session, err := idclient.SessionInstance(c.Session)

	if err != nil || !session.IsLogin() {
		rsp.Error = &types.ErrorMeta{"401", "Access Denied"}
		return
	}

	// profile
	if obj := store.BtAgent.ObjectGet("/global/ids/user-profile/" + session.UserID); obj.Error == nil {
		obj.JsonDecode(&rsp)
	}

	if rsp.Name == "" {

		rsp.Name = session.Name

		store.BtAgent.ObjectSet("/global/ids/user-profile/"+session.UserID, rsp, nil)
	}

	// login
	var user idsapi.User
	if obj := store.BtAgent.ObjectGet("/global/ids/user/" + session.UserID); obj.Error == nil {
		obj.JsonDecode(&user)
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
		req idsapi.UserProfile
		err error
	)

	defer c.RenderJson(&rsp)

	if err := c.Request.JsonDecode(&req); err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "Bad Request"}
		return
	}

	if req, err = profile.PutValidate(req); err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, err.Error()}
		return
	}

	session, err := idclient.SessionInstance(c.Session)

	if err != nil || !session.IsLogin() {
		rsp.Error = &types.ErrorMeta{"401", "Access Denied"}
		return
	}

	// login
	var user idsapi.User
	uobj := store.BtAgent.ObjectGet("/global/ids/user/" + session.UserID)
	if uobj.Error == nil {
		uobj.JsonDecode(&user)
	}

	if user.Meta.ID != session.UserID {
		rsp.Error = &types.ErrorMeta{"404", "User Not Found"}
		return
	}
	user.Name = req.Name

	store.BtAgent.ObjectSet("/global/ids/user/"+session.UserID, user, &btapi.ObjectWriteOptions{
		PrevVersion: uobj.Meta.Version,
	})

	// profile
	var profile idsapi.UserProfile
	pobj := store.BtAgent.ObjectGet("/global/ids/user-profile/" + session.UserID)
	if pobj.Error == nil {
		pobj.JsonDecode(&profile)
	}

	profile.Name = req.Name
	profile.Birthday = req.Birthday
	profile.About = req.About

	store.BtAgent.ObjectSet("/global/ids/user-profile/"+session.UserID, profile, &btapi.ObjectWriteOptions{
		PrevVersion: pobj.Meta.Version,
	})

	rsp.Kind = "UserProfile"
}

func (c User) PassSetAction() {

	var (
		rsp types.TypeMeta
		req idsapi.UserPasswordSet
	)

	defer c.RenderJson(&rsp)

	if err := c.Request.JsonDecode(&req); err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "Bad Request"}
		return
	}

	if err := login.PassSetValidate(req); err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, err.Error()}
		return
	}

	session, err := idclient.SessionInstance(c.Session)

	if err != nil || !session.IsLogin() {
		rsp.Error = &types.ErrorMeta{"401", "Access Denied"}
		return
	}

	var user idsapi.User
	uobj := store.BtAgent.ObjectGet("/global/ids/user/" + session.UserID)
	if uobj.Error == nil {
		uobj.JsonDecode(&user)
	}

	if user.Meta.ID != session.UserID {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "User Not Found"}
		return
	}

	if !pass.Check(req.CurrentPassword, user.Auth) {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "Current Password can not match"}
		return
	}

	user.Meta.Updated = utilx.TimeNow("atom")
	user.Auth, _ = pass.HashDefault(req.NewPassword)

	store.BtAgent.ObjectSet("/global/ids/user/"+user.Meta.ID, user, &btapi.ObjectWriteOptions{
		PrevVersion: uobj.Meta.Version,
	})

	rsp.Kind = "UserPassword"
}

func (c User) EmailSetAction() {

	var (
		rsp types.TypeMeta
		req idsapi.UserEmailSet
	)

	defer c.RenderJson(&rsp)

	if err := c.Request.JsonDecode(&req); err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "Bad Request"}
		return
	}

	if email, err := login.EmailSetValidate(req.Email); err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, err.Error()}
		return
	} else {
		req.Email = email
	}

	session, err := idclient.SessionInstance(c.Session)

	if err != nil || !session.IsLogin() {
		rsp.Error = &types.ErrorMeta{"401", "Access Denied"}
		return
	}

	var user idsapi.User
	uobj := store.BtAgent.ObjectGet("/global/ids/user/" + session.UserID)
	if uobj.Error == nil {
		uobj.JsonDecode(&user)
	}

	if user.Meta.ID != session.UserID {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "User Not Found"}
		return
	}

	if !pass.Check(req.Auth, user.Auth) {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "Password can not match"}
		return
	}

	user.Email = req.Email
	user.Meta.Updated = utilx.TimeNow("atom")

	store.BtAgent.ObjectSet("/global/ids/user/"+user.Meta.ID, user, &btapi.ObjectWriteOptions{
		PrevVersion: uobj.Meta.Version,
	})

	rsp.Kind = "UserEmail"
}

func (c User) PhotoSetAction() {

	var (
		rsp types.TypeMeta
		req idsapi.UserPhotoSet
	)

	defer c.RenderJson(&rsp)

	if err := c.Request.JsonDecode(&req); err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "Bad Request"}
		return
	}

	session, err := idclient.SessionInstance(c.Session)

	if err != nil || !session.IsLogin() {
		rsp.Error = &types.ErrorMeta{"401", "Access Denied"}
		return
	}

	//
	img64 := strings.SplitAfter(req.Data, ";base64,")
	if len(img64) != 2 {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "Bad Request"}
		return
	}
	imgreader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(img64[1]))
	imgsrc, _, err := image.Decode(imgreader)
	if err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, err.Error()}
		return
	}
	imgnew := imaging.Thumbnail(imgsrc, 96, 96, imaging.CatmullRom)

	var imgbuf bytes.Buffer
	err = png.Encode(&imgbuf, imgnew)
	if err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, err.Error()}
		return
	}
	imgphoto := base64.StdEncoding.EncodeToString(imgbuf.Bytes())

	// profile
	var profile idsapi.UserProfile
	pobj := store.BtAgent.ObjectGet("/global/ids/user-profile/" + session.UserID)
	if pobj.Error == nil {
		pobj.JsonDecode(&profile)
	}

	profile.Photo = "data:image/png;base64," + imgphoto
	profile.PhotoSource = req.Data

	store.BtAgent.ObjectSet("/global/ids/user-profile/"+session.UserID, profile, &btapi.ObjectWriteOptions{
		PrevVersion: pobj.Meta.Version,
	})

	rsp.Kind = "UserPhoto"
}

func (c User) RoleListAction() {

	ls := idsapi.UserRoleList{}

	defer c.RenderJson(&ls)

	session, err := idclient.SessionInstance(c.Session)

	if err != nil || !session.IsLogin() {
		ls.Error = &types.ErrorMeta{idsapi.ErrCodeUnauthorized, "Access Denied"}
		return
	}

	if objs := store.BtAgent.ObjectList("/global/ids/role/"); objs.Error == nil {

		for _, obj := range objs.Items {

			var role idsapi.UserRole
			if err := obj.JsonDecode(&role); err == nil {

				if role.IdxID <= 1000 || role.Meta.UserID == session.UserID {
					ls.Items = append(ls.Items, role)
				}
			}
		}
	}

	ls.Kind = "UserRoleList"
}
