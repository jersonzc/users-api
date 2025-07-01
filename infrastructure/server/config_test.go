package server

import (
	"errors"
	"testing"
	"time"
	errorspkg "users/domain/errors"
)

func TestNewConfig(t *testing.T) {
	readTimeout := time.Duration(1)
	writeTimeout := time.Duration(1)

	t.Run("on port number out of range", func(t *testing.T) {
		port := 65536
		prefix := "/api"

		_, err := NewConfig(port, prefix, readTimeout, writeTimeout)

		assertError(t, err, errorspkg.ServerInvalidPort)
	})

	t.Run("on missing prefix", func(t *testing.T) {
		port := 3001
		prefix := ""

		_, err := NewConfig(port, prefix, readTimeout, writeTimeout)

		assertError(t, err, errorspkg.ServerMissingPrefix)
	})

	t.Run("on OK", func(t *testing.T) {
		port := 3001
		prefix := "/api"

		_, err := NewConfig(port, prefix, readTimeout, writeTimeout)

		assertNoError(t, err)
	})
}

func assertError(t testing.TB, got, want error) {
	t.Helper()

	if !errors.Is(got, want) {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func assertNoError(t testing.TB, got error) {
	t.Helper()

	if got != nil {
		t.Fatal("got an error but didn't want one")
	}
}
