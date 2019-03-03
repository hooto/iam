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

package ctrl

import (
	"github.com/hooto/httpsrv"
	"github.com/hooto/iam/config"
	"github.com/hooto/iam/iamclient"
)

type Index struct {
	*httpsrv.Controller
}

func (c Index) IndexAction() {

	if !iamclient.SessionIsLogin(c.Session) {
		c.Redirect("/iam/service/login?continue=/iam")
		return
	}

	if c.Params.Get(iamclient.AccessTokenKey) != "" {
		c.Redirect("/iam")
		return
	}

	c.AutoRender = false
	c.Response.Out.Header().Set("Cache-Control", "no-cache")

	c.RenderString(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>` + config.Config.WebUiBannerTitle + `</title>
  <link rel="shortcut icon" type="image/x-icon" href="/iam/~/iam/img/iam-s2-32.png">
  <script src="/iam/~/lessui/js/sea.js?v=` + config.VersionHash + `"></script>
  <script src="/iam/~/iam/js/main.js?v=` + config.VersionHash + `"></script>
  <script type="text/javascript">
	iam.version = "` + config.VersionHash + `";
	iam.lang = "` + c.Request.Locale + `";
    window.onload = iam.Boot();
  </script>
</head>
<body id="body-content">
  <style>
  ._iam_loading {
    margin: 0;
    padding: 30px 40px;
    font-size: 48px;
    color: #000;
  }
  </style>
  <div class="_iam_loading">loading</div>
</body>
</html>`)
}
