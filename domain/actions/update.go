package actions

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"users/domain"
	"users/domain/errors"
)

type Update struct {
	getByID domain.GetByID
	update  domain.Update
	tracer  trace.Tracer
}

func NewUpdate(getByID domain.GetByID, update domain.Update, tracer trace.Tracer) (*Update, error) {
	return &Update{getByID: getByID, update: update, tracer: tracer}, nil
}

func (action *Update) Execute(ctx context.Context, id string, fields *map[string]interface{}) error {
	tracerCtx, span := action.tracer.Start(ctx, "Action-Update")
	defer span.End()

	result, err := action.getByID(tracerCtx, []string{id})
	if err != nil {
		return err
	}

	if len(result) == 0 {
		return errors.AppUserNotFound
	}

	if err = action.update(tracerCtx, id, fields); err != nil {
		return err
	}

	return nil
}
