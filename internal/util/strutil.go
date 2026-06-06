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

// Package util provides shared internal utility functions for the iam project.
package util

// containsSmallThreshold is the cutoff below which a brute-force nested loop
// is faster than allocating and probing a hash map.
const containsSmallThreshold = 16

// Contains checks whether slice a contains any element from slice b.
// It uses optimized strategies based on input sizes:
//   - Single-element b: linear scan of a.
//   - Small inputs (both <= 16): nested loop without allocation.
//   - Large inputs: hash-set built from the shorter slice for O(1) lookups.
func Contains(a, b []string) bool {
	if len(a) == 0 || len(b) == 0 {
		return false
	}

	// Fast path: single-element b avoids all overhead.
	if len(b) == 1 {
		s := b[0]
		for _, v := range a {
			if v == s {
				return true
			}
		}
		return false
	}

	// For small inputs, a simple nested loop is faster than map allocation.
	if len(a) <= containsSmallThreshold && len(b) <= containsSmallThreshold {
		for _, sa := range a {
			for _, sb := range b {
				if sa == sb {
					return true
				}
			}
		}
		return false
	}

	// For larger inputs, build a hash set from the shorter slice for O(1) lookups.
	if len(a) < len(b) {
		a, b = b, a
	}

	set := make(map[string]struct{}, len(b))
	for _, v := range b {
		set[v] = struct{}{}
	}

	for _, v := range a {
		if _, exists := set[v]; exists {
			return true
		}
	}

	return false
}
