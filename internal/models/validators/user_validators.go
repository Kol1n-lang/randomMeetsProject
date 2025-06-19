package validators

import (
	"github.com/go-playground/validator/v10"
)

type LoginUser struct {
	Username string `json:"username" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type User struct {
	Username string `json:"username" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	PhotoURL string `json:"photo_url" validate:"required,url"`
	Meet     MeetUp
}

func (l *LoginUser) Validate() error {
	validate := validator.New()
	return validate.Struct(l)
}
