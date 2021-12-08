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
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/go-units"
	pb "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

func filesInPath(path string) ([]string, error) {
	var files []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, "log") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (l LogResource) DiffString(newSize int64) ([]byte, error) {
	size := l.NewSize - l.OldSize
	if size < 0 {
		return nil, fmt.Errorf("Old size is bigger than the new log size old=%q, new=%q file=%q", l.OldSize, l.NewSize, l.Path)
	}
	res := make([]byte, size, size)
	openFile, _ := os.Open(l.Path)
	defer openFile.Close()
	_, err := openFile.ReadAt(res, l.OldSize-1)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (l LogResource) HasOldSize() bool {
	if l.OldSize != 0 {
		return true
	}
	return false
}

func watchDir(path string, fi os.FileInfo, err error) error {
	if strings.HasSuffix(path, "log") && strings.HasPrefix(path, "/var/log/containers/") {
		file, err := os.Stat(path)
		if err != nil {
			fmt.Println("ERROR", err)
			return err
		}
		fileSize := file.Size()
		LogResource := &LogResource{Path: path, OldSize: 0, NewSize: fileSize}
		dict[path] = *LogResource
		return watcher.Add(path)
	}
	return nil
}

func getContainerId(path string) (error, string) {
	path = strings.TrimRight(path, ".log")
	path = strings.TrimLeft(path, "/var/log/containers/")
	id := strings.Split(path, "-")
	if len(id[len(id)-1]) != 64 {
		return fmt.Errorf("malformed file name: %v", path), ""
	}
	return nil, id[len(id)-1]
}

func getContainerStats(containerId string, runtimeServiceClient pb.RuntimeServiceClient) (error, uint64, string) {
	filter := &pb.ContainerStatsFilter{}
	filter.Id = containerId
	m, err := runtimeServiceClient.ListContainerStats(context.Background(), &pb.ListContainerStatsRequest{
		Filter: filter,
	})
	//https://github.com/kubernetes-sigs/cri-tools/blob/a989838814805e1053688e0e94adf13d60e716c6/vendor/k8s.io/cri-api/pkg/apis/runtime/v1alpha2/api.pb.go#L6854-6865
	if err != nil {
		log.Fatal(err)
	}

	stats := m.GetStats()
	if len(stats) == 1 {
		mem := stats[0].GetMemory().GetWorkingSetBytes().GetValue()
		return nil, stats[0].Cpu.UsageCoreNanoSeconds.Value, units.HumanSize(float64(mem))
	}
	return fmt.Errorf("Container stats return more than one container"), 0, ""
}

func fileExist(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func containerRuntime(addr *string) error {
	if fileExist("/var/run/crio/crio.sock") {
		*addr = "unix:///var/run/crio/crio.sock"
		return nil
	} else if fileExist("/var/run/containerd/containerd.sock") {
		*addr = "unix:///var/run/containerd/containerd.sock"
		return nil
	} else {
		return fmt.Errorf("runtime socket not found")
	}
}
