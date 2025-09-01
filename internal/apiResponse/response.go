package apiresponse

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse represents the base structure for all API responses
type APIResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
	Status  bool   `json:"status" example:"true"`
}

// ErrorResponse represents error responses
type ErrorResponse struct {
	APIResponse
	Error string `json:"error,omitempty" example:"validation_error"`
}

// SuccessResponse represents successful responses with data
type SuccessResponse struct {
	Data interface{} `json:"data"`
	APIResponse
}

// PaginatedResponse for paginated data
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
	APIResponse
}

type Pagination struct {
	Page       int `json:"page" example:"1"`
	Limit      int `json:"limit" example:"10"`
	Total      int `json:"total" example:"100"`
	TotalPages int `json:"totalPages" example:"10"`
}

// Constructor functions
func NewErrorResponse(msg string, errorCode ...string) ErrorResponse {
	response := ErrorResponse{
		APIResponse: APIResponse{
			Message: msg,
			Status:  false,
		},
	}
	if len(errorCode) > 0 {
		response.Error = errorCode[0]
	}
	return response
}

func NewSuccessResponse(data interface{}, msg string) SuccessResponse {
	return SuccessResponse{
		Data: data,
		APIResponse: APIResponse{
			Message: msg,
			Status:  true,
		},
	}
}

// Helper methods for common responses
func SendError(c *gin.Context, statusCode int, message string, errorCode ...string) {
	c.JSON(statusCode, NewErrorResponse(message, errorCode...))
}

func SendSuccess(c *gin.Context, statusCode int, data interface{}, message string) {
	c.JSON(statusCode, NewSuccessResponse(data, message))
}

func SendValidationError(c *gin.Context, err error) {
	SendError(c, http.StatusBadRequest, err.Error(), "validation_error")
}

func SendInternalError(c *gin.Context, message string) {
	SendError(c, http.StatusInternalServerError, message, "internal_error")
}

func SendUnauthorized(c *gin.Context, message string) {
	SendError(c, http.StatusUnauthorized, message, "unauthorized")
}

func SendConflict(c *gin.Context, message string) {
	SendError(c, http.StatusConflict, message, "conflict")
}
