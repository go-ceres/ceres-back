package bootstrap

import (
    {{.imports}}
	"github.com/go-ceres/go-ceres"
	"github.com/go-ceres/go-ceres/logger"
	"github.com/go-ceres/go-ceres/client/grpc"{{if .Registry}}
    "github.com/go-ceres/go-ceres/registry/{{.Registry}}"{{end}}
)

type Bootstrap struct {
	ceres.Engine
}

func NewBootstrap() *Bootstrap {
	boot := new(Bootstrap)
	err := boot.SetInit(Init).MustSetup(
		boot.start{{.ServiceName}}Grpc,
	)
	if err != nil {
		logger.Panicd("must setup error",logger.FieldErr(err))
	}
	return boot
}

func Init() error {{.Extra.LeftBrackets}}{{if .Registry}}
    global.Registry = etcd.ScanConfig("{{.unTitleServiceName}}").Build(){{end}}
    global.Client = grpc.ScanConfig("{{.unTitleServiceName}}"){{if .Registry}}.WithRegistry(global.Registry){{end}}.Build(){{range .components}}
    global.{{.GlobalName}} = {{.InitStr}}{{end}}
    return nil
}
