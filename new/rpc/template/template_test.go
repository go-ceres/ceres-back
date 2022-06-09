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
	"github.com/go-ceres/ceres/new/rpc/config"
	"html/template"
	"os"
	"testing"
)

func DefaultConfig() *config.Config {
	conf := config.DefaultConfig()
	conf.Namespace = "github.com/go-ceres"
	conf.Registry = "etcd"
	conf.ConfigSource = "file"
	conf.Name = "user"
	conf.Dir = "."
	return conf
}

func TestMainTemplate(t *testing.T) {
	parse, err := template.New("f").Parse(MainTemplate)
	if err != nil {
		t.Error(err)
	}
	err1 := parse.Execute(os.Stdout, DefaultConfig())
	if err != nil {
		t.Error(err1)
	}
}

func TestEngineTemplate(t *testing.T) {
	parse, err := template.New("f").Parse(EngineTemplate)
	if err != nil {
		t.Error(err)
	}
	err1 := parse.Execute(os.Stdout, DefaultConfig())
	if err != nil {
		t.Error(err1)
	}
}

func TestContainerTemplate(t *testing.T) {
	parse, err := template.New("f").Parse(ContainerTemplate)
	if err != nil {
		t.Error(err)
	}
	err1 := parse.Execute(os.Stdout, DefaultConfig())
	if err != nil {
		t.Error(err1)
	}
}

func TestModuleTemplate(t *testing.T) {
	parse, err := template.New("f").Parse(ModuleTemplate)
	if err != nil {
		t.Error(err)
	}
	err1 := parse.Execute(os.Stdout, DefaultConfig())
	if err != nil {
		t.Error(err1)
	}
}
