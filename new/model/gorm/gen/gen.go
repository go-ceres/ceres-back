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
	"fmt"
	"github.com/go-ceres/ceres/new/model/gorm/parse"
	"github.com/go-ceres/ceres/new/model/gorm/template/model"
	"github.com/go-ceres/ceres/utils"
	"github.com/go-ceres/go-ceres/logger"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// Generator 代码构建器
type Generator struct {
	Statement *parse.Statement
	Config    *Config
}

// Start 开始生成
func (g *Generator) Start() error {
	m := make(map[string]string)
	for _, table := range g.Statement.Tables {
		code, err := g.genModel(table)
		if err != nil {
			return err
		}
		m[table.Table] = code
	}
	return g.writeFile(m)
}

// genModel 生成模块代码
func (g *Generator) genModel(table *parse.Table) (string, error) {
	// 如果表没有主键时
	if table.Primary == nil {
		// 如果设置了自动创建主键的配置
		if g.Config.AutoPrimary {
			field := &parse.Field{
				Name: "Id",
				Type: "int64",
				Tag:  "",
			}
			table.Primary = field
			table.Fields = append(table.Fields, field)
		}
	}
	// 导入包代码
	importCode, err := genImports(g.Config.Cache, table.Time)
	if err != nil {
		return "", err
	}
	// 变量代码
	varCode, err := genVars(table, g.Config.Cache)
	if err != nil {
		return "", err
	}
	// 结构体代码
	typesCode, err := genTypes(table, g.Config.Cache)
	if err != nil {
		return "", err
	}
	// 表名代码
	tablenameCode, err := genTablename(table, g.Config.Prefix)
	if err != nil {
		return "", err
	}
	// gorm设置实体代码
	dbCode, err := genDb(table)
	if err != nil {
		return "", err
	}
	// 数据迁移代码
	automigrateCode, err := genAutomigrate(table)
	if err != nil {
		return "", err
	}
	// 新建模型代码
	newCode, err := genNew(table, g.Config.Cache)
	if err != nil {
		return "", err
	}
	// 查询代码
	findCode, err := genFind(table)
	if err != nil {
		return "", err
	}
	// 查询列表
	queryCode, err := genQuery(table)
	if err != nil {
		return "", err
	}
	// 修改代码
	updateCode, err := genUpdate(table)
	if err != nil {
		return "", err
	}
	// 删除代码
	deleteCode, err := genDelete(table)
	if err != nil {
		return "", err
	}
	createCode, err := genCreate(table)
	if err != nil {
		return "", err
	}

	return g.execute(map[string]interface{}{
		"pkg":         g.Config.Pkg,
		"imports":     importCode,
		"vars":        varCode,
		"types":       typesCode,
		"tablename":   tablenameCode,
		"db":          dbCode,
		"automigrate": automigrateCode,
		"new":         newCode,
		"create":      createCode,
		"delete":      deleteCode,
		"update":      updateCode,
		"find":        findCode,
		"query":       queryCode,
	})
}

// genModel 生成模型代码
func (g *Generator) execute(data map[string]interface{}) (string, error) {
	tplText, err := utils.LoadTpl(category, modelTemplate, model.Model)
	if err != nil {
		return "", err
	}
	text, err := utils.NewTemplate("model").Parse(tplText).FormatGo(true).Execute(data)
	if err != nil {
		return "", err
	}
	return text.String(), nil
}

// writeFile 写出到文件
func (g *Generator) writeFile(fileList map[string]string) error {
	// 输出目录
	distDir, err := filepath.Abs(g.Config.Dist)
	if err != nil {
		return err
	}
	distDir = path.Join(distDir, "model")
	err = utils.MkdirIfNotFound(distDir)
	if err != nil {
		return err
	}
	// 输出所有的文件
	for tableName, code := range fileList {
		// 文件名为蛇形样式
		modelName := fmt.Sprintf("%s_model", utils.SnakeString(tableName))
		fileName := fmt.Sprintf("%s.go", modelName)
		fullPath := filepath.Join(distDir, fileName)
		if utils.FileExists(fullPath) {
			logger.Warnf("%s already exists, ignored.", fileName)
			continue
		}
		err = ioutil.WriteFile(fullPath, []byte(code), os.ModePerm)
		if err != nil {
			return err
		}
	}
	logger.Infof("code generation complete")
	return nil
}
