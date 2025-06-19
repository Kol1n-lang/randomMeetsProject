package validators

import (
	"github.com/go-playground/validator/v10"
)

type MeetUp struct {
	Title       string `json:"title" validate:"required,min=3"`
	Description string `json:"description" validate:"required,min=3"`
	Date        string `json:"date" validate:"required"`
	PeopleCount int    `json:"people_count" validate:"required,min=1"`
	PhotoURL    string `json:"photo_url"`
}

func (l *MeetUp) Validate() error {
	validate := validator.New()
	return validate.Struct(l)
}
