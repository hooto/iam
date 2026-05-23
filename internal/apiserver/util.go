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

package apiserver

import (
	"fmt"
	"strings"

	"github.com/hooto/httpsrv"

	"github.com/hooto/iam/v2/internal/data"
	"github.com/hooto/iam/v2/pkg/iamapi"
)

// remoteAddr extracts the IP address from a RemoteAddr string (host:port format).
func remoteAddr(ra string) string {
	if addridx := strings.Index(ra, ":"); addridx > 0 {
		return ra[:addridx]
	}
	return "127.0.0.1"
}

// userAuthDenyCheck checks if the user+IP combination has exceeded the auth
// deny threshold. Returns (denyCount, denyKey, error) where error is non-nil
// when the request should be rejected.
func userAuthDenyCheck(username string, req *httpsrv.Request) (int, []byte, error) {
	addr := remoteAddr(req.RemoteAddr)
	denyCount := 0
	denyKey := iamapi.NsUserAuthDeny(username, addr)
	if rs := data.Data.NewReader(denyKey).Exec(); rs.OK() {
		if denyCount = int(rs.Item().Int64Value()); denyCount >= 20 {
			return denyCount, denyKey, fmt.Errorf(
				"too many requests (%d), please try again in 1 day later", denyCount)
		}
	}
	return denyCount, denyKey, nil
}

// userAuthDenyIncr increments the deny counter for the given denyKey.
func userAuthDenyIncr(denyCount int, denyKey []byte) {
	data.Data.NewWriter(denyKey, []byte(fmt.Sprintf("%d", denyCount+1))).
		SetTTL(86400e3).Exec()
}
