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
	"errors"
	"fmt"
	"github.com/go-ceres/ceres/cli/model/sql/args"
	"github.com/go-ceres/ceres/cli/model/sql/parser"
	"github.com/go-ceres/ceres/ctx"
	"github.com/go-ceres/ceres/utils/formatx"
	"github.com/go-ceres/ceres/utils/pathx"
	"io/ioutil"
	"os"
	"path/filepath"
)

type CodeDescribe struct {
	FileName string // 文件名
	Update   bool   // 是否允许更新
	Content  string // 文件内容
}

// GeneratorFromDDl 生成模型根据DDl描述文件
func (g *Generator) GeneratorFromDDl(file string, args *args.DDlArgs) error {
	// 1.检查输出路径
	abs, err := filepath.Abs(args.Dist)
	if err != nil {
		return err
	}

	// 2.获取项目信息
	projectCtx, err := ctx.PrepareProject(abs)
	if err != nil {
		return err
	}

	// 3.获取entity输出路径
	var entityCtx *ctx.Project
	if len(args.Entity) > 0 {
		entityAbs, err := filepath.Abs(args.Entity)
		if err != nil {
			return err
		}
		if entityCtx, err = ctx.GetProjectInfo(entityAbs); err != nil {
			return err
		}
	}

	// 4.如果有entity输出目录则创建
	if entityCtx != nil {
		// 判断是否在同一个目录
		if projectCtx.Path != entityCtx.Path {
			return errors.New("the entity path and dist path are not in the same project")
		}

		err = pathx.MkdirIfNotExist(abs)
		if err != nil {
			return err
		}
	}

	// 5.创建输出文件夹
	err = pathx.MkdirIfNotExist(abs)
	if err != nil {
		return err
	}

	// 4.解析文件
	tables, err := parser.Parse(file, args.DataBase, args.Strict)
	if err != nil {
		return err
	}
	// 5.生成代码
	codes, err := g.genCodeDescribe(tables, projectCtx, entityCtx, args)
	if err != nil {
		return err
	}
	// 6.写出代码
	return g.createFile(codes)
}

// genFromDDl 生成代码描述
func (g *Generator) genCodeDescribe(tables []*parser.Table, projectCtx *ctx.Project, entityCtx *ctx.Project, args *args.DDlArgs) ([]*CodeDescribe, error) {
	var res = make([]*CodeDescribe, 0)
	for _, table := range tables {
		code, err := g.genModel(*table, projectCtx, entityCtx, args)
		if err != nil {
			return nil, err
		}
		res = append(res, code...)
	}
	return res, nil
}

// genModel 生成模型
func (g *Generator) genModel(table parser.Table, project *ctx.Project, entityCtx *ctx.Project, dlArgs *args.DDlArgs) ([]*CodeDescribe, error) {
	var res = make([]*CodeDescribe, 0)
	// 说明没有主键
	if len(table.Primary.Name.Source()) == 0 {
		return nil, fmt.Errorf("table %s: missing primary key", table.Name.Source())
	}
	// 生成实体
	if len(entityCtx.Path) > 0 {
		entity, err := g.genEntity(table, project, entityCtx, dlArgs)
		if err != nil {
			return nil, err
		}
		res = append(res, entity)
	}
	// 生成存储器
	repository, err := g.genRepository(table, project, entityCtx, dlArgs)
	if err != nil {
		return nil, err
	}
	// 生成自定义存储器
	res = append(res, repository)

	customRepository, err := g.genCustomRepository(table, project, entityCtx, dlArgs)
	if err != nil {
		return nil, err
	}
	res = append(res, customRepository)
	return res, nil
}

// genRepository 生成存储器
func (g *Generator) genRepository(table parser.Table, projectCtx *ctx.Project, entityCtx *ctx.Project, dlArgs *args.DDlArgs) (*CodeDescribe, error) {
	res := new(CodeDescribe)

	// 构建包导入代码
	importsCode, err := g.genImports(table, dlArgs.Cache, table.Time, projectCtx, entityCtx)
	if err != nil {
		return nil, err
	}

	// 构建结构体代码
	structCode, err := g.genStruct(table, projectCtx, entityCtx, dlArgs)
	if err != nil {
		return nil, err
	}

	// 生成表名
	tableNameCode, err := g.genTableName(table, dlArgs.Prefix)
	if err != nil {
		return nil, err
	}

	newCode, err := g.genNew(table, projectCtx, dlArgs.Cache)
	if err != nil {
		return nil, err
	}
	// 生成获取Db代码
	getDbCode, err := g.genGetDb(table)
	if err != nil {
		return nil, err
	}

	// 数据迁移代码
	autoMigrateCode, err := g.genAutoMigrate(table)
	if err != nil {
		return nil, err
	}

	// 新增代码
	createCode, err := g.GenCreate(table, entityCtx)
	if err != nil {
		return nil, err
	}

	// 删除代码
	deleteCode, err := g.genDelete(table, projectCtx, entityCtx, dlArgs)
	if err != nil {
		return nil, err
	}

	// 修改代码
	updateCode, err := g.GenUpdate(table, projectCtx, entityCtx, dlArgs)
	if err != nil {
		return nil, err
	}

	// 查询一条代码
	findCode, err := g.genFind(table, entityCtx)
	if err != nil {
		return nil, err
	}

	queryCode, err := g.genQueryListBytSql(table, entityCtx)
	if err != nil {
		return nil, err
	}

	content, err := g.genModelGenCode(map[string]interface{}{
		"pkg":         filepath.Base(projectCtx.WorkDir),
		"imports":     importsCode,
		"types":       structCode,
		"tablename":   tableNameCode,
		"db":          getDbCode,
		"automigrate": autoMigrateCode,
		"find":        findCode,
		"new":         newCode,
		"create":      createCode,
		"update":      updateCode,
		"delete":      deleteCode,
		"query":       queryCode,
	})
	if err != nil {
		return nil, err
	}
	modelFilename, err := formatx.FileNamingFormat(g.config.Style,
		fmt.Sprintf("%s_model", table.Name.Source()))
	if err != nil {
		return nil, err
	}
	res.Content = content
	res.Update = true
	res.FileName = filepath.Join(projectCtx.WorkDir, modelFilename+"_gen.go")
	return res, nil
}

// genEntity 生成实体代码
func (g *Generator) genEntity(table parser.Table, projectCtx *ctx.Project, entityCtx *ctx.Project, dlArgs *args.DDlArgs) (*CodeDescribe, error) {
	res := new(CodeDescribe)
	modelFilename, err := formatx.FileNamingFormat(g.config.Style,
		fmt.Sprintf("%s_entity", table.Name.Source()))
	if err != nil {
		return nil, err
	}
	// 生成字段代码
	fieldsStr, err := g.genFields(table.Fields, false)
	if err != nil {
		return nil, err
	}
	// 生成查询字段
	queryFieldStr, err := g.genQueryFields(table.Fields)
	if err != nil {
		return nil, err
	}
	code, err := g.genEntityCode(map[string]interface{}{
		"package":    filepath.Base(entityCtx.WorkDir),
		"camelName":  table.Name.ToCamel(),
		"fields":     fieldsStr,
		"queryField": queryFieldStr,
	})
	if err != nil {
		return nil, err
	}
	res.Content = code
	res.Update = false
	res.FileName = filepath.Join(entityCtx.WorkDir, modelFilename+".go")
	return res, nil
}

// createFile 创建文件
func (g *Generator) createFile(codes []*CodeDescribe) error {
	for _, code := range codes {
		exists := pathx.FileExists(code.FileName)
		// 如果文件存在并且是不允许更新的情况下提示信息
		if exists && !code.Update {
			g.log.Warning("%s already exists, ignored.", code.FileName)
			continue
		}
		// 写入文件
		if err := ioutil.WriteFile(code.FileName, []byte(code.Content), os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}
