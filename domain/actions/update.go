package actions

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"users/domain"
	"users/domain/entities"
	"users/domain/errors"
)

type Update struct {
	getByID domain.GetByID
	update  domain.Update
	tracer  trace.Tracer
}

func NewUpdate(getByID domain.GetByID, update domain.Update) (*Update, error) {
	return &Update{
		getByID: getByID,
		update:  update,
		tracer:  otel.Tracer("Action-Update")}, nil
}

func (action *Update) Execute(ctx context.Context, id string, fields map[string]interface{}) (*entities.User, error) {
	tracerCtx, span := action.tracer.Start(ctx, "Action-Update-Execute")
	defer span.End()

	result, err := action.getByID(tracerCtx, []string{id})
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.AppUserNotFound
	}

	return action.update(tracerCtx, id, fields)
}
