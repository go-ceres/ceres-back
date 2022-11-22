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
	"github.com/go-ceres/ceres/cli/rpc/config"
	"github.com/go-ceres/ceres/cli/rpc/parser"
	"github.com/go-ceres/ceres/cli/rpc/parser/model"
	"github.com/go-ceres/ceres/utils/formatc"
	"github.com/go-ceres/ceres/utils/pathc"
	"github.com/go-ceres/ceres/utils/templatec"
	"path/filepath"
)

//go:embed tpl/dto.tpl
var dtoTemplate string

// GenDto 生成application下的Dto文件
func (g *Generator) GenDto(ctx DirContext, proto model.Proto, conf *config.Config) error {
	// 判断是否生成dto
	if conf.DtoAndService {
		// 生成Dto 和 service
		if !conf.Multiple {
			return g.genDtoInCompatibility(ctx, proto, conf)
		}
		return g.genDtoGroup(ctx, proto, conf)
	}
	return nil
}

func (g *Generator) genDtoInCompatibility(ctx DirContext, proto model.Proto, conf *config.Config) error {
	dir := ctx.GetDto()
	for _, rpc := range proto.Service[0].RPC {
		// 生成request文件
		if err := g.genRequest(dir.Filename, dir.Base, rpc); err != nil {
			return err
		}
		// 生成Response文件
		if err := g.genResponse(dir.Filename, dir.Base, rpc); err != nil {
			return err
		}
	}

	return nil
}

// genRequest 生成request的Dto文件
func (g *Generator) genRequest(dirFileName string, packageName string, rpc *model.RPC) error {
	content, err := pathc.LoadTpl(category, dtoTemplateFileFile, dtoTemplate)
	if err != nil {
		return err
	}
	fileName, err := formatc.FileNamingFormat(g.config.Style, rpc.RequestType)
	if err != nil {
		return err
	}
	dtoFilename := filepath.Join(dirFileName, fileName+".go")
	return templatec.With("dto").GoFmt(true).Parse(content).SaveTo(map[string]interface{}{
		"packageName": packageName,
		"structName":  parser.CamelCase(rpc.RequestType),
	}, dtoFilename, false)
}

// genResponse 生成response的Dto文件
func (g *Generator) genResponse(dirFileName string, packageName string, rpc *model.RPC) error {
	content, err := pathc.LoadTpl(category, dtoTemplateFileFile, dtoTemplate)
	if err != nil {
		return err
	}
	fileName, err := formatc.FileNamingFormat(g.config.Style, rpc.ReturnsType)
	if err != nil {
		return err
	}
	dtoFilename := filepath.Join(dirFileName, fileName+".go")
	return templatec.With("dto").GoFmt(true).Parse(content).SaveTo(map[string]interface{}{
		"packageName": packageName,
		"structName":  parser.CamelCase(rpc.ReturnsType),
	}, dtoFilename, false)
}

func (g *Generator) genDtoGroup(ctx DirContext, proto model.Proto, conf *config.Config) error {
	dir := ctx.GetDto()
	for _, item := range proto.Service {
		serverName := item.Name
		for _, rpc := range item.RPC {
			// 当前服务下dto的包路径
			childPkg, err := dir.GetChildPackage(serverName)
			if err != nil {
				return err
			}
			// 包名称
			dtoDir := filepath.Base(childPkg)
			fileDir := filepath.Join(dir.Filename, dtoDir)
			// 生成request文件
			if err := g.genRequest(fileDir, dtoDir, rpc); err != nil {
				return err
			}
			// 生成Response文件
			if err := g.genResponse(fileDir, dtoDir, rpc); err != nil {
				return err
			}
		}
	}
	return nil
}
