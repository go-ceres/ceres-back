//    Copyright 2022. Go-Ceres
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

package generator

import (
	_ "embed"
	"fmt"
	"github.com/go-ceres/ceres/cli/rpc/config"
	"github.com/go-ceres/ceres/cli/rpc/parser"
	"github.com/go-ceres/ceres/cli/rpc/parser/model"
	"github.com/go-ceres/ceres/utils/formatc"
	"github.com/go-ceres/ceres/utils/pathc"
	"github.com/go-ceres/ceres/utils/stringc"
	"github.com/go-ceres/ceres/utils/templatec"
	"path/filepath"
	"strings"
)

const serviceFunctionTemplate = `{{if .hasComment}}{{.comment}}{{end}}
func (l *{{.logicName}}) {{.method}} (ctx context.Context,{{if .hasReq}}req {{.request}}{{if .stream}},stream {{.streamBody}}{{end}}{{else}}stream {{.streamBody}}{{end}}) ({{if .hasReply}}{{.response}},{{end}} error) {
	// todo: add your logic here and delete this line
	
	return {{if .hasReply}}&{{.responseType}}{},{{end}} nil
}
`

//go:embed tpl/service.tpl
var serviceTemplate string

// GenService 生成application下的service层
func (g *Generator) GenService(ctx DirContext, proto model.Proto, conf *config.Config) error {
	// 判断是否生成service
	if conf.DtoAndService {
		// 生成Dto 和 service
		if !conf.Multiple {
			return g.genServiceInCompatibility(ctx, proto, conf)
		}
		return g.genServiceGroup(ctx, proto, conf)
	}
	return nil
}

func (g *Generator) genServiceInCompatibility(ctx DirContext, proto model.Proto, conf *config.Config) error {
	dir := ctx.GetService()
	service := proto.Service[0].Service.Name
	for _, rpc := range proto.Service[0].RPC {
		serviceName := fmt.Sprintf("%sService", stringc.NewString(rpc.Name).ToCamel())
		fileName, err := formatc.FileNamingFormat(g.config.Style, rpc.Name+"_service")
		if err != nil {
			return err
		}
		serviceFilename := filepath.Join(dir.Filename, fileName+".go")
		functions, err := g.genServiceFunction(service, ctx.GetDto().Base, serviceName, rpc)
		if err != nil {
			return err
		}
		imports := []string{
			fmt.Sprintf(`"%v"`, ctx.GetDto().Package),
		}
		content, err := pathc.LoadTpl(category, serviceTemplateFileFile, serviceTemplate)
		if err != nil {
			return err
		}
		err = templatec.With("service").GoFmt(true).Parse(content).SaveTo(map[string]interface{}{
			"serviceName": fmt.Sprintf("%sService", stringc.NewString(rpc.Name).ToCamel()),
			"functions":   functions,
			"packageName": "service",
			"imports":     strings.Join(imports, "\n"),
		}, serviceFilename, false)
		if err != nil {
			return err
		}
	}

	return nil

}

func (g *Generator) genServiceGroup(ctx DirContext, proto model.Proto, conf *config.Config) error {
	dir := ctx.GetService()
	for _, item := range proto.Service {
		serverName := item.Name
		for _, rpc := range item.RPC {
			var (
				err             error
				filename        string
				serviceName     string
				serviceFilename string
				packageName     string
			)

			serviceName = fmt.Sprintf("%sService", stringc.NewString(rpc.Name).ToCamel())
			childPkg, err := dir.GetChildPackage(serverName)
			if err != nil {
				return err
			}

			dtoPackage, err := ctx.GetDto().GetChildPackage(serverName)
			if err != nil {
				return err
			}

			serviceDir := filepath.Base(childPkg)
			nameJoin := fmt.Sprintf("%s_service", serverName)
			packageName = strings.ToLower(stringc.NewString(nameJoin).ToCamel())
			serviceFilename, err = formatc.FileNamingFormat(g.config.Style, rpc.Name+"_service")
			if err != nil {
				return err
			}

			filename = filepath.Join(dir.Filename, serviceDir, serviceFilename+".go")
			functions, err := g.genLogicFunction(serverName, "dto", serviceName, rpc)
			if err != nil {
				return err
			}

			imports := []string{fmt.Sprintf(`dto "%v"`, dtoPackage)}
			text, err := pathc.LoadTpl(category, serviceTemplateFileFile, serviceTemplate)
			if err != nil {
				return err
			}

			if err = templatec.With("service").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
				"serviceName": serviceName,
				"functions":   functions,
				"packageName": packageName,
				"imports":     strings.Join(imports, "\n"),
			}, filename, false); err != nil {
				return err
			}
		}
	}
	return nil
}

// genServiceFunction 生成服务方法
func (g *Generator) genServiceFunction(serviceName, dtoPackage, logicName string, rpc *model.RPC) (string, error) {
	functions := make([]string, 0)
	text, err := pathc.LoadTpl(category, serviceFuncTemplateFileFile, serviceFunctionTemplate)
	if err != nil {
		return "", err
	}

	comment := parser.GetComment(rpc.Doc())
	streamServer := fmt.Sprintf("%s.%s_%s%s", dtoPackage, parser.CamelCase(serviceName),
		parser.CamelCase(rpc.Name), "Server")
	buffer, err := templatec.With("fun").Parse(text).Execute(map[string]interface{}{
		"logicName":    logicName,
		"method":       parser.CamelCase(rpc.Name),
		"hasReq":       !rpc.StreamsRequest,
		"request":      fmt.Sprintf("*%s.%s", dtoPackage, parser.CamelCase(rpc.RequestType)),
		"hasReply":     !rpc.StreamsRequest && !rpc.StreamsReturns,
		"response":     fmt.Sprintf("*%s.%s", dtoPackage, parser.CamelCase(rpc.ReturnsType)),
		"responseType": fmt.Sprintf("%s.%s", dtoPackage, parser.CamelCase(rpc.ReturnsType)),
		"stream":       rpc.StreamsRequest || rpc.StreamsReturns,
		"streamBody":   streamServer,
		"hasComment":   len(comment) > 0,
		"comment":      comment,
	})
	if err != nil {
		return "", err
	}

	functions = append(functions, buffer.String())
	return strings.Join(functions, "\n"), nil
}
