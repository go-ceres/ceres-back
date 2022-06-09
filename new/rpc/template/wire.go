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

import "github.com/go-ceres/ceres/new/rpc/config"

var WireTemplate = `
// +build wireinject

package container

import (
	"github.com/google/wire"
)

func NewContainer() (*Container,func(),error) {
	wire.Build({{if .Registry}}
		InitRegistry,{{end}}
		InitClient,
		ProviderSet,
	)
	return new(Container),nil,nil
}
`

// WriteWire 输出wire注入文件
func WriteWire(c *config.Config) error {
	outDir := c.DirContext.Wire
	err := write(c, outDir.FullName, WireTemplate)
	if err != nil {
		return err
	}
	return nil
}
