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

package config

// Config 新建项目时的配置文件
type Config struct {
	Namespace    string   // 项目空间，例如：github.com/go-ceres
	Name         string   // 项目名称，例如：demo
	Dir          string   // 项目输出路径，例如: .
	WorkDir      string   // 项目目录。绝对路径
	ImportPath   []string // proto文件中要包含文件的路径
	ConfigSource string   // 配置组件的类型，例如：file
	Registry     string   // 注册中心，例如：etcd
	ProtoFile    string   // 输入的文件
	Project      *Project
	Proto        *Proto
	DirContext   *DirContext // 各个文件的详细路径
	Extra        map[string]string
}

// File 输出的文件配置
type File struct {
	Base    string // 包名
	Path    string // 输出路径
	Package string // 包路径
	Tmpl    string // 魔板内容
}

type Proto struct {
	Package   string
	GoPackage string
	Imports   []string
	Name      string
	Src       string
	Services  []Service
}

type Service struct {
	Name    string
	Methods []Method
}

type Method struct {
	Name    string
	InType  string
	OutType string
}

// Project 项目解析（go mod）
type Project struct {
	WorkDir string
	Name    string
	Dir     string
	Path    string
}

// Dir 文件路径
type Dir struct {
	Package  string // 文件完整包名
	Base     string // 文件最后一级包名
	FullName string // 文件完整路径
	Dir      string // 文件夹完整路径
}

// DirContext 所有要输出的文件集合
type DirContext struct {
	Config    Dir            // 配置文件的路径
	Container Dir            // 服务的容器路径
	Wire      Dir            // wire注入
	Engine    Dir            // 服务的启动入口路径
	Main      Dir            // 程序入口路径
	Pb        Dir            // Pb文件路径
	Handler   Dir            // grpc处理程序路径
	Services  map[string]Dir // grpc服务处理程序路径
}

// DefaultConfig 默认配置信息
func DefaultConfig() *Config {
	return &Config{
		Extra: map[string]string{
			"LeftBrackets":  "{", // 左括号转义
			"RightBrackets": "}", // 右括号转义
		},
	}
}
