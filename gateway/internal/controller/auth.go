package controller

import (
	"log/slog"
	"net/http"
	"time"

	pb "github.com/SergeyBogomolovv/profile-manager/common/api/sso"
	"github.com/SergeyBogomolovv/profile-manager/common/httpx"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authController struct {
	logger   *slog.Logger
	validate *validator.Validate
	client   pb.SSOClient
}

func NewAuthController(logger *slog.Logger, client pb.SSOClient) *authController {
	validate := validator.New()
	return &authController{validate: validate, logger: logger, client: client}
}

// Init initializes authentication routes.
func (c *authController) Init(r *chi.Mux) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", c.HandleLogin)
		r.Post("/register", c.HandleRegister)
		r.Post("/refresh", c.HandleRefresh)
		r.Post("/logout", c.HandleLogout)
	})
}

// HandleLogin processes user login.
// @Summary User login
// @Description Authenticates a user with email and password and returns an access token.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "User login credentials"
// @Success 200 {object} AccessTokenResponse "Successful login, returns access token"
// @Failure 400 {object} httpx.ErrorResponse "Bad request"
// @Failure 401 {object} httpx.ErrorResponse "Invalid credentials"
// @Router /auth/login [post]
func (c *authController) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var body LoginRequest
	if err := httpx.DecodeBody(r, &body); err != nil {
		httpx.WriteError(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	if err := c.validate.Struct(body); err != nil {
		httpx.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := c.client.Login(r.Context(), &pb.LoginRequest{
		Email:    body.Email,
		Password: body.Password,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			httpx.WriteError(w, "Failed to login", http.StatusInternalServerError)
			return
		}
		switch st.Code() {
		case codes.InvalidArgument:
			httpx.WriteError(w, "Invalid parameters", http.StatusBadRequest)
		case codes.Unauthenticated:
			httpx.WriteError(w, "Invalid credentials", http.StatusUnauthorized)
		default:
			httpx.WriteError(w, "Failed to login", http.StatusInternalServerError)
		}
		return
	}

	http.SetCookie(w, newSecureCookie("refresh_token", resp.RefreshToken, time.Hour*24*7))
	http.SetCookie(w, newSecureCookie("access_token", resp.AccessToken, time.Hour*24))

	httpx.WriteSuccess(w, "Login successful", http.StatusOK)
}

// HandleRegister registers a new user.
// @Summary User registration
// @Description Registers a new user using email and password and returns the created user ID.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "User registration data"
// @Success 201 {object} RegisterResponse "User registered successfully, returns user ID"
// @Failure 400 {object} httpx.ErrorResponse "Invalid data or bad request"
// @Failure 409 {object} httpx.ErrorResponse "User with this email already exists"
// @Router /auth/register [post]
func (c *authController) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var body RegisterRequest
	if err := httpx.DecodeBody(r, &body); err != nil {
		httpx.WriteError(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	if err := c.validate.Struct(body); err != nil {
		httpx.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := c.client.Register(r.Context(), &pb.RegisterRequest{
		Email:    body.Email,
		Password: body.Password,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			httpx.WriteError(w, "Failed to register", http.StatusInternalServerError)
			return
		}
		switch st.Code() {
		case codes.InvalidArgument:
			httpx.WriteError(w, "Invalid parameters", http.StatusBadRequest)
		case codes.AlreadyExists:
			httpx.WriteError(w, "User with this email already exists", http.StatusConflict)
		default:
			httpx.WriteError(w, "Failed to register", http.StatusInternalServerError)
		}
		return
	}

	httpx.WriteJSON(w, RegisterResponse{UserID: resp.UserId, Message: "user registered successfully"}, http.StatusCreated)
}

// HandleRefresh refreshes an access token.
// @Summary Refresh access token
// @Description Refreshes the user's access token using the refresh token stored in cookies.
// @Tags auth
// @Accept json
// @Produce json
// @Param Cookie header string true "Refresh token stored in cookie (refresh_token=<token>)"
// @Success 200 {object} AccessTokenResponse "New access token generated successfully"
// @Failure 401 {object} httpx.ErrorResponse "Unauthorized"
// @Router /auth/refresh [post]
func (c *authController) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		httpx.WriteError(w, "Refresh token not found", http.StatusUnauthorized)
		return
	}
	resp, err := c.client.Refresh(r.Context(), &pb.RefreshRequest{
		RefreshToken: cookie.Value,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			httpx.WriteError(w, "Failed to refresh token", http.StatusInternalServerError)
			return
		}
		switch st.Code() {
		case codes.Unauthenticated:
			httpx.WriteError(w, "Invalid refresh token", http.StatusUnauthorized)
		default:
			httpx.WriteError(w, "Failed to refresh token", http.StatusInternalServerError)
		}
		return
	}

	http.SetCookie(w, newSecureCookie("access_token", resp.AccessToken, time.Hour*24))
	httpx.WriteSuccess(w, "Token refreshed", http.StatusOK)
}

// HandleLogout logs out the user.
// @Summary User logout
// @Description Logs out the user by invalidating the refresh token stored in cookies.
// @Tags auth
// @Accept json
// @Produce json
// @Param Cookie header string true "Refresh token stored in cookie (refresh_token=<token>)"
// @Success 200 {object} httpx.SuccessResponse "User logged out successfully"
// @Failure 401 {object} httpx.ErrorResponse "Unauthorized"
// @Router /auth/logout [post]
func (c *authController) HandleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		httpx.WriteError(w, "Refresh token not found", http.StatusUnauthorized)
		return
	}
	_, err = c.client.Logout(r.Context(), &pb.LogoutRequest{
		RefreshToken: cookie.Value,
	})
	if err != nil {
		httpx.WriteError(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, removeCookie("refresh_token"))
	http.SetCookie(w, removeCookie("access_token"))

	httpx.WriteSuccess(w, "Logout successful", http.StatusOK)
}
