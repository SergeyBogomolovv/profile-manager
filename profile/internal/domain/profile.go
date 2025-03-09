package domain

import "github.com/google/uuid"

type UserGender string

const (
	UserGenderMale         UserGender = "male"
	UserGenderFemale       UserGender = "female"
	UserGenderNotSpecified UserGender = "not specified"
)

type Profile struct {
	UserID    uuid.UUID
	Username  string
	FirstName string
	LastName  string
	BirthDate string
	Gender    UserGender
	Avatar    string
}
