package constant

import "time"

const (
	RequestID          = "X-Request-ID"
	RealIP             = "X-Real-IP"
	Account            = "account"
	RememberMe         = "remember_me"
	TokenIP            = "ip"
	Token              = "Token"
	RememberMeDuration = 7 * 24 * time.Hour
	RefreshLimit       = 24 * time.Hour
)
