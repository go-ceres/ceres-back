package main

import (
	{{.imports}}
	_ "github.com/go-ceres/go-ceres/cmd/plugins/config/source/{{.configSource}}"
	"github.com/go-ceres/go-ceres/logger"
)

func main() {
	boot := {{.bootBase}}.NewBootstrap()
	err := boot.Run()
	if err != nil {
		logger.Panic("run app error",logger.FieldErr(err))
	}
}
