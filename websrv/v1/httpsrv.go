// Copyright 2014 lessos Authors, All rights reserved.
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
	"github.com/lessos/lessgo/httpsrv"
)

func NewModule() httpsrv.Module {

	module := httpsrv.NewModule("iam_api")

	module.ControllerRegister(new(Service))

	module.ControllerRegister(new(User))
	module.ControllerRegister(new(MyApp))

	module.ControllerRegister(new(SysConfig))
	module.ControllerRegister(new(UserMgr))
	module.ControllerRegister(new(AppMgr))

	module.ControllerRegister(new(AppAuth))
	module.ControllerRegister(new(AccessKey))
	module.ControllerRegister(new(Status))

	return module
}
