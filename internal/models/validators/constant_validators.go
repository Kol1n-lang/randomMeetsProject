package validators

import (
	"github.com/go-playground/validator/v10"
)

type JWTTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type EmailPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     int    `json:"code"`
}

type PhotoURL struct {
	URL string `json:"url" validate:"required"`
}

func (url *PhotoURL) Validate() error {
	validate := validator.New()
	return validate.Struct(url)
}

func (emailPayload *EmailPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(emailPayload)
}
