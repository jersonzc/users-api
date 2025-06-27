package mappers

import (
	"log"
	"strconv"
	"time"
	"users/domain/entities"
)

func ToUserList(rows []map[string]string) []*entities.User {
	users := make([]*entities.User, len(rows))
	for i, row := range rows {
		users[i] = ToUser(row)
	}
	return users
}

func ToUser(row map[string]string) *entities.User {
	var user entities.User

	createdAt, err := time.Parse(time.DateTime, row["created_at"])
	if err != nil {
		log.Printf("Error while parsing 'created_at' field: %s", err.Error())
		return &user
	}

	updatedAt, err := time.Parse(time.DateTime, row["updated_at"])
	if err != nil {
		log.Printf("Error while parsing 'updated_at' field: %s", err.Error())
		return &user
	}

	active, err := strconv.ParseBool(row["active"])
	if err != nil {
		log.Printf("Error while parsing 'active' field: %s", err.Error())
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
	myTime, err := time.Parse(time.DateTime, result)
	if err != nil {
		return nil
	}
	return &myTime
}
