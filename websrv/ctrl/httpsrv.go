// Copyright 2014-2016 iam Author, All rights reserved.
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

package ctrl

import (
	"github.com/lessos/iam/config"
	"github.com/lessos/lessgo/httpsrv"
)

func NewModule() httpsrv.Module {

	module := httpsrv.NewModule("iam_ws")

	module.RouteSet(httpsrv.Route{
		Type:       httpsrv.RouteTypeStatic,
		Path:       "~",
		StaticPath: config.Prefix + "/webui",
	})

	module.TemplatePathSet(config.Prefix + "/websrv/views")

	module.ControllerRegister(new(Index))
	module.ControllerRegister(new(Service))
	module.ControllerRegister(new(Reg))
	module.ControllerRegister(new(User))
	module.ControllerRegister(new(AppAuth))

	// // TODO auto config
	// httpsrv.Config.I18n(config.Config.Prefix + "/i18n/en.json")
	// httpsrv.Config.I18n(config.Config.Prefix + "/i18n/zh_CN.json")

	return module
}
