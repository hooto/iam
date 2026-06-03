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

package open

import "github.com/hooto/httpsrv"

// NewModule creates an httpsrv.Module for the open API endpoints.
func NewModule() *httpsrv.Module {
	mod := httpsrv.NewModule()
	mod.RegisterAction("/app-auth/verify", AppAuth_Verify)
	mod.RegisterAction("/app-auth/token-exchange", AppAuth_TokenExchange)
	mod.RegisterAction("/app-auth/session", AppAuth_Session)
	mod.RegisterAction("/app-auth/update", AppAuth_Update)
	return mod
}
