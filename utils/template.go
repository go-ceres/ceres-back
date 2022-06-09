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

package utils

import (
	"bytes"
	"go/format"
	"text/template"
)

type Template struct {
	name   string
	text   string
	format bool
}

// NewTemplate 新建模板解析器
func NewTemplate(name string) *Template {
	return &Template{
		name:   name,
		format: false,
	}
}

// FormatGo 格式化为go代码
func (t *Template) FormatGo(f bool) *Template {
	t.format = f
	return t
}

// Parse 设置模板
func (t *Template) Parse(text string) *Template {
	t.text = text
	return t
}

// Execute 执行模板解析
func (t *Template) Execute(data interface{}) (*bytes.Buffer, error) {
	tem, err := template.New(t.name).Parse(t.text)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if err = tem.Execute(buf, data); err != nil {
		return nil, err
	}

	if !t.format {
		return buf, nil
	}

	formatOutput, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, err
	}

	buf.Reset()
	buf.Write(formatOutput)
	return buf, nil
}
