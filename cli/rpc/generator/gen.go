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
	"github.com/go-ceres/ceres/cli/rpc/config"
	"github.com/go-ceres/ceres/cli/rpc/parser"
	"github.com/go-ceres/ceres/ctx"
	"github.com/go-ceres/ceres/utils/pathc"
	"path/filepath"
)

// Generate 生成代码
func (g *Generator) Generate(conf *config.Config) error {
	// 1.检查输出路径
	abs, err := filepath.Abs(conf.Dist)
	if err != nil {
		return err
	}

	// 2.创建输出文件夹
	err = pathc.MkdirIfNotExist(abs)
	if err != nil {
		return err
	}

	// 3.前置检查
	err = g.Prepare()
	if err != nil {
		return err
	}

	// 4.获取项目信息
	projectCtx, err := ctx.PrepareProject(abs)
	if err != nil {
		return err
	}

	// 5.翻译proto原始文件为结构体
	p := parser.NewDefaultProtoParser()
	proto, err := p.Parse(conf.ProtoFile, conf.Multiple)
	if err != nil {
		return err
	}

	// 6.创建文件夹
	dirCtx, err := g.mkdir(projectCtx, proto, conf)
	if err != nil {
		return err
	}

	// 7.生成配置文件
	if err := g.GenConfig(dirCtx, proto, conf); err != nil {
		return err
	}

	// 8.生成pb文件
	if err := g.GenPb(dirCtx, conf); err != nil {
		return err
	}

	// 9.生成全局变量文件
	if err := g.GenGlobal(dirCtx, conf); err != nil {
		return err
	}

	// 10.生成dto
	if err := g.GenDto(dirCtx, proto, conf); err != nil {
		return err
	}

	// 11.生成service
	if err := g.GenService(dirCtx, proto, conf); err != nil {
		return err
	}

	// 12.生成logic
	if err := g.GenLogic(dirCtx, proto, conf); err != nil {
		return err
	}

	// 13.生成server
	if err := g.GenServer(dirCtx, proto, conf); err != nil {
		return err
	}

	// 14.生成启动框架
	if err := g.GenBoot(dirCtx, proto, conf); err != nil {
		return err
	}

	// 15.生成main文件
	if err := g.GenMain(dirCtx, conf); err != nil {
		return err
	}
	return nil
}
