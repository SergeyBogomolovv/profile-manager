package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/SergeyBogomolovv/profile-manager/common/httpx"
	"github.com/SergeyBogomolovv/profile-manager/sso/internal/config"
	"github.com/SergeyBogomolovv/profile-manager/sso/internal/models"
	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const userInfoUrl = "https://www.googleapis.com/oauth2/v2/userinfo"

type OAuthService interface {
	GoogleSignIn(ctx context.Context, user models.OAuthUserInfo) (models.Tokens, error)
}

type oauthController struct {
	state string
	oauth *oauth2.Config
	svc   OAuthService
}

func NewOAuthController(conf config.OAuth, svc OAuthService) *oauthController {
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
	return &oauthController{svc: svc, oauth: oauth, state: "23490rfdslmfjn34i0skldfj"}
}

// Инициализация роутов
func (c *oauthController) Init(router *chi.Mux) {
	router.Route("/auth", func(r chi.Router) {
		r.Get("/google", c.HandleLogin)
		r.Get("/google/callback", c.HandleCallback)
	})
}

// Перенаправление пользователя на Google OAuth
func (c *oauthController) HandleLogin(w http.ResponseWriter, r *http.Request) {
	url := c.oauth.AuthCodeURL(c.state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Обработка коллбэка от Google
func (c *oauthController) HandleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != c.state {
		httpx.WriteError(w, "Invalid OAuth state", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	token, err := c.oauth.Exchange(r.Context(), code)
	if err != nil {
		httpx.WriteError(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	client := c.oauth.Client(r.Context(), token)
	resp, err := client.Get(userInfoUrl)
	if err != nil {
		httpx.WriteError(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var user models.OAuthUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		httpx.WriteError(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	tokens, err := c.svc.GoogleSignIn(r.Context(), user)
	if err != nil {
		httpx.WriteError(w, "Failed to sign in", http.StatusInternalServerError)
		return
	}
	httpx.WriteJSON(w, tokens, http.StatusOK)
}
