package entity

import (
	helpers "darkness-awakens/helpers"

	"github.com/go-playground/validator/v10"
)

type RegisteringUser struct {
	Email       string `json:"email" validate:"email" message:"Invalid email"`
	DisplayName string `json:"name" validate:"required" message:"Name is required"`
	Password    string `json:"password" validate:"required" message:"Password is required"`
}

func (m *RegisteringUser) Validate(validate *validator.Validate) error {
	return helpers.ValidateFunc[RegisteringUser](*m, validate)
}

type LoggingInUser struct {
	Email    string `json:"email" validate:"email" message:"Invalid email"`
	Password string `json:"password" validate:"required" message:"Password is required"`
}

func (m *LoggingInUser) Validate(validate *validator.Validate) error {
	return helpers.ValidateFunc[LoggingInUser](*m, validate)
}

type NewOrLoggedInUser struct {
	Email       string `json:"email" validate:"email" message:"Invalid email"`
	DisplayName string `json:"name" validate:"required" message:"Name is required"`
}

func (m *NewOrLoggedInUser) Validate(validate *validator.Validate) error {
	return helpers.ValidateFunc[NewOrLoggedInUser](*m, validate)
}
