package {{.package}}

import "github.com/go-ceres/go-ceres/store/gorm"

type (
    {{.camelName}} struct{
        {{.fields}}
    }

    {{.camelName}}List []*{{.camelName}}

    QueryParam struct {
        gorm.PaginationParam
        gorm.QueryOptions
    }

    {{.camelName}}QueryResult struct {
        PageResult *gorm.PaginationResult
        List                  {{.camelName}}List
    }
)
