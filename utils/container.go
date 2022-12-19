package utils

import (
	"fmt"
	"math/rand"
	"reflect"
)

func AddToMap[K comparable, V any, M map[K]any](m M, v *V, k K) {
	if v == nil {
		return
	}
	m[k] = *v
}

func GetFromAnyMap[T any, K comparable](m map[K]any, k K) (T, error) {
	var result T
	var ok bool
	if val, exist := m[k]; !exist {
		return result, fmt.Errorf("key %v does not exist", k)
	} else if result, ok = val.(T); !ok {
		return result, fmt.Errorf("value does not has type %v, it's %v", reflect.TypeOf(result), reflect.TypeOf(val))
	} else {
		return result, nil
	}
}

func Sample[T any](samples []T, count int) []T {
	if len(samples) <= count {
		return samples
	}
	result := make([]T, 0, count)
	for _, v := range rand.Perm(len(samples))[:count] {
		result = append(result, samples[v])
	}
	return result
}

func Contain[T comparable](items []T, value T) bool {
	for _, item := range items {
		if item == value {
			return true
		}
	}
	return false
}
