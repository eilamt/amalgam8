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
	Actions     []FancyAction   `json:"actions,omitempty"`
}

type ActionInterface interface {
	GetType() string
}

type Route struct {
	Backends []Backend `json:"backends"`
}

type Backend struct {
	Name    string   `json:"name,omitempty"`
	Tags    []string `json:"tags"`
	Weight  float64  `json:"weight,omitempty"`
	Timeout float64  `json:"timeout,omitempty"`
	Retries int      `json:"retries,omitempty"` // FIXME: this BREAKS disabling retries by setting them to 0!
}

type FancyAction struct {
	internal   interface{}
	actionType string
}

func (a *FancyAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.internal)
}

func (a *FancyAction) UnmarshalJSON(data []byte) error {
	action := Action{}
	err := json.Unmarshal(data, &action)
	if err != nil {
		return err
	}

	a.actionType = action.Type

	switch action.Type {
	case "delay":
		delay := DelayAction{}
		if err = json.Unmarshal(data, &delay); err != nil {
			return err
		}
		a.internal = delay
	case "abort":
		abort := AbortAction{}
		if err = json.Unmarshal(data, &abort); err != nil {
			return err
		}
		a.internal = abort
	}
	return nil
}

func (a *FancyAction) GetType() string {
	return a.actionType
}

func (a *FancyAction) Internal() interface{} {
	return a.internal
}

type Action struct {
	Type string `json:"action"`
}

type DelayAction struct {
	Type        string
	Probability int      `json:"probability"`
	Tags        []string `json:"tags"`
	Duration    int      `json:"duration"`
}

type AbortAction struct {
	Type        string
	Probability int
	Tags        []string
	ReturnCode  int
}
