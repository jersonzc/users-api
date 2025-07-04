package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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

	rows, err := repo.client.queries.GetUsers(tracerCtx, ids)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(attribute.Int("repo.postgres.rows.count", len(rows)))

	return toUserList(rows), nil
}

func (repo *Repository) Save(ctx context.Context, user *entities.User) (*entities.User, error) {
	tracerCtx, span := repo.tracer.Start(ctx, "PostgresRepository-Save")
	defer span.End()

	arg, err := toSaveUserParams(user)
	if err != nil {
		return nil, err
	}

	row, err := repo.client.queries.CreateUser(tracerCtx, arg)
	if err != nil {
		return nil, err
	}

	return toUser(row), nil
}

func (repo *Repository) Update(ctx context.Context, id string, fields map[string]interface{}) (*entities.User, error) {
	tracerCtx, span := repo.tracer.Start(ctx, "PostgresRepository-Update")
	defer span.End()

	arg, err := toUpdateUserParams(id, fields)
	if err != nil {
		return nil, err
	}

	row, err := repo.client.queries.UpdateUser(tracerCtx, arg)
	if err != nil {
		return nil, err
	}

	return toUser(row), nil
}

func (repo *Repository) Remove(ctx context.Context, id string) error {
	tracerCtx, span := repo.tracer.Start(ctx, "PostgresRepository-Remove")
	defer span.End()

	return repo.client.queries.DeleteUser(tracerCtx, id)
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

	var birth *time.Time
	if row.Birth.Valid {
		birth = &row.Birth.Time
	}

	var email *string
	if row.Email.Valid {
		email = &row.Email.String
	}

	var location *string
	if row.Location.Valid {
		location = &row.Location.String
	}

	user.ID = row.ID
	user.Name = row.Name
	user.Birth = birth
	user.Email = email
	user.Location = location
	user.CreatedAt = row.CreatedAt.Time
	user.UpdatedAt = row.UpdatedAt.Time
	user.Active = row.Active

	return &user
}

func toSaveUserParams(user *entities.User) (CreateUserParams, error) {
	var birth pgtype.Date
	if user.Birth != nil {
		if err := birth.Scan(*user.Birth); err != nil {
			return CreateUserParams{}, err
		}
	} else {
		birth.Valid = false
	}

	var email pgtype.Text
	if user.Email != nil {
		email.String = *user.Email
		email.Valid = true
	} else {
		email.Valid = false
	}

	var location pgtype.Text
	if user.Location != nil {
		location.String = *user.Location
		location.Valid = true
	} else {
		location.Valid = false
	}
	return CreateUserParams{
		ID:       user.ID,
		Name:     user.Name,
		Birth:    birth,
		Email:    email,
		Location: location,
		Active:   user.Active,
	}, nil
}

func toUpdateUserParams(id string, fields map[string]interface{}) (UpdateUserParams, error) {
	var row UpdateUserParams

	// ID
	row.ID = id

	// Name
	if value, ok := fields["name"]; ok {
		row.NameDoUpdate = true
		row.Name = value.(string)
	}

	// Birth
	if value, ok := fields["birth"]; ok {
		row.BirthDoUpdate = true
		if value != nil {
			if err := row.Birth.Scan(value); err != nil {
				return UpdateUserParams{}, err
			}
		} else {
			row.Birth.Valid = false
		}
	}

	// Email
	if value, ok := fields["email"]; ok {
		row.EmailDoUpdate = true
		if value != nil {
			row.Email.String = value.(string)
			row.Email.Valid = true
		} else {
			row.Email.Valid = false
		}
	}

	// Location
	if value, ok := fields["location"]; ok {
		row.LocationDoUpdate = true
		if value != nil {
			row.Location.String = value.(string)
			row.Location.Valid = true
		} else {
			row.Location.Valid = false
		}
	}

	// Active
	if value, ok := fields["active"]; ok {
		row.ActiveDoUpdate = true
		row.Active = value.(bool)
	}

	return row, nil
}
