package {{.packageName}}

import (
	"context"
	{{.imports}}
)

type {{.logicName}} struct {
	ctx    context.Context{{if .dtoAndService}}
	{{.unTitleServiceName}} *{{.servicePackageName}}.{{.serviceName}}{{end}}
}

func New{{.logicName}}() *{{.logicName}} {
	return &{{.logicName}}{{.extra.LeftBrackets}}{{if .dtoAndService}}
        {{.unTitleServiceName}}: {{.servicePackageName}}.New{{.serviceName}}(),{{end}}
	}
}
{{.functions}}
