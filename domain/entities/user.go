package entities

import "time"

type User struct {
	ID        string
	Name      string
	Birth     *time.Time
	Email     *string
	Location  *string
	CreatedAt time.Time
	UpdatedAt time.Time
	Active    bool
}
