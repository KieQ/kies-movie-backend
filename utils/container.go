package utils

func AddToMap[K comparable, V any, M map[K]any](m M, v *V, k K) {
	if v == nil {
		return
	}
	m[k] = *v
}
