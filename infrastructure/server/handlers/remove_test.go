package handlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
	"users/infrastructure/dependencies"
)

type RemoveMock struct {
	execute func(context.Context, string) error
	err     error
}

func NewRemoveMock(err error) *RemoveMock {
	mock := &RemoveMock{
		err: err,
	}

	mock.execute = func(ctx context.Context, id string) error {
		return err
	}

	return mock
}

func TestRemove(t *testing.T) {
	tests := []struct {
		name         string
		remove       *RemoveMock
		expectedCode int
		expectedBody string
	}{
		{
			name:         "on OK execution",
			remove:       NewRemoveMock(nil),
			expectedCode: http.StatusNoContent,
			expectedBody: "",
		},
		{
			name:         "on repository error",
			remove:       NewRemoveMock(errors.New("an error occurred")),
			expectedCode: http.StatusInternalServerError,
			expectedBody: "{\"errors\":\"an error occurred\"}",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url := "/1"

			actions := dependencies.Actions{Remove: test.remove.execute}
			handler := New(&actions)

			request, _ := http.NewRequest(http.MethodDelete, url, nil)
			response := httptest.NewRecorder()

			router := gin.New()
			router.DELETE(url, handler.Remove)
			router.ServeHTTP(response, request)

			assertInt(t, response.Code, test.expectedCode)
			assertString(t, response.Body.String(), test.expectedBody)
		})
	}
}
