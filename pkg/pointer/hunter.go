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
	"fmt"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/golang/glog"

	"log"
	"time"

	"google.golang.org/grpc"
	pb "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

func Hunt() error {
	q := &Queue{nodes: make([]*LogForQueue, 3)}
	err := containerRuntime(&addr)
	if err != nil {
		glog.Error(err)
		return err
	}
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(2*time.Second))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	runtimeClient = pb.NewRuntimeServiceClient(conn)
	go func() {
		for {
			if q.count > 0 {
				l := q.Pop()
				//err, cpu, mem := getContainerStats(l.ContainerId, runtimeClient)
				//if err != nil {
				//	fmt.Println(err)
				//}
				fmt.Println(l.Content)
			}
		}
	}()
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	if err := filepath.Walk("/var/log/containers/", watchDir); err != nil {
		fmt.Println("ERROR", err)
	}
	watcher.Add("/var/log/containers/")
	done := make(chan bool)

	go dispatchEvents(q)
	//fmt.Println("debug 1")
	//go func() {
	//	fmt.Println("debug 2")
	//	for {
	//		var logsfiles []string
	//		files, err := ioutil.ReadDir("/var/log/containers/")
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//		for _, file := range files {
	//			logsfiles = append(logsfiles, file.Name())
	//		}
	//		logs <- logsfiles
	//		time.Sleep(time.Second * 4)
	//	}
	//}()
	<-done
	return nil

}
