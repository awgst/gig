package repository

var RepositoryTemplate = `package repository

import (
	"gorm.io/gorm"
)

type {{.LowerName}}Repository struct {
	db *gorm.DB
}

func New{{.Name}}Repository(db *gorm.DB) {{.Name}}Repository {
	return &{{.LowerName}}Repository{db}
}
`

var RepositoryInterfaceTemplate = `package repository

type {{.Name}}Repository interface {
}`

var RepositoryInterfaceCRUDTemplate = `package repository

import (
	"{{.ModelPath}}/model"
)

type {{.Name}}Repository interface {
	Create({{.CamelCaseName}} *model.{{.ModelName}}) error
	Update({{.CamelCaseName}} *model.{{.ModelName}}) error
	Delete(id uint) error
	GetByID(id uint) (model.{{.ModelName}}, error)
	GetAll() ([]model.{{.ModelName}}, error)
}`

var RepositoryCRUDTemplate = `package repository

import (
	"{{.ModelPath}}/model"
	"gorm.io/gorm"
)

type {{.CamelCaseName}}Repository struct {
	db *gorm.DB
}

func New{{.Name}}Repository(db *gorm.DB) {{.Name}}Repository {
	return &{{.CamelCaseName}}Repository{db}
}

func (r *{{.CamelCaseName}}Repository) Create({{.CamelCaseName}} *model.{{.ModelName}}) error {
	return r.db.Create({{.CamelCaseName}}).Error
}

func (r *{{.CamelCaseName}}Repository) Update({{.CamelCaseName}} *model.{{.ModelName}}) error {
	return r.db.Save({{.CamelCaseName}}).Error
}

func (r *{{.CamelCaseName}}Repository) Delete(id uint) error {
	return r.db.Delete(&model.{{.ModelName}}{}, id).Error
}

func (r *{{.CamelCaseName}}Repository) GetByID(id uint) (model.{{.ModelName}}, error) {
	var {{.CamelCaseName}} model.{{.ModelName}}
	err := r.db.First(&{{.CamelCaseName}}, id).Error
	return {{.CamelCaseName}}, err
}

func (r *{{.CamelCaseName}}Repository) GetAll() ([]model.{{.ModelName}}, error) {
	var {{.CamelCaseName}}s []model.{{.ModelName}}
	err := r.db.Find(&{{.CamelCaseName}}s).Error
	return {{.CamelCaseName}}s, err
}
`
