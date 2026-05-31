// Copyright 2014 Eryx <evorui at gmail dot com>, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses-2.0
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
	"log"
	"log/slog"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hooto/httpsrv"
	"github.com/sysinner/incore/v2/pkg/inlog"
	"github.com/sysinner/incore/v2/pkg/signals"

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

	httpsrv.DefaultService.Config.HttpPort = config.Config.HttpPort

	httpsrv.DefaultService.HandleModule("/iam/v2/auth", auth.NewModule())
	httpsrv.DefaultService.HandleModule("/iam/v2/user", user.NewModule())
	httpsrv.DefaultService.HandleModule("/iam/v2/open", open.NewModule())
	httpsrv.DefaultService.HandleModule("/iam/v2/admin", admin.NewModule())

	httpsrv.DefaultService.HandleModule("/iam", newUIModule())

	signals.Go(func() {
		httpsrv.DefaultService.Start()
	}, func() {
		httpsrv.DefaultService.Stop()
		data.Data.Close()
		slog.Info("server kill done")
	})

	signals.Wait()
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

// spaFileSystem 包装 http.FileSystem，增加 SPA fallback
type spaFileSystem struct {
	base http.FileSystem
}

// 静态资源扩展名，这些文件 404 时不 fallback 到 index.html
var assetExts = map[string]bool{
	".js": true, ".mjs": true, ".css": true,
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".svg": true, ".ico": true, ".webp": true,
	".woff": true, ".woff2": true, ".ttf": true, ".eot": true, ".otf": true,
	".map": true, ".json": true, ".xml": true, ".txt": true,
	".wasm": true,
}

func (sfs *spaFileSystem) Open(name string) (http.File, error) {

	name = strings.TrimPrefix(filepath.Clean(name), "/")

	// 根路径返回 index.html
	if name == "" || name == "/" {
		name = "index.html"
	}

	f, err := sfs.base.Open(name)
	if err == nil {
		return f, nil
	}

	// SPA fallback: 非 assets 扩展名的路径返回 index.html
	ext := strings.ToLower(filepath.Ext(name))
	if !assetExts[ext] {
		return sfs.base.Open("index.html")
	}

	return nil, err
}
