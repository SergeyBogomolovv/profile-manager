package events

type UserRegister struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

const RegisterTopic = "register"

const (
	ProfileRegisterQueue      = "profile_register_queue"
	NotificationRegisterQueue = "notification_register_queue"
)
