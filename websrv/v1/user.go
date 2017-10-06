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

	"github.com/eryx/imaging"
	"github.com/hooto/httpsrv"
	"github.com/lessos/lessgo/pass"
	"github.com/lessos/lessgo/types"
	"github.com/lynkdb/iomix/skv"

	"github.com/hooto/iam/base/login"
	"github.com/hooto/iam/base/profile"
	"github.com/hooto/iam/iamapi"
	"github.com/hooto/iam/iamclient"
	"github.com/hooto/iam/store"
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

	var set struct {
		types.TypeMeta
		iamapi.UserProfile
	}
	defer c.RenderJson(&set)

	// profile
	if obj := store.Data.ProgGet(iamapi.DataUserProfileKey(c.us.UserName)); obj.OK() {
		obj.Decode(&set.UserProfile)
	}

	set.Login = nil
	if set.Login == nil || set.Login.Name == "" {

		// login
		var user iamapi.User
		if obj := store.Data.ProgGet(iamapi.DataUserKey(c.us.UserName)); obj.OK() {
			obj.Decode(&user)
		}

		if user.Name != c.us.UserName {
			set.Error = types.NewErrorMeta(iamapi.ErrCodeNotFound, "User Not Found")
			return
		}

		set.Login = &user
		store.Data.ProgPut(iamapi.DataUserProfileKey(c.us.UserName), skv.NewProgValue(set), nil)
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
	uobj := store.Data.ProgGet(iamapi.DataUserKey(c.us.UserName))
	if uobj.OK() {
		uobj.Decode(&user)
	}

	if user.Name != c.us.UserName {
		set.Error = types.NewErrorMeta("404", "User Not Found")
		return
	}
	user.DisplayName = req.Login.DisplayName

	store.Data.ProgPut(iamapi.DataUserKey(c.us.UserName), skv.NewProgValue(user), nil)

	// profile
	var profile iamapi.UserProfile
	pobj := store.Data.ProgGet(iamapi.DataUserProfileKey(c.us.UserName))
	if pobj.OK() {
		pobj.Decode(&profile)
	}

	profile.Birthday = req.Birthday
	profile.About = req.About

	profile.Login = &user
	profile.Login.Keys = nil

	store.Data.ProgPut(iamapi.DataUserProfileKey(c.us.UserName), skv.NewProgValue(profile), nil)

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
	uobj := store.Data.ProgGet(iamapi.DataUserKey(c.us.UserName))
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

	store.Data.ProgPut(iamapi.DataUserKey(c.us.UserName), skv.NewProgValue(user), nil)

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
	uobj := store.Data.ProgGet(iamapi.DataUserKey(c.us.UserName))
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
	store.Data.ProgPut(iamapi.DataUserKey(c.us.UserName), skv.NewProgValue(user), nil)

	if rs := store.Data.ProgGet(iamapi.DataUserProfileKey(c.us.UserName)); rs.OK() {
		var preprofile iamapi.UserProfile
		if err := rs.Decode(&preprofile); err == nil {
			preprofile.Login = &user
			store.Data.ProgPut(iamapi.DataUserProfileKey(c.us.UserName), skv.NewProgValue(preprofile), nil)
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
	pobj := store.Data.ProgGet(iamapi.DataUserProfileKey(c.us.UserName))
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

	store.Data.ProgPut(iamapi.DataUserProfileKey(c.us.UserName), skv.NewProgValue(profile), nil)

	set.Kind = "UserPhoto"
}

func (c User) RoleListAction() {

	sets := iamapi.UserRoleList{}
	defer c.RenderJson(&sets)

	// TODO page
	if objs := store.Data.ProgScan(iamapi.DataRoleKey(0), iamapi.DataRoleKey(99999999), 1000); objs.OK() {

		rss := objs.KvList()
		for _, obj := range rss {

			var role iamapi.UserRole
			if err := obj.Decode(&role); err == nil {

				if role.Id == 1 {
					continue
				}

				if role.Id <= 1000 || role.User == c.us.UserName {
					sets.Items = append(sets.Items, role)
				}
			}
		}
	}

	sets.Kind = "UserRoleList"
}
