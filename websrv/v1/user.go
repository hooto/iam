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
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/hooto/httpsrv"
	"github.com/lessos/lessgo/pass"
	"github.com/lessos/lessgo/types"

	"github.com/hooto/iam/base/login"
	"github.com/hooto/iam/base/profile"
	"github.com/hooto/iam/data"
	"github.com/hooto/iam/iamapi"
	"github.com/hooto/iam/iamclient"
)

type User struct {
	*httpsrv.Controller
	us iamapi.UserSession
	//
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

	var set struct {
		types.TypeMeta
		iamapi.UserProfile
	}
	defer c.RenderJson(&set)

	// profile
	if obj := data.Data.NewReader(iamapi.ObjKeyUserProfile(c.us.UserName)).Query(); obj.OK() {
		obj.Decode(&set.UserProfile)
	}

	set.Login = nil
	if set.Login == nil || set.Login.Name == "" {

		// login
		var user iamapi.User
		if obj := data.Data.NewReader(iamapi.ObjKeyUser(c.us.UserName)).Query(); obj.OK() {
			obj.Decode(&user)
		}

		if user.Name != c.us.UserName {
			set.Error = types.NewErrorMeta(iamapi.ErrCodeNotFound, "User Not Found")
			return
		}

		set.Login = &user
		data.Data.NewWriter(iamapi.ObjKeyUserProfile(c.us.UserName), set).Commit()
	}

	set.Photo = ""
	set.PhotoSource = ""

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
	uobj := data.Data.NewReader(iamapi.ObjKeyUser(c.us.UserName)).Query()
	if uobj.OK() {
		uobj.Decode(&user)
	}

	if user.Name != c.us.UserName {
		set.Error = types.NewErrorMeta("404", "User Not Found")
		return
	}
	user.DisplayName = req.Login.DisplayName

	data.Data.NewWriter(iamapi.ObjKeyUser(c.us.UserName), user).Commit()

	// profile
	var profile iamapi.UserProfile
	pobj := data.Data.NewReader(iamapi.ObjKeyUserProfile(c.us.UserName)).Query()
	if pobj.OK() {
		pobj.Decode(&profile)
	}

	profile.Birthday = req.Birthday
	profile.About = req.About

	profile.Login = &user
	profile.Login.Keys = nil

	data.Data.NewWriter(iamapi.ObjKeyUserProfile(c.us.UserName), profile).Commit()

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
	uobj := data.Data.NewReader(iamapi.ObjKeyUser(c.us.UserName)).Query()
	if uobj.OK() {
		uobj.Decode(&user)
	}

	if user.Name != c.us.UserName {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "User Not Found")
		return
	}

	if auth := user.Keys.Get(iamapi.UserKeyDefault); auth == nil ||
		!pass.Check(req.CurrentPassword, auth.String()) {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Current Password can not match")
		return
	}

	user.Updated = types.MetaTimeNow()
	auth_key, _ := pass.HashDefault(req.NewPassword)
	user.Keys.Set(iamapi.UserKeyDefault, auth_key)

	data.Data.NewWriter(iamapi.ObjKeyUser(c.us.UserName), user).Commit()

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
	uobj := data.Data.NewReader(iamapi.ObjKeyUser(c.us.UserName)).Query()
	if uobj.OK() {
		uobj.Decode(&user)
	}

	if user.Name != c.us.UserName {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "User Not Found")
		return
	}

	if auth := user.Keys.Get(iamapi.UserKeyDefault); auth == nil ||
		!pass.Check(req.Auth, auth.String()) {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "Password can not match")
		return
	}

	user.Email = req.Email
	user.Updated = types.MetaTimeNow()
	data.Data.NewWriter(iamapi.ObjKeyUser(c.us.UserName), user).Commit()

	if rs := data.Data.NewReader(iamapi.ObjKeyUserProfile(c.us.UserName)).Query(); rs.OK() {
		var preprofile iamapi.UserProfile
		if err := rs.Decode(&preprofile); err == nil {
			preprofile.Login = &user
			data.Data.NewWriter(iamapi.ObjKeyUserProfile(c.us.UserName), preprofile).Commit()
		}
	}

	set.Kind = "UserEmail"
}

func (c User) PhotoSetAction() {

	var (
		set types.TypeMeta
		req iamapi.UserPhotoSet
	)
	defer c.RenderJson(&set)

	if err := c.Request.JsonDecode(&req); err != nil {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "InvalidArgument")
		return
	}

	//
	img64 := strings.SplitAfter(req.Data, ";base64,")
	if len(img64) != 2 {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "InvalidArgument")
		return
	}

	// profile
	var profile iamapi.UserProfile
	pobj := data.Data.NewReader(iamapi.ObjKeyUserProfile(c.us.UserName)).Query()
	if pobj.OK() {
		pobj.Decode(&profile)
	}

	if profile.Login != nil && profile.Login.Name != "" && profile.Login.Name != c.us.UserName {
		set.Error = types.NewErrorMeta(iamapi.ErrCodeInvalidArgument, "InvalidArgument")
		return
	}

	//
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

	profile.Photo = "data:image/png;base64," + imgphoto
	profile.PhotoSource = req.Data

	data.Data.NewWriter(iamapi.ObjKeyUserProfile(c.us.UserName), profile).Commit()

	set.Kind = "UserPhoto"
}

func (c User) RoleListAction() {

	sets := iamapi.UserRoleList{}
	defer c.RenderJson(&sets)

	if rs := data.Data.NewReader(nil).KeyRangeSet(
		iamapi.ObjKeyRole(""), iamapi.ObjKeyRole("")).LimitNumSet(1000).Query(); rs.OK() {

		for _, obj := range rs.Items {

			var role iamapi.UserRole
			if err := obj.DataValue().Decode(&role, nil); err == nil {

				if obj.Meta.IncrId == 1 {
					continue
				}

				role.Id = uint32(obj.Meta.IncrId)

				if role.Id <= 1000 || role.User == c.us.UserName {
					sets.Items = append(sets.Items, role)
				}
			}
		}
	}

	sets.Kind = "UserRoleList"
}
