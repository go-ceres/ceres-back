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

// genAutomigrate 生成迁移代码
func genAutomigrate(table *parse.Table) (string, error) {
	tplText, err := utils.LoadTpl(category, optionTemplate, model.Automigrate)
	if err != nil {
		return "", err
	}
	text, err := utils.NewTemplate("option").Parse(tplText).Execute(map[string]interface{}{
		"camelName": utils.CamelString(table.Table),
		"options":   strings.Join(table.Options, " "),
	})
	if err != nil {
		return "", err
	}
	return text.String(), nil
}
