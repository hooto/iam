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
	"github.com/lessos/bigtree/btagent"
	"github.com/lessos/bigtree/btapi"
)

var (
	BtAgent *btagent.Agent
)

func init() {
	BtAgent, _ = btagent.NewAgent(btapi.DataAccessConfig{
		PathPoint: "/sys/ids",
	})
}
