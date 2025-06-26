package postgres

import (
	"context"
	"github.com/doug-martin/goqu/v9"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"users/domain/entities"
	"users/infrastructure/mappers"
)

type Repository struct {
	execModify   func(context.Context, string) error
	execRetrieve func(context.Context, string) ([]map[string]string, error)
	tracer       trace.Tracer
}

func NewRepository(
	execModify func(context.Context, string) error,
	execRetrieve func(context.Context, string) ([]map[string]string, error),
	tracer trace.Tracer,
) (*Repository, error) {
	return &Repository{execModify: execModify, execRetrieve: execRetrieve, tracer: tracer}, nil
}

func (repo *Repository) Get(ctx context.Context) ([]*entities.User, error) {
	tracerCtx, span := repo.tracer.Start(ctx, "PostgresRepository-Get")
	defer span.End()

	query, _, _ := goqu.From("users").ToSQL()

	rows, err := repo.execRetrieve(tracerCtx, query)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(attribute.Int("repo.postgres.rows.count", len(rows)))

	return mappers.ToUserList(rows), nil
}

func (repo *Repository) GetByID(ctx context.Context, ids []string) ([]*entities.User, error) {
	tracerCtx, span := repo.tracer.Start(ctx, "PostgresRepository-GetByID")
	defer span.End()

	query, _, _ := goqu.From("users").Where(goqu.C("id").In(ids)).ToSQL()

	rows, err := repo.execRetrieve(tracerCtx, query)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(attribute.Int("repo.postgres.rows.count", len(rows)))

	return mappers.ToUserList(rows), nil
}

func (repo *Repository) Save(ctx context.Context, user *entities.User) error {
	tracerCtx, span := repo.tracer.Start(ctx, "PostgresRepository-Save")
	defer span.End()

	query, _, _ := goqu.Insert("users").
		Cols("id", "name", "birth", "active", "location").
		Vals(
			goqu.Vals{user.ID, user.Name, user.Birth, user.Active, user.Location},
		).ToSQL()

	return repo.execModify(tracerCtx, query)
}

func (repo *Repository) Update(ctx context.Context, user *entities.User) error {
	tracerCtx, span := repo.tracer.Start(ctx, "PostgresRepository-Update")
	defer span.End()

	query, _, _ := goqu.Update("users").Set(user).Where(goqu.C("id").Eq(user.ID)).ToSQL()

	return repo.execModify(tracerCtx, query)
}

func (repo *Repository) Remove(ctx context.Context, id string) error {
	tracerCtx, span := repo.tracer.Start(ctx, "PostgresRepository-Remove")
	defer span.End()

	query, _, _ := goqu.Delete("users").Where(goqu.C("id").Eq(id)).ToSQL()

	return repo.execModify(tracerCtx, query)
}
