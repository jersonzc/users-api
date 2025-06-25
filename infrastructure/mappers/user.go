package mappers

import (
	"strconv"
	"time"
	"users/domain/entities"
)

const dateLayout = "2006-01-02 15:04:05 -0700 MST"

func ToUserList(rows []map[string]string) []*entities.User {
	users := make([]*entities.User, len(rows))
	for i, row := range rows {
		users[i] = ToUser(row)
	}
	return users
}

func ToUser(row map[string]string) *entities.User {
	var user entities.User

	id, err := strconv.Atoi(row["id"])
	if err != nil {
		return &user
	}

	birth, err := time.Parse(dateLayout, row["birth"])
	if err != nil {
		return &user
	}

	active, err := strconv.ParseBool(row["active"])
	if err != nil {
		return &user
	}

	user.ID = id
	user.Name = row["name"]
	user.Birth = birth
	user.Active = active
	user.Location = StringFromNullableColumn(row, "location")

	return &user
}

func StringFromNullableColumn(row map[string]string, column string) *string {
	result, ok := row[column]
	if !ok {
		return nil
	}
	return &result
}
