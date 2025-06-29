package utils

import (
	"time"
)

func GetCurrentTimestamp() string {
	return time.Now().Format(time.RFC3339)
}
