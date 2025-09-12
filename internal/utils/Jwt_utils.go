package utils

import "time"

func GetExpiryTime() time.Time {
	return time.Now().Add(time.Hour * 24)
}
