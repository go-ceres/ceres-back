//    Copyright 2021. Go-Ceres
//    Author https://github.com/go-ceres/go-ceres
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package template

import (
	"github.com/alecthomas/template"
	"github.com/go-ceres/ceres/new/rpc/config"
	"os"
	"path/filepath"
)

var protoTemplate = `
syntax = "proto3";

package proto;

service Greeter {
  rpc SayHello (HelloRequest) returns (HelloReply) {}
}

message HelloRequest {
  string name = 1;
}

message HelloReply {
  string message = 1;
}
`

// WriteProto 写出默认的文件
func WriteProto(c *config.Config) error {
	dir := filepath.Join(c.WorkDir, "proto")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	file := filepath.Join(dir, "greeter.proto")
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	t, err := template.New("f").Parse(protoTemplate)
	if err != nil {
		return err
	}
	c.ProtoFile = file
	return t.Execute(f, map[string]string{})
}
