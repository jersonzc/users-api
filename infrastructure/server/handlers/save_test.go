package handlers

import (
	"bytes"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"users/domain/entities"
	"users/infrastructure/dependencies"
)

type SaveMock struct {
	execute func(context.Context, *entities.User) (*entities.User, error)
	answer  *entities.User
	err     error
}

func NewSaveMock(answer *entities.User, err error) *SaveMock {
	mock := &SaveMock{
		answer: answer,
		err:    err,
	}

	mock.execute = func(ctx context.Context, user *entities.User) (*entities.User, error) {
		if err != nil {
			return nil, err
		}
		return answer, nil
	}

	return mock
}

func TestSave(t *testing.T) {
	tests := []struct {
		name         string
		save         *SaveMock
		body         io.Reader
		expectedCode int
		expectedBody string
	}{
		{
			name:         "on OK execution",
			save:         NewSaveMock(&entities.User{ID: "2"}, nil),
			body:         bytes.NewReader([]byte(`{"name":"test"}`)),
			expectedCode: http.StatusCreated,
			expectedBody: "{\"data\":{\"id\":\"2\",\"name\":\"\",\"birth\":\"\",\"email\":\"\",\"location\":null,\"created_at\":\"0001-01-01 00:00:00 +0000 UTC\",\"updated_at\":\"0001-01-01 00:00:00 +0000 UTC\",\"active\":false}}",
		},
		{
			name:         "on nil payload",
			save:         NewSaveMock(&entities.User{ID: "2"}, nil),
			body:         nil,
			expectedCode: http.StatusBadRequest,
			expectedBody: "{\"errors\":\"cannot read nil body\"}",
		},
		{
			name:         "on missing name field",
			save:         NewSaveMock(&entities.User{ID: "2"}, nil),
			body:         bytes.NewReader([]byte(`{"email":"test@test.com"}`)),
			expectedCode: http.StatusBadRequest,
			expectedBody: "{\"errors\":\"Key: 'SaveUser.Name' Error:Field validation for 'Name' failed on the 'required' tag\"}",
		},
		{
			name:         "on invalid birth",
			save:         NewSaveMock(&entities.User{ID: "2"}, nil),
			body:         bytes.NewReader([]byte(`{"name":"test","birth":"23/09/92"}`)),
			expectedCode: http.StatusBadRequest,
			expectedBody: "{\"errors\":\"error while parsing 'birth' field from \\\"23/09/92\\\": parsing time \\\"23/09/92\\\" as \\\"02/01/2006\\\": cannot parse \\\"92\\\" as \\\"2006\\\"\"}",
		},
		{
			name:         "on save repository error",
			save:         NewSaveMock(nil, errors.New("an error occurred")),
			body:         bytes.NewReader([]byte(`{"name":"test"}`)),
			expectedCode: http.StatusInternalServerError,
			expectedBody: "{\"errors\":\"an error occurred\"}",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url := "/"

			actions := dependencies.Actions{Save: test.save.execute}
			handler := New(&actions)

			request, _ := http.NewRequest(http.MethodPost, url, test.body)
			response := httptest.NewRecorder()

			router := gin.New()
			router.POST(url, handler.Save)
			router.ServeHTTP(response, request)

			assertInt(t, response.Code, test.expectedCode)
			assertString(t, response.Body.String(), test.expectedBody)
		})
	}
}
