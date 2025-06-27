package actions

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"users/domain"
	"users/domain/entities"
)

type GetByID struct {
	getByID domain.GetByID
	tracer  trace.Tracer
}

func NewGetByID(getByID domain.GetByID) (*GetByID, error) {
	return &GetByID{
		getByID: getByID,
		tracer:  otel.Tracer("Action-GetByID")}, nil
}

func (action *GetByID) Execute(ctx context.Context, ids []string) ([]*entities.User, error) {
	tracerCtx, span := action.tracer.Start(ctx, "Action-GetByID-Execute")
	defer span.End()

	return action.getByID(tracerCtx, ids)
}
