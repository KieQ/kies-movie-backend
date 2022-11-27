package constant

type StatusCode int32

const (
	ServiceError          StatusCode = 10000
	RequestParameterError            = 10001
)

func (s StatusCode) String() string {
	switch s {
	case ServiceError:
		return "Service Error Happened"
	case RequestParameterError:
		return "RequestParameter has error"
	default:
		return "Unknown Error"
	}
}
