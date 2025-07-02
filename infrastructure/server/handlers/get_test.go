package handlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
	"users/domain/entities"
	"users/infrastructure/dependencies"
)

type GetMock struct {
	execute func(context.Context) ([]*entities.User, error)
	answer  []*entities.User
	err     error
}

func NewGetMock(answer []*entities.User, err error) *GetMock {
	mock := &GetMock{
		answer: answer,
		err:    err,
	}

	mock.execute = func(ctx context.Context) ([]*entities.User, error) {
		if err != nil {
			return nil, err
		}
		return answer, nil
	}

	return mock
}

func TestGet(t *testing.T) {
	tests := []struct {
		name         string
		get          *GetMock
		expectedCode int
		expectedBody string
	}{
		{
			name:         "on OK execution",
			get:          NewGetMock([]*entities.User{{ID: "1"}, {ID: "2"}}, nil),
			expectedCode: http.StatusOK,
			expectedBody: "{\"data\":[{\"id\":\"1\",\"name\":\"\",\"birth\":\"\",\"email\":\"\",\"location\":null,\"created_at\":\"0001-01-01 00:00:00 +0000 UTC\",\"updated_at\":\"0001-01-01 00:00:00 +0000 UTC\",\"active\":false},{\"id\":\"2\",\"name\":\"\",\"birth\":\"\",\"email\":\"\",\"location\":null,\"created_at\":\"0001-01-01 00:00:00 +0000 UTC\",\"updated_at\":\"0001-01-01 00:00:00 +0000 UTC\",\"active\":false}]}",
		},
		{
			name:         "on repository error",
			get:          NewGetMock(nil, errors.New("an error occurred")),
			expectedCode: http.StatusInternalServerError,
			expectedBody: "{\"errors\":\"an error occurred\"}",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url := "/"

			actions := dependencies.Actions{Get: test.get.execute}
			handler := New(&actions)

			request, _ := http.NewRequest(http.MethodGet, url, nil)
			response := httptest.NewRecorder()

			router := gin.New()
			router.GET(url, handler.Get)
			router.ServeHTTP(response, request)

			assertInt(t, response.Code, test.expectedCode)
			assertString(t, response.Body.String(), test.expectedBody)
		})
	}
}
