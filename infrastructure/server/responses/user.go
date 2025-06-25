package responses

import (
	"users/domain/entities"
)

const dateLayout = "02/01/2006"

type UserResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Birth    string `json:"birth"`
	Active   bool   `json:"active"`
	Location string `json:"location"`
}

func FromUser(user *entities.User) *UserResponse {
	return &UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Birth:    user.Birth.Format(dateLayout),
		Active:   user.Active,
		Location: fromNullableString(user.Location),
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
