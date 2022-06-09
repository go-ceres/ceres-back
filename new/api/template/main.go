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

package template

var MainTemplate = `package main

import (
	{{.imports}}
	_ "github.com/go-ceres/go-ceres/cmd/plugins/config/source/{{.configSource}}"
	"github.com/go-ceres/go-ceres/logger"
)

func main() {
	eng := engine.NewEngine()
	err := eng.Run()
	if err != nil {
		logger.Panic("run app error",logger.FieldErr(err))
	}
}
`
