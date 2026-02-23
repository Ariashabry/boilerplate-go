package product

import (
	"net/http"
	"strconv"

	"github.com/ariashabry/boilerplate-go/helpers/log"
	"github.com/ariashabry/boilerplate-go/internal/domains/product/model/dto"
	"github.com/ariashabry/boilerplate-go/internal/domains/product/service"
	"github.com/ariashabry/boilerplate-go/transport/http/response"
	"github.com/gin-gonic/gin"
)

// Ensure dto is used for swag type resolution
var _ dto.Product

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
		public.GET("/:id", c.GetProductById)
	}
}

// @Summary Get Product
// @Description Get All Data Product
// @Tags Product
// @Accept  json
// @Produce json
// @Success 200 {object} response.BaseResponse{data=[]dto.Product}
// @Failure 500 {object} response.BaseResponseError "Internal server error"
// @Router /product/ [get]
func (h *ProductHandler) GetProduct(ctx *gin.Context) {
	data, err := h.svc.GetList(ctx.Request.Context())
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

	resp := response.BaseResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    data,
		Success: true,
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary Get Product By Id
// @Description Get Product By Id
// @Tags Product
// @Accept  json
// @Produce json
// @Success 200 {object} response.BaseResponse{data=dto.Product}
// @Failure 500 {object} response.BaseResponseError "Internal server error"
// @Router /product/:id [get]
func (h *ProductHandler) GetProductById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.log.Error(err.Error())
		response.WithError(ctx, err)
		return
	}

	data, err := h.svc.GetProductById(ctx.Request.Context(), id)
	if err != nil {
		h.log.Error(err.Error())
		response.WithError(ctx, err)
		return
	}

	if data.Name == "" {
		resp := gin.H{
			"data":    nil,
			"error":   false,
			"message": "No products found",
		}
		ctx.JSON(http.StatusNotFound, resp)
		return
	}

	resp := response.BaseResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    data,
	}

	ctx.JSON(http.StatusOK, resp)
}
