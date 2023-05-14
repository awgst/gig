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
	Create({{.LowerName}} *model.{{.Name}}) error
	Update({{.LowerName}} *model.{{.Name}}) error
	Delete({{.LowerName}} *model.{{.Name}}) error
	GetByID(id uint) (model.{{.Name}}, error)
	GetAll() ([]model.{{.Name}}, error)
}`

var RepositoryCRUDTemplate = `package repository

import (
	"{{.ModelPath}}/model"
	"gorm.io/gorm"
)

type {{.LowerName}}Repository struct {
	db *gorm.DB
}

func New{{.Name}}Repository(db *gorm.DB) {{.Name}}Repository {
	return &{{.LowerName}}Repository{db}
}

func (r *{{.LowerName}}Repository) Create({{.LowerName}} *model.{{.Name}}) error {
	return r.db.Create({{.LowerName}}).Error
}

func (r *{{.LowerName}}Repository) Update({{.LowerName}} *model.{{.Name}}) error {
	return r.db.Save({{.LowerName}}).Error
}

func (r *{{.LowerName}}Repository) Delete({{.LowerName}} *model.{{.Name}}) error {
	return r.db.Delete({{.LowerName}}).Error
}

func (r *{{.LowerName}}Repository) GetByID(id uint) (model.{{.Name}}, error) {
	var {{.LowerName}} model.{{.Name}}
	err := r.db.First(&{{.LowerName}}, id).Error
	return {{.LowerName}}, err
}

func (r *{{.LowerName}}Repository) GetAll() ([]model.{{.Name}}, error) {
	var {{.LowerName}}s []model.{{.Name}}
	err := r.db.Find(&{{.LowerName}}s).Error
	return {{.LowerName}}s, err
}
`
