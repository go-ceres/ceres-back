package bootstrap

import (
	{{.imports}}
	"github.com/go-ceres/go-ceres/server/grpc"
)

func (boot *Bootstrap) start{{.ServiceName}}Grpc() (func(),error) {
	srv := grpc.ScanConfig("{{.BaseServiceName}}").Build()
	{{if .Registry}}
	boot.SetRegistry(global.Registry){{end}}
	{{range .serverNames}}{{.ServerPkg}}.Register{{.Server}}Server(srv.Server)
	{{end}}
	return func() {},boot.Server(srv)
}
