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

package args

import "github.com/go-ceres/ceres/config"

type DDlArgs struct {
	Src      string         // 输入的文件
	Dist     string         // 输出文件目录
	Prefix   string         // 数据库前缀
	Cache    bool           // 是否使用缓存
	Strict   bool           // 严格模式
	Style    string         // 文件样式
	DataBase string         // 数据库名
	Entity   string         // 实体文件夹
	Home     string         // home文件夹
	Remote   string         // 远程地址
	Branch   string         // 远程分支
	Config   *config.Config // 样式配置
}

func NewDefaultDDlArgs() *DDlArgs {
	return &DDlArgs{}
}
