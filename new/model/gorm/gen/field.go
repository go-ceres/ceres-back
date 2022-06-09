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
	"github.com/go-ceres/ceres/new/model/gorm/parse"
	"github.com/go-ceres/ceres/new/model/gorm/template/model"
	"github.com/go-ceres/ceres/utils"
	"strings"
)

// genFields 生成所有字段代码
func genFields(fields []*parse.Field) (string, error) {
	var fieldList []string
	for _, field := range fields {
		res, err := genField(field)
		if err != nil {
			return "", err
		}
		fieldList = append(fieldList, res)
	}
	return strings.Join(fieldList, "\n\t"), nil
}

// genQueryFields 生成查询字段
func genQueryFields(fields []*parse.Field) (string, error) {
	var queryFieldList []string
	for _, field := range fields {
		if !(field.Unique || field.Primary || field.Fulltext) {
			continue
		}
		this := new(parse.Field)
		this.OriginalName = field.OriginalName
		// 主键配置，用于多个主键查询
		if field.Primary {
			this.Name = field.Name + "s"
			this.OriginalName = field.OriginalName + "s"
			this.Type = "[]" + field.Type
		} else if field.Unique || field.Fulltext {
			this.Name = field.Name
			switch field.Type {
			case "int64", "int32", "int", "uint64", "uint32", "uint", "float64", "float32":
				this.Type = "*" + field.Type
			default:
				this.Type = field.Type
			}
		} else {
			this.Name = field.Name
			this.Type = field.Type
		}
		res, err := genQueryField(this)
		if err != nil {
			return "", err
		}
		queryFieldList = append(queryFieldList, res)

	}
	return strings.Join(queryFieldList, "\n\t"), nil
}

// genField 生成查询字段
func genQueryField(field *parse.Field) (string, error) {
	tag, err := genQueryTag(field)
	if err != nil {
		return "", err
	}
	tplText, err := utils.LoadTpl(category, fieldTemplate, model.Field)
	if err != nil {
		return "", err
	}
	text, err := utils.NewTemplate("field").Parse(tplText).Execute(map[string]interface{}{
		"name": utils.CamelString(field.Name),
		"type": field.Type,
		"tag":  tag,
	})
	if err != nil {
		return "", err
	}
	return text.String(), nil
}

// genField 生成带个字段代码
func genField(field *parse.Field) (string, error) {
	tag, err := genTag(field.Tag, field.OriginalName)
	if err != nil {
		return "", err
	}
	tplText, err := utils.LoadTpl(category, fieldTemplate, model.Field)
	if err != nil {
		return "", err
	}
	text, err := utils.NewTemplate("field").Parse(tplText).Execute(map[string]interface{}{
		"name": utils.CamelString(field.Name),
		"type": field.Type,
		"tag":  tag,
	})
	if err != nil {
		return "", err
	}
	return text.String(), nil
}
