package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	contentTypeHeader = "Content-Type"
	applicationJSON   = "application/json"
)

type TestRequest struct {
	Name string `json:"name" binding:"required"`
	Age  int    `json:"age" binding:"gte=18"`
}

func setupRouter(handler gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/test", handler)
	return r
}

func TestWithBody(t *testing.T) {
	// 1. Success Case
	t.Run("Valid Body", func(t *testing.T) {
		called := false
		mockHandler := func(c *gin.Context, req TestRequest) {
			called = true
			assert.Equal(t, "John", req.Name)
			assert.Equal(t, 25, req.Age)
			c.Status(http.StatusOK)
		}

		r := setupRouter(WithBody(mockHandler))
		body, _ := json.Marshal(TestRequest{Name: "John", Age: 25})
		req, _ := http.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(body))
		req.Header.Set(contentTypeHeader, applicationJSON)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.True(t, called)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	// 2. Validation Error
	t.Run("Invalid Body - Missing Required Field", func(t *testing.T) {
		called := false
		mockHandler := func(c *gin.Context, req TestRequest) {
			called = true
		}

		r := setupRouter(WithBody(mockHandler))
		// Missing Name
		body, _ := json.Marshal(map[string]interface{}{"age": 25})
		req, _ := http.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(body))
		req.Header.Set(contentTypeHeader, applicationJSON)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.False(t, called)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// 3. Validation Error - Constraint Failed
	t.Run("Invalid Body - Constraint Failed", func(t *testing.T) {
		called := false
		mockHandler := func(c *gin.Context, req TestRequest) {
			called = true
		}

		r := setupRouter(WithBody(mockHandler))
		// Age < 18
		body, _ := json.Marshal(TestRequest{Name: "John", Age: 10})
		req, _ := http.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(body))
		req.Header.Set(contentTypeHeader, applicationJSON)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.False(t, called)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// 4. Malformed JSON
	t.Run("Malformed JSON", func(t *testing.T) {
		called := false
		mockHandler := func(c *gin.Context, req TestRequest) {
			called = true
		}

		r := setupRouter(WithBody(mockHandler))
		req, _ := http.NewRequest(http.MethodPost, "/test", bytes.NewBufferString("{invalid-json"))
		req.Header.Set(contentTypeHeader, applicationJSON)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.False(t, called)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
