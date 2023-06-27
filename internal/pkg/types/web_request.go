package types

import (
	"time"

	"regexp"

	"github.com/go-playground/validator/v10"
)

const (
	patternAlphaNumeric            = "^[a-zA-Z0-9_-]+$"
	patternAlphaNumericSpecialChar = `^[a-zA-Z0-9~!@#$%^&*()-_=+{}\|;:'",<.>/?]+$`
)

type LoginRequest struct {
	UserName string `json:"username" validate:"required,username"`
	Password string `json:"password" validate:"required,password"`
}

type CreateUserRequest struct {
	UserName  string `json:"username" validate:"required,username"`
	Password  string `json:"password" validate:"required,password"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	DoB       string `json:"dateofbirth" validate:"required,dob"`
	Email     string `json:"email" validate:"required,email"`
}

type EditUserRequest struct {
	Password  string `json:"password" validate:"omitempty,password"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	DoB       string `json:"dateofbirth" validate:"omitempty,dob"`
	Email     string `json:"email" validate:"omitempty,email"`
}

type CreatePostRequest struct {
	ContentText      string   `json:"content_text" validate:"required"`
	ContentImagePath []string `json:"content_image_path" validate:"omitempty,dive,url"`
}

type EditPostRequest struct {
	ContentText      string   `json:"content_text" validate:"omitempty"`
	ContentImagePath []string `json:"content_image_path" validate:"omitempty,dive,url"`
}

type CommentPostRequest struct {
	Content string `json:"content"`
}

func NewValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("dob", validateDOB)
	validate.RegisterValidation("username", validateUsername)
	validate.RegisterValidation("password", validatePassword)

	return validate
}

func validateDOB(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()

	// Define the expected date format
	dateFormat := time.DateOnly

	// Parse the date string into a time.Time value
	_, err := time.Parse(dateFormat, dateStr)

	return err == nil
}

func validatePassword(fl validator.FieldLevel) bool {
	if len(fl.Field().String()) < 4 || len(fl.Field().String()) > 200 {
		return false
	}

	alphaRegex, err := regexp.Compile(patternAlphaNumericSpecialChar)
	if err != nil {
		return false
	}
	return alphaRegex.MatchString(fl.Field().String())
}

func validateUsername(fl validator.FieldLevel) bool {
	if len(fl.Field().String()) < 4 || len(fl.Field().String()) > 200 {
		return false
	}

	alphaNumRegex, err := regexp.Compile(patternAlphaNumeric)
	if err != nil {
		return false
	}
	return alphaNumRegex.MatchString(fl.Field().String())
}