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
	"github.com/hooto/hlang4g/hlang"
	"github.com/hooto/httpsrv"
	"github.com/hooto/iam/bindata"
	"github.com/hooto/iam/config"
)

func NewModule() httpsrv.Module {

	module := httpsrv.NewModule("iam_api")

	module.ControllerRegister(new(Service))

	module.ControllerRegister(new(User))
	module.ControllerRegister(new(UserGroup))
	module.ControllerRegister(new(Account))
	module.ControllerRegister(new(AccountCharge))
	module.ControllerRegister(new(App))

	module.ControllerRegister(new(AppAuth))
	module.ControllerRegister(new(AccessKey))
	module.ControllerRegister(new(Status))

	module.ControllerRegister(new(UserMgr))
	module.ControllerRegister(new(AppMgr))
	module.ControllerRegister(new(AccountMgr))
	module.ControllerRegister(new(SysConfig))
	module.ControllerRegister(new(SysMsg))

	// TODO auto config
	if fs := bindata.NewFs("iam_i18n"); fs != nil {
		hlang.StdLangFeed.LoadMessageWithFs(fs)
	} else if config.Prefix != "" {
		hlang.StdLangFeed.LoadMessages(config.Prefix+"/i18n/en.json", true)
		hlang.StdLangFeed.LoadMessages(config.Prefix+"/i18n/zh-CN.json", true)
	}

	if hlang.StdLangFeed.Init() {
		module.ControllerRegister(new(hlang.Langsrv))
	}

	return module
}
