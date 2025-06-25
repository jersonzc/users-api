package requests

import (
	"fmt"
	"strconv"
	"time"
	"users/domain/entities"
)

type UpdateUser struct {
	Name     string `json:"name"`
	Birth    string `json:"birth"`
	Active   string `json:"active"`
	Location string `json:"location"`
}

func (p *UpdateUser) ToUser() (*entities.User, error) {
	birth, err := time.Parse(dateLayout, p.Birth)
	if err != nil {
		return nil, fmt.Errorf("error while parsing 'birth' field from %q: %w", p.Birth, err)
	}

	active, err := strconv.ParseBool(p.Active)
	if err != nil {
		return nil, fmt.Errorf("error while parsing 'active' field from %q: %w", p.Active, err)
	}

	return &entities.User{
		Name:     p.Name,
		Birth:    birth,
		Active:   active,
		Location: toNullableString(p.Location),
	}, nil
}
