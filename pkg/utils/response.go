package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type FieldError struct {
	Field   string
	Message string
}

type ValidationErrorResponse struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Details map[string]string `json:"details"`
}

func Response(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, gin.H{
		"status":  "success",
		"message": message,
		"data":    data,
	})
}

func ResponseError(c *gin.Context, code int, message string, details interface{}) {
	c.JSON(code, gin.H{
		"status":  "error",
		"message": message,
		"details": details,
	})
}

func ResponseValidationFailed(c *gin.Context, errors map[string]string) {
	ResponseError(c, http.StatusBadRequest, "Validation failed", errors)
}

func ResponseValidationError(c *gin.Context, err error) {
	ResponseError(c, http.StatusBadRequest, "Invalid request payload", err.Error())
}

func ResponseParameterErrors(c *gin.Context, fieldErrors []FieldError) {
	details := make(map[string]string)
	for _, fe := range fieldErrors {
		details[fe.Field] = fe.Message
	}

	c.JSON(http.StatusBadRequest, ValidationErrorResponse{
		Status:  "error",
		Message: "Parameter failed",
		Details: details,
	})
}
