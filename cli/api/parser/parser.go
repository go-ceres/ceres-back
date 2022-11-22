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

import (
	"errors"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-ceres/ceres/cli/api/config"
	"net/url"
	"reflect"
)

type Parser struct {
	spec *openapi3.T
}

// NewParser 创建解析器
func NewParser(c *config.Config) (*Parser, error) {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	u, err := url.Parse(c.SwaggerFile)
	var swagger *openapi3.T
	var loadErr error
	if err == nil && u.Scheme != "" && u.Host != "" {
		swagger, loadErr = loader.LoadFromURI(u)
	} else {
		swagger, loadErr = loader.LoadFromFile(c.SwaggerFile)
	}
	if loadErr != nil {
		return nil, err
	}
	fmt.Println(swagger)
	// 解析info
	return &Parser{
		spec: swagger,
	}, err
}

// Parse 解析为swagger
func (p *Parser) Parse() (*Swagger, error) {
	res := &Swagger{}
	// 解析info
	info, err := p.parseInfo()
	if err != nil {
		return nil, err
	}
	res.Info = info

	// 解析Schemas（对应go里面的结构体）
	schemas, err := p.parseSchemas()
	if err != nil {
		return nil, err
	}
	res.Routers = schemas

	return nil, nil
}

// parseInfo 解析info
func (p *Parser) parseInfo() (*Info, error) {
	res := &Info{
		Title:       p.spec.Info.Title,
		Version:     p.spec.Info.Version,
		Description: p.spec.Info.Description,
	}
	return res, nil
}

// parseSchemas 解析schema
func (p *Parser) parseSchemas() (map[string][]*Router, error) {
	var res = map[string][]*Router{}
	var methods = []string{"Post", "Get", "Put", "Delete", "Option"}
	// 设置一个默认组
	res["default"] = []*Router{}
	for s, ref := range p.spec.Paths {
		value := reflect.ValueOf(ref).Elem()
		for _, method := range methods {
			Operation, ok := value.FieldByName(method).Interface().(*openapi3.Operation)
			if ok && Operation != nil {
				group := "default"
				// 定义路由
				router := &Router{
					Path: s,
				}
				// 如果有分组
				if len(Operation.Tags) > 0 {
					group = Operation.Tags[0]
				}
				router.Method = method
				// 如果没有OperationID，表示没有handler
				if Operation.OperationID == "" {
					return nil, errors.New(fmt.Sprintf("no OperationID in current path %s method %s", s, method))
				}
				router.OperationID = Operation.OperationID
				// 简介
				router.Summary = Operation.Summary
				// 描述
				router.Description = Operation.Description

				// 获取当前组已经存在的所有路由
				routers, ok := res[group]
				if !ok {
					routers = []*Router{}
				}

				// 添加路由到当前路由组
				routers = append(routers, router)
				res[group] = routers
			}
		}
	}
	return res, nil
}
