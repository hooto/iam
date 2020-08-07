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
	"errors"
	"fmt"

	"github.com/hooto/hauth/go/hauth/v1"
	"github.com/hooto/iam/iamapi"
	"github.com/lessos/lessgo/crypto/idhash"
	"github.com/lessos/lessgo/encoding/json"
	"github.com/lessos/lessgo/net/httpclient"
	"github.com/lessos/lessgo/types"
)

func SysMsgPost(req iamapi.MsgItem, ak *iamapi.AccessKey) error {

	req.Id = idhash.RandHexString(16)

	js, _ := json.Encode(req, "")

	hc := httpclient.Post(fmt.Sprintf(
		"%s/v1/sys-msg/post",
		service_url_global(),
	))
	defer hc.Close()

	hc.Header("contentType", "application/json; charset=utf-8")

	ac := hauth.NewAppCredential(ak.AuthKey())
	ac.SignHttpToken(hc.Req, js)

	hc.Body(js)

	var rsp types.TypeMeta
	if err := hc.ReplyJson(&rsp); err != nil {
		return err
	} else if rsp.Error != nil {
		return errors.New(rsp.Error.Message + ", " + rsp.Error.Code)
	} else if rsp.Kind != "MsgItem" {
		return errors.New("unknown error")
	}

	return nil
}
