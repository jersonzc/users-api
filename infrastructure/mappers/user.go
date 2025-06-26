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

	createdAt, err := time.Parse(dateLayout, row["createdAt"])
	if err != nil {
		return &user
	}

	updatedAt, err := time.Parse(dateLayout, row["updatedAt"])
	if err != nil {
		return &user
	}

	active, err := strconv.ParseBool(row["active"])
	if err != nil {
		return &user
	}

	user.ID = row["id"]
	user.Name = row["name"]
	user.Birth = timeFromNullableColumn(row, "birth")
	user.Email = stringFromNullableColumn(row, "email")
	user.Location = stringFromNullableColumn(row, "location")
	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt
	user.Active = active

	return &user
}

func stringFromNullableColumn(row map[string]string, column string) *string {
	result, ok := row[column]
	if !ok {
		return nil
	}
	return &result
}

func timeFromNullableColumn(rows map[string]string, column string) *time.Time {
	result, ok := rows[column]
	if !ok {
		return nil
	}
	myTime, err := time.Parse(dateLayout, result)
	if err != nil {
		return nil
	}
	return &myTime
}
