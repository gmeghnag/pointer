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
	"fmt"
	"io/ioutil"

	"github.com/gmeghnag/pointer/types"
	"github.com/gmeghnag/pointer/types/vars"
	"github.com/golang/glog"
	yaml "gopkg.in/yaml.v2"
)

func readConf(filename string) (*types.Conf, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &types.Conf{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}

	return c, nil
}

func parseConf(types.Conf) error {
	c, err := readConf(vars.ConfigPath)
	if err != nil {
		glog.Error(err)
		return err
	}
	vars.Config = *c
	return nil
}
