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
	strings2 "github.com/go-ceres/ceres/utils"
	"path/filepath"
	"strings"
)

// handlerTemplate 处理程序的模板
var handlerTemplate = `package handler

import (
	"context"
	{{.imports}}
)

type {{.ServiceName}}Handler struct {
	container *container.Container
}

func New{{.ServiceName}}Handler(ctn *container.Container) *{{.ServiceName}}Handler {
	return &{{.ServiceName}}Handler{
		container: ctn,
	}
}
{{range $i,$v := .Methods}}
func (h *{{$.ServiceName}}Handler) {{$v.Name}}(ctx context.Context, request *{{$.Package}}.{{$v.InType}}) (*{{$.Package}}.{{$v.OutType}}, error) {
	return &{{$.Package}}.{{$v.OutType}}{
	},nil
}
{{end}}
`

// registerTemplate 注册处理程序模板
var registerTemplate = `package handler

import (
	{{.imports}}
	"google.golang.org/grpc"
)

func Register(srv *grpc.Server,ctn *container.Container)  {

	{{range $i,$v := .services}}{{$.Package}}.Register{{$v.Name}}Server(srv, New{{$v.Name}}Handler(ctn))
	{{end}}
}
`

// WriteHandler 输出grpc处理程序
func WriteHandler(c *config.Config) error {
	// 输出处理程序
	outDir := c.DirContext.Handler // 处理文件输出路径
	containerImport := fmt.Sprintf(`"%s"`, c.DirContext.Container.Package)
	pbImport := fmt.Sprintf(`"%s"`, c.DirContext.Pb.Package)
	imports := strings.Join([]string{containerImport, pbImport}, "\n\t")
	registerService := c.Proto.Services
	for _, service := range c.Proto.Services {
		// 保存文件名
		saveFile := filepath.Join(outDir.Dir, strings2.SnakeString(service.Name)+".go")
		serviceName := strings2.CamelString(service.Name)
		err := write(map[string]interface{}{
			"imports":     imports,
			"ServiceName": serviceName,
			"Package":     c.DirContext.Pb.Base,
			"Methods":     service.Methods,
		}, saveFile, handlerTemplate)
		if err != nil {
			return err
		}
	}
	// 修改服务名为驼峰
	for i, service := range registerService {
		registerService[i].Name = strings2.CamelString(service.Name)
	}
	// 输出注册程序
	err := write(map[string]interface{}{
		"imports":  imports,
		"services": registerService,
		"Package":  c.DirContext.Pb.Base,
	}, c.DirContext.Handler.FullName, registerTemplate)
	if err != nil {
		return err
	}
	return nil
}
