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

const logicFunctionTemplate = `{{if .hasComment}}{{.comment}}{{end}}
func (l *{{.logicName}}) {{.method}} (ctx context.Context,{{if .hasReq}}req {{.request}}{{if .stream}},stream {{.streamBody}}{{end}}{{else}}stream {{.streamBody}}{{end}}) ({{if .hasReply}}{{.response}},{{end}} error) {
	// todo: add your logic here and delete this line
	
	return {{if .hasReply}}&{{.responseType}}{},{{end}} nil
}
`

//go:embed tpl/logic.tpl
var logicTemplate string

// GenLogic 生成逻辑层
func (g *Generator) GenLogic(ctx DirContext, proto model.Proto, conf *config.Config) error {
	if !conf.Multiple {
		return g.genLogicInCompatibility(ctx, proto, conf)
	}
	return g.genLogicGroup(ctx, proto, conf)
}

// genLogicInCompatibility 生成单服务下的逻辑处理
func (g *Generator) genLogicInCompatibility(ctx DirContext, proto model.Proto, conf *config.Config) error {
	dir := ctx.GetLogic()
	service := proto.Service[0].Service.Name
	for _, rpc := range proto.Service[0].RPC {
		logicName := fmt.Sprintf("%sLogic", stringc.NewString(rpc.Name).ToCamel())
		fileName, err := formatc.FileNamingFormat(g.config.Style, rpc.Name+"_logic")
		if err != nil {
			return err
		}
		logicFilename := filepath.Join(dir.Filename, fileName+".go")
		functions, err := g.genLogicFunction(service, proto.PbPackage, logicName, rpc)
		if err != nil {
			return err
		}
		imports := []string{
			fmt.Sprintf(`"%v"`, ctx.GetPb().Package),
		}
		if conf.DtoAndService {
			imports = append(imports, fmt.Sprintf(`"%v"`, ctx.GetService().Package))
		}
		content, err := pathc.LoadTpl(category, logicTemplateFileFile, logicTemplate)
		if err != nil {
			return err
		}
		err = templatec.With("logic").GoFmt(true).Parse(content).SaveTo(map[string]interface{}{
			"logicName":          fmt.Sprintf("%sLogic", stringc.NewString(rpc.Name).ToCamel()),
			"functions":          functions,
			"packageName":        "logic",
			"serviceName":        fmt.Sprintf("%sService", stringc.NewString(rpc.Name).ToCamel()),
			"unTitleServiceName": fmt.Sprintf("%sService", stringc.NewString(stringc.NewString(rpc.Name).ToCamel()).UnTitle()),
			"dtoAndService":      conf.DtoAndService,
			"servicePackageName": ctx.GetService().Base,
			"extra":              extra,
			"imports":            strings.Join(imports, "\n"),
		}, logicFilename, false)
		if err != nil {
			return err
		}
	}

	return nil
}

// genLogicGroup 生成多服务下的逻辑处理层
func (g *Generator) genLogicGroup(ctx DirContext, proto model.Proto, conf *config.Config) error {
	dir := ctx.GetLogic()
	for _, item := range proto.Service {
		serviceName := item.Name
		for _, rpc := range item.RPC {
			var (
				err           error
				filename      string
				logicName     string
				logicFilename string
				packageName   string
			)

			logicName = fmt.Sprintf("%sLogic", stringc.NewString(rpc.Name).ToCamel())
			childPkg, err := dir.GetChildPackage(serviceName)
			if err != nil {
				return err
			}

			serviceDir := filepath.Base(childPkg)
			nameJoin := fmt.Sprintf("%s_logic", serviceName)
			packageName = strings.ToLower(stringc.NewString(nameJoin).ToCamel())
			logicFilename, err = formatc.FileNamingFormat(g.config.Style, rpc.Name+"_logic")
			if err != nil {
				return err
			}

			filename = filepath.Join(dir.Filename, serviceDir, logicFilename+".go")
			functions, err := g.genLogicFunction(serviceName, proto.PbPackage, logicName, rpc)
			if err != nil {
				return err
			}

			imports := []string{fmt.Sprintf(`"%v"`, ctx.GetPb().Package)}
			text, err := pathc.LoadTpl(category, logicTemplateFileFile, logicTemplate)
			if err != nil {
				return err
			}

			if err = templatec.With("logic").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
				"logicName":   logicName,
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

func (g *Generator) genLogicFunction(serviceName, goPackage, logicName string, rpc *model.RPC) (string, error) {
	functions := make([]string, 0)
	text, err := pathc.LoadTpl(category, logicFuncTemplateFileFile, logicFunctionTemplate)
	if err != nil {
		return "", err
	}

	comment := parser.GetComment(rpc.Doc())
	streamServer := fmt.Sprintf("%s.%s_%s%s", goPackage, parser.CamelCase(serviceName),
		parser.CamelCase(rpc.Name), "Server")
	buffer, err := templatec.With("fun").Parse(text).Execute(map[string]interface{}{
		"logicName":    logicName,
		"method":       parser.CamelCase(rpc.Name),
		"hasReq":       !rpc.StreamsRequest,
		"request":      fmt.Sprintf("*%s.%s", goPackage, parser.CamelCase(rpc.RequestType)),
		"hasReply":     !rpc.StreamsRequest && !rpc.StreamsReturns,
		"response":     fmt.Sprintf("*%s.%s", goPackage, parser.CamelCase(rpc.ReturnsType)),
		"responseType": fmt.Sprintf("%s.%s", goPackage, parser.CamelCase(rpc.ReturnsType)),
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
