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

const (
	category                    = "rpc"
	configTemplateFileFile      = "config.tpl"
	logicTemplateFileFile       = "logic.tpl"
	serviceTemplateFileFile     = "service.tpl"
	dtoTemplateFileFile         = "dto.tpl"
	serviceFuncTemplateFileFile = "service-func.tpl"
	serverTemplateFile          = "server.tpl"
	serverFuncTemplateFile      = "server-func.tpl"
	globalTemplateFile          = "global.tpl"
	logicFuncTemplateFileFile   = "logic-func.tpl"
	bootTemplateFileFile        = "boot.tpl"
	serverBootTemplateFileFile  = "server-boot.tpl"
	mainTemplateFileFile        = "main.tpl"
)

var extra = map[string]string{
	"LeftBrackets":  "{", // 左括号转义
	"RightBrackets": "}", // 右括号转义
}
