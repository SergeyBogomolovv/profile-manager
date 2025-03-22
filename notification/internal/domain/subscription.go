package domain

type SubscriptionType string

const (
	SubscriptionTypeEmail    SubscriptionType = "email"
	SubscriptionTypeTelegram SubscriptionType = "telegram"
)

type Subscription struct {
	User    User
	Type    SubscriptionType
	Enabled bool
}
