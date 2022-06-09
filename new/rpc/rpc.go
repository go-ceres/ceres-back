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

package rpc

import (
	"errors"
	"fmt"
	"github.com/go-ceres/ceres/new/rpc/config"
	"github.com/go-ceres/ceres/new/rpc/module"
	tmpl "github.com/go-ceres/ceres/new/rpc/template"
	strings2 "github.com/go-ceres/ceres/utils"
	exec2 "github.com/go-ceres/ceres/utils/exec"
	"github.com/go-ceres/cli/v2"
	"github.com/gookit/gcli/v3/interact"
	"github.com/jhump/protoreflect/desc/protoparse"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Flags 创建项目的参数
var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:    "proto",
		Aliases: []string{"p"},
		Usage:   "proto file",
	},
	&cli.StringFlag{
		Name:    "namespace",
		Aliases: []string{"n"},
		Usage:   "you can add a space name to a project. For example: github.com/ceres",
	},
	&cli.StringSliceFlag{
		Name:    "proto_path",
		Aliases: []string{"i"},
		Usage:   "native command of protoc, specify the directory in which to search for imports.",
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
	conf.WorkDir = filepath.Join(projectPath, strings2.SnakeString(conf.Name))
	// proto文件
	conf.ProtoFile = ctx.String("proto")
	// proto文件路径
	conf.ImportPath = ctx.StringSlice("proto_path")
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
	// 如果不是
	if registry != "none" {
		conf.Registry = registry
	}

	// 创建文件
	if err := create(conf); err != nil {
		return err
	}

	// 构建wire
	if err := wire(conf); err != nil {
		return err
	}

	return nil
}

// create 主方法
func create(c *config.Config) error {
	// 1. 检查环境
	if err := checkEnvironment(); err != nil {
		return err
	}

	// 2. 创建项目文件夹
	if err := strings2.MkdirIfNotFound(c.WorkDir); err != nil {
		return err
	}

	// 3. 解析go mod
	if err := parseMod(c); err != nil {
		return err
	}

	// 4. 解析proto文件
	if err := parseProto(c); err != nil {
		return err
	}

	// 5. 创建文件夹
	if err := createDir(c); err != nil {
		return err
	}

	// 6. 输出配置文件
	if err := tmpl.WriteConfig(c); err != nil {
		return err
	}

	// 7. 输出pb文件
	if err := createPb(c); err != nil {
		return err
	}

	// 8. 输出handler文件
	if err := tmpl.WriteHandler(c); err != nil {
		return err
	}

	// 9. 输出container文件
	if err := tmpl.WriteContainer(c); err != nil {
		return err
	}

	// 10.输出注入构建文件
	if err := tmpl.WriteWire(c); err != nil {
		return err
	}

	// 11. 输出engine文件
	if err := tmpl.WriteEngine(c); err != nil {
		return err
	}

	// 12. 输出main文件
	if err := tmpl.WriteMain(c); err != nil {
		return err
	}

	return nil
}

// checkEnvironment 检查环境
func checkEnvironment() error {
	// 1.检查是否安装go
	if _, err := exec.LookPath("go"); err != nil {
		return err
	}

	// 2.检查是否安装protoc
	if _, err := exec.LookPath("protoc"); err != nil {
		return err
	}

	// 3.检查是否安装了proto-gen-go插件
	if _, err := exec.LookPath("protoc-gen-go"); err != nil {
		return err
	}

	// 4.检查是否安装了wire插件
	if _, err := exec.LookPath("wire"); err != nil {
		return err
	}
	return nil
}

// parseMod 解析go mod
func parseMod(c *config.Config) error {
	// 初始化go.mod
	project, err := module.InitMod(c)
	if err != nil {
		return err
	}
	c.Project = project
	return nil
}

// ParseProto 输出proto文件
func parseProto(c *config.Config) error {
	// 如果没有设置proto文件，则表示使用默认文件
	if len(c.ProtoFile) == 0 {
		// 输出proto文件
		if err := tmpl.WriteProto(c); err != nil {
			return err
		}
	}
	abs, err := filepath.Abs(c.ProtoFile)
	if err != nil {
		return err
	}
	// 解析proto文件
	Parser := protoparse.Parser{}
	//加载并解析 proto文件,得到一组 FileDescriptor

	descs, err := Parser.ParseFiles(abs)
	if err != nil {
		return err
	}

	desc := descs[0].AsFileDescriptorProto()
	ret := config.Proto{
		Name:    filepath.Base(abs),
		Package: desc.GetPackage(),
		Imports: desc.GetDependency(),
		Src:     abs,
	}
	ret.GoPackage = desc.GetOptions().GetGoPackage()
	if len(desc.GetOptions().GetGoPackage()) == 0 {
		ret.GoPackage = desc.GetPackage()
	}
	descsTwo, err := Parser.ParseFilesButDoNotLink(abs)
	if err != nil {
		return err
	}
	desc = descsTwo[0]
	services := desc.GetService()
	for _, service := range services {
		srv := config.Service{
			Name: service.GetName(),
		}
		methods := service.GetMethod()
		for _, method := range methods {
			m := config.Method{
				Name:    method.GetName(),
				InType:  method.GetInputType(),
				OutType: method.GetOutputType(),
			}
			srv.Methods = append(srv.Methods, m)
		}
		ret.Services = append(ret.Services, srv)
	}
	c.Proto = &ret
	return nil
}

// wire 依赖注入
func wire(c *config.Config) error {
	dir := c.DirContext.Container
	cs := "wire"
	_, err := exec2.Command(cs, dir.Dir)
	if err != nil {
		return err
	}
	return nil
}

// createPb 创建pb文件
func createPb(c *config.Config) error {
	dir := c.DirContext.Pb
	// 自己proto文件所在的路径
	mp := filepath.Dir(c.Proto.Src)
	cs := "protoc "
	for _, s := range c.ImportPath {
		cs += " -I=" + s
	}
	cs += " -I=" + mp
	cs += " " + c.Proto.Name
	if strings.Contains(c.Proto.GoPackage, "/") {
		cs += " --go_out=plugins=grpc:" + c.DirContext.Main.Dir
	} else {
		cs += " --go_out=plugins=grpc:" + dir.Dir
	}
	fmt.Println(cs)
	_, err := exec2.Command(cs, "")
	if err != nil {
		return err
	}
	return nil
}

// createDir 创建目录
func createDir(c *config.Config) error {
	ctx := c.Project
	proto := c.Proto
	// 要创建的文件夹
	var paths []string
	// 设置文件夹路径
	configDir := filepath.Join(ctx.WorkDir, "config")
	containerDir := filepath.Join(ctx.WorkDir, "container")
	wireDir := filepath.Join(ctx.WorkDir, "container")
	engineDir := filepath.Join(ctx.WorkDir, "engine")
	pbDir := filepath.Join(ctx.WorkDir, proto.GoPackage)
	handlerDir := filepath.Join(ctx.WorkDir, "handler")
	// 添加进要创建的文件夹列表
	paths = append(paths, configDir, containerDir, engineDir, pbDir, handlerDir)
	// 创建返回
	res := &config.DirContext{
		Main: config.Dir{
			Dir:      ctx.WorkDir,
			FullName: filepath.Join(ctx.WorkDir, "main.go"),
			Package:  filepath.ToSlash(filepath.Join(ctx.Path, strings.TrimPrefix(ctx.WorkDir, ctx.Dir))),
			Base:     filepath.Base(ctx.WorkDir),
		},
		Config: config.Dir{
			Dir:      configDir,
			FullName: filepath.Join(configDir, "config.toml"),
			Package:  filepath.ToSlash(filepath.Join(ctx.Path, strings.TrimPrefix(configDir, ctx.Dir))),
			Base:     filepath.Base(configDir),
		},
		Container: config.Dir{
			Dir:      containerDir,
			FullName: filepath.Join(containerDir, "container.go"),
			Package:  filepath.ToSlash(filepath.Join(ctx.Path, strings.TrimPrefix(containerDir, ctx.Dir))),
			Base:     filepath.Base(containerDir),
		},
		Wire: config.Dir{
			Dir:      wireDir,
			FullName: filepath.Join(wireDir, "wire.go"),
			Package:  filepath.ToSlash(filepath.Join(ctx.Path, strings.TrimPrefix(containerDir, ctx.Dir))),
			Base:     filepath.Base(containerDir),
		},
		Engine: config.Dir{
			Dir:      engineDir,
			FullName: filepath.Join(engineDir, "engine.go"),
			Package:  filepath.ToSlash(filepath.Join(ctx.Path, strings.TrimPrefix(engineDir, ctx.Dir))),
			Base:     filepath.Base(engineDir),
		},
		Pb: config.Dir{
			Dir:      pbDir,
			FullName: filepath.Join(pbDir, proto.Name+".pb.go"),
			Package:  filepath.ToSlash(filepath.Join(ctx.Path, strings.TrimPrefix(pbDir, ctx.Dir))),
			Base:     filepath.Base(pbDir),
		},
		Handler: config.Dir{
			Dir:      handlerDir,
			FullName: filepath.Join(handlerDir, "handler.go"),
			Package:  filepath.ToSlash(filepath.Join(ctx.Path, strings.TrimPrefix(handlerDir, ctx.Dir))),
			Base:     filepath.Base(handlerDir),
		},
		Services: map[string]config.Dir{},
	}
	// 服务文件
	for _, service := range proto.Services {
		this := config.Dir{
			Dir:      handlerDir,
			FullName: filepath.Join(handlerDir, strings2.SnakeString(service.Name)),
			Package:  filepath.ToSlash(filepath.Join(ctx.Path, strings.TrimPrefix(handlerDir, ctx.Dir))),
			Base:     filepath.Base(handlerDir),
		}
		res.Services[service.Name] = this
	}
	// 新建文件夹
	for _, path := range paths {
		err := strings2.MkdirIfNotFound(path)
		if err != nil {
			return err
		}
	}
	c.DirContext = res
	return nil
}
