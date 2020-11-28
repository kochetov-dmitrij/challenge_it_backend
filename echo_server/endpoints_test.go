package main

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test HealthCheck endpoint
func TestHealthCheck(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expectedBody := "{\"data\":\"Server is up and running\"}\n"

	// Assertions
	if assert.NoError(t, HealthCheck(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expectedBody, rec.Body.String())
	}
}
