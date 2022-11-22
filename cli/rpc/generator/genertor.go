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

import (
	config2 "github.com/go-ceres/ceres/config"
	"github.com/go-ceres/ceres/env"
	"github.com/go-ceres/ceres/utils/console"
	"github.com/go-ceres/go-ceres/logger"
)

// Generator 生成器结构体
type Generator struct {
	config  *config2.Config
	log     console.Console
	verbose bool
}

// NewGenerator 创建生成器
func NewGenerator(style string, verbose bool) *Generator {
	conf, err := config2.NewConfig(style)
	if err != nil {
		logger.DPanic(err)
	}
	return &Generator{
		config:  conf,
		log:     console.NewColorConsole(verbose),
		verbose: verbose,
	}
}

// Prepare 前置步骤安装没有安装的依赖
func (g *Generator) Prepare() error {
	return env.Prepare(true, true, g.verbose)
}
