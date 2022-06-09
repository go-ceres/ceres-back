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

package rpc

import (
	"github.com/go-ceres/ceres/new/rpc/config"
	tpl "github.com/go-ceres/ceres/new/rpc/template"
	"path/filepath"
	"testing"
)

func defaultConfig() *config.Config {
	conf := config.DefaultConfig()
	conf.Namespace = "github.com/go-ceres"
	conf.Registry = "etcd"
	conf.ConfigSource = "file"
	conf.Name = "user"
	conf.Dir = "."
	abs, _ := filepath.Abs(conf.Dir)
	conf.WorkDir = abs
	return conf
}

func TestCreateProto(t *testing.T) {
	f := "user.proto"
	conf := defaultConfig()
	conf.ProtoFile = f
	err := parseProto(conf)
	if err != nil {
		t.Error(err)
	}
	//conf.ImportPath = append(conf.ImportPath, "/Users/liuqin/Documents/workspaces/private/go/pkg")
	err = createPb(conf)
	if err != nil {
		t.Error(err)
	}
}

func TestHandlerTemplate(t *testing.T) {
	f := "user.proto"
	conf := defaultConfig()
	conf.ProtoFile = f
	err := parseProto(conf)
	if err != nil {
		t.Error(err)
	}
	err1 := tpl.WriteHandler(conf)
	if err1 != nil {
		t.Error(err1)
	}
}

func TestWriteProto(t *testing.T) {
	f := "user.proto"
	conf := defaultConfig()
	conf.ProtoFile = f
	err := tpl.WriteProto(conf)
	if err != nil {
		t.Error(err)
	}
}

func TestCreate(t *testing.T) {
	conf := defaultConfig()
	if err := create(conf); err != nil {
		t.Error(err)
	}

}
