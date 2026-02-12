package product

import (
	"net/http"

	"github.com/ariashabry/boilerplate-go/helpers/log"
	"github.com/ariashabry/boilerplate-go/internal/domains/product/service"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	svc service.ProductService
	log *log.AppLog
}

func ProvideProductHandler(svc service.ProductService, log *log.AppLog) ProductHandler {
	return ProductHandler{svc: svc, log: log}
}

func (c *ProductHandler) Router(group *gin.RouterGroup) {
	public := group.Group("/product")
	{
		public.GET("/", c.GetProduct)
	}
}

func (h *ProductHandler) GetProduct(ctx *gin.Context) {
	data, err := h.svc.GetProduct(ctx.Request.Context())
	if err != nil {
		h.log.Error(err.Error())
		resp := gin.H{
			"data":    nil,
			"error":   true,
			"message": err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	if len(data) == 0 {
		resp := gin.H{
			"data":    nil,
			"error":   false,
			"message": "No products found",
		}
		ctx.JSON(http.StatusNotFound, resp)
		return
	}

	resp := gin.H{
		"data":    data,
		"error":   false,
		"message": "Success",
	}

	ctx.JSON(http.StatusOK, resp)
}
