func (m *default{{.camelName}}Model) Create(ctx context.Context, param *{{if .entity}}{{.packageName}}.{{end}}{{.camelName}}) error {
	var po = new({{.camelName}})
	_ = objectx.Copy(po,param)
	result := Get{{.camelName}}Db(ctx,m.db).Create(po)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
