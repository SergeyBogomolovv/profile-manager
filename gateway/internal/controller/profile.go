package controller

import (
	"io"
	"log/slog"
	"net/http"

	pb "github.com/SergeyBogomolovv/profile-manager/common/api/profile"
	"github.com/SergeyBogomolovv/profile-manager/common/httpx"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type profileController struct {
	logger   *slog.Logger
	validate *validator.Validate
	client   pb.ProfileClient
}

func NewProfileController(logger *slog.Logger, client pb.ProfileClient) *profileController {
	validate := validator.New()
	return &profileController{validate: validate, logger: logger, client: client}
}

// Init initializes profile routes.
func (c *profileController) Init(r *chi.Mux) {
	r.Route("/profile", func(r chi.Router) {
		r.Post("/update", c.HandleUpdate)
		r.Get("/my", c.HandleGet)
	})
}

// HandleUpdate updates the user's profile.
// @Summary Update user profile
// @Description Updates the authenticated user's profile information using multipart/form-data.
// @Tags profile
// @Accept multipart/form-data
// @Produce json
// @Param username formData string false "Username"
// @Param first_name formData string false "First name"
// @Param last_name formData string false "Last name"
// @Param birth_date formData string false "Birth date (YYYY-MM-DD)"
// @Param gender formData string false "Gender (male or female)"
// @Param avatar formData file false "Profile avatar"
// @Success 200 {object} ProfileResponse "Profile updated successfully"
// @Failure 400 {object} httpx.ErrorResponse "Validation error or bad request"
// @Failure 401 {object} httpx.ErrorResponse "Unauthorized"
// @Router /profile/update [post]
// @Security BearerAuth
func (c *profileController) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		httpx.WriteError(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}
	req := UpdateProfileRequest{
		Username:  r.FormValue("username"),
		FirstName: r.FormValue("first_name"),
		LastName:  r.FormValue("last_name"),
		BirthDate: r.FormValue("birth_date"),
		Gender:    r.FormValue("gender"),
	}
	// process avatar
	file, _, err := r.FormFile("avatar")
	if err == nil {
		req.Avatar, err = io.ReadAll(file)
		if err != nil {
			httpx.WriteError(w, "Failed to read avatar file", http.StatusBadRequest)
			return
		}
		file.Close()
	} else if err != http.ErrMissingFile {
		httpx.WriteError(w, "Failed to process avatar", http.StatusBadRequest)
		return
	}
	// validate
	if err := c.validate.Struct(req); err != nil {
		httpx.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := c.client.UpdateProfile(authCtx(r), &pb.UpdateProfileRequest{
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		BirthDate: req.BirthDate,
		Gender:    req.Gender,
		Avatar:    req.Avatar,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			httpx.WriteError(w, "Failed to update profile", http.StatusInternalServerError)
			return
		}
		switch st.Code() {
		case codes.InvalidArgument:
			httpx.WriteError(w, "Invalid request", http.StatusBadRequest)
		case codes.AlreadyExists:
			httpx.WriteError(w, "Username already exists", http.StatusConflict)
		case codes.Unauthenticated:
			httpx.WriteError(w, "Unauthorized", http.StatusUnauthorized)
		case codes.NotFound:
			httpx.WriteError(w, "Profile not found", http.StatusNotFound)
		default:
			httpx.WriteError(w, "Failed to update profile", http.StatusInternalServerError)
		}
		return
	}

	httpx.WriteJSON(w, profileResponse(resp), http.StatusOK)
}

// HandleGet retrieves the authenticated user's profile.
// @Summary Get user profile
// @Description Retrieves the profile of the authenticated user
// @Tags profile
// @Accept json
// @Produce json
// @Success 200 {object} ProfileResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /profile/my [get]
// @Security BearerAuth
func (c *profileController) HandleGet(w http.ResponseWriter, r *http.Request) {
	resp, err := c.client.GetProfile(authCtx(r), &pb.GetProfileRequest{})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			httpx.WriteError(w, "Failed to get profile", http.StatusInternalServerError)
			return
		}
		switch st.Code() {
		case codes.InvalidArgument:
			httpx.WriteError(w, "Invalid user id", http.StatusBadRequest)
		case codes.NotFound:
			httpx.WriteError(w, "Profile not found", http.StatusNotFound)
		case codes.Unauthenticated:
			httpx.WriteError(w, "Unauthorized", http.StatusUnauthorized)
		default:
			httpx.WriteError(w, "Failed to get profile", http.StatusInternalServerError)
		}
		return
	}

	httpx.WriteJSON(w, profileResponse(resp), http.StatusOK)
}

func profileResponse(profile *pb.ProfileResponse) ProfileResponse {
	return ProfileResponse{
		UserID:    profile.UserId,
		Username:  profile.Username,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		BirthDate: profile.BirthDate,
		Gender:    profile.Gender,
		Avatar:    profile.Avatar,
	}
}
