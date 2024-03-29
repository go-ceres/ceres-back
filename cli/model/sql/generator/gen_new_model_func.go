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
	"io/ioutil"
	"path/filepath"
	"strings"
)

//go:embed tpl/new.tpl
var newTemplate string

// genNew 生成new方法代码
func (g *Generator) genNew(table parser.Table, projectCtx *ctx.Project, cache bool) (string, error) {
	tplText, err := pathx.LoadTpl(category, newTemplateFile, newTemplate)
	if err != nil {
		return "", err
	}
	// 查看是否有指定目录和文件
	hasGormDb := false
	hasCache := false
	globalFile := filepath.Join(projectCtx.Dir, "global", "global.go")
	if pathx.FileExists(globalFile) {
		fileContent, _ := ioutil.ReadFile(globalFile)
		if strings.Contains(string(fileContent), "gorm.DB") {
			hasGormDb = true
		}
		if strings.Contains(string(fileContent), "cache.Cache") {
			hasCache = true
		}
	}

	text, err := utils.NewTemplate("import").Parse(tplText).Execute(map[string]interface{}{
		"camelName": table.Name.ToCamel(),
		"cache":     cache,
		"hasGormDb": hasGormDb,
		"hasCache":  hasCache,
		"extra":     extra,
	})
	if err != nil {
		return "", err
	}
	return text.String(), nil
}
