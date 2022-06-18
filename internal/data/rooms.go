package data

import (
	"time"
)

type Room struct {
	ID            int64       `json:"id"` // Unique integer ID for the Room
	Data          time.Time   `json:"-"`  // Timestamp when Room was created for our database
	Description   string      `json:"description,omitempty"`
	Title         string      `json:"title"` // Custom Title for Room created by user
	Width         int64       `json:"width"`
	Height        int64       `json:"height"`
	FurnitureList []Furniture `json:"furniture_list,omitempty"` // Furniture inside room
}
