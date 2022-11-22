func new{{.camelName}}Model() *default{{.camelName}}Model {
	return &default{{.camelName}}Model{{.extra.LeftBrackets}}{{if .hasGormDb}}
		db: global.Db,{{end}}{{if .hasCache}}{{if .cache}}
		cache: global.Cache,{{end}}{{end}}
	}
}
