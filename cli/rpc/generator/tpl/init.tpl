package boot

import (
    {{.imports}}
    "github.com/go-ceres/go-ceres/client/grpc"{{if .Registry}}
    "github.com/go-ceres/go-ceres/registry"
    "github.com/go-ceres/go-ceres/registry/{{.Registry}}"{{end}}
)

func Init() error {{.Extra.LeftBrackets}}{{if .Registry}}
    // 注册中心
    global.Registry = etcd.ScanConfig("{{.Registry}}").Build(){{end}}
    // Grpc连接客户端
    global.Client = grpc.ScanConfig("default"){{if .Registry}}.WithRegistry(global.Registry){{end}}.Build()
}
