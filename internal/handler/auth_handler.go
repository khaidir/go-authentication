package handler

import (
	"net/http"

	"auth-services/internal/dto"
	"auth-services/internal/usecase"
	"auth-services/pkg/logger"
	"auth-services/pkg/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Uc usecase.UserUsecase
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Response(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
		return
	}

	err := h.Uc.Register(req)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	res := dto.RegisterResponse{
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
	}
	utils.Response(c, http.StatusOK, "Successfully registered user", res)
}

func (h *UserHandler) Login(c *gin.Context) {

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseValidationError(c, err)
		return
	}
	if err := utils.ValidateStruct(&req); err != nil {
		utils.ResponseValidationFailed(c, utils.FormatValidationErrors(err))
		return
	}

	resp, err := h.Uc.Login(req)
	if err != nil {
		utils.ResponseError(c, http.StatusUnauthorized, "Invalid credentials", err.Error())
		return
	}
	utils.Response(c, http.StatusOK, "Logged", resp)
}

func (h *UserHandler) Profile(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		logger.LogWithTrace(c, "Parameter not consistent", nil, "warn")
		utils.ResponseParameterErrors(c, []utils.FieldError{
			{Field: "id", Message: "Id is required"},
		})
		return
	}

	// Validate UUID format if needed
	if !utils.IsValidUUID(id) {
		logger.LogWithTrace(c, "Invalid UUID format", nil, "warn")
		utils.ResponseParameterErrors(c, []utils.FieldError{
			{Field: "id", Message: "Invalid UUID format"},
		})
		return
	}

	res, err := h.Uc.GetProfile(id)
	if err != nil {
		logger.LogWithTrace(c, "Failed to get user profile", err, "warn")
		utils.ResponseError(c, http.StatusNotFound, "Not Found", err.Error())
		return
	}
	utils.Response(c, http.StatusOK, "OK", res)
}
