//    Copyright 2021. Go-Ceres
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

package gen

import (
	"errors"
	"github.com/go-ceres/ceres/new/model/gorm/parse"
	"github.com/go-ceres/ceres/new/model/gorm/template/model"
	"github.com/go-ceres/ceres/utils"
	"strings"
)

// genPrimaryWhere 生成主键查询
func genPrimaryWhere(field *parse.Field) (string, error) {
	if field == nil {
		return "", errors.New("primaryKey is nil")
	}
	tplText, err := utils.LoadTpl(category, primaryTemplate, model.PrimaryWhere)
	if err != nil {
		return "", err
	}
	text, err := utils.NewTemplate("primary_where").Parse(tplText).Execute(map[string]interface{}{
		"camelPrimary": utils.CamelString(field.Name),
		"primary":      field.OriginalName,
	})
	if err != nil {
		return "", err
	}
	return text.String(), nil
}

// genOtherWheres 生成其他查询
func genOtherWheres(table *parse.Table) (string, error) {
	var queryList []string
	for _, field := range table.Fields {
		if field == table.Primary {
			continue
		} else if !(field.Unique || field.Fulltext) {
			continue
		}
		code, err := genOtherWhere(field)
		if err != nil {
			return "", err
		}
		queryList = append(queryList, code)
	}
	return strings.Join(queryList, "\n"), nil
}

// genOtherWhere 生成单个查询
func genOtherWhere(field *parse.Field) (string, error) {
	tplText, err := utils.LoadTpl(category, otherTemplate, model.OtherWhere)
	if err != nil {
		return "", err
	}
	text, err := utils.NewTemplate("other_where").Parse(tplText).Execute(map[string]interface{}{
		"camelName": utils.CamelString(field.Name),
		"name":      field.OriginalName,
		"number":    field.Number,
	})
	if err != nil {
		return "", err
	}
	return text.String(), nil
}
