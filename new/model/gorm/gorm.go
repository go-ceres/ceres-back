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

package gorm

import (
	"errors"
	"github.com/go-ceres/ceres/new/model/gorm/gen"
	"github.com/go-ceres/cli/v2"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:  "src",
		Usage: "sql file path",
	},
	&cli.StringFlag{
		Name:  "dist",
		Usage: "model file output directory",
	},
	&cli.StringFlag{
		Name:  "prefix",
		Usage: "table prefix",
	},
	&cli.BoolFlag{
		Name:  "cache",
		Usage: "Use cache",
	},
}

func Run(ctx *cli.Context) error {
	src := ctx.String("src")
	if len(src) == 0 {
		return errors.New("expected path or path globbing patterns, but nothing found")
	}
	conf := gen.DefaultConfig()
	conf.Src = src
	conf.Dist = ctx.String("dist")
	conf.Prefix = ctx.String("prefix")
	conf.Cache = ctx.Bool("cache")
	generator, err := conf.BuildGenerator()
	if err != nil {
		return err
	}
	return generator.Start()
}
