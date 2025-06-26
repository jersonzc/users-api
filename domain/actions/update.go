package actions

import (
	"context"
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

func NewUpdate(getByID domain.GetByID, update domain.Update, tracer trace.Tracer) (*Update, error) {
	return &Update{getByID: getByID, update: update, tracer: tracer}, nil
}

func (action *Update) Execute(ctx context.Context, user *entities.User) (*entities.User, error) {
	tracerCtx, span := action.tracer.Start(ctx, "Action-Update")
	defer span.End()

	result, err := action.getByID(tracerCtx, []string{user.ID})
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.AppUserNotFound
	}

	if err = action.update(tracerCtx, user); err != nil {
		return nil, err
	}

	return user, nil
}
