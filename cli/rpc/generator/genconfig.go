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
	"github.com/go-ceres/ceres/cli/rpc/parser/model"
	"github.com/go-ceres/ceres/utils/formatc"
	"github.com/go-ceres/ceres/utils/pathc"
	"github.com/go-ceres/ceres/utils/templatec"
	"path/filepath"
)

//go:embed tpl/config.tpl
var configTemplate string

// GenConfig 生成配置文件
func (g *Generator) GenConfig(ctx DirContext, _ model.Proto, conf *config.Config) error {
	dir := ctx.GetConfig()
	configFilename, err := formatc.FileNamingFormat(g.config.Style, "config")
	if err != nil {
		return err
	}
	// 文件名
	fileName := filepath.Join(dir.Filename, fmt.Sprintf("%v.toml", configFilename))
	// 获取模板内容
	context, err := pathc.LoadTpl(category, configTemplateFileFile, configTemplate)
	if err != nil {
		return err
	}
	return templatec.With("etc").Parse(context).SaveTo(map[string]interface{}{
		"serviceName": ctx.GetServiceName().UnTitle(),
		"registry":    conf.Registry,
		"orm":         conf.Orm,
	}, fileName, false)
}
