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
	"fmt"
	"github.com/go-ceres/ceres/new/rpc/config"
	"strings"
)

var engineTemplate = `package engine

import (
	{{.imports}}
	"github.com/go-ceres/go-ceres"
	"github.com/go-ceres/go-ceres/logger"
	"github.com/go-ceres/go-ceres/server/grpc"
)

type Engine struct {
	ceres.Engine
}

func NewEngine() *Engine {
	eng := &Engine{}
	err := eng.MustSetup(
		eng.startGrpc,
	)
	if err != nil {
		logger.Panicd("must setup error",logger.FieldErr(err))
	}
	return eng
}

func (e *Engine) startGrpc() (func(),error) {
	clear := func() {}
	server := grpc.ScanConfig("grpc").Build()
	ctn,fn,err := container.NewContainer()
	if err != nil {
		return clear,err
	}{{if .Registry}}
	e.SetRegistry(ctn.Registry){{end}}
	handler.Register(server.Server, ctn)
	return func() {
		fn()
	},e.Server(server)
}
`

// WriteEngine 输出启动引擎
func WriteEngine(c *config.Config) error {
	outDir := c.DirContext.Engine
	containerImport := fmt.Sprintf(`"%s"`, c.DirContext.Container.Package)
	handlerImport := fmt.Sprintf(`"%s"`, c.DirContext.Handler.Package)
	imports := strings.Join([]string{containerImport, handlerImport}, "\n\t")
	err := write(map[string]interface{}{
		"imports":  imports,
		"Registry": c.Registry,
	}, outDir.FullName, engineTemplate)
	if err != nil {
		return err
	}
	return nil
}
