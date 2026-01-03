package handlers

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func RespondOK(c *gin.Context, status int, payload interface{}) {
	c.JSON(status, payload)
}

func RespondError(
	c *gin.Context,
	status int,
	code string,
	message string,
	details interface{},
) {
	c.JSON(status, ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	})
}
