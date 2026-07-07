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

package admin

import (
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/hooto/httpsrv"
	"github.com/lessos/lessgo/pass"
	"github.com/sysinner/innerstack/v2/pkg/inauth"

	"github.com/hooto/iam/v2/internal/data"
	"github.com/hooto/iam/v2/internal/util"
	"github.com/hooto/iam/v2/pkg/iamapi"
)

var errPasswordLength = errors.New("Password must be between 8 and 30 characters long")

// UserList returns all user accounts (group-type entries are excluded),
// optionally filtered by a case-insensitive match on name/email.
// Password keys are never included in the response.
func UserList(ctx httpsrv.Ctx) error {

	if authAdmin(ctx) == nil {
		return nil
	}

	var rsp AdminUserListResponse
	defer ctx.JSON(&rsp)

	qry := strings.ToLower(strings.TrimSpace(ctx.Params().Value("qry_text")))

	users := data.UserList()
	names := make([]string, 0, len(users))
	for name := range users {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {

		u := users[name]
		if u == nil || u.Type == iamapi.UserTypeGroup {
			continue
		}

		if qry != "" &&
			!strings.Contains(strings.ToLower(u.Name), qry) &&
			!strings.Contains(strings.ToLower(u.Email), qry) {
			continue
		}

		// shallow copy + detach keys so cached entries are not mutated
		// and password hashes are never leaked.
		item := *u
		item.Keys = nil
		rsp.Items = append(rsp.Items, &item)
	}

	rsp.Status = inauth.NewServiceStatus("200", "ok")
	return nil
}

type AdminUserListResponse struct {
	Status inauth.ServiceStatus `json:"status"`
	Items  []*iamapi.User       `json:"items,omitempty"`
}

// UserEntry returns a single user (with profile extras) for editing.
// The password is masked with userMgrPasswdHidden.
func UserEntry(ctx httpsrv.Ctx) error {

	if authAdmin(ctx) == nil {
		return nil
	}

	var rsp AdminUserEntryResponse
	defer ctx.JSON(&rsp)

	username := strings.ToLower(strings.TrimSpace(ctx.Params().Value("username")))
	if username == "" {
		rsp.Status = inauth.NewServiceStatus("404", "User Not Found")
		return nil
	}

	u := data.UserGet(username)
	if u == nil || u.Type == iamapi.UserTypeGroup {
		rsp.Status = inauth.NewServiceStatus("404", "User Not Found")
		return nil
	}

	// copy + detach keys, then set the masked sentinel.
	cu := *u
	cu.Keys = nil
	cu.Keys.Set(iamapi.UserKeyDefault, userMgrPasswdHidden)

	item := &AdminUserEntry{User: &cu}

	// merge profile extras (birthday / about)
	var profile iamapi.UserProfile
	if rs := data.Data.NewReader(iamapi.NsUserProfile(username)).Exec(); rs.OK() {
		rs.Item().JsonDecode(&profile)
	}
	item.Birthday = profile.Birthday
	item.About = profile.About

	rsp.Status = inauth.NewServiceStatus("200", "ok")
	rsp.Item = item
	return nil
}

type AdminUserEntry struct {
	*iamapi.User
	Birthday string `json:"birthday,omitempty"`
	About    string `json:"about,omitempty"`
}

type AdminUserEntryResponse struct {
	Status inauth.ServiceStatus `json:"status"`
	Item   *AdminUserEntry      `json:"item,omitempty"`
}

// UserSet creates or updates a user. Create-vs-update is decided by whether
// the named user already exists. On update, a password equal to the masked
// sentinel (or empty) leaves the stored hash untouched.
func UserSet(ctx httpsrv.Ctx) error {

	if authAdmin(ctx) == nil {
		return nil
	}

	var (
		req AdminUserSetRequest
		rsp AdminUserSetResponse
	)
	defer ctx.JSON(&rsp)

	if err := ctx.Request().JsonDecode(&req); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", "Bad Request")
		return nil
	}

	// normalize + validate username
	req.Name = strings.ToLower(strings.TrimSpace(req.Name))
	if err := iamapi.UsernameValid(req.Name); err != nil {
		rsp.Status = inauth.NewServiceStatus("400", err.Error())
		return nil
	}

	// validate roles
	for _, r := range req.Roles {
		if err := iamapi.RoleValid(r); err != nil {
			rsp.Status = inauth.NewServiceStatus("400", "Invalid Role: "+r)
			return nil
		}
	}

	// validate status
	if req.Status != iamapi.UserStatusActive && req.Status != iamapi.UserStatusBanned {
		rsp.Status = inauth.NewServiceStatus("400", "Invalid Status")
		return nil
	}

	tn := time.Now().Unix()
	existing := data.UserGet(req.Name)

	if existing == nil {

		// -- CREATE --
		email := strings.ToLower(strings.TrimSpace(req.Email))
		if err := iamapi.EmailValid(email); err != nil {
			rsp.Status = inauth.NewServiceStatus("400", err.Error())
			return nil
		}
		if err := validatePassword(req.Password); err != nil {
			rsp.Status = inauth.NewServiceStatus("400", err.Error())
			return nil
		}

		roles := req.Roles
		if len(roles) == 0 {
			roles = []string{iamapi.Role_User}
		}

		auth, _ := pass.HashDefault(req.Password)

		user := iamapi.User{
			Name:        req.Name,
			Email:       email,
			DisplayName: strings.TrimSpace(req.DisplayName),
			Roles:       roles,
			Status:      iamapi.UserStatusActive, // new accounts always start active
			Created:     tn,
			Updated:     tn,
		}
		user.Keys.Set(iamapi.UserKeyDefault, auth)

		if rs := data.Data.NewWriter(iamapi.NsUser(user.Name), nil).
			SetJsonValue(&user).Exec(); !rs.OK() {
			rsp.Status = inauth.NewServiceStatus("500", rs.ErrorMessage())
			return nil
		}
		data.UserSet(&user)

		saveProfileExtras(req.Name, req.Birthday, req.About, tn)

	} else {

		// -- UPDATE --
		if existing.Type == iamapi.UserTypeGroup {
			rsp.Status = inauth.NewServiceStatus("400", "Cannot edit a group user")
			return nil
		}

		email := strings.ToLower(strings.TrimSpace(req.Email))
		if err := iamapi.EmailValid(email); err != nil {
			rsp.Status = inauth.NewServiceStatus("400", err.Error())
			return nil
		}

		// operate on the cached pointer (same pattern as user/profile.go),
		// then persist + refresh the cache.
		u := existing
		u.Email = email
		u.DisplayName = strings.TrimSpace(req.DisplayName)
		u.Roles = req.Roles
		if len(u.Roles) == 0 {
			u.Roles = []string{iamapi.Role_User}
		}
		u.Status = req.Status

		// lockout guard: the sysadmin account always keeps the sa role
		// and an active status.
		if u.Name == iamapi.UserSysadmin {
			if !util.Contains(u.Roles, []string{iamapi.Role_Sysadmin}) {
				u.Roles = append(u.Roles, iamapi.Role_Sysadmin)
			}
			u.Status = iamapi.UserStatusActive
		}

		// password: only change when a real new value is provided
		if req.Password != "" && req.Password != userMgrPasswdHidden {
			if err := validatePassword(req.Password); err != nil {
				rsp.Status = inauth.NewServiceStatus("400", err.Error())
				return nil
			}
			auth, _ := pass.HashDefault(req.Password)
			u.Keys.Set(iamapi.UserKeyDefault, auth)
		}

		u.Updated = tn

		if rs := data.Data.NewWriter(iamapi.NsUser(u.Name), nil).
			SetJsonValue(u).Exec(); !rs.OK() {
			rsp.Status = inauth.NewServiceStatus("500", rs.ErrorMessage())
			return nil
		}
		data.UserSet(u)

		saveProfileExtras(req.Name, req.Birthday, req.About, tn)
	}

	rsp.Status = inauth.NewServiceStatus("200", "ok")
	return nil
}

type AdminUserSetRequest struct {
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	DisplayName string   `json:"display_name"`
	Roles       []string `json:"roles"`
	Status      uint8    `json:"status"`
	Birthday    string   `json:"birthday"`
	About       string   `json:"about"`
}

type AdminUserSetResponse struct {
	Status inauth.ServiceStatus `json:"status"`
}

// validatePassword enforces the same length rules as auth/sign-up.go.
func validatePassword(p string) error {
	if len(p) < 8 || len(p) > 30 {
		return errPasswordLength
	}
	return nil
}

// saveProfileExtras merges birthday/about into the user's profile record,
// preserving other fields (e.g. photo) set via other flows.
func saveProfileExtras(name, birthday, about string, tn int64) {
	var profile iamapi.UserProfile
	if rs := data.Data.NewReader(iamapi.NsUserProfile(name)).Exec(); rs.OK() {
		rs.Item().JsonDecode(&profile)
	}
	profile.Birthday = birthday
	profile.About = about
	profile.Updated = tn
	data.Data.NewWriter(iamapi.NsUserProfile(name), nil).SetJsonValue(profile).Exec()
}
