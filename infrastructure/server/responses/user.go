package responses

import (
	"time"
	"users/domain/entities"
)

const dateLayout = "02/01/2006"

type UserResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Birth     string  `json:"birth"`
	Email     string  `json:"email"`
	Location  *string `json:"location"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	Active    bool    `json:"active"`
}

func FromUser(user *entities.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Birth:     fromNullableTime(user.Birth),
		Email:     fromNullableString(user.Email),
		Location:  user.Location,
		CreatedAt: user.CreatedAt.Format(time.DateTime),
		UpdatedAt: user.UpdatedAt.Format(time.DateTime),
		Active:    user.Active,
	}
}

func FromUserList(users []*entities.User) []*UserResponse {
	result := make([]*UserResponse, len(users))
	for i, user := range users {
		result[i] = FromUser(user)
	}
	return result
}

func fromNullableString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func fromNullableTime(t *time.Time) string {
	if t != nil {
		return (*t).Format(dateLayout)
	}
	return ""
}
