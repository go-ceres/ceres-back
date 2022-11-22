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
	"github.com/go-ceres/ceres/utils/formatx"
	"github.com/go-ceres/ceres/utils/pathx"
	"github.com/go-ceres/ceres/utils/templatex"
	"path/filepath"
	"strings"
)

//go:embed tpl/main.tpl
var mainTemplate string

// ServerTemplateData 入口文件
type ServerTemplateData struct {
	Server    string
	ServerPkg string
}

// GenMain 生成main文件
func (g *Generator) GenMain(ctx DirContext, conf *config.Config) error {
	mainFilename, err := formatx.FileNamingFormat(g.config.Style, "main")
	if err != nil {
		return err
	}
	fileName := filepath.Join(ctx.GetMain().Filename, fmt.Sprintf("%v.go", mainFilename))
	imports := make([]string, 0)
	bootImport := fmt.Sprintf(`"%s"`, ctx.GetBoot().Package)
	imports = append(imports, bootImport)

	content, err := pathx.LoadTpl(category, mainTemplateFileFile, mainTemplate)
	if err != nil {
		return err
	}
	return templatex.With("main").GoFmt(true).Parse(content).SaveTo(map[string]interface{}{
		"imports":      strings.Join(imports, "\n"),
		"configSource": conf.ConfigSource,
		"bootBase":     ctx.GetBoot().Base,
	}, fileName, false)
}
