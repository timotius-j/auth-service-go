package util

import (
	"github.com/TimX-21/auth-service-go/internal/dto"
	"github.com/gin-gonic/gin"
)

func HandleResponse(data any, code int, c *gin.Context) {
	res := dto.Response[any]{}
	res.Success = true
	res.Data = data

	c.JSON(code, res)
}

func HandlePagedResponse(data any, paging *dto.PageMetadata, code int, c *gin.Context) {
	res := dto.Response[any]{}
	res.Success = true
	res.Data = data
	res.Paging = paging

	c.JSON(code, res)
}

func HandleAbortOAuth(c *gin.Context, err *gin.Error) {
	metaMap, ok := err.Meta.(map[string]int)
	if !ok {
		return
	}

	statusCode := metaMap["status_code"]

	e := &dto.Error{
		Message: err.Error(),
	}

	c.AbortWithStatusJSON(statusCode, dto.Response[any]{
		Success: false,
		Error:   e,
	})
}
