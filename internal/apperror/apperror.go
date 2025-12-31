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
	ErrTypeConflict       int = http.StatusConflict
)

var (
	ErrDatabase       = InternalServerError("Error Database")
	ErrInvalidRequest = BadRequestError("Invalid Request")
	ErrRecordNotFound = NotFoundError("Record Not Found")
	ErrUnauthorized   = UnauthorizedError("Unauthorized")
	ErrUserNotFound   = NotFoundError("User Not Found")
	ErrInternalServer = InternalServerError("Internal Server Error")
	ErrTimeout        = TimeoutError("Request Timeout")
	ErrCanceled       = CanceledError("Request Canceled")

	ErrTooManyRequests   = BadRequestError("Too Many Requests")
	ErrInvalidOTP        = BadRequestError("Invalid OTP")
	ErrInvalidResetToken = BadRequestError("Invalid Reset Token")
	ErrPasswordNotMatch  = BadRequestError("Password doesn't match")

	ErrAlreadyVerified        = ConflictError("account already verified")
	ErrInvalidOrdExpiredToken = BadRequestError("invalid or expired verification token")
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

func ConflictError(errMsg string) *gin.Error {
	err := errors.New(errMsg)

	return &gin.Error{
		Meta: map[string]int{
			"status_code": ErrTypeConflict,
		},
		Err: err,
	}
}
