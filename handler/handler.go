package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/leifarriens/go-microservice/service"
)

type Handler struct {
	ProductService service.ProductService
}

type HandlerConfig struct {
	E              *echo.Echo
	ProductService service.ProductService
}

func NewHandler(c *HandlerConfig) *Handler {
	h := &Handler{
		ProductService: c.ProductService,
	}

	c.E.POST("/products", h.CreateProduct)
	c.E.GET("/products", h.GetAllProducts)
	c.E.GET("/products/:id", h.GetProductByID)

	return h
}
