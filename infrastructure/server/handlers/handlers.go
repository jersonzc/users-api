package handlers

import (
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"strings"
	"users/infrastructure/dependencies"
)

const xAppID = "X-Application-ID"

type Handlers struct {
	actions *dependencies.Actions
	tracer  trace.Tracer
}

func New(actions *dependencies.Actions) *Handlers {
	return &Handlers{
		actions: actions,
		tracer:  otel.Tracer("Handler"),
	}
}

func mapToString(myMap map[string][]string) string {
	var sb strings.Builder

	for key, values := range myMap {
		for _, v := range values {
			sb.WriteString(fmt.Sprintf("%q:%q", key, v))
			sb.WriteString(" ")
		}
	}

	result := sb.String()
	if len(result) > 0 {
		result = result[:len(result)-1]
	}

	return result
}
