func AutoMigrateGorm{{.camelName}}() error {
	return global.Db{{if .options}}.Set("gorm:table_options","{{.options}}"){{end}}.AutoMigrate(new({{.camelName}}))
}
