func (m *default{{.camelName}}Model) FindOne(ctx context.Context,params {{if .entity}}{{.packageName}}.{{end}}{{.camelName}}) (*{{if .entity}}{{.packageName}}.{{end}}{{.camelName}},error) {
    var {{if .entity}}po{{else}}en{{end}} = new({{.camelName}}){{if .entity}}
    var en = new({{.packageName}}.{{.camelName}}){{end}}
	db := Get{{.camelName}}Db(ctx,m.db).Where(&params)
	_,err := gorm.FindOne(db,{{if .entity}}po{{else}}en{{end}})
	if err != nil {
		return nil,err
	}{{if .entity}}
	_ = structure.Copy(en,po)
	return en,nil{{else}}
	return po,nil
	{{end}}
}
