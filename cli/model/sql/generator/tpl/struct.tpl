type (
	default{{.camelName}}Model struct {
		db *gorm.DB{{if .cache}}
		cache cache.Cache{{end}}
	}

	{{.camelName}} struct{
		{{.fields}}
	}

	{{.camelName}}List []*{{.camelName}}{{if .noEntity}}

	QueryParam struct {
        gorm.PaginationParam
        gorm.QueryOptions
    }

	{{.camelName}}QueryResult struct {
        PageResult *gorm.PaginationResult
        List                  {{.camelName}}List
    }{{end}}
)
