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

package pointer

import (
	"github.com/fsnotify/fsnotify"
	pb "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

var addr string

// "unix:///var/run/containerd/containerd.sock"
// "unix:///var/run/crio/crio.sock"

var watcher *fsnotify.Watcher
var dict = map[string]LogResource{}
var runtimeClient pb.RuntimeServiceClient

type LogResource struct {
	OldSize     int64
	NewSize     int64
	Container   string
	ContainerId string
	Path        string
	Pod         string
	Namespace   string
}

type LogForQueue struct {
	OldSize     int64
	NewSize     int64
	Container   string
	ContainerId string
	Path        string
	Pod         string
	Namespace   string
	Content     string
}
