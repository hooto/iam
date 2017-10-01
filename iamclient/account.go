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
	"github.com/lessos/lessgo/encoding/json"
	"github.com/lessos/lessgo/net/httpclient"
	"github.com/lessos/lessgo/types"
)

func AccountChargePrepay(req iamapi.AccountChargePrepay) iamapi.AccountChargePrepay {

	js, _ := json.Encode(req, "")

	hc := httpclient.Post(fmt.Sprintf(
		"%s/v1/account-charge/prepay",
		ServiceUrl,
	))
	hc.Header("contentType", "application/json; charset=utf-8")
	hc.Body(js)
	defer hc.Close()

	var rsp iamapi.AccountChargePrepay
	if err := hc.ReplyJson(&rsp); err != nil && rsp.Error == nil {
		rsp.Error = types.NewErrorMeta("400", "Network Error")
	}
	return rsp
}

func AccountChargePayout(req iamapi.AccountChargePayout) iamapi.AccountChargePayout {

	js, _ := json.Encode(req, "")

	hc := httpclient.Post(fmt.Sprintf(
		"%s/v1/account-charge/payout",
		ServiceUrl,
	))
	hc.Header("contentType", "application/json; charset=utf-8")
	hc.Body(js)
	defer hc.Close()

	var rsp iamapi.AccountChargePayout
	if err := hc.ReplyJson(&rsp); err != nil && rsp.Error == nil {
		rsp.Error = types.NewErrorMeta("400", "Network Error")
	}
	return rsp
}
