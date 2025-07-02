package handlers

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
	"users/domain/entities"
	"users/infrastructure/dependencies"
)

func TestGetMultiple(t *testing.T) {
	tests := []struct {
		name         string
		getByID      *GetByIDMock
		body         []byte
		expectedCode int
		expectedBody string
	}{
		{
			name:         "on OK execution",
			getByID:      NewGetByIDMock([]*entities.User{{ID: "1"}, {ID: "2"}}, nil),
			body:         []byte(`{"users":["1","2"]}`),
			expectedCode: http.StatusOK,
			expectedBody: "{\"data\":[{\"id\":\"1\",\"name\":\"\",\"birth\":\"\",\"email\":\"\",\"location\":null,\"created_at\":\"0001-01-01 00:00:00 +0000 UTC\",\"updated_at\":\"0001-01-01 00:00:00 +0000 UTC\",\"active\":false},{\"id\":\"2\",\"name\":\"\",\"birth\":\"\",\"email\":\"\",\"location\":null,\"created_at\":\"0001-01-01 00:00:00 +0000 UTC\",\"updated_at\":\"0001-01-01 00:00:00 +0000 UTC\",\"active\":false}]}",
		},
		{
			name:         "on invalid json",
			getByID:      NewGetByIDMock(nil, errors.New("an error occurred")),
			body:         []byte(`{"users":["1","2"]`),
			expectedCode: http.StatusBadRequest,
			expectedBody: "{\"errors\":\"unexpected EOF\"}",
		},
		{
			name:         "on repository error",
			getByID:      NewGetByIDMock(nil, errors.New("an error occurred")),
			body:         []byte(`{"users":["1","2"]}`),
			expectedCode: http.StatusInternalServerError,
			expectedBody: "{\"errors\":\"an error occurred\"}",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url := "/search"

			actions := dependencies.Actions{GetByID: test.getByID.execute}
			handler := New(&actions)

			request, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(test.body))
			response := httptest.NewRecorder()

			router := gin.New()
			router.POST(url, handler.GetMultiple)
			router.ServeHTTP(response, request)

			assertInt(t, response.Code, test.expectedCode)
			assertString(t, response.Body.String(), test.expectedBody)
		})
	}
}
