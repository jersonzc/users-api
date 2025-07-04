package actions

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"users/domain"
	"users/domain/entities"
	"users/domain/errors"
)

type Save struct {
	getByID domain.GetByID
	save    domain.Save
	tracer  trace.Tracer
}

func NewSave(getByID domain.GetByID, save domain.Save) (*Save, error) {
	return &Save{
		getByID: getByID,
		save:    save,
		tracer:  otel.Tracer("Action-Save")}, nil
}

func (action *Save) Execute(ctx context.Context, user *entities.User) (*entities.User, error) {
	tracerCtx, span := action.tracer.Start(ctx, "Action-Save-Execute")
	defer span.End()

	result, err := action.getByID(tracerCtx, []string{user.ID})
	if err != nil {
		return nil, err
	}

	if len(result) == 1 {
		return nil, errors.AppUserExists
	}

	return action.save(tracerCtx, user)
}
