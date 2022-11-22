package {{.packageName}}

import (
	"context"
	{{.imports}}
)

type {{.serviceName}} struct {
}

func New{{.serviceName}}() *{{.serviceName}} {
	return &{{.serviceName}}{}
}
{{.functions}}
