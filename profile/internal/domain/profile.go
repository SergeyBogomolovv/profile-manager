package domain

type UserGender string

const (
	UserGenderMale         UserGender = "male"
	UserGenderFemale       UserGender = "female"
	UserGenderNotSpecified UserGender = "not specified"
)

type Profile struct {
	UserID    int64
	Username  string
	FirstName string
	LastName  string
	BirthDate string
	Gender    UserGender
	Avatar    string
}
