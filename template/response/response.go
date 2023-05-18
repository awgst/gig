package response

var ResponseTemplate = `package response

import (
	"{{.ModelPath}}/model"
	"time"
)

type {{.Name}}Response struct {
	ID uint ` + "`json:\"id\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
}

func (r *{{.Name}}Response) MakeOne({{.CamelCaseName}} model.{{.ModelName}}) {{.Name}}Response {
	return {{.Name}}Response{}
}

func (r *{{.Name}}Response) Make({{.CamelCaseName}} []model.{{.ModelName}}) []{{.Name}}Response {
	{{.CamelCaseName}}Responses := make([]{{.Name}}Response, len({{.CamelCaseName}}))
	for k, v := range {{.CamelCaseName}} {
		{{.CamelCaseName}}Responses[k] = r.MakeOne(v)
	}
	return {{.CamelCaseName}}Responses
}
`

var PlainResponseTemplate = `package response

type {{.Name}}Response struct {
}`
