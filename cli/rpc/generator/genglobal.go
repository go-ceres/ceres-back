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
	"fmt"
	"github.com/go-ceres/ceres/cli/rpc/config"
	"github.com/go-ceres/ceres/utils/formatc"
	"github.com/go-ceres/ceres/utils/pathc"
	"github.com/go-ceres/ceres/utils/templatec"
	"path/filepath"
	"strings"
)

//go:embed tpl/global.tpl
var globalTemplate string

// GenGlobal 生成全局上下文
func (g *Generator) GenGlobal(ctx DirContext, conf *config.Config) error {
	dir := ctx.GetGlobal()
	globalFilename, err := formatc.FileNamingFormat(g.config.Style, "global")
	if err != nil {
		return err
	}
	imports := make([]string, 0)
	// orm框架
	if conf.Orm != nil {
		imports = append(imports, fmt.Sprintf(`"%s"`, filepath.Join("github.com/go-ceres/go-ceres/store", conf.Orm["name"])))
	}

	fileName := filepath.Join(dir.Filename, globalFilename+".go")
	text, err := pathc.LoadTpl(category, globalTemplateFile, globalTemplate)
	if err != nil {
		return err
	}

	return templatec.With("global").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"Registry": conf.Registry,
		"Extra":    extra,
		"imports":  strings.Join(imports, "\n"),
		"orm":      conf.Orm,
	}, fileName, false)
}
