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

package config

import (
	"github.com/lessos/lessgo/data/rdo"
	"github.com/lessos/lessgo/data/rdo/base"
	"github.com/lessos/lessgo/utils"
)

func (c *ConfigCommon) DatabaseInstance() (*base.Client, error) {

	dc, err := rdo.NewClient("def", c.Database)
	if err != nil {
		return dc, err
	}

	ds, err := base.LoadDataSetFromString(databaseSchema)
	err = dc.Dialect.SchemaSync(c.Database.Dbname, ds)
	if err != nil {
		return dc, err
	}

	timenow := base.TimeNow("datetime")

	uid := utils.StringEncode16("sysadmin", 8)

	dc.Base.InsertIgnore("ids_role", map[string]interface{}{
		"rid":        1,
		"uid":        uid,
		"status":     1,
		"name":       "Administrator",
		"desc":       "Root System Administrator",
		"privileges": "1,2",
		"created":    timenow,
		"updated":    timenow,
	})
	dc.Base.InsertIgnore("ids_role", map[string]interface{}{
		"rid":        100,
		"uid":        uid,
		"status":     1,
		"name":       "Member",
		"desc":       "Universal Member",
		"privileges": "",
		"created":    timenow,
		"updated":    timenow,
	})

	dc.Base.InsertIgnore("ids_role", map[string]interface{}{
		"rid":        101,
		"uid":        uid,
		"status":     0,
		"name":       "Anonymous",
		"desc":       "Anonymous Member",
		"privileges": "",
		"created":    timenow,
		"updated":    timenow,
	})

	dc.Base.InsertIgnore("ids_instance", map[string]interface{}{
		"id":        "lessids",
		"uid":       uid,
		"status":    1,
		"app_id":    "lessids",
		"app_title": "less Identity Server",
		"version":   Version,
		"created":   timenow,
		"updated":   timenow,
	})

	dc.Base.InsertIgnore("ids_privilege", map[string]interface{}{
		"pid":       "1",
		"instance":  "lessids",
		"uid":       uid,
		"privilege": "user.admin",
		"desc":      "User Management",
		"created":   timenow,
	})
	dc.Base.InsertIgnore("ids_privilege", map[string]interface{}{
		"pid":       "2",
		"instance":  "lessids",
		"uid":       uid,
		"privilege": "sys.admin",
		"desc":      "System Management",
		"created":   timenow,
	})

	dc.Base.InsertIgnore("ids_sysconfig", map[string]interface{}{
		"key":     "service_name",
		"value":   "less Identity Service",
		"created": timenow,
		"updated": timenow,
	})
	dc.Base.InsertIgnore("ids_sysconfig", map[string]interface{}{
		"key":     "webui_banner_title",
		"value":   "Account Center",
		"created": timenow,
		"updated": timenow,
	})

	return dc, err
}
