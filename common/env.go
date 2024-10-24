package common

import "syscall"


func EnvString(key string, fallback string) string {
	if value, exists := syscall.Getenv(key); exists {
		return value
	}
	return fallback
}