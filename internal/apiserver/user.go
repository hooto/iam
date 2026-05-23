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

package apiserver

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/hooto/httpsrv"
	"github.com/lessos/lessgo/pass"
	"github.com/sysinner/incore/v2/pkg/inauth"

	"github.com/hooto/iam/v2/internal/data"
	"github.com/hooto/iam/v2/pkg/iamapi"
)

// User provides authenticated user self-service actions:
//   - Profile       GET  - get current user profile
//   - ProfileSet    POST - update display name, birthday, about
//   - PassSet       POST - change password
//   - EmailSet      POST - change email
//   - PhotoSet      POST - upload avatar
//   - RoleList      GET  - list roles visible to current user
type User struct {
	*httpsrv.Controller
}

type UserProfileResponse struct {
	Status inauth.ServiceStatus `json:"status"`
	Item   *iamapi.UserProfile  `json:"item,omitempty"`
}

func (c User) ProfileAction() {

	user := userAuth(c.Controller)
	if user == nil {
		return
	}

	var rsp UserProfileResponse
	defer c.RenderJson(&rsp)

	entry := iamapi.UserProfile{
		DisplayName: user.DisplayName,
		Email:       user.Email,
	}

	// profile extras
	var profile iamapi.UserProfile
	if rs := data.Data.NewReader(iamapi.NsUserProfile(user.Name)).Exec(); rs.OK() {
		rs.Item().JsonDecode(&profile)
	}

	entry.Gender = profile.Gender
	entry.Birthday = profile.Birthday
	entry.About = profile.About
	entry.Photo = profile.Photo

	rsp.Status = inauth.NewServiceStatus("200", "ok")
	rsp.Item = &entry
}

type UserProfileSetRequest struct {
	DisplayName string `json:"display_name"`
	Birthday    string `json:"birthday"`
	About       string `json:"about"`
}

type UserProfileSetResponse struct {
	Status inauth.ServiceStatus `json:"status"`
}

func (c User) ProfileSetAction() {

	user := userAuth(c.Controller)
	if user == nil {
		return
	}

	var (
		req UserProfileSetRequest
		rsp UserProfileSetResponse
	)
	defer c.RenderJson(&rsp)

	if err := c.Request.JsonDecode(&req); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", "Bad Request")
		return
	}

	// validate display name
	req.DisplayName = strings.TrimSpace(req.DisplayName)
	if len(req.DisplayName) < 1 || len(req.DisplayName) > 30 {
		rsp.Status = inauth.NewServiceStatus("400", "DisplayName must be between 1 and 30 characters long")
		return
	}

	// validate birthday
	if _, err := time.Parse("2006-01-02", req.Birthday); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", "Birthday is not valid")
		return
	}

	if req.About == "" {
		rsp.Status = inauth.NewServiceStatus("400", "About Me can not be null")
		return
	}

	user.DisplayName = req.DisplayName
	user.Updated = time.Now().Unix()
	data.Data.NewWriter(iamapi.NsUser(user.Name), nil).SetJsonValue(user).Exec()
	data.UserSet(user)

	// update profile extras
	var profile iamapi.UserProfile
	if rs := data.Data.NewReader(iamapi.NsUserProfile(user.Name)).Exec(); rs.OK() {
		rs.Item().JsonDecode(&profile)
	}

	profile.Birthday = req.Birthday
	profile.About = req.About

	data.Data.NewWriter(iamapi.NsUserProfile(user.Name), nil).SetJsonValue(profile).Exec()

	rsp.Status = inauth.NewServiceStatus("200", "ok")
}

type UserPassSetRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type UserPassSetResponse struct {
	Status inauth.ServiceStatus `json:"status"`
}

func (c User) PassSetAction() {

	user := userAuth(c.Controller)
	if user == nil {
		return
	}

	var (
		req UserPassSetRequest
		rsp UserPassSetResponse
	)
	defer c.RenderJson(&rsp)

	if err := c.Request.JsonDecode(&req); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", "Bad Request")
		return
	}

	if len(req.NewPassword) < 8 || len(req.NewPassword) > 30 {
		rsp.Status = inauth.NewServiceStatus("400", "Password must be between 8 and 30 characters long")
		return
	}

	if auth := user.Keys.Get(iamapi.UserKeyDefault); auth == nil ||
		!pass.Check(req.CurrentPassword, auth.String()) {
		rsp.Status = inauth.NewServiceStatus("400", "Current Password can not match")
		return
	}

	user.Updated = time.Now().Unix()
	authKey, _ := pass.HashDefault(req.NewPassword)
	user.Keys.Set(iamapi.UserKeyDefault, authKey)

	data.Data.NewWriter(iamapi.NsUser(user.Name), nil).SetJsonValue(user).Exec()
	data.UserSet(user)

	rsp.Status = inauth.NewServiceStatus("200", "ok")
}

type UserEmailSetRequest struct {
	Auth  string `json:"auth"`
	Email string `json:"email"`
}

type UserEmailSetResponse struct {
	Status inauth.ServiceStatus `json:"status"`
}

func (c User) EmailSetAction() {

	user := userAuth(c.Controller)
	if user == nil {
		return
	}

	var (
		req UserEmailSetRequest
		rsp UserEmailSetResponse
	)
	defer c.RenderJson(&rsp)

	if err := c.Request.JsonDecode(&req); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", "Bad Request")
		return
	}

	// validate email
	email := strings.TrimSpace(strings.ToLower(req.Email))
	if err := iamapi.EmailValid(email); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", err.Error())
		return
	}

	if auth := user.Keys.Get(iamapi.UserKeyDefault); auth == nil ||
		!pass.Check(req.Auth, auth.String()) {
		rsp.Status = inauth.NewServiceStatus("400", "Password can not match")
		return
	}

	user.Email = email
	user.Updated = time.Now().Unix()
	data.Data.NewWriter(iamapi.NsUser(user.Name), nil).SetJsonValue(user).Exec()
	data.UserSet(user)

	rsp.Status = inauth.NewServiceStatus("200", "ok")
}

type UserPhotoSetRequest struct {
	Data string `json:"data"`
}

type UserPhotoSetResponse struct {
	Status inauth.ServiceStatus `json:"status"`
}

func (c User) PhotoSetAction() {

	user := userAuth(c.Controller)
	if user == nil {
		return
	}

	var (
		req UserPhotoSetRequest
		rsp UserPhotoSetResponse
	)
	defer c.RenderJson(&rsp)

	if err := c.Request.JsonDecode(&req); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", "Bad Request")
		return
	}

	img64 := strings.SplitAfter(req.Data, ";base64,")
	if len(img64) != 2 {
		rsp.Status = inauth.NewServiceStatus("400", "Bad Request")
		return
	}

	// decode and resize image
	imgReader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(img64[1]))
	imgSrc, _, err := image.Decode(imgReader)
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("400", err.Error())
		return
	}

	imgNew := imaging.Thumbnail(imgSrc, 96, 96, imaging.CatmullRom)

	var imgBuf bytes.Buffer
	if err := png.Encode(&imgBuf, imgNew); err != nil {
		rsp.Status = inauth.NewServiceStatus("500", err.Error())
		return
	}

	// load or create profile
	var profile iamapi.UserProfile
	if rs := data.Data.NewReader(iamapi.NsUserProfile(user.Name)).Exec(); rs.OK() {
		rs.Item().JsonDecode(&profile)
	}

	if profile.Login != nil && profile.Login.Name != "" &&
		profile.Login.Name != user.Name {
		rsp.Status = inauth.NewServiceStatus("400", "Unauthorized")
		return
	}

	profile.Photo = "data:image/png;base64," + base64.StdEncoding.EncodeToString(imgBuf.Bytes())
	profile.PhotoSource = req.Data

	data.Data.NewWriter(iamapi.NsUserProfile(user.Name), nil).SetJsonValue(profile).Exec()

	rsp.Status = inauth.NewServiceStatus("200", "ok")
}

type UserRoleListResponse struct {
	Status inauth.ServiceStatus `json:"status"`
	Items  []iamapi.UserRole    `json:"items,omitempty"`
}

func (c User) RoleListAction() {

	user := userAuth(c.Controller)
	if user == nil {
		return
	}

	var rsp UserRoleListResponse
	defer c.RenderJson(&rsp)

	if rs := data.Data.NewRanger(
		iamapi.NsRole(""), iamapi.NsRole("")).SetLimit(1000).Exec(); rs.OK() {

		for _, obj := range rs.Items {

			var role iamapi.UserRole
			if err := obj.JsonDecode(&role); err != nil {
				continue
			}

			if role.Name != "" || role.User == user.Name {
				rsp.Items = append(rsp.Items, role)
			}
		}
	}

	rsp.Status = inauth.NewServiceStatus("200", "ok")
}
