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

package common

import "github.com/go-ceres/cli/v2"

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:  "home",
		Value: "",
		Usage: "The ceres home path of the template",
	},
	&cli.StringFlag{
		Name:  "branch",
		Value: "",
		Usage: "The branch of the remote repo，it does work with --remote",
	},
	&cli.BoolFlag{
		Name:  "verbose",
		Value: false,
		Usage: "Enable log output",
	},
	&cli.StringFlag{
		Name:  "remote",
		Value: "",
		Usage: "The remote git repo of the template",
	},
}
