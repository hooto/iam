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

package worker

import (
	"sync"
	"time"

	"github.com/hooto/hlog4g/hlog"
	"github.com/hooto/iam/iamapi"
	"github.com/hooto/iam/store"
	"github.com/lessos/lessgo/encoding/json"
	"github.com/lessos/lessgo/net/email"
)

var (
	msgMu           sync.Mutex
	msgQueuePending        = false
	msgQueueTimeout uint32 = 864000
)

func MsgQueueRefresh() {

	msgMu.Lock()
	if msgQueuePending {
		msgMu.Unlock()
		return
	}
	msgQueuePending = true
	msgMu.Unlock()

	defer func() {
		msgQueuePending = false
	}()

	var (
		offset = iamapi.DataMsgQueue("")
		cutset = iamapi.DataMsgQueue("zzzz")
		limit  = 100
	)

	for {

		mailer, err := email.MailerPull("def")
		if err != nil {
			hlog.Printf("warn", "mailer setup err %s", err.Error())
			break
		}

		rs := store.Data.KvScan(offset, cutset, limit)
		if !rs.OK() {
			hlog.Printf("info", "mailer scan err")
			break
		}

		rss := rs.KvList()

		for _, v := range rss {

			var item iamapi.MsgItem
			if err := v.Decode(&item); err != nil {
				hlog.Printf("info", "mailer err %s", err.Error())
				continue
			}

			if rs := store.Data.KvProgGet(iamapi.DataUserKey(item.ToUser)); rs.OK() {
				var userLogin iamapi.User
				rs.Decode(&userLogin)
				if iamapi.UserEmailRe2.MatchString(userLogin.Email) {
					item.ToEmail = userLogin.Email
				}
			}

			item.Updated = uint32(time.Now().Unix())

			if item.ToEmail != "" {
				if err := msgPost(mailer, item); err != nil {
					item.Retry += 1
					if item.Retry < 10 {
						continue
					}
					item.Action = iamapi.MsgActionPostTimeout
				} else {
					item.Action = iamapi.MsgActionPostOK
				}
			} else {
				item.Action = iamapi.MsgActionPostError
			}

			if iamapi.OpActionAllow(item.Action, iamapi.MsgActionPostOK) ||
				iamapi.OpActionAllow(item.Action, iamapi.MsgActionPostError) ||
				iamapi.OpActionAllow(item.Action, iamapi.MsgActionPostTimeout) {
				if item.Posted < 1 {
					item.Posted = item.Updated
				}
				js, _ := json.Encode(item, "")
				if rs := store.Data.KvPut(iamapi.DataMsgSent(item.SentId()), js, nil); rs.OK() {
					store.Data.KvDel(iamapi.DataMsgQueue(item.Id))
					hlog.Printf("info", "mailer post %s, to %s, retry %d, ok", item.Id, item.ToEmail, item.Retry)
				}
			} else {
				js, _ := json.Encode(item, "")
				store.Data.KvPut(iamapi.DataMsgQueue(item.Id), js, nil)
				hlog.Printf("warn", "mailer post %s, retry %d", item.Id, item.ToEmail, item.Retry)
			}
		}

		if len(rss) < limit {
			break
		}
	}
}

func msgPost(mailer *email.Mailer, msg iamapi.MsgItem) error {
	return mailer.SendMail(msg.ToEmail, msg.Title, msg.Body+"\n")
}
