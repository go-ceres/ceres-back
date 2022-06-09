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

package gen

import (
	"fmt"
	"github.com/go-ceres/ceres/new/api/config"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc/protoparse"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/types/descriptorpb"
	"os"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	conf := config.DefaultConfig()
	conf.SwaggerFile = "../example/swagger/user.yaml"
	generator, err := NewGenerator(conf)
	if err != nil {
		return
	}
	err = generator.Start()
	if err != nil {
		return
	}
	fmt.Println()
}

func TestProto(t *testing.T) {
	// 解析proto文件
	Parser := protoparse.Parser{
		ImportPaths: []string{os.Getenv("GOPATH") + "/src", os.Getenv("GOPATH") + "/src/github.com/googleapis/googleapis"},
	}
	//加载并解析 proto文件,得到一组 FileDescriptor

	filenames, err := protoparse.ResolveFilenames([]string{os.Getenv("GOPATH") + "/src"}, "../example/proto/sayhello.proto")
	if err != nil {
		return
	}
	descs, err := Parser.ParseFiles(filenames...)
	if err != nil {
		return
	}
	desc := descs[0]
	dpb := desc.AsFileDescriptorProto()

	opt := dpb.Service[0].Method[0].Options

	msg := desc.GetMessageTypes()[0].GetFields()[0].GetType()

	if msg == descriptorpb.FieldDescriptorProto_TYPE_STRING {
		fmt.Println("进来了")
	}

	if opt != nil && proto.HasExtension(opt, annotations.E_Http) {
		if ext, _ := proto.GetExtension(opt, annotations.E_Http); ext != nil {
			if x, _ := ext.(*annotations.HttpRule); x != nil {
				fmt.Println(x)
			}
		}
	}
	fmt.Println(dpb, filenames, msg)
}
