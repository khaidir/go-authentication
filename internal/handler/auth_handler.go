package handler

import (
	"net/http"
	"strconv"

	"auth-services/internal/dto"
	"auth-services/internal/usecase"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Uc usecase.UserUsecase
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.Uc.Register(req)
	c.JSON(http.StatusOK, gin.H{"success": err == nil})
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
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) Profile(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	res, err := h.Uc.GetProfile(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
