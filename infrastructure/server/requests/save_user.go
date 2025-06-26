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
	uuidV4 := uuid.New()

	birth, err := time.Parse(dateLayout, p.Birth)
	if err != nil {
		return nil, fmt.Errorf("error while parsing 'birth' field from %q: %w", p.Birth, err)
	}

	now := time.Now().UTC()

	return &entities.User{
		ID:        uuidV4.String(),
		Name:      p.Name,
		Birth:     toNullableTime(birth),
		Email:     toNullableString(p.Email),
		Location:  toNullableString(p.Location),
		CreatedAt: now,
		UpdatedAt: now,
		Active:    true,
	}, nil
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
