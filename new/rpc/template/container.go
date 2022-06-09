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
)

var containerTemplate = `package container

import (
	"github.com/go-ceres/go-ceres/client/grpc"{{if .Registry}}
	"github.com/go-ceres/go-ceres/registry"
	"github.com/go-ceres/go-ceres/registry/{{.Registry}}"{{end}}
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(wire.Struct(new(Container),"*"))

type Container struct {{.Extra.LeftBrackets}}{{if .Registry}}
	Registry registry.Registry{{end}}
	Client	*grpc.Client
}

{{if .Registry}}
func InitRegistry() (registry.Registry,func(),error) {
	reg := etcd.ScanConfig("{{.Registry}}").Build()
	return reg, func() {
		err := reg.Close()
		if err != nil {
			return 
		}
	},nil
}
{{end}}

func InitClient({{if .Registry}}registry registry.Registry{{end}}) (*grpc.Client,func(),error) {
	client := grpc.ScanConfig("default"){{if .Registry}}.WithRegistry(registry){{end}}.Build()
	return client, func() {
		client.Close()
	},nil
}
`

// WriteContainer 输出container文件
func WriteContainer(c *config.Config) error {
	outDir := c.DirContext.Container
	err := write(c, outDir.FullName, containerTemplate)
	if err != nil {
		return err
	}
	return nil
}
