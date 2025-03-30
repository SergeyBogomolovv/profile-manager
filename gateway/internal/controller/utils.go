package controller

import (
	"context"
	"net/http"
	"time"

	"google.golang.org/grpc/metadata"
)

func newSecureCookie(name, value string, duration time.Duration) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(duration.Seconds()),
		Expires:  time.Now().Add(duration),
		Secure:   true,
		HttpOnly: true,
	}
}

func removeCookie(name string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
}

func authCtx(r *http.Request) context.Context {
	cookie, err := r.Cookie("access_token")
	if err != nil {
		return r.Context()
	}
	ctx := metadata.NewOutgoingContext(r.Context(), metadata.Pairs("authorization", "Bearer "+cookie.Value))
	return ctx
}
