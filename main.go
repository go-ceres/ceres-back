package main

import (
	"fmt"
	"github.com/go-ceres/ceres/cli/build"
	"github.com/go-ceres/ceres/cli/common"
	"github.com/go-ceres/ceres/cli/model"
	"github.com/go-ceres/ceres/cli/rpc"
	"github.com/go-ceres/ceres/internal/version"
	"github.com/go-ceres/cli/v2"
	"github.com/go-ceres/go-ceres/logger"
	"os"
	"runtime"
)

var (
	commands = []*cli.Command{
		// 构建相关
		{
			Name:   "build",
			Usage:  "compile go-ceres application",
			Flags:  build.Flags,
			Action: build.Run,
		},
		{
			Name:  "api",
			Usage: "generate api code",
		},
		{
			Name:        "rpc",
			Usage:       "generate rpc code",
			Flags:       append(common.Flags, rpc.Flags...),
			Subcommands: rpc.Commands,
		},
		{
			Name:        "model",
			Usage:       "generate model code",
			Flags:       model.Flags,
			Subcommands: model.Commands,
		},
	}
)

func main() {
	app := cli.NewApp()
	app.Usage = "a cli tool for go-ceres"
	app.Description = "a cli tool for go-ceres"
	app.UseShortOptionHandling = true
	app.Version = fmt.Sprintf("%s %s/%s", version.BuildVersion, runtime.GOOS, runtime.GOARCH)
	app.Commands = commands
	app.Flags = append(app.Flags)
	app.ExitErrHandler = func(context *cli.Context, err error) {
		_ = cli.ShowCommandHelp(context, context.Command.Name)
		logger.Error(err)
	}
	_ = app.Run(os.Args)
}
