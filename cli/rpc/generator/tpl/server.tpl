{{.head}}

package server

import (
	{{if .notStream}}"context"{{end}}

	{{.imports}}
	"google.golang.org/grpc"
)

type {{.server}}Server struct {
	{{.unimplementedServer}}{{range .logicList}}
	{{.UnTitleName}} *{{.LogicPackage}}.{{.LogicName}}{{end}}
}

func Register{{.server}}Server(srv *grpc.Server) {
	{{.pbPackage}}.Register{{.server}}Server(srv, New{{.server}}Server())
}

func New{{.server}}Server() *{{.server}}Server {
	return &{{.server}}Server{{.extra.LeftBrackets}}{{range .logicList}}
        {{.UnTitleName}}:{{.LogicPackage}}.New{{.LogicName}}(),{{end}}
	}
}

{{.funcs}}
