package constant

import "time"

const (
	TokenIP            = "ip"    // in JWT
	Token              = "Token" // in cookie
	Lang               = "lang"  // in context
	RememberMeDuration = 14 * 24 * time.Hour
)
