package handler

import (
	"net/http"

	"github.com/TimX-21/auth-service-go/internal/apperror"
	"github.com/TimX-21/auth-service-go/internal/auth/dto"
	"github.com/TimX-21/auth-service-go/internal/auth/model"
	"github.com/TimX-21/auth-service-go/internal/auth/service"
	"github.com/TimX-21/auth-service-go/internal/util"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service service.AuthServiceItf
}

func NewAuthHandler(s service.AuthServiceItf) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) GetUserDataHandler(c *gin.Context) {
	
	request := dto.GetUserDataRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(apperror.ErrInvalidRequest)
		return
	}
	if err := util.ValidateStruct(request); err != nil {
		c.Error(apperror.ErrInvalidRequest)
		return
	}
	
	user := model.User{Email: request.Email}
	userData, err := h.service.GetUserDataService(c, user)
	if err != nil {
		c.Error(err)
		return
	}
	
	response := dto.GetUserDataResponse{
		ID:        userData.ID,
		Email:     userData.Email,
		IsActive:  userData.IsActive,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
	}

	util.HandleResponse(response, http.StatusOK, c)
}
