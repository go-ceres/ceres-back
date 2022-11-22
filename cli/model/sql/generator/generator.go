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
	"github.com/go-ceres/ceres/config"
	"github.com/go-ceres/ceres/env"
	"github.com/go-ceres/ceres/utils/console"
	"github.com/go-ceres/go-ceres/logger"
)

type Generator struct {
	config  *config.Config
	log     console.Console
	verbose bool
}

const (
	category                = "model"
	importTemplateFile      = "import.tpl"
	structTemplateFile      = "struct.tpl"
	fieldTemplateFile       = "field.tpl"
	tagTemplateFile         = "tag.tpl"
	queryTagTemplateFile    = "query-tag.tpl"
	getDbTemplateFile       = "db.tpl"
	newTemplateFile         = "new.tpl"
	modelTemplateFile       = "model.tpl"
	modelGenTemplateFile    = "model-gen.tpl"
	autoMigrateTemplateFile = "automigrate.tpl"
	tableNameTemplateFile   = "tablename.tpl"
	entityTemplateFIle      = "entity.tpl"
	createTemplateFile      = "create.tpl"
	deleteTemplateFile      = "delete.tpl"
	updateTemplateFile      = "update.tpl"
	findTemplateFile        = "find.tpl"
	queryTemplateFile       = "query.tpl"
)

var extra = map[string]string{
	"LeftBrackets":  "{", // 左括号转义
	"RightBrackets": "}", // 右括号转义
}

// NewGenerator 创建生成器
func NewGenerator(style string, verbose bool) *Generator {
	conf, err := config.NewConfig(style)
	if err != nil {
		logger.DPanic(err)
	}
	return &Generator{
		config:  conf,
		log:     console.NewColorConsole(verbose),
		verbose: verbose,
	}
}

// Prepare 前置步骤安装没有安装的依赖
func (g *Generator) Prepare() error {
	return env.Prepare(true, true, g.verbose)
}
