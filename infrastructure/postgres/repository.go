package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"strings"
	"time"
	"users/domain/entities"
)

type Repository struct {
	client *Client
	tracer trace.Tracer
}

func NewRepository(
	client *Client,
) (*Repository, error) {
	return &Repository{
		client: client,
		tracer: otel.Tracer("PostgresRepository")}, nil
}

func (repo *Repository) Get(ctx context.Context) ([]*entities.User, error) {
	tracerCtx, span := repo.tracer.Start(ctx, "PostgresRepository-Get")
	defer span.End()

	rows, err := repo.client.queries.ListActiveUsers(tracerCtx)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(attribute.Int("repo.postgres.rows.count", len(rows)))

	return toUserList(rows), nil
}

func (repo *Repository) GetByID(ctx context.Context, ids []string) ([]*entities.User, error) {
	tracerCtx, span := repo.tracer.Start(ctx, "PostgresRepository-GetByID")
	defer span.End()

	rows, err := repo.client.queries.GetUsers(tracerCtx, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}

	span.SetAttributes(attribute.Int("repo.postgres.rows.count", len(rows)))

	return toUserList(rows), nil
}

func (repo *Repository) Save(ctx context.Context, user *entities.User) error {
	tracerCtx, span := repo.tracer.Start(ctx, "PostgresRepository-Save")
	defer span.End()

	_, err := repo.client.queries.CreateUser(tracerCtx, toSaveUserParams(user))
	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) Update(ctx context.Context, id string, fields map[string]interface{}) error {
	tracerCtx, span := repo.tracer.Start(ctx, "PostgresRepository-Update")
	defer span.End()

	err := repo.client.queries.UpdateUser(tracerCtx, toUpdateUserParams(id, fields))
	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) Remove(ctx context.Context, id string) error {
	tracerCtx, span := repo.tracer.Start(ctx, "PostgresRepository-Remove")
	defer span.End()

	err := repo.client.queries.DeleteUser(tracerCtx, id)
	if err != nil {
		return err
	}

	return nil
}

func toUserList(rows []User) []*entities.User {
	users := make([]*entities.User, len(rows))
	for i, row := range rows {
		users[i] = toUser(row)
	}
	return users
}

func toUser(row User) *entities.User {
	var user entities.User

	user.ID = row.ID
	user.Name = row.Name
	user.Birth = &row.Birth.Time
	user.Email = &row.Email.String
	user.Location = &row.Location.String
	user.CreatedAt = row.CreatedAt.Time
	user.UpdatedAt = row.UpdatedAt.Time
	user.Active = row.Active

	return &user
}

func toSaveUserParams(user *entities.User) CreateUserParams {
	var row CreateUserParams

	row.ID = user.ID
	row.Name = user.Name
	row.Birth = pgtype.Date{Time: fromNullableTime(user.Birth), Valid: true}
	row.Email = pgtype.Text{String: fromNullableString(user.Email), Valid: true}
	row.Location = pgtype.Text{String: fromNullableString(user.Location), Valid: true}
	row.CreatedAt = pgtype.Timestamp{Time: user.CreatedAt, Valid: true}
	row.UpdatedAt = pgtype.Timestamp{Time: user.UpdatedAt, Valid: true}
	row.Active = user.Active

	return row
}

func toUpdateUserParams(id string, fields map[string]interface{}) UpdateUserParams {
	var row UpdateUserParams

	if value, ok := fields["name"]; ok {
		row.NameDoUpdate = true
		row.Name = value.(string)
	}

	if value, ok := fields["birth"]; ok {
		row.BirthDoUpdate = true
		row.Birth = value.(pgtype.Date)
	}

	if value, ok := fields["email"]; ok {
		row.EmailDoUpdate = true
		row.Email = value.(string)
	}

	if value, ok := fields["location"]; ok {
		row.LocationDoUpdate = true
		row.Location = value.(string)
	}

	if value, ok := fields["active"]; ok {
		row.ActiveDoUpdate = true
		row.Active = value.(bool)
	}

	row.ID = id
	row.UpdatedAt = pgtype.Timestamp{Time: fields["updated_at"].(time.Time), Valid: true}

	return row
}

func fromNullableString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func fromNullableTime(t *time.Time) time.Time {
	if t != nil {
		return *t
	}
	return time.Time{}
}
