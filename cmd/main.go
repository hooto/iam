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

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/lessos/lessgo/httpsrv"

	"github.com/lessos/lessids/config"
	"github.com/lessos/lessids/websrv/ctrl"
	"github.com/lessos/lessids/websrv/v1"
)

var flagPrefix = flag.String("prefix", "", "the prefix folder path")

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	//
	flag.Parse()
	if err := config.Init(*flagPrefix); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//
	httpsrv.GlobalService.Config.HttpPort = config.Config.Port
	httpsrv.GlobalService.Config.InstanceID = "lessids"
	httpsrv.GlobalService.Config.LessIdsServiceUrl = fmt.Sprintf("http://127.0.0.1:%d/ids", config.Config.Port)

	httpsrv.GlobalService.ModuleRegister("/ids/v1", v1.NewModule())
	httpsrv.GlobalService.ModuleRegister("/ids/", ctrl.NewModule())

	//
	httpsrv.GlobalService.Start()
}
