package global

import ({{.imports}}{{if .Registry}}
    "github.com/go-ceres/go-ceres/registry"{{end}}
    "github.com/go-ceres/go-ceres/client/grpc"
)

var({{if .Registry}}
    Registry registry.Registry{{end}}
    Client *grpc.Client{{if .orm}}
   Db *{{.orm.name}}.DB
   {{end}}
)
