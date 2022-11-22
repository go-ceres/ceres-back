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

package action

import (
	"errors"
	"fmt"
	"github.com/go-ceres/ceres/cli/rpc/config"
	"github.com/go-ceres/ceres/cli/rpc/generator"
	"github.com/go-ceres/ceres/utils"
	"github.com/go-ceres/ceres/utils/pathx"
	"github.com/go-ceres/cli/v2"
	"github.com/gookit/gcli/v3/interact"
	"os"
	"path/filepath"
	"strings"
)

var (
	errInvalidDistOutput = errors.New("ceres: missing --dist")
	// Flags 创建项目的参数
	Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "proto_path",
			Usage: "proto path",
		},
		&cli.BoolFlag{
			Name:  "multiple",
			Value: false,
			Usage: "describes whether support generating multiple rpc services or not",
		},
		&cli.StringSliceFlag{
			Name:  "go_opt",
			Value: nil,
			Usage: "protoc args for go_opt",
		},
		&cli.StringSliceFlag{
			Name:  "go-grpc_opt",
			Value: nil,
			Usage: "protoc args for go-grpc_opt",
		},
		&cli.StringSliceFlag{
			Name:  "plugin",
			Value: nil,
			Usage: "protoc args for plugin",
		},
		&cli.StringFlag{
			Name:  "dist",
			Usage: "code output directory",
		},
	}
)

// Protoc 创建根据proto文件适配的项目
func Protoc(ctx *cli.Context) error {
	// 创建默认配置
	conf := config.DefaultConfig()
	// 是否支持生成多个服务
	conf.Multiple = ctx.Bool("multiple")
	// grpc_out目录
	conf.GrpcOut = "./interfaces/grpc/proto"
	// go_out目录
	conf.GoOut = "./interfaces/grpc/proto"
	// proto_path目录
	conf.ProtoPath = ctx.String("proto_path")
	// 额外的go_opt参数
	conf.GoOpt = ctx.StringSlice("go_opt")
	// 额外的go-grpc_opt参数
	conf.GoGrpcOpt = ctx.StringSlice("go-grpc_opt")
	// 代码输出目录
	conf.Dist = ctx.String("dist")
	// 是否打印复杂日志
	verbose := ctx.Bool("verbose")
	// 没有设置dist
	if len(conf.Dist) == 0 {
		return errInvalidDistOutput
	}
	// 检查goOut目录是否有效
	var err error
	conf.GoOut, err = filepath.Abs(conf.GoOut)
	if err != nil {
		return err
	}
	// 检查goGrpcOut目录是否有效
	conf.GrpcOut, err = filepath.Abs(conf.GrpcOut)
	if err != nil {
		return err
	}
	// 创建goOut输出目录
	if err := pathx.MkdirIfNotFound(conf.GoOut); err != nil {
		return err
	}
	// 创建grpc输出目录
	if err := pathx.MkdirIfNotFound(conf.GrpcOut); err != nil {
		return err
	}
	// 组装protoc参数
	protoArgs := wrapProtocCmd(conf, ctx.Args().Slice())
	// 获取当前工作目录
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	// 模板相关
	home := ctx.String("home")
	remote := ctx.String("remote")
	branch := ctx.String("branch")
	if len(remote) > 0 {
		repo, _ := utils.CloneIntoGitHome(remote, branch)
		if len(repo) > 0 {
			home = repo
		}
	}
	if len(home) > 0 {
		// 注册home变量
		pathx.RegisterCeresHome(home)
	}
	// 设置输出目录
	if !filepath.IsAbs(conf.Dist) {
		conf.Dist = filepath.Join(pwd, conf.Dist)
	}
	// 配置输出路径
	conf.Dist, err = filepath.Abs(conf.Dist)
	if err != nil {
		return err
	}
	// proto文件
	conf.ProtoFile = ctx.Args().First()
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
	// 如果不是
	if registry != "none" {
		conf.Registry = registry
	}
	// 数据库框架
	orm := interact.SelectOne(
		"please select database orm",
		[]string{"none", "gorm"},
		"0",
	)
	if orm != "none" {
		conf.Components = append(conf.Components, &config.Component{
			Name:          orm,
			GlobalImport:  []string{fmt.Sprintf(`"github.com/go-ceres/go-ceres/store/%s"`, orm)},
			ImportPackage: []string{fmt.Sprintf(`"github.com/go-ceres/go-ceres/store/%s"`, orm)},
			InitFunc: func(name string) string {
				return orm + `.ScanConfig("` + name + `").Build()`
			},
			GlobalName: "Db",
			ConfigFunc: func(name string) string {
				return `[ceres.store.` + orm + `.` + name + `]
    dns="root:root@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
`
			},
			TypeName: "*gorm.DB",
		})
	}
	// 缓存框架
	cache := interact.SelectOne(
		"please select cache center!",
		[]string{"none", "redis", "memcached"},
		"0",
	)
	if cache != "none" {
		conf.Components = append(conf.Components, &config.Component{
			Name:          cache,
			GlobalImport:  []string{`"github.com/go-ceres/go-ceres/cache"`},
			ImportPackage: []string{fmt.Sprintf(`"github.com/go-ceres/go-ceres/cache/%s"`, cache)},
			InitFunc: func(name string) string {
				return cache + `.ScanConfig("` + name + `").Build()`
			},
			GlobalName: "Cache",
			ConfigFunc: func(name string) string {
				return `[ceres.cache.` + name + `]
    addr=["127.0.0.1"]
`
			},
			TypeName: "cache.Cache",
		})
	}

	// 是否生成对应的dto与service文件
	conf.DtoAndService = interact.Confirm("Generate dto and service？", true)
	// 设置protoc命令
	conf.ProtocCmd = strings.Join(protoArgs, " ")
	// 创建生成器
	g := generator.NewGenerator("go_ceres", verbose)
	return g.Generate(conf)
}

// wrapProtocCmd 包装protoc命令
func wrapProtocCmd(conf *config.Config, args []string) []string {
	res := append([]string{"protoc"}, args...)
	if len(conf.ProtoPath) > 0 {
		res = append(res, "--proto_path", conf.ProtoPath)
	}
	for _, goOpt := range conf.GoOpt {
		res = append(res, "--go_opt", goOpt)
	}
	for _, goGrpcOpt := range conf.GoGrpcOpt {
		res = append(res, "--go-grpc_opt", goGrpcOpt)
	}
	// go数据目录
	res = append(res, "--go_out", conf.GoOut)
	// grpc输出目录
	res = append(res, "--go-grpc_out", conf.GrpcOut)
	// 插件
	for _, plugin := range conf.Plugins {
		res = append(res, "--plugin="+plugin)
	}
	return res
}
