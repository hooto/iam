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

package iamclient

import (
	"fmt"

	"github.com/hooto/iam/iamapi"
	"github.com/lessos/lessgo/net/httpclient"
	"github.com/lessos/lessgo/types"
)

func PublicUserEntry(user string) iamapi.UserEntry {

	hc := httpclient.Get(fmt.Sprintf(
		"%s/user/entry?user=%s",
		ServiceUrl,
		user,
	))
	defer hc.Close()

	var rsp iamapi.UserEntry
	if err := hc.ReplyJson(&rsp); err != nil {
		rsp.Error = types.NewErrorMeta("400", err.Error())
	}

	return rsp
}
