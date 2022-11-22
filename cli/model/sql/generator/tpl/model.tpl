package {{.pkg}}

type {{.upperStartCamelObject}}Model struct {
	*defaultMembersModel
}

func New{{.upperStartCamelObject}}Model() *{{.upperStartCamelObject}}Model {
	return &{{.upperStartCamelObject}}Model{
		default{{.upperStartCamelObject}}Model: new{{.upperStartCamelObject}}Model(),
	}
}
