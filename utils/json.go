package utils

import "encoding/json"

func JSONStatus(message string) []byte {
	m, _ := json.Marshal(struct {
		Message string `json:"message"`
	}{
		message,
	})
	return m
}
