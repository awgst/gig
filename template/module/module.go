package module

var ModuleTemplate = `package {{.PackageName}}

import (
	"{{.ModulePath}}/http/handler"
	"{{.ModulePath}}/service"
	"{{.ModulePath}}/repository"
	"gorm.io/gorm"
)

func BuildModule(db *gorm.DB) (handler.{{.HandlerName}}, service.{{.ServiceName}}, repository.{{.RepositoryName}}) {
	{{.CamelCaseName}}Repository := repository.New{{.RepositoryName}}(db)
	{{.CamelCaseName}}Service := service.New{{.ServiceName}}()
	{{.CamelCaseName}}Handler := handler.New{{.HandlerName}}()
	return {{.CamelCaseName}}Handler, {{.CamelCaseName}}Service, {{.CamelCaseName}}Repository
}
`

var ModuleCRUDTemplate = `package {{.PackageName}}

import (
	"{{.ProjectName}}/pkg/database"
	"{{.ModulePath}}/http/handler"
	"{{.ModulePath}}/service"
	"{{.ModulePath}}/repository"
)

func BuildModule(db database.Connection) (handler.{{.HandlerName}}, service.{{.ServiceName}}, repository.{{.RepositoryName}}) {
	{{.CamelCaseName}}Repository := repository.New{{.RepositoryName}}(db.SQL)
	{{.CamelCaseName}}Service := service.New{{.ServiceName}}({{.CamelCaseName}}Repository)
	{{.CamelCaseName}}Handler := handler.New{{.HandlerName}}({{.CamelCaseName}}Service)
	return {{.CamelCaseName}}Handler, {{.CamelCaseName}}Service, {{.CamelCaseName}}Repository
}`
