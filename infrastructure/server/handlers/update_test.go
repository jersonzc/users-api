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
	"users/infrastructure/dependencies"
)

type UpdateMock struct {
	execute func(context.Context, string, map[string]interface{}) error
	err     error
}

func NewUpdateMock(err error) *UpdateMock {
	mock := &UpdateMock{
		err: err,
	}

	mock.execute = func(ctx context.Context, id string, fields map[string]interface{}) error {
		return err
	}

	return mock
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name         string
		update       *UpdateMock
		body         io.Reader
		expectedCode int
		expectedBody string
	}{
		{
			name:         "on OK execution",
			update:       NewUpdateMock(nil),
			body:         bytes.NewReader([]byte(`{"name":"test"}`)),
			expectedCode: http.StatusNoContent,
			expectedBody: "",
		},
		{
			name:         "on invalid type",
			update:       NewUpdateMock(nil),
			body:         bytes.NewReader([]byte(`{"name":1}`)),
			expectedCode: http.StatusBadRequest,
			expectedBody: "{\"errors\":\"json: cannot unmarshal number into Go struct field UpdateUser.name of type string\"}",
		},
		{
			name:         "on nil payload",
			update:       NewUpdateMock(nil),
			body:         nil,
			expectedCode: http.StatusBadRequest,
			expectedBody: "{\"errors\":\"cannot read nil body\"}",
		},
		{
			name:         "on invalid birth",
			update:       NewUpdateMock(nil),
			body:         bytes.NewReader([]byte(`{"name":"test","birth":"23/09/92"}`)),
			expectedCode: http.StatusBadRequest,
			expectedBody: "{\"errors\":\"error while parsing 'birth' field from \\\"23/09/92\\\": parsing time \\\"23/09/92\\\" as \\\"02/01/2006\\\": cannot parse \\\"92\\\" as \\\"2006\\\"\"}",
		},
		{
			name:         "on save repository error",
			update:       NewUpdateMock(errors.New("an error occurred")),
			body:         bytes.NewReader([]byte(`{"name":"test"}`)),
			expectedCode: http.StatusInternalServerError,
			expectedBody: "{\"errors\":\"an error occurred\"}",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url := "/1"

			actions := dependencies.Actions{Update: test.update.execute}
			handler := New(&actions)

			request, _ := http.NewRequest(http.MethodPut, url, test.body)
			response := httptest.NewRecorder()

			router := gin.New()
			router.PUT(url, handler.Update)
			router.ServeHTTP(response, request)

			assertInt(t, response.Code, test.expectedCode)
			assertString(t, response.Body.String(), test.expectedBody)
		})
	}
}
