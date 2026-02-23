package response

import (
	"github.com/ariashabry/boilerplate-go/shared/failure"
	"github.com/gin-gonic/gin"
)

type Base struct {
	Data     *interface{} `json:"data,omitempty"`
	MetaData *interface{} `json:"metadata,omitempty"`
	Error    *string      `json:"error,omitempty"`
	Message  *string      `json:"message,omitempty"`
	Success  *bool        `json:"success,omitempty"`
}

type BaseResponse struct {
	Code    int         `json:"code" example:"200"`
	Message string      `json:"message" example:"Success"`
	Data    interface{} `json:"data"`
	Success bool        `json:"success,omitempty" example:"true"`
}

type BaseResponseError struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"Internal Server Error"`
	Success bool   `json:"success" example:"false"`
}

func WithError(ctx *gin.Context, err error) {
	code := failure.GetCode(err)
	errMsg := err.Error()
	ctx.JSON(code, Base{Error: &errMsg})
}

func WithSuccess(ctx *gin.Context, code int, data interface{}) {
	ctx.JSON(code, BaseResponse{
		Code:    code,
		Message: "Success",
		Data:    data,
	})
}
