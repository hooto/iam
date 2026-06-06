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

import (
	"github.com/hooto/httpsrv"
	"github.com/hooto/iam/v2/internal/config"
	"github.com/sysinner/incore/v2/pkg/inauth"
)

type Sys_InfoResponse struct {
	Status          inauth.ServiceStatus `json:"status"`
	InstanceId      string               `json:"instance_id"`
	AllowUserSignUp bool                 `json:"allow_user_sign_up"`
}

func Sys_Info(ctx httpsrv.Ctx) error {
	rsp := Sys_InfoResponse{
		Status:          inauth.NewServiceStatus("200", "ok"),
		InstanceId:      config.Config.InstanceID,
		AllowUserSignUp: config.AllowUserSignUp,
	}
	return ctx.JSON(&rsp)
}
