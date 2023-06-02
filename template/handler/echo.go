package handler

var handlerEchoCRUDTemplate = `package handler

import (
	"strconv"
	"net/http"
	"github.com/labstack/echo/v4"
	"{{.ServicePath}}/service"
	"{{.ResponsePath}}/response"
	"{{.RequestPath}}/request"
)

type {{.Name}}Handler interface {
	Create(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	Show(c echo.Context) error
	Index(c echo.Context) error
}

type {{.CamelCaseName}}Handler struct {
	service service.{{.Name}}Service
}

var {{.CamelCaseName}}Response response.{{.Name}}Response

func New{{.Name}}Handler(service service.{{.Name}}Service) {{.Name}}Handler {
	return &{{.CamelCaseName}}Handler{service}
}

func (h *{{.CamelCaseName}}Handler) Create(c echo.Context) error {
	var {{.LowerName}}Request request.Create{{.Name}}Request
	if err := c.Bind({{.LowerName}}Request); err != nil {
		return c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if err := h.service.Create(&{{.LowerName}}Request); err != nil {
		return c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, gin.H{})
}

func (h *{{.CamelCaseName}}Handler) Update(c echo.Context) error {
	var {{.LowerName}}Request request.Update{{.Name}}Request
	if err := c.Bind({{.LowerName}}Request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(uint(id), &{{.LowerName}}Request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *{{.CamelCaseName}}Handler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *{{.CamelCaseName}}Handler) Show(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	{{.LowerName}}, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, {{.CamelCaseName}}Response.MakeOne({{.CamelCaseName}}))
}

func (h *{{.CamelCaseName}}Handler) Index(c echo.Context) error {
	{{.CamelCaseName}}s, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, {{.CamelCaseName}}Response.Make({{.CamelCaseName}}s))
}
`
