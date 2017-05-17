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
	us iamapi.UserSession
}

func (c *User) Init() int {

	//
	c.us, _ = iamclient.SessionInstance(c.Session)

	if !c.us.IsLogin() {
		c.Response.Out.WriteHeader(401)
		c.RenderJson(types.NewTypeErrorMeta(iamapi.ErrCodeUnauthorized, "Unauthorized"))
		return 1
	}

	return 0
}

func (c User) ProfileAction() {

	set := iamapi.UserProfile{}

	defer c.RenderJson(&set)

	// profile
	if obj := store.PvGet("user-profile/" + c.us.UserID); obj.OK() {
		obj.Decode(&set)
	}

	if set.Name == "" {

		set.Name = c.us.Name

		store.PvPut("user-profile/"+c.us.UserID, set, nil)
	}

	// login
	var user iamapi.User
	if obj := store.PvGet("user/" + c.us.UserID); obj.OK() {
		obj.Decode(&user)
	}

	if user.Meta.ID != c.us.UserID {
		set.Error = types.NewErrorMeta("401", "Access Denied")
		return
	}

	set.Login.Meta = user.Meta
	set.Login.Name = user.Name
	set.Login.Email = user.Email
	set.Name = user.Name

	set.Kind = "UserProfile"
}

func (c User) ProfileSetAction() {

	var (
		set types.TypeMeta
		req iamapi.UserProfile
		err error
	)

	defer c.RenderJson(&set)

	if err := c.Request.JsonDecode(&req); err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Bad Request")
		return
	}

	if req, err = profile.PutValidate(req); err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, err.Error())
		return
	}

	// login
	var user iamapi.User
	uobj := store.PvGet("user/" + c.us.UserID)
	if uobj.OK() {
		uobj.Decode(&user)
	}

	if user.Meta.ID != c.us.UserID {
		set.Error = types.NewErrorMeta("404", "User Not Found")
		return
	}
	user.Name = req.Name

	store.PvPut("user/"+c.us.UserID, user, &skv.PvWriteOptions{
		PrevVersion: uobj.Meta().Version,
	})

	// profile
	var profile iamapi.UserProfile
	pobj := store.PvGet("user-profile/" + c.us.UserID)
	if pobj.OK() {
		pobj.Decode(&profile)
	}

	profile.Name = req.Name
	profile.Birthday = req.Birthday
	profile.About = req.About

	store.PvPut("user-profile/"+c.us.UserID, profile, &skv.PvWriteOptions{
		PrevVersion: pobj.Meta().Version,
	})

	set.Kind = "UserProfile"
}

func (c User) PassSetAction() {

	var (
		set types.TypeMeta
		req iamapi.UserPasswordSet
	)

	defer c.RenderJson(&set)

	if err := c.Request.JsonDecode(&req); err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Bad Request")
		return
	}

	if err := login.PassSetValidate(req); err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, err.Error())
		return
	}

	var user iamapi.User
	uobj := store.PvGet("user/" + c.us.UserID)
	if uobj.OK() {
		uobj.Decode(&user)
	}

	if user.Meta.ID != c.us.UserID {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "User Not Found")
		return
	}

	if !pass.Check(req.CurrentPassword, user.Auth) {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Current Password can not match")
		return
	}

	user.Meta.Updated = utilx.TimeNow("atom")
	user.Auth, _ = pass.HashDefault(req.NewPassword)

	store.PvPut("user/"+user.Meta.ID, user, &skv.PvWriteOptions{
		PrevVersion: uobj.Meta().Version,
	})

	set.Kind = "UserPassword"
}

func (c User) EmailSetAction() {

	var (
		set types.TypeMeta
		req iamapi.UserEmailSet
	)

	defer c.RenderJson(&set)

	if err := c.Request.JsonDecode(&req); err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Bad Request")
		return
	}

	if email, err := login.EmailSetValidate(req.Email); err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, err.Error())
		return
	} else {
		req.Email = email
	}

	var user iamapi.User
	uobj := store.PvGet("user/" + c.us.UserID)
	if uobj.OK() {
		uobj.Decode(&user)
	}

	if user.Meta.ID != c.us.UserID {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "User Not Found")
		return
	}

	if !pass.Check(req.Auth, user.Auth) {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Password can not match")
		return
	}

	user.Email = req.Email
	user.Meta.Updated = utilx.TimeNow("atom")

	store.PvPut("user/"+user.Meta.ID, user, &skv.PvWriteOptions{
		PrevVersion: uobj.Meta().Version,
	})

	set.Kind = "UserEmail"
}

func (c User) PhotoSetAction() {

	var (
		set types.TypeMeta
		req iamapi.UserPhotoSet
	)

	defer c.RenderJson(&set)

	if err := c.Request.JsonDecode(&req); err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Bad Request")
		return
	}

	//
	img64 := strings.SplitAfter(req.Data, ";base64,")
	if len(img64) != 2 {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Bad Request")
		return
	}
	imgreader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(img64[1]))
	imgsrc, _, err := image.Decode(imgreader)
	if err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, err.Error())
		return
	}
	imgnew := imaging.Thumbnail(imgsrc, 96, 96, imaging.CatmullRom)

	var imgbuf bytes.Buffer
	err = png.Encode(&imgbuf, imgnew)
	if err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, err.Error())
		return
	}
	imgphoto := base64.StdEncoding.EncodeToString(imgbuf.Bytes())

	// profile
	var profile iamapi.UserProfile
	pobj := store.PvGet("user-profile/" + c.us.UserID)
	if pobj.OK() {
		pobj.Decode(&profile)
	}

	profile.Photo = "data:image/png;base64," + imgphoto
	profile.PhotoSource = req.Data

	store.PvPut("user-profile/"+c.us.UserID, profile, &skv.PvWriteOptions{
		PrevVersion: pobj.Meta().Version,
	})

	set.Kind = "UserPhoto"
}

func (c User) RoleListAction() {

	sets := iamapi.UserRoleList{}
	defer c.RenderJson(&sets)

	// TODO page
	if objs := store.PvScan("role/", "", "", 1000); objs.OK() {

		rss := objs.KvList()
		for _, obj := range rss {

			var role iamapi.UserRole
			if err := obj.Decode(&role); err == nil {

				if role.Id == 1 {
					continue
				}

				if role.Id <= 1000 || role.Meta.UserID == c.us.UserID {
					sets.Items = append(sets.Items, role)
				}
			}
		}
	}

	sets.Kind = "UserRoleList"
}
