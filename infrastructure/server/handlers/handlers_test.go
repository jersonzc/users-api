package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestMapToString(t *testing.T) {
	tests := []struct {
		name  string
		myMap map[string][]string
		want  string
	}{
		{
			name:  "on data",
			myMap: http.Header{"x-app": {"test"}},
			want:  "\"x-app\":\"test\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, mapToString(tt.myMap), "mapToString(%v)", tt.myMap)
		})
	}
}

func assertInt(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("got '%d', want '%d'", got, want)
	}
}

func assertString(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}
