package requests

import (
	"fmt"
	"github.com/google/uuid"
	"time"
	"users/domain/entities"
)

const dateLayout = "02/01/2006"

type SaveUser struct {
	Name     string `json:"name" binding:"required"`
	Birth    string `json:"birth"`
	Email    string `json:"email"`
	Location string `json:"location"`
}

func (p *SaveUser) ToUser() (*entities.User, error) {
	birth, err := parseTime(dateLayout, p.Birth)
	if err != nil {
		return nil, fmt.Errorf("error while parsing 'birth' field from %q: %w", p.Birth, err)
	}

	return &entities.User{
		ID:       uuid.New().String(),
		Name:     p.Name,
		Birth:    toNullableTime(birth),
		Email:    toNullableString(p.Email),
		Location: toNullableString(p.Location),
		Active:   true,
	}, nil
}

func parseTime(layout string, value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}

	result, err := time.Parse(layout, value)
	if err != nil {
		return time.Time{}, err
	}

	return result, nil
}

func toNullableString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func toNullableTime(value time.Time) *time.Time {
	if value.IsZero() {
		return nil
	}
	return &value
}
