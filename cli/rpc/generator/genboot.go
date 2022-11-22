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
	"github.com/go-ceres/ceres/utils/templatec"
	"path/filepath"
	"strings"
)

//go:embed tpl/bootstrap.tpl
var bootTemplate string

//go:embed tpl/server-boot.tpl
var serverBootTemplate string

// GenBoot 生成引擎入口
func (g *Generator) GenBoot(ctx DirContext, proto model.Proto, conf *config.Config) error {
	// 生成engine文件
	if err := g.genBase(ctx, conf); err != nil {
		return err
	}
	return g.genServerBoot(ctx, proto, conf)
}

// genBase 生成基础
func (g *Generator) genBase(ctx DirContext, conf *config.Config) error {
	fileName := filepath.Join(ctx.GetBoot().Filename, fmt.Sprintf("%v.go", "bootstrap"))
	content, err := pathc.LoadTpl(category, bootTemplateFileFile, bootTemplate)
	if err != nil {
		return err
	}
	imports := make([]string, 0)
	globalImport := fmt.Sprintf(`"%s"`, ctx.GetGlobal().Package)
	imports = append(imports, globalImport)
	// orm框架
	if conf.Orm != nil {
		imports = append(imports, fmt.Sprintf(`"%s"`, filepath.Join("github.com/go-ceres/go-ceres/store", conf.Orm["name"])))
		conf.Orm["newFunc"] = conf.Orm["name"] + `.ScanConfig("` + ctx.GetServiceName().UnTitle() + `").Build()`
	}
	return templatec.With("boot").GoFmt(true).Parse(content).SaveTo(map[string]interface{}{
		"Registry":    conf.Registry,
		"Extra":       extra,
		"imports":     strings.Join(imports, "\n"),
		"ServiceName": ctx.GetServiceName().UnTitle(),
		"orm":         conf.Orm,
	}, fileName, false)
}

// genServerBoot 初始化服务启动代码
func (g *Generator) genServerBoot(ctx DirContext, proto model.Proto, conf *config.Config) error {
	bootFilename, err := formatc.FileNamingFormat(g.config.Style, ctx.GetServiceName().Source())
	if err != nil {
		return err
	}
	fileName := filepath.Join(ctx.GetBoot().Filename, fmt.Sprintf("%v.go", bootFilename))
	imports := make([]string, 0)
	var serverNames []ServerTemplateData
	for _, e := range proto.Service {
		var (
			remoteImport string
			serverPkg    string
		)
		if !conf.Multiple {
			serverPkg = "server"
			remoteImport = fmt.Sprintf(`"%v"`, ctx.GetServer().Package)
		} else {
			childPkg, err := ctx.GetServer().GetChildPackage(e.Name)
			if err != nil {
				return err
			}

			serverPkg = filepath.Base(childPkg + "Server")
			remoteImport = fmt.Sprintf(`%s "%v"`, serverPkg, childPkg)
		}
		imports = append(imports, remoteImport)
		serverNames = append(serverNames, ServerTemplateData{
			Server:    parser.CamelCase(e.Name),
			ServerPkg: serverPkg,
		})
	}

	content, err := pathc.LoadTpl(category, serverBootTemplateFileFile, serverBootTemplate)
	if err != nil {
		return err
	}
	return templatec.With("boot").GoFmt(true).Parse(content).SaveTo(map[string]interface{}{
		"imports":         strings.Join(imports, "\n"),
		"Registry":        conf.Registry,
		"serverNames":     serverNames,
		"BaseServiceName": ctx.GetServiceName().Source(),
		"ServiceName":     ctx.GetServiceName().ToCamel(),
	}, fileName, true)
}
