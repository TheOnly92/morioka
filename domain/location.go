package domain

import (
	"time"
)

func GetLocation() *time.Location {
	location, _ := time.LoadLocation("Asia/Tokyo")
	return location
}
