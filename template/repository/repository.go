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

var RepositoryCRUDGormTemplate = `package repository

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

var RepositoryCRUDTemplate = `package repository

import (
	"{{.ModelPath}}/model"
	"database/sql"
)

const (
	create{{.Name}}Query = "INSERT INTO {{.TableName}} () VALUES () RETURNING id"
	update{{.Name}}Query = "UPDATE {{.TableName}} SET () WHERE id = ?"
	delete{{.Name}}Query = "DELETE FROM {{.TableName}} WHERE id = ?"
	get{{.Name}}ByIDQuery = "SELECT * FROM {{.TableName}} WHERE id = ?"
	getAll{{.Name}}Query = "SELECT * FROM {{.TableName}}"
)

type {{.CamelCaseName}}Repository struct {
	db *sql.DB
}

func New{{.Name}}Repository(db *sql.DB) {{.Name}}Repository {
	return &{{.CamelCaseName}}Repository{db}
}

func (r *{{.CamelCaseName}}Repository) Create({{.CamelCaseName}} *model.{{.ModelName}}) error {
	_, err := r.db.Exec(create{{.Name}}Query)
	return err
}

func (r *{{.CamelCaseName}}Repository) Update({{.CamelCaseName}} *model.{{.ModelName}}) error {
	_, err := r.db.Exec(update{{.Name}}Query, {{.CamelCaseName}}.ID)
	return err
}

func (r *{{.CamelCaseName}}Repository) Delete(id uint) error {
	_, err := r.db.Exec(delete{{.Name}}Query, id)
	return err
}

func (r *{{.CamelCaseName}}Repository) GetByID(id uint) (model.{{.ModelName}}, error) {
	var {{.CamelCaseName}} model.{{.ModelName}}
	err := r.db.
		QueryRow(get{{.Name}}ByIDQuery, id).
		Scan(
			&{{.CamelCaseName}}.ID,
			&{{.CamelCaseName}}.CreatedAt,
			&{{.CamelCaseName}}.UpdatedAt,
		)

	return {{.CamelCaseName}}, err
}

func (r *{{.CamelCaseName}}Repository) GetAll() ([]model.{{.ModelName}}, error) {
	var {{.CamelCaseName}}s []model.{{.ModelName}}
	rows, err := r.db.Query(getAll{{.Name}}Query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var {{.CamelCaseName}} model.{{.ModelName}}
		err := rows.Scan(
			&{{.CamelCaseName}}.ID,
			&{{.CamelCaseName}}.CreatedAt,
			&{{.CamelCaseName}}.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		{{.CamelCaseName}}s = append({{.CamelCaseName}}s, {{.CamelCaseName}})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return {{.CamelCaseName}}s, nil
}
`
