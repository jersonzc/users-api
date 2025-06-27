package requests

import (
	"fmt"
	"strconv"
	"time"
)

type UpdateUser struct {
	Name     *string `json:"name"`
	Birth    *string `json:"birth"`
	Email    *string `json:"email"`
	Location *string `json:"location"`
	Active   *string `json:"active"`
}

func (p *UpdateUser) ToMap() (*map[string]interface{}, error) {
	fields := make(map[string]interface{})

	if p.Name != nil {
		fields["name"] = *p.Name
	}

	if p.Birth != nil {
		birth, err := time.Parse(dateLayout, *p.Birth)
		if err != nil {
			return nil, fmt.Errorf("error while parsing 'birth' field from %q: %w", p.Birth, err)
		}
		fields["birth"] = birth
	}

	if p.Email != nil {
		fields["email"] = *p.Email
	}

	if p.Location != nil {
		if *p.Location == "" {
			fields["location"] = nil
		} else {
			fields["location"] = *p.Location
		}
	}

	if p.Active != nil {
		active, err := strconv.ParseBool(*p.Active)
		if err != nil {
			return nil, fmt.Errorf("error while parsing 'active' field from %q: %w", p.Active, err)
		}
		fields["active"] = active
	}

	fields["updated_at"] = time.Now().UTC()

	return &fields, nil
}
