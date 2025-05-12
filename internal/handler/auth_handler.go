package handler

import (
	"net/http"
	"strconv"

	"auth-services/internal/dto"
	"auth-services/internal/usecase"
	"auth-services/pkg/logger"
	"auth-services/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	err := h.Uc.Register(req)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	utils.Response(c, http.StatusOK, "Successfully registered user", nil)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.Uc.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	logged := gin.H{
		"email": req.Email,
		"token": token,
	}
	utils.Response(c, http.StatusOK, "Logged", logged)
	// c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) Profile(c *gin.Context) {
	traceID := c.GetString(logger.TraceIDKey)
	id, _ := strconv.Atoi(c.Param("id"))
	res, err := h.Uc.GetProfile(uint(id))
	if err != nil {
		logger.Log.Error("failed to get profile", zap.String("trace_id", traceID), zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	utils.Response(c, http.StatusOK, "OK", res)
}
