package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shamank/warehouse-service/internal/domain/schemas"
	"log/slog"
	"net/http"
)

//go:generate mockery --name=Service
type Service interface {
	GetRemainingProducts(warehouseUUID string) ([]schemas.Product, error)
	ReserveProducts(productsToReserve []string) error
	ReleaseProducts(productsToRelease []string) error
}

type Handler struct {
	service Service
	logger  *slog.Logger
}

func NewHandler(service Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) InitAPIRoutes() *gin.Engine {
	r := gin.Default()

	r.Use(CORS)

	api := r.Group("/api")

	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
		api.GET("/getRemainingProducts", h.getRemainingProducts)
		api.POST("/reserveProducts", h.reserveProducts)
		api.POST("/releaseProducts", h.releaseProducts)
	}

	return r

}

func (h *Handler) getRemainingProducts(c *gin.Context) {
	params := c.Request.URL.Query()

	warehouseUUID := params.Get("warehouse_uuid")

	if warehouseUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "warehouse_uuid is required",
		})

		return
	}

	result, err := h.service.GetRemainingProducts(warehouseUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unkown error",
		})

		return
	}

	c.JSON(http.StatusOK, result)

}

type reserveProductRequest []string

func (h *Handler) reserveProducts(c *gin.Context) {
	var productsToReserve reserveProductRequest

	err := c.BindJSON(&productsToReserve)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return

	}

	err = h.service.ReserveProducts(productsToReserve)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})

}

type releaseProductRequest []string

func (h *Handler) releaseProducts(c *gin.Context) {
	var productsToRelease releaseProductRequest

	err := c.BindJSON(&productsToRelease)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = h.service.ReleaseProducts(productsToRelease)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})

}
