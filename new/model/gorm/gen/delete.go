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
)

// genDelete 生成删除代码
func genDelete(table *parse.Table) (string, error) {
	tplText, err := utils.LoadTpl(category, deleteTemplate, model.Delete)
	if err != nil {
		return "", err
	}
	text, err := utils.NewTemplate("update").Parse(tplText).Execute(map[string]interface{}{
		"camelName":    utils.CamelString(table.Table),
		"originalName": table.Primary.OriginalName,
		"fieldType":    table.Primary.Type,
		"fieldName":    table.Primary.Name,
	})
	if err != nil {
		return "", err
	}
	return text.String(), nil
}
