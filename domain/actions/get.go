package actions

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"users/domain"
	"users/domain/entities"
)

type Get struct {
	get    domain.Get
	tracer trace.Tracer
}

func NewGet(get domain.Get, tracer trace.Tracer) (*Get, error) {
	return &Get{get: get, tracer: tracer}, nil
}

func (action *Get) Execute(ctx context.Context) ([]*entities.User, error) {
	tracerCtx, span := action.tracer.Start(ctx, "Action-Get")
	defer span.End()

	return action.get(tracerCtx)
}
