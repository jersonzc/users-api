package dependencies

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"users/domain/actions"
	"users/domain/entities"
	"users/infrastructure/postgres"
)

type Actions struct {
	Get     func(context.Context) ([]*entities.User, error)
	GetByID func(context.Context, []string) ([]*entities.User, error)
	Save    func(context.Context, *entities.User) (*entities.User, error)
	Update  func(context.Context, *entities.User) (*entities.User, error)
	Remove  func(context.Context, string) error
}

func NewActions(postgresClient *postgres.Client, tracer trace.Tracer) (*Actions, error) {
	postgresRepo, err := postgres.NewRepository(postgresClient.Modify, postgresClient.Retrieve, tracer)
	if err != nil {
		return nil, err
	}

	get, err := actions.NewGet(postgresRepo.Get, tracer)
	if err != nil {
		return nil, err
	}

	getByID, err := actions.NewGetByID(postgresRepo.GetByID, tracer)
	if err != nil {
		return nil, err
	}

	save, err := actions.NewSave(postgresRepo.GetByID, postgresRepo.Save, tracer)
	if err != nil {
		return nil, err
	}

	update, err := actions.NewUpdate(postgresRepo.GetByID, postgresRepo.Update, tracer)
	if err != nil {
		return nil, err
	}

	remove, err := actions.NewRemove(postgresRepo.GetByID, postgresRepo.Remove, tracer)
	if err != nil {
		return nil, err
	}

	return &Actions{
		Get:     get.Execute,
		GetByID: getByID.Execute,
		Save:    save.Execute,
		Update:  update.Execute,
		Remove:  remove.Execute,
	}, nil
}
