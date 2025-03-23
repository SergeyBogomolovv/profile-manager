package domain

import (
	"errors"
	"time"
)

const TokenTTL = time.Minute * 15

var ErrInvalidToken = errors.New("invalid token")
