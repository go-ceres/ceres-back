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
	"github.com/go-ceres/ceres/utils"
	"github.com/go-ceres/ceres/utils/pathx"
)

//go:embed tpl/tablename.tpl
var tableNameTemplate string

// genTableName 生成获取表名方法代码
func (g *Generator) genTableName(table parser.Table, prefix string) (string, error) {
	tplText, err := pathx.LoadTpl(category, tableNameTemplateFile, tableNameTemplate)
	if err != nil {
		return "", err
	}
	text, err := utils.NewTemplate("import").Parse(tplText).Execute(map[string]interface{}{
		"camelName": table.Name.ToCamel(),
		"name":      table.Name.Source(),
		"prefix":    prefix,
	})
	if err != nil {
		return "", err
	}
	return text.String(), nil
}
