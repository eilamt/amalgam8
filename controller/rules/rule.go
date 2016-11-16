// Copyright 2016 IBM Corporation
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package rules

import "encoding/json"

// Rule represents an individual rule.
type Rule struct {
	ID          string          `json:"id"`
	Priority    int             `json:"priority"`
	Tags        []string        `json:"tags,omitempty"`
	Destination string          `json:"destination"`
	Match       json.RawMessage `json:"match,omitempty"`
	Route       *Route          `json:"route,omitempty"`
	Actions     json.RawMessage `json:"actions,omitempty"`
}

type Route struct {
	Backends []Backend `json:"backends"`
}

type Backend struct {
	Name    string   `json:"name"`
	Tags    []string `json:"tags,omitempty"`
	Weight  float64  `json:"weight,omitempty"`
	Timeout float64  `json:"timeout,omitempty"`
	Retries int      `json:"retries,omitempty"` // FIXME: this BREAKS disabling retries by setting them to 0!
}
