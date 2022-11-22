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

package parser

// Info API 元数据信息
type Info struct {
	Title       string // 标题
	Version     string // 版本
	Description string // 描述

}

// Contact 所开放的 API 的联系人信息
type Contact struct {
	Name  string
	Url   string
	Email string
}

// License 开源协议
type License struct {
	Name string
	Url  string
}
