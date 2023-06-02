package model

var ModelTemplate = `package model

import (
	"time"
)

type {{.Name}} struct {
	ID uint
	CreatedAt time.Time
	UpdatedAt time.Time
}`

var PlainModelTemplate = `package model

type {{.Name}} struct {
}`
