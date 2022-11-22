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

package model

import (
	"github.com/go-ceres/ceres/cli/common"
	"github.com/go-ceres/ceres/cli/model/sql/action"
	"github.com/go-ceres/cli/v2"
)

var Flags []cli.Flag

var Commands = []*cli.Command{
	{
		Name:  "mysql",
		Usage: "generate model code for mysql",
		Subcommands: []*cli.Command{
			{
				Name:   "ddl",
				Flags:  append(action.MysqlDDlFlags, common.Flags...),
				Action: action.MysqlDDl,
			},
		},
	},
}
