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
	"github.com/go-ceres/ceres/utils"
	"github.com/go-ceres/ceres/utils/pathx"
)

//go:embed tpl/model-gen.tpl
var modelGenTemplate string

// genModel 生成模型代码
func (g *Generator) genModelGenCode(data map[string]interface{}) (string, error) {
	tplText, err := pathx.LoadTpl(category, modelGenTemplateFile, modelGenTemplate)
	if err != nil {
		return "", err
	}
	text, err := utils.NewTemplate("model").Parse(tplText).FormatGo(true).Execute(data)
	if err != nil {
		return "", err
	}
	return text.String(), nil
}
