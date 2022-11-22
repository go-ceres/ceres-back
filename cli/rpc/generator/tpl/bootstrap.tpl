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
    // registry
    global.Registry = etcd.ScanConfig("{{.ServiceName}}").Build(){{end}}
    // grpc client
    global.Client = grpc.ScanConfig("{{.ServiceName}}"){{if .Registry}}.WithRegistry(global.Registry){{end}}.Build(){{if .orm}}
    // orm
    global.Db = {{.orm.newFunc}}
    {{end}}
    return nil
}
