package jsonmapper

import "encoding/json"

func FromJSON[T any](data []byte) T {
	var v T
	_ = json.Unmarshal(data, &v)
	return v
}

func ToJSON[T any](v T) []byte {
	bytes, _ := json.Marshal(v)
	return bytes
}
