package constant

type StatusCode int32

const (
	ServiceError          StatusCode = 10000
	UserNotLogin                     = 10001
	UserIPChanged                    = 10002
	RequestParameterError            = 10003
	FailedToProcess                  = 10004
	NoAuthority                      = 10005
)
