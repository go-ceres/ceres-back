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
	_ "embed"
)

////go:embed tpl/init.tpl
//var bootTemplate string
//
//// GenBoot 生成初始化文件
//func (g *Generator) GenBoot(ctx DirContext, proto model.Proto, conf *config.Config) error {
//	dir := ctx.GetBoot()
//	bootFilename, err := formatc.FileNamingFormat(g.config.Style, "boot")
//	if err != nil {
//		return err
//	}
//	imports := make([]string, 0)
//	globalImport := fmt.Sprintf(`"%s"`, ctx.GetGlobal().Package)
//	imports = append(imports, globalImport)
//	fileName := filepath.Join(dir.Filename, bootFilename+".go")
//	text, err := pathc.LoadTpl(category, initTemplateFile, bootTemplate)
//	if err != nil {
//		return err
//	}
//
//	return templatec.With("boot").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
//		"Registry": conf.Registry,
//		"Extra":    extra,
//		"imports":  strings.Join(imports, "\n"),
//	}, fileName, false)
//}
