package service

var ServiceTemplate = `package service

type {{.Name}}Service interface {
}

type {{.CamelCaseName}}Service struct {
}

func New{{.Name}}Service() {{.Name}}Service {
	return &{{.CamelCaseName}}Service{}
}
`

var ServiceCRUDTemplate = `package service

import (
	"{{.RequestPath}}/request"
	"{{.RepositoryPath}}/repository"
	"{{.ModelPath}}/model"
)

type {{.Name}}Service interface {
	Create({{.CamelCaseName}}Request *request.Create{{.Name}}Request) error
	Update(id uint, {{.CamelCaseName}}Request *request.Update{{.Name}}Request) error
	Delete(id uint) error
	GetByID(id uint) (model.{{.Name}}, error)
	GetAll() ([]model.{{.Name}}, error)
}

type {{.CamelCaseName}}Service struct {
	repository repository.{{.Name}}Repository
}

func New{{.Name}}Service(repository repository.{{.Name}}Repository) {{.Name}}Service {
	return &{{.CamelCaseName}}Service{repository}
}

func (s *{{.CamelCaseName}}Service) Create({{.CamelCaseName}}Request *request.Create{{.Name}}Request) error {
	{{.CamelCaseName}} := model.{{.Name}}{
	}
	return s.repository.Create(&{{.CamelCaseName}})
}

func (s *{{.CamelCaseName}}Service) Update(id uint, {{.CamelCaseName}}Request *request.Update{{.Name}}Request) error {
	{{.CamelCaseName}} := model.{{.Name}}{
		ID: id,
	}
	return s.repository.Update(&{{.CamelCaseName}})
}

func (s *{{.CamelCaseName}}Service) Delete(id uint) error {
	return s.repository.Delete(id)
}

func (s *{{.CamelCaseName}}Service) GetByID(id uint) (model.{{.ModelName}}, error) {
	return s.repository.GetByID(id)
}

func (s *{{.CamelCaseName}}Service) GetAll() ([]model.{{.ModelName}}, error) {
	return s.repository.GetAll()
}
`
