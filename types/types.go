/*
Copyright © 2021 The Pointer Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package types

type Conf struct {
	Skip struct {
		Namespaces struct {
			Startswith string   `json:"startswith"`
			Endswith   string   `json:"endswith"`
			Exact      []string `json:"exact"`
		} `json:"namespaces"`
	} `json:"skip"`
	Include struct {
		Namespaces struct {
			Startswith string   `json:"startswith"`
			Endswith   string   `json:"endswith"`
			Exact      []string `json:"exact"`
		} `json:"namespaces"`
	} `json:"include"`
}
