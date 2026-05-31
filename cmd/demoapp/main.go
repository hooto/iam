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

package main

import (
	"embed"
	"flag"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/hooto/htoml4g/htoml"
	"github.com/hooto/httpsrv"

	"github.com/hooto/iam/v2/pkg/iamserver"
)

//go:embed dist
var embedDir embed.FS

var (
	configPath string

	cfg demoConfig
)

type demoConfig struct {
	HttpPort uint16                   `json:"http_port" toml:"http_port"`
	IAM      *iamserver.AppAuthConfig `json:"iam" toml:"iam"`
}

func main() {
	prefix := flag.String("prefix", "", "path prefix for config and data")
	flag.Parse()

	if err := loadConfig(*prefix); err != nil {
		slog.Error("config load failed", "error", err)
	}

	iamserver.AppConfig = cfg.IAM

	httpsrv.DefaultService.HandleModule("/demoapp", newUIModule())
	httpsrv.DefaultService.HandleModule("/demoapp/api", newAPIModule())

	httpsrv.DefaultService.Start()
}

func loadConfig(prefix string) error {
	if prefix == "" {
		prefix, _ = filepath.Abs(filepath.Dir(os.Args[0]) + "/..")
	}

	configPath = filepath.Join(prefix, "etc/demoapp_config.toml")

	if err := htoml.DecodeFromFile(configPath, &cfg); err != nil {
		if !os.IsNotExist(err) {
			slog.Warn("config decode warning", "error", err)
		}
	}

	if cfg.HttpPort == 0 {
		cfg.HttpPort = 3001
	}

	if cfg.IAM == nil {
		cfg.IAM = &iamserver.AppAuthConfig{}
	}

	cfg.IAM.SaveFunc = func() error {
		return htoml.EncodeToFile(cfg, configPath, nil)
	}

	httpsrv.DefaultService.Config.HttpPort = cfg.HttpPort

	return cfg.IAM.SaveFunc()
}

func newUIModule() *httpsrv.Module {
	mod := httpsrv.NewModule()

	distFS, err := fs.Sub(embedDir, "dist")
	if err != nil {
		panic(err)
	}

	mod.RegisterFileServer("/", "/", &spaFileSystem{base: http.FS(distFS)})

	return mod
}

func newAPIModule() *httpsrv.Module {
	mod := httpsrv.NewModule()
	mod.RegisterController(new(iamserver.UserAuth))
	return mod
}

type spaFileSystem struct {
	base http.FileSystem
}

func (sfs *spaFileSystem) Open(name string) (http.File, error) {
	name = strings.TrimPrefix(filepath.Clean(name), "/")

	if name == "" {
		name = "index.html"
	}

	f, err := sfs.base.Open(name)
	if err != nil {
		return sfs.base.Open("index.html")
	}
	return f, nil
}
