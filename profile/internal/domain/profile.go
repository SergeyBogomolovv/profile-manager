package domain

import "errors"

type UserGender string

const (
	UserGenderMale         UserGender = "male"
	UserGenderFemale       UserGender = "female"
	UserGenderNotSpecified UserGender = "not specified"
)

type Profile struct {
	UserID    string
	Username  string
	FirstName string
	LastName  string
	BirthDate string
	Gender    UserGender
	Avatar    string
}

type UpdateProfileDTO struct {
	Username  string
	FirstName string
	LastName  string
	BirthDate string     `validate:"omitempty,datetime=2006-01-02"`
	Gender    UserGender `validate:"omitempty,oneof=male female"`
	Avatar    []byte
}

var (
	ErrProfileNotFound = errors.New("profile not found")
	ErrUsernameExists  = errors.New("username already exists")
)
