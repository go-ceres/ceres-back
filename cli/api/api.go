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

package api

import (
	"errors"
	"github.com/go-ceres/ceres/cli/api/config"
	"github.com/go-ceres/ceres/cli/api/gen"
	"github.com/go-ceres/ceres/utils/stringx"
	"github.com/go-ceres/cli/v2"
	"github.com/gookit/gcli/v3/interact"
	"os"
	"os/exec"
	"path/filepath"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:    "file",
		Aliases: []string{"f"},
		Usage:   "swagger file path",
	},
	&cli.StringFlag{
		Name:    "namespace",
		Aliases: []string{"n"},
		Usage:   "you can add a space name to a project. For example: github.com/ceres",
	},
	&cli.StringFlag{
		Name:    "dir",
		Aliases: []string{"d"},
		Value:   ".",
		Usage:   "code output path",
	},
}

// Run new 的相关命令
func Run(ctx *cli.Context) error {
	if ctx.Args().First() == "" {
		return errors.New("specify service name")
	}
	if err := checkEnvironment(); err != nil {
		return err
	}
	// 创建默认配置
	conf := config.DefaultConfig()
	// 项目名称
	conf.Name = ctx.Args().First()
	// 服务空间名
	conf.Namespace = ctx.String("namespace")
	// 输出路径
	conf.Dir = ctx.String("dir")
	// 项目工作路径
	projectPath, err := filepath.Abs(conf.Dir)
	if err != nil {
		return err
	}
	conf.WorkDir = filepath.Join(projectPath, stringx.SnakeString(conf.Name))
	// swagger文件
	conf.SwaggerFile = ctx.String("file")
	// 如果工作空间已经存在相应代码是否进行覆盖
	if _, err := os.Stat(conf.WorkDir); err == nil {
		cover := interact.Confirm("code already exists,do you want to overwrite the existing code", true)
		if !cover {
			return nil
		}
	}
	// 配置组件类型
	conf.ConfigSource = interact.SelectOne(
		"please select configuration center!",
		[]string{"file", "etcd"},
		"0",
	)
	// 注册中心
	registry := interact.SelectOne(
		"please select registration Center!",
		[]string{"none", "etcd"},
		"0",
	)
	// 如果有选择注册中心
	if registry != "none" {
		conf.Registry = registry
	}

	generator, err := gen.NewGenerator(conf)
	if err != nil {
		return err
	}
	return generator.Start()
}

// checkEnvironment 检查环境
func checkEnvironment() error {
	// 1.检查是否安装go
	if _, err := exec.LookPath("go"); err != nil {
		return err
	}

	// 4.检查是否安装了wire插件
	if _, err := exec.LookPath("wire"); err != nil {
		return err
	}
	return nil
}
