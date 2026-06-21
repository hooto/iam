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

package iamserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/sysinner/incore/v2/pkg/inauth"
)

// iamPost sends a JSON POST to an IAM endpoint and decodes the response.
func iamPost(baseUrl string, endpoint string, auth string, reqBody interface{}, rspBody interface{}) error {

	var (
		iamURL  = urlJoinPath(baseUrl, endpoint)
		body, _ = json.Marshal(reqBody)
	)

	req, err := http.NewRequest("POST", iamURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if auth != "" {
		req.Header.Set("Authorization", "Bearer "+auth)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if rspBody != nil {
		if err := json.Unmarshal(respBytes, rspBody); err != nil {
			return err
		}
	}

	type statusHolder struct {
		Status inauth.ServiceStatus `json:"status"`
	}
	var holder statusHolder
	json.Unmarshal(respBytes, &holder)
	if holder.Status.Code != "200" {
		return errors.New(holder.Status.Message)
	}

	return nil
}

func urlJoinPath(basePath, addon string) string {
	s, _ := url.JoinPath(basePath, addon)
	return s
}

// isSameSite reports whether the given URL belongs to the same host as the
// incoming HTTP request. It validates Referer-based redirects to mitigate
// open-redirect attacks. Relative URLs with an empty host are treated as
// same-origin and considered safe.
func isSameSite(u *url.URL, r *http.Request) bool {
	if u == nil || r == nil {
		return false
	}
	if u.Host == "" {
		return true
	}
	return strings.EqualFold(u.Host, r.Host)
}

// isBrowserNavigation reports whether the request looks like a full-page
// browser navigation rather than an AJAX/XHR call.
//
// It prefers the modern Sec-Fetch-Mode hint, but falls back to classic
// browser signals (Upgrade-Insecure-Requests + Accept: text/html) because
// Sec-Fetch-* headers may be stripped by proxies or missing in some clients.
// An explicit X-Requested-With: XMLHttpRequest always means an AJAX call.
func isBrowserNavigation(r *http.Request) bool {
	if r == nil {
		return false
	}

	// explicit AJAX indicator (jQuery/axios etc.)
	if strings.EqualFold(r.Header.Get("X-Requested-With"), "XMLHttpRequest") {
		return false
	}

	// modern standard, preferred when present
	if r.Header.Get("Sec-Fetch-Mode") == "navigate" {
		return true
	}

	// classic browser navigation hints
	if r.Header.Get("Upgrade-Insecure-Requests") == "1" &&
		strings.Contains(strings.ToLower(r.Header.Get("Accept")), "text/html") {
		return true
	}

	return false
}
