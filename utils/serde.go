package utils

import "encoding/json"

func ToJSON(data interface{}) string {
	d, _ := json.Marshal(data)
	return string(d)
}

func FromJSON[T any](jsonStr string) T {
	var result T
	_ = json.Unmarshal([]byte(jsonStr), &result)
	return result
}
