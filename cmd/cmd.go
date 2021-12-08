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

package cmd

import (
	"flag"
	"os"

	"github.com/gmeghnag/pointer/pkg/pointer"
	"github.com/gmeghnag/pointer/types/vars"
)

func parseFlags() {
	flag.StringVar(&vars.ConfigPath, "config", "/pointer/etc/pointer.conf", "configuration file path.")
	flag.Parse()
}

func Init() {
	parseFlags()
	//fmt.Printf("%q", vars.Config)
	err := parseConf(vars.Config)
	if err != nil {
		os.Exit(1)
	}
	//fmt.Printf("%q", vars.Config)
	err = pointer.Hunt()
	if err != nil {
		os.Exit(1)
	}
}
