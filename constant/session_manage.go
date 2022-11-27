package constant

import "time"

const (
	LogID              = "LOG_ID"
	Account            = "account"
	Token              = "Token"
	RememberMe         = "remember_me"
	RememberMeDuration = 7 * 24 * time.Hour
	RefreshLimit       = 24 * time.Hour
)
