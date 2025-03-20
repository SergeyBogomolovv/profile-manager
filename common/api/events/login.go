package events

import "time"

type UserLogin struct {
	ID   string    `json:"id"`
	IP   string    `json:"ip"`
	Time time.Time `json:"time"`
	Type string    `json:"type"`
}

const LoginTopic = "login"

const (
	NotificationLoginQueue = "notification_login_queue"
)

const (
	LoginTypeCredentials = "credentials"
	LoginTypeGoogle      = "google"
)
