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

var (
	ErrProfileNotFound = errors.New("profile not found")
)
