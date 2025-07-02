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

type GetByIDMock struct {
	execute func(context.Context, []string) ([]*entities.User, error)
	answer  []*entities.User
	err     error
}

func NewGetByIDMock(answer []*entities.User, err error) *GetByIDMock {
	mock := &GetByIDMock{
		answer: answer,
		err:    err,
	}

	mock.execute = func(ctx context.Context, ids []string) ([]*entities.User, error) {
		if err != nil {
			return nil, err
		}
		return answer, nil
	}

	return mock
}

func TestGetSingle(t *testing.T) {
	tests := []struct {
		name         string
		getByID      *GetByIDMock
		expectedCode int
		expectedBody string
	}{
		{
			name:         "on OK execution",
			getByID:      NewGetByIDMock([]*entities.User{{ID: "1"}}, nil),
			expectedCode: http.StatusOK,
			expectedBody: "{\"data\":[{\"id\":\"1\",\"name\":\"\",\"birth\":\"\",\"email\":\"\",\"location\":null,\"created_at\":\"0001-01-01 00:00:00 +0000 UTC\",\"updated_at\":\"0001-01-01 00:00:00 +0000 UTC\",\"active\":false}]}",
		},
		{
			name:         "on repository error",
			getByID:      NewGetByIDMock(nil, errors.New("an error occurred")),
			expectedCode: http.StatusInternalServerError,
			expectedBody: "{\"errors\":\"an error occurred\"}",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url := "/search/1"

			actions := dependencies.Actions{GetByID: test.getByID.execute}
			handler := New(&actions)

			request, _ := http.NewRequest(http.MethodGet, url, nil)
			response := httptest.NewRecorder()

			router := gin.New()
			router.GET(url, handler.GetSingle)
			router.ServeHTTP(response, request)

			assertInt(t, response.Code, test.expectedCode)
			assertString(t, response.Body.String(), test.expectedBody)
		})
	}
}
