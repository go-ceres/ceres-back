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

import (
	"bytes"
	"github.com/alecthomas/template"
	"go/format"
	"os"
)

// write 写入文件
func write(data interface{}, file, tmpl string) error {
	fns := template.FuncMap{}
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	parse, err := template.New("f").Funcs(fns).Parse(tmpl)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err := parse.Execute(buf, data); err != nil {
		return err
	}
	formatOutput, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}
	_, err = f.Write(formatOutput)
	if err != nil {
		return err
	}
	return nil
}
