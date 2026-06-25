package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// envelope wraps all successful responses as {"data": ...}.
type envelope[T any] struct {
	Data T `json:"data"`
}

// errDetail is the inner object of all error responses.
type errDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// errBody is the shape of all error responses: {"error": {"code": "...", "message": "..."}}.
type errBody struct {
	Error errDetail `json:"error"`
}

// JSON writes a 200 response with the data wrapped in {"data": ...}.
func JSON[T any](c *gin.Context, data T) {
	c.JSON(http.StatusOK, envelope[T]{Data: data})
}

// JSONStatus writes any status code with the data wrapped in {"data": ...}.
func JSONStatus[T any](c *gin.Context, status int, data T) {
	c.JSON(status, envelope[T]{Data: data})
}

// JSONError writes an error response as {"error": {"code": "...", "message": "..."}}.
func JSONError(c *gin.Context, status int, code, message string) {
	c.JSON(status, errBody{Error: errDetail{Code: code, Message: message}})
}
