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

type Config struct {
	Namespace    string // 项目空间，例如：github.com/go-ceres
	SwaggerFile  string // swagger 描述文件
	Name         string // 项目名称，例如：demo
	Dir          string // 代码输出，例如: .
	WorkDir      string // 项目目录。绝对路径
	ConfigSource string // 配置组件的类型，例如：file
	Registry     string // 注册中心，例如：etcd
	Extra        map[string]string
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
