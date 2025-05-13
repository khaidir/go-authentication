package handler

import (
	"net/http"
	"strconv"

	"auth-services/internal/dto"
	"auth-services/internal/usecase"
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
	err := h.Uc.Register(req)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	res := dto.RegisterResponse{
		Name:  req.Name,
		Email: req.Email,
	}
	utils.Response(c, http.StatusOK, "Successfully registered user", res)
}

func (h *UserHandler) Login(c *gin.Context) {

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogError(c, http.StatusBadRequest, "Internal Server Error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.Uc.Login(req)
	if err != nil {
		utils.LogError(c, http.StatusInternalServerError, "Internal Server Error", err)
		utils.Response(c, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}
	logged := gin.H{
		"email": req.Email,
		"token": token,
	}
	utils.Response(c, http.StatusOK, "Logged", logged)
}

func (h *UserHandler) Profile(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	res, err := h.Uc.GetProfile(uint(id))
	if err != nil {
		utils.LogError(c, 500, "failed to get profile", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	utils.Response(c, http.StatusOK, "OK", res)
}
