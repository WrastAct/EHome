package data

import (
	"time"

	"github.com/WrastAct/EHome/internal/validator"
)

type Room struct {
	ID            int64       `json:"id"` // Unique integer ID for the Room
	Date          time.Time   `json:"-"`  // Timestamp when Room was created for our database
	Description   string      `json:"description,omitempty"`
	Title         string      `json:"title"` // Custom Title for Room created by user
	Width         int64       `json:"width"`
	Height        int64       `json:"height"`
	FurnitureList []Furniture `json:"furniture_list,omitempty"` // Furniture inside room
}

func ValidateRoom(v *validator.Validator, room *Room) {
	v.Check(room.Title != "", "title", "must be provided")
	v.Check(len(room.Title) <= 30, "title", "must not be more than 30 bytes long")

	v.Check(len(room.Description) <= 400, "description", "must not be more than 400 bytes long")

	v.Check(room.Width != 0, "width", "must be provided")
	v.Check(room.Height != 0, "height", "must be provided")

	for _, val := range room.FurnitureList {
		ValidateFurniture(v, &val)
	}
}
