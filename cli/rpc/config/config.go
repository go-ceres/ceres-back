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
	Multiple        bool              // 是否支持生成多个rpc服务
	Dist            string            // 项目输出路径，例如: .
	Home            string            // ceres的模板目录
	WorkDir         string            // 项目目录。绝对路径
	GrpcOut         string            // grpc_out 路径
	GoOut           string            // go_out 路径
	ConfigSource    string            // 配置组件的类型，例如：file
	Registry        string            // 注册中心，例如：etcd
	ProtoPath       string            // proto_path 参数
	ProtoGenGrpcDir string            // proto的grpc输出目录
	ProtoGenGoDir   string            // proto的go文件输出夹
	GoOpt           []string          // protoc 的 opt参数
	GoGrpcOpt       []string          // protoc 的go-grpc_opt 参数
	Plugins         []string          // protoc的插件
	ProtoFile       string            // proto文件
	IsGooglePlugin  bool              // 指示proto文件是否由google插件生成的标志
	ProtocCmd       string            // protoc的命令
	DtoAndService   bool              // 是否生成对应的dto与service
	Orm             map[string]string // orm框架
}

// DefaultConfig 默认配置信息
func DefaultConfig() *Config {
	return &Config{}
}
