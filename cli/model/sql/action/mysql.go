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

package action

import (
	"errors"
	"github.com/go-ceres/ceres/cli/model/sql/args"
	"github.com/go-ceres/ceres/cli/model/sql/generator"
	"github.com/go-ceres/ceres/config"
	"github.com/go-ceres/ceres/utils"
	"github.com/go-ceres/ceres/utils/pathc"
	"github.com/go-ceres/cli/v2"
)

var MysqlDDlFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "src",
		Value:   "",
		Aliases: []string{"s"},
		Usage:   "The path or path globbing patterns of the ddl",
	},
	&cli.StringFlag{
		Name:    "dist",
		Value:   "",
		Aliases: []string{"d"},
		Usage:   "The target dir",
	},
	&cli.BoolFlag{
		Name:    "cache",
		Value:   false,
		Aliases: []string{"c"},
		Usage:   "The target dir",
	},
	&cli.BoolFlag{
		Name:  "strict",
		Value: false,
		Usage: "The strict mode is enabled",
	},
	&cli.StringFlag{
		Name:    "entity",
		Value:   "",
		Aliases: []string{"e"},
		Usage:   "entity out path",
	},
	&cli.StringFlag{
		Name:  "prefix",
		Value: "",
		Usage: "table name prefix",
	},
	&cli.StringFlag{
		Name:  "style",
		Value: "",
		Usage: "The target dir",
	},
}

// MysqlDDl 根据sql文件创建model
func MysqlDDl(ctx *cli.Context) error {
	ddlArgs := args.NewDefaultDDlArgs()
	ddlArgs.Src = ctx.String("src")
	ddlArgs.Dist = ctx.String("dist")
	ddlArgs.Cache = ctx.Bool("cache")
	ddlArgs.Entity = ctx.String("entity")
	ddlArgs.Style = ctx.String("style")
	ddlArgs.DataBase = ctx.String("database")
	ddlArgs.Home = ctx.String("home")
	ddlArgs.Remote = ctx.String("remote")
	ddlArgs.Branch = ctx.String("branch")
	ddlArgs.Strict = ctx.Bool("strict")
	ddlArgs.Prefix = ctx.String("prefix")
	verbose := ctx.Bool("verbose")
	if len(ddlArgs.Src) == 0 {
		return errors.New("expected path or path globbing patterns, but nothing found")
	}
	if len(ddlArgs.Remote) > 0 {
		repo, _ := utils.CloneIntoGitHome(ddlArgs.Remote, ddlArgs.Branch)
		if len(repo) > 0 {
			ddlArgs.Home = repo
		}
	}
	if len(ddlArgs.Home) > 0 {
		pathc.RegisterCeresHome(ddlArgs.Home)
	}
	conf, err := config.NewConfig(ddlArgs.Style)
	if err != nil {
		return err
	}
	ddlArgs.Config = conf
	files, err := pathc.MatchFiles(ddlArgs.Src)
	if err != nil {
		return err
	}
	// 创建生成器
	g := generator.NewGenerator("go_ceres", verbose)
	for _, file := range files {
		if err := g.GeneratorFromDDl(file, ddlArgs); err != nil {
			return err
		}
	}
	return nil
}
