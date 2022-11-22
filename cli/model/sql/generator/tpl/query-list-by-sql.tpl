func (m *default{{.camelName}}Model) QueryListBySql(ctx context.Context, params *{{if .entity}}{{.packageName}}.{{end}}QueryParam, sql string,args ...interface{} ) (*{{if .entity}}{{.packageName}}.{{end}}{{.camelName}}QueryResult,error) {
	db := Get{{.camelName}}Db(ctx, m.db)
	if len(sql) > 0 {
        db.Where(sql, args)
    }
    opt := gorm.GetQueryOption(params.QueryOptions)
	opt.OrderFields = append(opt.OrderFields, gorm.NewOrderField("{{.primary}}", gorm.OrderByDESC))
	db = db.Order(gorm.ParseOrder(opt.OrderFields))
	if len(opt.OrderFields) > 0 {
        db.Select(opt.SelectFields)
    }
	var po {{.camelName}}List
	pr, err := gorm.WrapPageQuery(ctx, db, params.PaginationParam, &po)
	if err != nil {
        return nil, err
    }
    var en {{if .entity}}{{.packageName}}.{{end}}{{.camelName}}List
    _ = objectx.Copy(&en, &po)
    qr := &{{if .entity}}{{.packageName}}.{{end}}{{.camelName}}QueryResult{
        List:       en,
        PageResult: pr,
    }
    return qr, nil
}
