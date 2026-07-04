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

package user

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/hooto/httpsrv"
	"github.com/sysinner/innerstack/v2/pkg/inauth"

	"github.com/hooto/iam/v2/internal/data"
	"github.com/hooto/iam/v2/pkg/iamapi"
)

type UserProfileResponse struct {
	Status inauth.ServiceStatus `json:"status"`
	Item   *iamapi.UserProfile  `json:"item,omitempty"`
}

// ProfileEntry returns the current user profile.
func ProfileEntry(ctx httpsrv.Ctx) error {

	u := authCtx(ctx)
	if u == nil {
		return nil
	}

	var rsp UserProfileResponse
	defer ctx.JSON(&rsp)

	entry := iamapi.UserProfile{
		DisplayName: u.DisplayName,
		Email:       u.Email,
	}

	// profile extras
	var profile iamapi.UserProfile
	if rs := data.Data.NewReader(iamapi.NsUserProfile(u.Name)).Exec(); rs.OK() {
		rs.Item().JsonDecode(&profile)
	}

	entry.Gender = profile.Gender
	entry.Birthday = profile.Birthday
	entry.About = profile.About
	entry.Photo = profile.Photo

	rsp.Status = inauth.NewServiceStatus("200", "ok")
	rsp.Item = &entry
	return nil
}

type UserProfileSetRequest struct {
	DisplayName string `json:"display_name"`
	Birthday    string `json:"birthday"`
	About       string `json:"about"`
}

type UserProfileSetResponse struct {
	Status inauth.ServiceStatus `json:"status"`
}

// ProfileSet updates the current user profile (display name, birthday, about).
func ProfileSet(ctx httpsrv.Ctx) error {

	u := authCtx(ctx)
	if u == nil {
		return nil
	}

	var (
		req UserProfileSetRequest
		rsp UserProfileSetResponse
	)
	defer ctx.JSON(&rsp)

	if err := ctx.Request().JsonDecode(&req); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", "Bad Request")
		return nil
	}

	// validate display name
	req.DisplayName = strings.TrimSpace(req.DisplayName)
	if len(req.DisplayName) < 1 || len(req.DisplayName) > 30 {
		rsp.Status = inauth.NewServiceStatus("400", "DisplayName must be between 1 and 30 characters long")
		return nil
	}

	// validate birthday
	if _, err := time.Parse("2006-01-02", req.Birthday); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", "Birthday is not valid")
		return nil
	}

	if req.About == "" {
		rsp.Status = inauth.NewServiceStatus("400", "About Me can not be null")
		return nil
	}

	u.DisplayName = req.DisplayName
	u.Updated = time.Now().Unix()
	data.Data.NewWriter(iamapi.NsUser(u.Name), nil).SetJsonValue(u).Exec()
	data.UserSet(u)

	// update profile extras
	var profile iamapi.UserProfile
	if rs := data.Data.NewReader(iamapi.NsUserProfile(u.Name)).Exec(); rs.OK() {
		rs.Item().JsonDecode(&profile)
	}

	profile.Birthday = req.Birthday
	profile.About = req.About

	data.Data.NewWriter(iamapi.NsUserProfile(u.Name), nil).SetJsonValue(profile).Exec()

	rsp.Status = inauth.NewServiceStatus("200", "ok")
	return nil
}

type UserPhotoSetRequest struct {
	Data string `json:"data"`
}

type UserPhotoSetResponse struct {
	Status inauth.ServiceStatus `json:"status"`
}

// ProfilePhotoSet uploads and sets the user avatar image.
func ProfilePhotoSet(ctx httpsrv.Ctx) error {

	u := authCtx(ctx)
	if u == nil {
		return nil
	}

	var (
		req UserPhotoSetRequest
		rsp UserPhotoSetResponse
	)
	defer ctx.JSON(&rsp)

	if err := ctx.Request().JsonDecode(&req); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", "Bad Request")
		return nil
	}

	img64 := strings.SplitAfter(req.Data, ";base64,")
	if len(img64) != 2 {
		rsp.Status = inauth.NewServiceStatus("400", "Bad Request")
		return nil
	}

	// decode and resize image
	imgReader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(img64[1]))
	imgSrc, _, err := image.Decode(imgReader)
	if err != nil {
		rsp.Status = inauth.NewServiceStatus("400", err.Error())
		return nil
	}

	imgNew := imaging.Thumbnail(imgSrc, 96, 96, imaging.CatmullRom)

	var buf []byte
	imgBuf := bytes.NewBuffer(buf)
	if err := png.Encode(imgBuf, imgNew); err != nil {
		rsp.Status = inauth.NewServiceStatus("500", err.Error())
		return nil
	}

	// load or create profile
	var profile iamapi.UserProfile
	if rs := data.Data.NewReader(iamapi.NsUserProfile(u.Name)).Exec(); rs.OK() {
		rs.Item().JsonDecode(&profile)
	}

	if profile.Login != nil && profile.Login.Name != "" &&
		profile.Login.Name != u.Name {
		rsp.Status = inauth.NewServiceStatus("400", "Unauthorized")
		return nil
	}

	profile.Photo = "data:image/png;base64," + base64.StdEncoding.EncodeToString(imgBuf.Bytes())
	profile.PhotoSource = req.Data

	data.Data.NewWriter(iamapi.NsUserProfile(u.Name), nil).SetJsonValue(profile).Exec()

	rsp.Status = inauth.NewServiceStatus("200", "ok")
	return nil
}
