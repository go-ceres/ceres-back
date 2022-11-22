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
	"github.com/go-ceres/ceres/utils/templatex"
	"strings"
)

//go:embed tpl/tag.tpl
var tagTemplate string

//go:embed tpl/query-tag.tpl
var queryTagTemplate string

// genTag 生成字段标签
func genTag(tag string) (string, error) {
	tplText, err := pathx.LoadTpl(category, tagTemplateFile, strings.ReplaceAll(tagTemplate, "\n", ""))
	if err != nil {
		return "", err
	}
	text, err := templatex.With("tag").Parse(tplText).Execute(map[string]interface{}{
		"tag": tag,
	})
	if err != nil {
		return "", err
	}
	return text.String(), nil
}

// genQueryTag 生成查询字段标签
func genQueryTag(field *parser.Field) (string, error) {
	tplText, err := pathx.LoadTpl(category, queryTagTemplateFile, strings.ReplaceAll(queryTagTemplate, "\n", ""))
	if err != nil {
		return "", err
	}
	tag := "json:\"" + field.OriginalName + "\" form:\"" + field.OriginalName + "\""
	text, err := utils.NewTemplate("querytag").Parse(tplText).Execute(map[string]interface{}{
		"tag": tag,
	})
	if err != nil {
		return "", err
	}
	return text.String(), nil
}
