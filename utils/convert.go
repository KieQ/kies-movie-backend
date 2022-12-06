package utils

func DowncastWithDefault[T any](val any, defaultVal T) T {
	if val == nil {
		return defaultVal
	}
	if v, ok := val.(T); ok {
		return v
	}
	return defaultVal
}
