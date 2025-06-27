package actions

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"users/domain"
	"users/domain/entities"
)

type GetByID struct {
	getByID domain.GetByID
	tracer  trace.Tracer
}

func NewGetByID(getByID domain.GetByID, tracer trace.Tracer) (*GetByID, error) {
	return &GetByID{getByID: getByID, tracer: tracer}, nil
}

func (action *GetByID) Execute(ctx context.Context, ids []string) ([]*entities.User, error) {
	tracerCtx, span := action.tracer.Start(ctx, "Action-GetByID")
	defer span.End()

	return action.getByID(tracerCtx, ids)
}
