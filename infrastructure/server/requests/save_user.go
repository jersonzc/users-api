package requests

import (
	"fmt"
	"strconv"
	"time"
	"users/domain/entities"
)

const dateLayout = "02/01/2006"

type SaveUser struct {
	ID       string `json:"id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Birth    string `json:"birth" binding:"required"`
	Active   string `json:"active" binding:"required"`
	Location string `json:"location"`
}

func (p *SaveUser) ToUser() (*entities.User, error) {
	id, err := strconv.Atoi(p.ID)
	if err != nil {
		return nil, fmt.Errorf("error while parsing 'id' field from %q: %w", p.ID, err)
	}

	birth, err := time.Parse(dateLayout, p.Birth)
	if err != nil {
		return nil, fmt.Errorf("error while parsing 'birth' field from %q: %w", p.Birth, err)
	}

	active, err := strconv.ParseBool(p.Active)
	if err != nil {
		return nil, fmt.Errorf("error while parsing 'active' field from %q: %w", p.Active, err)
	}

	return &entities.User{
		ID:       id,
		Name:     p.Name,
		Birth:    birth,
		Active:   active,
		Location: toNullableString(p.Location),
	}, nil
}

func toNullableString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
