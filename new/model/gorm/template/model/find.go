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

var Find = `
func (m *{{.camelName}}Model) FindOne(ctx context.Context,params {{.camelName}}) (*{{.camelName}},error) {
	var resp {{.camelName}}
	db := Get{{.camelName}}Db(ctx,m.DB).Where(&params)
	_,err := gorm.FindOne(db,&resp)
	if err != nil {
		return nil,err
	}
	return &resp,nil
}
`
