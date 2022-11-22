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
	"github.com/go-ceres/ceres/cli/model/sql/args"
	"github.com/go-ceres/ceres/cli/model/sql/parser"
	"github.com/go-ceres/ceres/ctx"
	"github.com/go-ceres/ceres/utils"
	"github.com/go-ceres/ceres/utils/pathc"
	"github.com/go-ceres/ceres/utils/stringc"
)

//go:embed tpl/update.tpl
var updateTemplate string

// GenUpdate 生成修改代码
func (g *Generator) GenUpdate(table parser.Table, projectCtx *ctx.Project, entityCtx *ctx.Project, args *args.DDlArgs) (string, error) {
	tplText, err := pathc.LoadTpl(category, updateTemplateFile, updateTemplate)
	if err != nil {
		return "", err
	}
	text, err := utils.NewTemplate("update").Parse(tplText).Execute(map[string]interface{}{
		"camelName":    table.Name.ToCamel(),
		"originalName": table.Primary.OriginalName,
		"entity":       len(entityCtx.Path) > 0,
		"packageName":  stringc.NewString(table.Name.ToCamel()).UnTitle() + "Entity",
		"fieldType":    table.Primary.Type,
		"fieldName":    table.Primary.Name.UnTitle(),
	})
	if err != nil {
		return "", err
	}
	return text.String(), nil
}
