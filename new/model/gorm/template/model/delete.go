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

package model

var Delete = `
func (m *{{.camelName}}Model) Delete(ctx context.Context,{{.fieldName}}s []{{.fieldType}}) error {
	result:=Get{{.camelName}}Db(ctx,m.DB).Where("{{.originalName}} IN (?)",{{.fieldName}}s).Delete({{.camelName}}{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
`
