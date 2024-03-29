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
	"github.com/go-ceres/ceres/cli/model/sql/parser"
	"github.com/go-ceres/ceres/ctx"
	"github.com/go-ceres/ceres/utils"
	"github.com/go-ceres/ceres/utils/pathx"
	"github.com/go-ceres/ceres/utils/stringx"
)

//go:embed tpl/create.tpl
var createTemplate string

// GenCreate 生成新增代码
func (g *Generator) GenCreate(table parser.Table, entityCtx *ctx.Project) (string, error) {
	tplText, err := pathx.LoadTpl(category, createTemplateFile, createTemplate)
	if err != nil {
		return "", err
	}
	text, err := utils.NewTemplate("create").Parse(tplText).Execute(map[string]interface{}{
		"camelName":   table.Name.ToCamel(),
		"entity":      len(entityCtx.Path) > 0,
		"packageName": stringx.NewString(table.Name.ToCamel()).UnTitle() + "Entity",
	})
	if err != nil {
		return "", err
	}
	return text.String(), nil
}
