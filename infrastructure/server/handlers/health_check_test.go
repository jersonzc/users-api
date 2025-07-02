package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	url := "/health"

	request, _ := http.NewRequest("GET", url, nil)
	response := httptest.NewRecorder()

	router := gin.New()
	router.GET(url, HealthCheck)
	router.ServeHTTP(response, request)

	got := response.Code
	want := http.StatusOK

	if got != want {
		t.Errorf("got '%d' want '%d'", got, want)
	}
}
