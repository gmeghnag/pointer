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
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/golang/glog"
)

func dispatchEvents(q *Queue) {
	for {
		select {
		case event := <-watcher.Events:
			//switch {
			//case event.Op&fsnotify.Write == fsnotify.Write:
			//	log.Printf("Write:  %s: %s", event.Op, event.Name)
			//case event.Op&fsnotify.Create == fsnotify.Create:
			//	log.Printf("Create: %s: %s", event.Op, event.Name)
			//case event.Op&fsnotify.Remove == fsnotify.Remove:
			//	log.Printf("Remove: %s: %s", event.Op, event.Name)
			//case event.Op&fsnotify.Rename == fsnotify.Rename:
			//	log.Printf("Rename: %s: %s", event.Op, event.Name)
			//case event.Op&fsnotify.Chmod == fsnotify.Chmod:
			//	log.Printf("Chmod:  %s: %s", event.Op, event.Name)
			//}
			go procesEvent(event, q)
			//err := processEvent(event)
			//if err != nil {
			//	fmt.Println("ERROR", err)
			//}
		case err := <-watcher.Errors:
			fmt.Println("ERROR", err)
		}
	}
}

//pointer_openshift-image-registry_pointer-44fc621d3bd6744a88200361703390ce81384899e55b2cc8e01958891c7bf962.log
// <pod_name>_<namespace>_<container_name>-<container_id>.log

func parseLogPath(path string) (err error, container_id string) {
	path = strings.TrimRight(path, ".log")
	path = strings.TrimLeft(path, "/var/log/containers/")
	id := path[len(path)-65:] // -<64chars_id>
	if id[:1] != "-" {
		glog.Warning("malformed path: %v", path)
		return fmt.Errorf("malformed path: %v", path), ""
	}
	return nil, id[1:]
}

func processEvent(event fsnotify.Event) error {
	if event.Op.String() == "WRITE" {
		//fmt.Println(event.Name, event.Op, event.Op.String() == "WRITE")
		//err, containerId := parseLogPath(event.Name)
		//if err != nil {
		//	fmt.Printf(err.Error())
		//	return err
		//}
		logFile := dict[event.Name]
		file, err := os.Stat(logFile.Path)
		if err != nil {
			fmt.Printf(err.Error())
			return err
		}
		newfileSize := file.Size()
		logFile.OldSize = logFile.NewSize
		logFile.NewSize = newfileSize
		dict[event.Name] = LogResource{Path: event.Name, OldSize: logFile.OldSize, NewSize: logFile.NewSize}
		lbytes, err := logFile.DiffString(newfileSize)
		if err != nil {
			fmt.Println("ERROR", err)
			return err
		}
		err, id := getContainerId(event.Name)
		if err != nil {
			fmt.Println("ERROR", err)
			return err
		}
		text := ""
		if len(lbytes) > 0 {
			lbytes = lbytes[1:]
			text = string(lbytes)
		}

		err, cpu, mem := getContainerStats(id, runtimeClient)
		if err != nil {
			fmt.Println("ERROR", err)
			return err
		}

		fmt.Println(id, cpu, mem, text)

	}
	if event.Op.String() == "CREATE" {
		err := watcher.Add(event.Name)
		file, err := os.Stat(event.Name)
		if err != nil {
			fmt.Println("ERROR", err)
			return err
		}
		fileSize := file.Size()
		logFile := &LogResource{Path: event.Name, OldSize: 0, NewSize: fileSize}
		dict[event.Name] = *logFile
		if err != nil {
			fmt.Println("ERROR", err)
			return err
		}
	}
	if event.Op.String() == "REMOVE" {
		err := watcher.Remove(event.Name)
		if err != nil {
			fmt.Println("ERROR", err)
			return err
		}
	}
	return nil
}

func procesEvent(event fsnotify.Event, q *Queue) {
	if event.Op.String() == "WRITE" {
		logFile := dict[event.Name]
		file, err := os.Stat(logFile.Path)
		if err != nil {
			fmt.Printf(err.Error())
			return
		}
		newfileSize := file.Size()
		logFile.OldSize = logFile.NewSize
		logFile.NewSize = newfileSize
		dict[event.Name] = LogResource{Path: event.Name, OldSize: logFile.OldSize, NewSize: logFile.NewSize}
		lbytes, err := logFile.DiffString(newfileSize)
		if err != nil {
			fmt.Println("ERROR", err)
			return
		}
		err, id := getContainerId(event.Name)
		if err != nil {
			fmt.Println("ERROR", err)
			return
		}
		text := ""
		if len(lbytes) > 0 {
			lbytes = lbytes[1:]
			text = string(lbytes)
		}
		q.Push(&LogForQueue{Path: event.Name, Content: text, ContainerId: id})
		fmt.Println(q.count)
	}
}
