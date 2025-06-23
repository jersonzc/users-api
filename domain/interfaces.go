package domain

import (
	"context"
	"users/domain/entities"
)

type Get func(context.Context) ([]*entities.User, error)

type GetByID func(context.Context, []int) ([]*entities.User, error)

type Save func(context.Context, *entities.User) error

type Update func(context.Context, *entities.User) error

type Remove func(context.Context, int) error
