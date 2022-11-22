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

const functionTemplate = `
{{if .hasComment}}{{.comment}}{{end}}
func (s *{{.server}}Server) {{.method}} ({{if .notStream}}ctx context.Context,{{if .hasReq}} req {{.request}}{{end}}{{else}}{{if .hasReq}} req {{.request}},{{end}}stream {{.streamBody}}{{end}}) ({{if .notStream}}{{.response}},{{end}}error) {
	return s.{{.unTitleLogicName}}.{{.method}}(ctx,{{if .hasReq}}req{{if .stream}} ,stream{{end}}{{else}}{{if .stream}}stream{{end}}{{end}})
}
`

//go:embed tpl/server.tpl
var serverTemplate string

type LogicDesc struct {
	LogicName    string
	UnTitleName  string
	LogicPackage string
}

// GenServer 生成接口门面
func (g *Generator) GenServer(ctx DirContext, proto model.Proto, conf *config.Config) error {
	if !conf.Multiple {
		return g.genServerCompatibility(ctx, proto)
	}
	return g.genServerGroup(ctx, proto)
}

// genServerCompatibility 生成单个服务
func (g *Generator) genServerCompatibility(ctx DirContext, proto model.Proto) error {
	dir := ctx.GetServer()
	logicImport := fmt.Sprintf(`"%v"`, ctx.GetLogic().Package)
	pbImport := fmt.Sprintf(`"%v"`, ctx.GetPb().Package)
	imports := []string{logicImport, pbImport}

	head := templatec.GetHead(proto.Name)

	service := proto.Service[0]
	serverFilename, err := formatc.FileNamingFormat(g.config.Style, service.Name+"_server")
	if err != nil {
		return err
	}

	serverFile := filepath.Join(dir.Filename, serverFilename+".go")
	funcList, logicList, err := g.genFunctions(proto.PbPackage, service, false)
	if err != nil {
		return err
	}

	text, err := pathc.LoadTpl(category, serverTemplateFile, serverTemplate)
	if err != nil {
		return err
	}

	notStream := false
	for _, rpc := range service.RPC {
		if !rpc.StreamsRequest && !rpc.StreamsReturns {
			notStream = true
			break
		}
	}

	return templatec.With("server").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"head":      head,
		"pbPackage": proto.PbPackage,
		"unimplementedServer": fmt.Sprintf("%s.Unimplemented%sServer", proto.PbPackage,
			stringc.NewString(service.Name).ToCamel()),
		"server":    stringc.NewString(service.Name).ToCamel(),
		"imports":   strings.Join(imports, "\n"),
		"funcs":     strings.Join(funcList, "\n"),
		"logicList": logicList,
		"notStream": notStream,
		"extra":     extra,
	}, serverFile, true)
}

// genServerGroup 生成多个服务
func (g *Generator) genServerGroup(ctx DirContext, proto model.Proto) error {
	dir := ctx.GetServer()
	for _, service := range proto.Service {
		var (
			serverFile  string
			logicImport string
		)

		serverFilename, err := formatc.FileNamingFormat(g.config.Style, service.Name+"_server")
		if err != nil {
			return err
		}

		serverChildPkg, err := dir.GetChildPackage(service.Name)
		if err != nil {
			return err
		}

		logicChildPkg, err := ctx.GetLogic().GetChildPackage(service.Name)
		if err != nil {
			return err
		}

		serverDir := filepath.Base(serverChildPkg)
		logicImport = fmt.Sprintf(`"%v"`, logicChildPkg)
		serverFile = filepath.Join(dir.Filename, serverDir, serverFilename+".go")

		pbImport := fmt.Sprintf(`"%v"`, ctx.GetPb().Package)

		imports := []string{logicImport, pbImport}

		head := templatec.GetHead(proto.Name)

		funcList, logicList, err := g.genFunctions(proto.PbPackage, service, true)
		if err != nil {
			return err
		}

		text, err := pathc.LoadTpl(category, serverTemplateFile, serverTemplate)
		if err != nil {
			return err
		}

		notStream := false
		for _, rpc := range service.RPC {
			if !rpc.StreamsRequest && !rpc.StreamsReturns {
				notStream = true
				break
			}
		}

		if err = templatec.With("server").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
			"head":      head,
			"pbPackage": proto.PbPackage,
			"unimplementedServer": fmt.Sprintf("%s.Unimplemented%sServer", proto.PbPackage,
				stringc.NewString(service.Name).ToCamel()),
			"server":    stringc.NewString(service.Name).ToCamel(),
			"imports":   strings.Join(imports, "\n"),
			"funcs":     strings.Join(funcList, "\n"),
			"logicList": logicList,
			"notStream": notStream,
			"extra":     extra,
		}, serverFile, true); err != nil {
			return err
		}
	}
	return nil
}

// genFunctions 生成server方法
func (g *Generator) genFunctions(goPackage string, service model.Service, multiple bool) ([]string, []*LogicDesc, error) {
	var (
		functionList []string
		logicList    []*LogicDesc
		logicPkg     string
	)
	for _, rpc := range service.RPC {
		text, err := pathc.LoadTpl(category, serverFuncTemplateFile, functionTemplate)
		if err != nil {
			return nil, nil, err
		}

		var logicName string
		if !multiple {
			logicPkg = "logic"
			logicName = fmt.Sprintf("%sLogic", stringc.NewString(rpc.Name).ToCamel())
		} else {
			nameJoin := fmt.Sprintf("%s_logic", service.Name)
			logicPkg = strings.ToLower(stringc.NewString(nameJoin).ToCamel())
			logicName = fmt.Sprintf("%sLogic", stringc.NewString(rpc.Name).ToCamel())
		}

		comment := parser.GetComment(rpc.Doc())
		streamServer := fmt.Sprintf("%s.%s_%s%s", goPackage, parser.CamelCase(service.Name),
			parser.CamelCase(rpc.Name), "Server")
		buffer, err := templatec.With("func").Parse(text).Execute(map[string]interface{}{
			"server":           stringc.NewString(service.Name).ToCamel(),
			"logicName":        logicName,
			"unTitleLogicName": stringc.NewString(logicName).UnTitle(),
			"method":           parser.CamelCase(rpc.Name),
			"request":          fmt.Sprintf("*%s.%s", goPackage, parser.CamelCase(rpc.RequestType)),
			"response":         fmt.Sprintf("*%s.%s", goPackage, parser.CamelCase(rpc.ReturnsType)),
			"hasComment":       len(comment) > 0,
			"comment":          comment,
			"hasReq":           !rpc.StreamsRequest,
			"stream":           rpc.StreamsRequest || rpc.StreamsReturns,
			"notStream":        !rpc.StreamsRequest && !rpc.StreamsReturns,
			"streamBody":       streamServer,
			"logicPkg":         logicPkg,
		})
		if err != nil {
			return nil, nil, err
		}

		temp := &LogicDesc{
			LogicPackage: logicPkg,
			LogicName:    logicName,
			UnTitleName:  stringc.NewString(logicName).UnTitle(),
		}

		functionList = append(functionList, buffer.String())
		// 添加到logic列表
		logicList = append(logicList, temp)
	}
	return functionList, logicList, nil
}
