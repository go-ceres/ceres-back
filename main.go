package main

import (
	"fmt"
	"github.com/go-ceres/ceres/build"
	new2 "github.com/go-ceres/ceres/new"
	"github.com/go-ceres/cli/v2"
	"github.com/go-ceres/go-ceres/logger"
	"os"
	"runtime"
)

var (
	version  = "0.1.0"
	commands = []*cli.Command{
		// 构建相关
		{
			Name:   "build",
			Usage:  "compile go-ceres application",
			Flags:  build.Flags,
			Action: build.Run,
		},
		{
			Name:        "new",
			Usage:       "create go-ceres application",
			Subcommands: new2.Command,
		},
	}
)

func main() {
	app := cli.NewApp()
	app.Usage = "a cli tool for go-ceres"
	app.Description = "a cli tool for go-ceres"
	app.Version = fmt.Sprintf("%s %s/%s", version, runtime.GOOS, runtime.GOARCH)
	app.Commands = commands
	if err := app.Run(os.Args); err != nil {
		logger.Error(err)
		os.Exit(0)
	}
}
