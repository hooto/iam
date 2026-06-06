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

package auth

import "github.com/hooto/httpsrv"

// NewModule creates an httpsrv.Module for the auth API endpoints.
func NewModule() *httpsrv.Module {
	mod := httpsrv.NewModule()
	mod.RegisterAction("/sys/info", Sys_Info)
	mod.RegisterAction("/sign-in", SignIn)
	mod.RegisterAction("/sign-out", SignOut)
	mod.RegisterAction("/sign-up", SignUp)
	mod.RegisterAction("/session", Session)
	mod.RegisterAction("/password/reset-ticket", Password_ResetTicket)
	mod.RegisterAction("/password/reset-confirm", Password_ResetConfirm)
	mod.RegisterAction("/photo/:username", Photo)
	return mod
}
