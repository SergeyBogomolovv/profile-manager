package domain

import (
	"errors"
	"time"
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

var (
	ErrInvalidToken = errors.New("invalid token")
)

const AccessTokenTTL = time.Hour
const RefreshTokenTTL = time.Hour * 24 * 7
