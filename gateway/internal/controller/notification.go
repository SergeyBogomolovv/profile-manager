package controller

import (
	"net/http"

	pb "github.com/SergeyBogomolovv/profile-manager/common/api/notification"
	"github.com/SergeyBogomolovv/profile-manager/common/httpx"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type notiController struct {
	client pb.NotificationClient
}

func NewNotificationController(client pb.NotificationClient) *notiController {
	return &notiController{client: client}
}

func (c *notiController) Init(r *chi.Mux) {
	r.Route("/notification", func(r chi.Router) {
		r.Post("/token", c.HandleToken)
	})
}

// HandleToken generates a token for Telegram notifications.
// @Summary Generate Telegram token
// @Description Generates a token for Telegram notifications
// @Tags notification
// @Accept json
// @Produce json
// @Success 200 {object} TokenResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /notification/token [post]
// @Security BearerAuth
func (c *notiController) HandleToken(w http.ResponseWriter, r *http.Request) {
	resp, err := c.client.GenerateTelegramToken(authCtx(r), &pb.GenerateTelegramTokenRequest{})

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			httpx.WriteError(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}
		switch st.Code() {
		case codes.Unauthenticated:
			httpx.WriteError(w, "Unauthorized", http.StatusUnauthorized)
		default:
			httpx.WriteError(w, "Failed to generate token", http.StatusInternalServerError)
		}
		return
	}

	httpx.WriteJSON(w, TokenResponse{Token: resp.Token}, http.StatusOK)
}
