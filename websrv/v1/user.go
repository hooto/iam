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
	"html"
	"image"
	"image/png"
	"strings"

	"github.com/eryx/imaging"
	"github.com/lessos/lessgo/data/rdo"
	"github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/httpsrv"
	"github.com/lessos/lessgo/pass"
	"github.com/lessos/lessgo/types"

	"../../base/login"
	"../../base/profile"
	"../../base/session"
	"../../idsapi"
)

type User struct {
	*httpsrv.Controller
}

func (c User) ProfileAction() {

	rsp := idsapi.UserProfile{}

	defer c.RenderJson(&rsp)

	s := session.GetSession(c.Request)
	if s.Uid == "" {
		rsp.Error = &types.ErrorMeta{"401", "Access Denied"}
		return
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Error = &types.ErrorMeta{"401", "Access Denied"}
		return
	}

	// login
	q := base.NewQuerySet().From("ids_login").Limit(1)
	q.Where.And("uid", s.Uid)
	rslogin, err := dcn.Base.Query(q)
	if err != nil || len(rslogin) != 1 {
		rsp.Error = &types.ErrorMeta{"401", "Access Denied"}
		return
	}

	rsp.Login.Meta.ID = rslogin[0].Field("uid").String()
	rsp.Login.Meta.Name = rslogin[0].Field("uname").String()
	rsp.Login.Name = rslogin[0].Field("name").String()
	rsp.Login.Email = rslogin[0].Field("email").String()

	rsp.Name = rslogin[0].Field("name").String()

	//
	q = base.NewQuerySet().From("ids_profile").Limit(1)
	q.Where.And("uid", s.Uid)
	rs, err := dcn.Base.Query(q)
	if err != nil || len(rs) != 1 {

		item := map[string]interface{}{
			"uid":     s.Uid,
			"gender":  0,
			"created": base.TimeNow("datetime"), // TODO
			"updated": base.TimeNow("datetime"), // TODO
		}

		dcn.Base.Insert("ids_profile", item)

	} else {
		rsp.Birthday = rs[0].Field("birthday").String()
		rsp.About = html.EscapeString(rs[0].Field("aboutme").String())
	}

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

	s := session.GetSession(c.Request)
	if s.Uid == "" {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInternalError, "Server Error"}
		return
	}

	itemlogin := map[string]interface{}{
		"name":    req.Name,
		"updated": base.TimeNow("datetime"),
	}
	ft := base.NewFilter()
	ft.And("uid", s.Uid)
	dcn.Base.Update("ids_login", itemlogin, ft)

	itemprofile := map[string]interface{}{
		"birthday": req.Birthday,
		"aboutme":  req.About,
		"updated":  base.TimeNow("datetime"), // TODO
	}
	dcn.Base.Update("ids_profile", itemprofile, ft)

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

	s := session.GetSession(c.Request)
	if s.Uid == "" {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInternalError, "Server Error"}
		return
	}

	q := base.NewQuerySet().From("ids_login").Limit(1)
	q.Where.And("uid", s.Uid)
	rsu, err := dcn.Base.Query(q)
	if err == nil && len(rsu) == 0 {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "User Not Found"}
		return
	}

	if !pass.Check(req.CurrentPassword, rsu[0].Field("pass").String()) {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "Current Password can not match"}

		return
	}

	pstr, _ := pass.HashDefault(req.NewPassword)

	itemlogin := map[string]interface{}{
		"pass":    pstr,
		"updated": base.TimeNow("datetime"),
	}
	ft := base.NewFilter()
	ft.And("uid", s.Uid)
	dcn.Base.Update("ids_login", itemlogin, ft)

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

	s := session.GetSession(c.Request)
	if s.Uid == "" {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
		return
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInternalError, "Server Error"}
		return
	}

	q := base.NewQuerySet().From("ids_login").Limit(1)
	q.Where.And("uid", s.Uid)
	rsu, err := dcn.Base.Query(q)
	if err == nil && len(rsu) == 0 {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "User can not found"}
		return
	}

	if !pass.Check(req.Auth, rsu[0].Field("pass").String()) {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInvalidArgument, "Password can not match"}
		return
	}

	itemlogin := map[string]interface{}{
		"email":   req.Email,
		"updated": base.TimeNow("datetime"),
	}

	ft := base.NewFilter()
	ft.And("uid", s.Uid)
	dcn.Base.Update("ids_login", itemlogin, ft)

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

	s := session.GetSession(c.Request)
	if s.Uid == "" {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeAccessDenied, "Access Denied"}
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

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Error = &types.ErrorMeta{idsapi.ErrCodeInternalError, "Server Error"}
		return
	}

	itemprofile := map[string]interface{}{
		"photo":    "data:image/png;base64," + imgphoto,
		"photosrc": req.Data,
		"updated":  base.TimeNow("datetime"),
	}
	ft := base.NewFilter()
	ft.And("uid", s.Uid)
	dcn.Base.Update("ids_profile", itemprofile, ft)

	rsp.Kind = "UserPhoto"
}
