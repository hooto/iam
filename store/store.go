// Copyright 2015 lessOS.com, All rights reserved.
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

package store

import (
	"time"

	"github.com/lessos/bigtree/btagent"
	"github.com/lessos/bigtree/btapi"
)

var (
	Ready   bool
	BtAgent btagent.ApiAgent
)

func init() {

	BtAgent = btagent.ApiAgent{}

	// BtAgent, _ = btagent.NewAgent(btapi.DataAccessConfig{
	// 	PathPoint: "/sys/ids",
	// })

	go func() {

		for {

			BtAgent.ObjectSet("global/ids/ids-test", "test", &btapi.ObjectWriteOptions{
				Ttl: 3000,
			})

			if rs := BtAgent.ObjectGet("global/ids/ids-test"); rs.Data == "test" {

				in := InitNew{}
				in.Init()

				Ready = true

				break
			}

			time.Sleep(1e9)
		}
	}()
}
