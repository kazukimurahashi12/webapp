package dto

import "github.com/go-playground/validator/v10"

type Session struct {
	UserID string `validate:"required,uuid4"`
	Token  string `validate:"required,min=32"`
}

func (s *Session) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}
