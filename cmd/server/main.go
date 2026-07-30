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
	"context"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hooto/httpsrv/v2"
	"github.com/hooto/httpsrv/v2/middleware/static"
	"github.com/sysinner/innerstack/v2/pkg/inlog"
	"github.com/sysinner/innerstack/v2/pkg/signals"

	"github.com/hooto/iam/v2/internal/apiserver/admin"
	"github.com/hooto/iam/v2/internal/apiserver/auth"
	"github.com/hooto/iam/v2/internal/apiserver/open"
	"github.com/hooto/iam/v2/internal/apiserver/user"
	"github.com/hooto/iam/v2/internal/config"
	"github.com/hooto/iam/v2/internal/data"
)

//go:embed dist
var embedDir embed.FS

var (
	flagPrefix = flag.String("prefix", "", "the prefix folder path for config/data")
)

func main() {

	inlog.Setup()

	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()
	if err := config.Setup(*flagPrefix); err != nil {
		log.Fatal(err)
	}

	if err := data.Setup(); err != nil {
		log.Fatal(err)
	}

	app := httpsrv.New()

	auth.Register(app.Group("/iam/v2/auth"))
	user.Register(app.Group("/iam/v2/user"))
	open.Register(app.Group("/iam/v2/open"))
	admin.Register(app.Group("/iam/v2/admin"))

	{
		distFS, err := fs.Sub(embedDir, "dist")
		if err != nil {
			log.Fatal(err)
		}

		spa := static.New(".", static.Config{
			FS: &spaFileSystem{fs: distFS},
		})
		app.Get("/iam{*}", spa)
	}

	signals.Go(func() {
		if err := app.Run(fmt.Sprintf(":%d", config.Config.HttpPort)); err != nil {
			slog.Error("httpsrv run error", "error", err)
		}
	}, func() {
		_ = app.Shutdown(context.Background())
		data.Data.Close()
		slog.Info("server kill done")
	})

	signals.Wait()
}

// spaFileSystem 包装 fs.FS，增加 SPA fallback
type spaFileSystem struct {
	fs fs.FS
}

// 静态资源扩展名，这些文件 404 时不 fallback 到 index.html
var assetExts = map[string]bool{
	".js": true, ".mjs": true, ".css": true,
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".svg": true, ".ico": true, ".webp": true,
	".woff": true, ".woff2": true, ".ttf": true, ".eot": true, ".otf": true,
	".map": true, ".json": true, ".xml": true, ".txt": true,
	".wasm": true,
}

func (sfs *spaFileSystem) Open(name string) (fs.File, error) {

	name = strings.TrimPrefix(filepath.Clean(name), "/")

	// 根路径返回 index.html
	if name == "" || name == "/" || name == "." {
		name = "index.html"
	}

	f, err := sfs.fs.Open(name)
	if err == nil {
		return f, nil
	}

	// SPA fallback: 非 assets 扩展名的路径返回 index.html
	ext := strings.ToLower(filepath.Ext(name))
	if !assetExts[ext] {
		return sfs.fs.Open("index.html")
	}

	return nil, err
}
