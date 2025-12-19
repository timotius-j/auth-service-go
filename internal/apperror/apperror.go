package apperror

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ErrTypeInternalServer int = http.StatusInternalServerError
	ErrTypeBadRequest     int = http.StatusBadRequest
	ErrTypeNotFound       int = http.StatusNotFound
	ErrTypeUnauthorized   int = http.StatusUnauthorized
	ErrTypeTimeout        int = http.StatusRequestTimeout
	ErrTypeCancelled      int = http.StatusConflict
)

var (
	ErrDatabase = InternalServerError("Error Database")
	ErrInvalidRequest = BadRequestError("Invalid Request")
	ErrRecordNotFound = NotFoundError("Record Not Found")
	ErrUnauthorized   = UnauthorizedError("Unauthorized")
)

func TimeoutError(errMsg string) *gin.Error {
	err := errors.New(errMsg)

	return &gin.Error{
		Meta: map[string]int{
			"status_code": ErrTypeTimeout,
		},
		Err: err,
	}
}

func CanceledError(errMsg string) *gin.Error {
	err := errors.New(errMsg)

	return &gin.Error{
		Meta: map[string]int{
			"status_code": ErrTypeCancelled,
		},
		Err: err,
	}
}

func InternalServerError(errMsg string) *gin.Error {
	err := errors.New(errMsg)

	return &gin.Error{
		Meta: map[string]int{
			"status_code": ErrTypeInternalServer,
		},
		Err: err,
	}
}

func NotFoundError(errMsg string) *gin.Error {
	err := errors.New(errMsg)

	return &gin.Error{
		Meta: map[string]int{
			"status_code": ErrTypeNotFound,
		},
		Err: err,
	}
}

func BadRequestError(errMsg string) *gin.Error {
	err := errors.New(errMsg)

	return &gin.Error{
		Meta: map[string]int{
			"status_code": ErrTypeBadRequest,
		},
		Err: err,
	}
}

func UnauthorizedError(errMsg string) *gin.Error {
	err := errors.New(errMsg)

	return &gin.Error{
		Meta: map[string]int{
			"status_code": ErrTypeUnauthorized,
		},
		Err: err,
	}
}
