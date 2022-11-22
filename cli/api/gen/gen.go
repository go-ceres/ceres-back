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

package gen

import (
	"fmt"
	"github.com/go-ceres/ceres/cli/api/config"
	"github.com/go-ceres/ceres/cli/api/parser"
)

type Generator struct {
	config *config.Config // 配置信息
	parser *parser.Parser // swagger解析器
}

// NewGenerator 创建代码生成器
func NewGenerator(conf *config.Config) (*Generator, error) {
	newParser, err := parser.NewParser(conf)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	fmt.Println(newParser)
	return &Generator{
		parser: newParser,
		config: conf,
	}, nil
}

// Start 开始生成代码
func (g *Generator) Start() error {
	swagger, err := g.parser.Parse()
	if err != nil {
		return err
	}
	fmt.Println(swagger)
	return nil
}
