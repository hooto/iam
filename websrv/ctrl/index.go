// Copyright 2014 lessos Authors, All rights reserved.
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
	"code.hooto.com/lessos/iam/config"
	"code.hooto.com/lessos/iam/iamclient"
	"github.com/lessos/lessgo/httpsrv"
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

	c.RenderString(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>` + config.Config.WebUiBannerTitle + `</title>
  <link rel="shortcut icon" type="image/x-icon" href="/iam/~/iam/img/favicon.ico">
  <script src="/iam/~/lessui/js/sea.js"></script>
  <script src="/iam/~/iam/js/main.js"></script>
  <script type="text/javascript">
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
