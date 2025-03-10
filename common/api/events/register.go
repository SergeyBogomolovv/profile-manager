package events

type UserRegister struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

const (
	RegisterExchange = "register_exchange"
)

const (
	ProfileQueue = "profile_queue"
)
