package controller

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/SergeyBogomolovv/profile-manager/common/httpx"
	"github.com/SergeyBogomolovv/profile-manager/sso/internal/config"
	"github.com/SergeyBogomolovv/profile-manager/sso/internal/domain"
	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const userInfoUrl = "https://www.googleapis.com/oauth2/v2/userinfo"

type OAuthService interface {
	OAuth(ctx context.Context, info domain.OAuthUserInfo, provider domain.AccountType) (domain.Tokens, error)
}

type httpController struct {
	logger *slog.Logger
	state  string
	oauth  *oauth2.Config
	svc    OAuthService
}

func NewHTTPController(logger *slog.Logger, conf config.OAuth, svc OAuthService) *httpController {
	oauth := &oauth2.Config{
		ClientID:     conf.ClientID,
		ClientSecret: conf.ClientSecret,
		RedirectURL:  conf.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
	return &httpController{svc: svc, oauth: oauth, state: "23490rfdslmfjn34i0skldfj", logger: logger}
}

// Инициализация роутов
func (c *httpController) Init(router *chi.Mux) {
	router.Route("/auth", func(r chi.Router) {
		r.Get("/google", c.HandleLogin)
		r.Get("/google/callback", c.HandleCallback)
	})
}

// Перенаправление пользователя на Google OAuth
func (c *httpController) HandleLogin(w http.ResponseWriter, r *http.Request) {
	url := c.oauth.AuthCodeURL(c.state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Обработка коллбэка от Google
func (c *httpController) HandleCallback(w http.ResponseWriter, r *http.Request) {
	const op = "oauth.HandleCallback"
	logger := c.logger.With(slog.String("op", op))

	state := r.FormValue("state")
	if state != c.state {
		httpx.WriteError(w, "Invalid OAuth state", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	token, err := c.oauth.Exchange(r.Context(), code)
	if err != nil {
		logger.Error("failed to exchange token", "err", err)
		httpx.WriteError(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	client := c.oauth.Client(r.Context(), token)
	resp, err := client.Get(userInfoUrl)
	if err != nil {
		logger.Error("failed to get user info", "err", err)
		httpx.WriteError(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var user domain.OAuthUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		logger.Error("failed to decode user info", "err", err)
		httpx.WriteError(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	tokens, err := c.svc.OAuth(r.Context(), user, domain.AccountTypeGoogle)
	if err != nil {
		logger.Error("failed to sign in", "err", err)
		httpx.WriteError(w, "Failed to sign in", http.StatusInternalServerError)
		return
	}
	httpx.WriteJSON(w, tokens, http.StatusOK)
}
