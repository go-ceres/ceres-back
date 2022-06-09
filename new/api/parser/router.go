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

type Router struct {
	Path         string      // 路径
	Method       string      // 请求方法
	Summary      string      // 简介
	Description  string      // 描述
	OperationID  string      // 操作方法
	HeaderParams []Parameter // 头请求参数
	BodyParams   []Parameter // body请求参数
}
