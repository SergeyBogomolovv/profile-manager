package controller

type RegisterRequest struct {
	Email    string `json:"email" validate:"email" example:"xLb3u@example.com"`
	Password string `json:"password" validate:"min=6" example:"password"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"email" example:"xLb3u@example.com"`
	Password string `json:"password" validate:"min=6" example:"password"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token" example:"access_token"`
}

type RegisterResponse struct {
	Message string `json:"message" example:"user registered successfully"`
	UserID  string `json:"user_id" example:"user_id"`
}

type ProfileResponse struct {
	UserID    string `json:"user_id" example:"user_id"`
	Username  string `json:"username,omitempty" example:"username"`
	FirstName string `json:"first_name,omitempty" example:"first_name"`
	LastName  string `json:"last_name,omitempty" example:"last_name"`
	BirthDate string `json:"birth_date,omitempty" example:"birth_date"`
	Gender    string `json:"gender,omitempty" example:"gender"`
	Avatar    string `json:"avatar,omitempty" example:"avatar"`
}

type UpdateProfileRequest struct {
	Username  string `form:"username" json:"username" example:"username" validate:"omitempty"`
	FirstName string `form:"first_name" json:"first_name" example:"John" validate:"omitempty"`
	LastName  string `form:"last_name" json:"last_name" example:"Doe" validate:"omitempty"`
	BirthDate string `form:"birth_date" json:"birth_date" example:"2000-01-01" validate:"omitempty,datetime=2006-01-02"`
	Gender    string `form:"gender" json:"gender" example:"male" validate:"omitempty,oneof=male female"`
	Avatar    []byte `form:"avatar" json:"avatar" swaggerignore:"true" validate:"omitempty"`
}

type TokenResponse struct {
	Token string `json:"token" example:"sf34fdsfsdf-sdf3ef"`
}
