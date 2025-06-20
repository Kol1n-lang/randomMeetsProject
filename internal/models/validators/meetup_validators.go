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

type PutMeetUp struct {
	Title       string `json:"title" validate:"min=3"`
	Description string `json:"description" validate:"min=3"`
	Date        string `json:"date" validate:"required"`
	PeopleCount int    `json:"people_count" validate:"min=1"`
}

type GetMeetUp struct {
	Title       string `json:"title" validate:"required,min=3"`
	Description string `json:"description" validate:"required,min=3"`
	Date        string `json:"date" validate:"required"`
	PeopleCount int    `json:"people_count" validate:"required,min=1"`
	PhotoURL    string `json:"photo_url"`
	TGContact   string `json:"tg_user"`
}

func (l *MeetUp) Validate() error {
	validate := validator.New()
	return validate.Struct(l)
}

func (l *PutMeetUp) Validate() error {
	validate := validator.New()
	return validate.Struct(l)
}
