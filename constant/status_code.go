package constant

type StatusCode int32

const (
	ServiceError          StatusCode = 10000
	UserNotLogin                     = 10001
	RequestParameterError            = 10002
	FailedToProcess                  = 10003
	NoAuthority                      = 10004
)

func (s StatusCode) String() string {
	switch s {
	case ServiceError:
		return "Service Error Happened"
	case UserNotLogin:
		return "User has not logged in"
	case RequestParameterError:
		return "RequestParameter has error"
	case FailedToProcess:
		return "Failed to process this request"
	case NoAuthority:
		return "User has no authority"
	default:
		return "Unknown Error"
	}
}
