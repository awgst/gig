package handler

var HandlerPlainTemplate = `package handler

type {{.Name}}Handler interface {
}

type {{.CamelCaseName}}Handler struct {
}

func New{{.Name}}Handler() {{.Name}}Handler {
	return &{{.CamelCaseName}}Handler{}
}
`