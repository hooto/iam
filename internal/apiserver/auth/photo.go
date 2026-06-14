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

package auth

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/hooto/httpsrv"
	"github.com/hooto/iam/v2/internal/data"
	"github.com/hooto/iam/v2/pkg/iamapi"
)

// emptyPhoto is the default avatar SVG returned when the user has no photo.
const emptyPhoto = `<svg xmlns="http://www.w3.org/2000/svg" width="128" height="128" fill="currentColor" class="bi bi-person" viewBox="0 0 16 16">
  <path d="M8 8a3 3 0 1 0 0-6 3 3 0 0 0 0 6m2-3a2 2 0 1 1-4 0 2 2 0 0 1 4 0m4 8c0 1-1 1-1 1H3s-1 0-1-1 1-4 6-4 6 3 6 4m-1-.004c-.001-.246-.154-.986-.832-1.664C11.516 10.68 10.289 10 8 10s-3.516.68-4.168 1.332c-.678.678-.83 1.418-.832 1.664z"/>
</svg>`

// dataURIPrefix is the expected prefix for a base64-encoded data URI.
const dataURIPrefix = "data:"

// Photo returns the profile photo for the given username.
// If the user has a photo stored as a data URI (e.g. "data:image/png;base64,..."),
// it is decoded and served with the correct Content-Type.
// Otherwise, a default SVG avatar is returned.
func Photo(ctx httpsrv.Ctx) error {

	username := strings.ToLower(ctx.Params().Value("username"))
	if username == "" {
		username = strings.ToLower(ctx.Params().Value("action"))
		if username == "" {
			return writeEmptyPhoto(ctx)
		}
	}

	if err := iamapi.UsernameValid(username); err != nil {
		return writeEmptyPhoto(ctx)
	}

	var profile iamapi.UserProfile
	rs := data.Data.NewReader(iamapi.NsUserProfile(username)).Exec()
	if !rs.OK() {
		return writeEmptyPhoto(ctx)
	}
	rs.Item().JsonDecode(&profile)

	if profile.Photo == "" || !strings.HasPrefix(profile.Photo, dataURIPrefix) {
		return writeEmptyPhoto(ctx)
	}

	contentType, imageData, err := decodeDataURI(profile.Photo)
	if err != nil {
		return writeEmptyPhoto(ctx)
	}

	ctx.Response().Header().Set("Content-Type", contentType)
	ctx.Response().Header().Set("Cache-Control", "private, max-age=3600")
	ctx.Response().Write(imageData)
	return nil
}

// writeEmptyPhoto responds with the default SVG avatar.
func writeEmptyPhoto(ctx httpsrv.Ctx) error {
	ctx.Response().Header().Set("Content-Type", "image/svg+xml")
	ctx.Response().Header().Set("Cache-Control", "public, max-age=300")
	ctx.Response().Write([]byte(emptyPhoto))
	return nil
}

// decodeDataURI parses a data URI of the form "data:<mediatype>;base64,<data>"
// and returns the MIME type and decoded bytes.
func decodeDataURI(dataURI string) (string, []byte, error) {
	// data:<mediatype>;base64,<data>
	s := strings.TrimPrefix(dataURI, dataURIPrefix)

	commaIdx := strings.Index(s, ",")
	if commaIdx < 0 {
		return "", nil, errors.New("invalid data URI: missing comma separator")
	}

	meta := s[:commaIdx]
	encoded := s[commaIdx+1:]

	// Extract media type from meta (e.g. "image/png;base64")
	mediaType := "application/octet-stream"
	if meta != "" {
		parts := strings.SplitN(meta, ";", 2)
		if parts[0] != "" {
			mediaType = parts[0]
		}
	}

	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		// Try URL-safe encoding as fallback
		data, err = base64.URLEncoding.DecodeString(encoded)
		if err != nil {
			return "", nil, err
		}
	}

	return mediaType, data, nil
}
