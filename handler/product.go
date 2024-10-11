package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/leifarriens/go-microservice/model"
	"github.com/leifarriens/go-microservice/service"
	_ "gorm.io/gorm"
)

// CreateProduct godoc
//
//	@Summary		Create product
//	@Description	Create product
//	@Tags			product
//	@Accept			json
//	@Produce		json
//
//	@Param			product	body		model.ProductDto	true	"The input product struct"
//	@Success		200		{object}	model.Product
//
//	@failure		400		{string}	string	"error"
//	@failure		404		{string}	string	"error"
//	@failure		500		{string}	string	"error"
//
//	@Router			/products [post]
func (h *Handler) CreateProduct(c echo.Context) error {
	var p model.ProductDto

	if err := c.Bind(&p); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(p); err != nil {
		return err
	}

	product, err := h.ProductService.Add(c.Request().Context(), &model.ProductDto{
		Name:      p.Name,
		Price:     p.Price,
		Available: p.Available,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if product == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	return c.JSON(http.StatusOK, product)
}

type GetAllProductsParams struct {
	model.PageableRequest
}

type GetAllProductsResponse struct {
	model.PageableResponse
	Items []*model.Product `json:"items"`
}

// GetAllProducts godoc
//
//	@Summary		Get all products
//	@Description	Get all products
//	@Tags			product
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int	false	"limit"
//	@Param			offset	query		int	false	"offset"
//	@Success		200		{object}	GetAllProductsResponse
//
//	@failure		400		{string}	string	"error"
//	@failure		404		{string}	string	"error"
//	@failure		500		{string}	string	"error"
//
//	@Router			/products [get]
func (h *Handler) GetAllProducts(c echo.Context) error {
	var p GetAllProductsParams

	if err := c.Bind(&p); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if p.Limit == 0 {
		p.Limit = 20
	}

	if p.Offset == 0 {
		p.Offset = 0
	}

	if err := c.Validate(p); err != nil {
		return err
	}

	products, err := h.ProductService.Get(c.Request().Context(), p.Limit, p.Offset)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &GetAllProductsResponse{
		PageableResponse: model.PageableResponse{
			Limit:  p.Limit,
			Offset: p.Offset,
		},
		Items: products,
	})
}

type GetByIDParams struct {
	ID string `param:"id" validate:"required"`
}

// GetProductByID godoc
//
//	@Summary		Get products by id
//	@Description	Get products by id
//	@Tags			product
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Product ID"
//	@Success		200	{object}	model.Product
//
//	@failure		400	{string}	string	"error"
//	@failure		404	{string}	string	"error"
//	@failure		500	{string}	string	"error"
//
//	@Router			/products/{id} [get]
func (h *Handler) GetProductByID(c echo.Context) error {
	var p GetByIDParams

	if err := c.Bind(&p); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if err := c.Validate(p); err != nil {
		return err
	}

	product, err := h.ProductService.GetByID(c.Request().Context(), p.ID)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Product not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, product)
}
