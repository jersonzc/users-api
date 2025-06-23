package entities

import "time"

type User struct {
	ID       int
	Name     string
	Birth    time.Time
	Active   bool
	Location *string
}
