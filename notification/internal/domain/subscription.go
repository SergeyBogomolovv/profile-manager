package domain

type SubscriptionType string

const (
	SubscriptionTypeEmail    SubscriptionType = "email"
	SubscriptionTypeTelegram SubscriptionType = "telegram"
)

type Subscription struct {
	UserID  string
	Type    SubscriptionType
	Enabled bool
}
