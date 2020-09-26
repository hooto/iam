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
	"strings"
	"sync"
	"time"

	"github.com/hooto/hlog4g/hlog"
	"github.com/hooto/hmsg/go/hmsg/v1"
	"github.com/hooto/iam/data"
	"github.com/hooto/iam/iamapi"
	"github.com/lessos/lessgo/net/email"
)

var (
	msgMu           sync.Mutex
	msgQueuePending        = false
	msgQueueTimeout uint32 = 864000
	msgQueueCheckN  int64  = 0
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
		offset = iamapi.ObjKeyMsgQueue("")
		cutset = iamapi.ObjKeyMsgQueue("")
		limit  = 100
	)

	for {

		mailer, err := email.MailerPull("def")
		if err != nil {
			if msgQueueCheckN == 1 {
				hlog.Printf("warn", "mailer setup err %s", err.Error())
			}
			msgQueueCheckN += 1
			break
		}

		rs := data.Data.NewReader(nil).KeyRangeSet(offset, cutset).
			LimitNumSet(int64(limit)).Query()
		if !rs.OK() {
			hlog.Printf("info", "mailer scan err")
			break
		}

		for _, v := range rs.Items {

			var item hmsg.MsgItem
			if err := v.DataValue().Decode(&item, nil); err != nil {
				hlog.Printf("info", "mailer err %s", err.Error())
				continue
			}

			toMail := []string{}
			if item.ToEmail != "" {
				toMail = strings.Split(item.ToEmail, ";")
			}

			if len(toMail) == 0 {

				if u := data.UserGet(item.ToUser); u != nil {

					if u.Type == iamapi.UserTypeGroup {

						for _, ov := range u.Owners {
							if ou := data.UserGet(ov); ou != nil {
								if iamapi.EmailRE.MatchString(ou.Email) {
									toMail = append(toMail, ou.Email)
								}
							}
						}

					} else if iamapi.EmailRE.MatchString(u.Email) {
						toMail = append(toMail, u.Email)
					}
				}
			}

			/**
			if rs := data.Data.NewReader(iamapi.ObjKeyUser(item.ToUser)).Query(); rs.OK() {
				var userLogin iamapi.User
				rs.Decode(&userLogin)

				if userLogin.Type == iamapi.UserTypeGroup {
					//
				} else {

					if iamapi.EmailRE.MatchString(userLogin.Email) {
						// item.ToEmail = userLogin.Email
						toMail = append(toMail, userLogin.Email)
					}
				}
			}
			*/

			item.Updated = uint32(time.Now().Unix())

			if len(toMail) > 0 {

				item.ToEmail = strings.Join(toMail, ";")

				if err := msgPost(mailer, item); err != nil {
					item.Retry += 1
					if item.Retry < 10 {
						continue
					}
					item.Action = hmsg.MsgActionPostTimeout
				} else {
					item.Action = hmsg.MsgActionPostOK
				}
			} else {
				item.Action = hmsg.MsgActionPostError
			}

			if iamapi.OpActionAllow(item.Action, hmsg.MsgActionPostOK) ||
				iamapi.OpActionAllow(item.Action, hmsg.MsgActionPostError) ||
				iamapi.OpActionAllow(item.Action, hmsg.MsgActionPostTimeout) {
				if item.Posted < 1 {
					item.Posted = item.Updated
				}
				if rs := data.Data.NewWriter(iamapi.ObjKeyMsgSent(item.SentId()), item).Commit(); rs.OK() {
					data.Data.NewWriter(v.Meta.Key, nil).ModeDeleteSet(true).Commit()
					hlog.Printf("info", "mailer post %s, to %s, retry %d, ok", item.Id, item.ToEmail, item.Retry)
				}
			} else {
				data.Data.NewWriter(v.Meta.Key, item).Commit()
				hlog.Printf("warn", "mailer post %s, retry %d", item.Id, item.ToEmail, item.Retry)
			}
		}

		if !rs.Next {
			break
		}
	}
}

func msgPost(mailer *email.Mailer, msg hmsg.MsgItem) error {
	return mailer.SendMail(msg.ToEmail, msg.Title, msg.Body+"\n")
}
