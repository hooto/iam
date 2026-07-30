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
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/hooto/htoml4g/htoml"
	"github.com/hooto/httpsrv/v2"
	"github.com/hooto/httpsrv/v2/middleware/static"
	"github.com/sysinner/innerstack/v2/pkg/inlog"

	"github.com/hooto/iam/v2/pkg/iamapi"
	"github.com/hooto/iam/v2/pkg/iamserver"
)

//go:embed dist
var embedDir embed.FS

var (
	configPath string

	cfg demoConfig

	version = "1.0.0"
)

var permissions = []*iamapi.AppPermission{
	{
		Permission: "list",
		Roles:      []string{iamapi.Role_User, iamapi.Role_Guest},
		Summary:    "List access",
	},
	{
		Permission: "read",
		Roles:      []string{iamapi.Role_User, iamapi.Role_Guest},
		Summary:    "Read access",
	},
	{
		Permission: "write",
		Roles:      []string{iamapi.Role_User},
		Summary:    "Write access",
	},
	{
		Permission: "delete",
		Roles:      []string{iamapi.Role_User},
		Summary:    "Delete access",
	},
}

type demoConfig struct {
	HttpPort uint16                   `json:"http_port" toml:"http_port"`
	IAM      *iamserver.AppAuthConfig `json:"iam" toml:"iam"`
}

func main() {

	inlog.Setup()

	prefix := flag.String("prefix", "", "path prefix for config and data")
	flag.Parse()

	if err := loadConfig(*prefix); err != nil {
		slog.Error("config load failed", "error", err)
	}

	if err := iamserver.AppVerifier.Setup(cfg.IAM); err != nil {
		slog.Error("app auth verifier setup failed", "error", err)
	} else if err = iamserver.AppVerifier.Update(&iamapi.AppInstance{
		ID:          cfg.IAM.AppId,
		Version:     version,
		Status:      1,
		Permissions: permissions,
	}); err != nil {
		slog.Error("app auth update failed", "error", err)
	}

	app := httpsrv.New()

	api := app.Group("/demoapp/api")
	iamserver.RegisterRoutes(api)

	{
		distFS, err := fs.Sub(embedDir, "dist")
		if err != nil {
			slog.Error("fs setup failed", "error", err)
		}
		spa := static.New(".", static.Config{
			FS: &spaFileSystem{fs: distFS},
		})
		app.Get("/demoapp{*}", spa)
	}

	if err := app.Run(fmt.Sprintf(":%d", cfg.HttpPort)); err != nil {
		slog.Error("httpsrv run error", "error", err)
	}
}

func loadConfig(prefix string) error {
	if prefix == "" {
		prefix, _ = filepath.Abs(filepath.Dir(os.Args[0]) + "/..")
	}

	configPath = filepath.Join(prefix, "etc/demoapp_config.toml")

	if err := htoml.DecodeFromFile(configPath, &cfg); err != nil {
		if !os.IsNotExist(err) {
			slog.Warn("config decode warning", "error", err)
		} else {
			return err
		}
	}

	if cfg.HttpPort == 0 {
		cfg.HttpPort = 3001
	}

	{
		if cfg.IAM == nil {
			cfg.IAM = &iamserver.AppAuthConfig{}
		}

		cfg.IAM.Flush = func() error {
			return htoml.EncodeToFile(cfg, configPath, nil)
		}
	}

	return nil
}

type spaFileSystem struct {
	fs fs.FS
}

func (sfs *spaFileSystem) Open(name string) (fs.File, error) {
	name = strings.TrimPrefix(filepath.Clean(name), "/")

	if name == "" || name == "/" || name == "." {
		name = "index.html"
	}

	f, err := sfs.fs.Open(name)
	if err != nil {
		return sfs.fs.Open("index.html")
	}
	return f, nil
}
