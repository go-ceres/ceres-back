func (m *default{{.camelName}}Model) Update(ctx context.Context,{{.fieldName}} {{.fieldType}},param *{{if .entity}}{{.packageName}}.{{end}}{{.camelName}}) error {
    var po = new({{.camelName}})
    _ = structure.Copy(po, param)
	result:=Get{{.camelName}}Db(ctx,m.db).Where("{{.originalName}} = ?",{{.fieldName}}).Updates(po)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
