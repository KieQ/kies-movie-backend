package constant

import "time"

const (
	RequestID          = "X-Request-ID"
	RealIP             = "X-Real-IP"
	Account            = "account"
	TokenIP            = "ip"
	Token              = "Token"
	Lang               = "lang"
	RememberMeDuration = 14 * 24 * time.Hour
)
