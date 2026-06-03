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
