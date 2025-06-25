package actions

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"users/domain"
	"users/domain/errors"
)

type Remove struct {
	getByID domain.GetByID
	remove  domain.Remove
	tracer  trace.Tracer
}

func NewRemove(getByID domain.GetByID, remove domain.Remove, tracer trace.Tracer) (*Remove, error) {
	return &Remove{getByID: getByID, remove: remove, tracer: tracer}, nil
}

func (action *Remove) Execute(ctx context.Context, id int) error {
	tracerCtx, span := action.tracer.Start(ctx, "Action-Remove")
	defer span.End()

	result, err := action.getByID(tracerCtx, []int{id})
	if err != nil {
		return err
	}

	if len(result) == 0 {
		return errors.AppUserNotFound
	}

	return action.remove(tracerCtx, id)
}
