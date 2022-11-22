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
	"github.com/go-ceres/ceres/cli/model/sql/args"
	"github.com/go-ceres/ceres/cli/model/sql/parser"
	"github.com/go-ceres/ceres/ctx"
	"github.com/go-ceres/ceres/utils/formatx"
	"github.com/go-ceres/ceres/utils/pathx"
	"github.com/go-ceres/ceres/utils/stringx"
	"github.com/go-ceres/ceres/utils/templatex"
	"path/filepath"
)

//go:embed tpl/model.tpl
var modelTemplate string

// genCustomRepository 生成自定义代码
func (g *Generator) genCustomRepository(table parser.Table, projectCtx *ctx.Project, entityCtx *ctx.Project, dlArgs *args.DDlArgs) (*CodeDescribe, error) {
	res := new(CodeDescribe)
	content, err := pathx.LoadTpl(category, modelTemplateFile, modelTemplate)
	if err != nil {
		return nil, err
	}

	out, err := templatex.With("model").Parse(content).Execute(map[string]interface{}{
		"pkg":                   filepath.Base(projectCtx.WorkDir),
		"cache":                 dlArgs.Cache,
		"upperStartCamelObject": table.Name.ToCamel(),
		"lowerStartCamelObject": stringx.NewString(table.Name.ToCamel()).UnTitle(),
	})
	modelFilename, err := formatx.FileNamingFormat(g.config.Style,
		fmt.Sprintf("%s_model", table.Name.Source()))
	if err != nil {
		return nil, err
	}
	res.Content = out.String()
	res.Update = true
	res.FileName = filepath.Join(projectCtx.WorkDir, modelFilename+".go")
	return res, nil
}
