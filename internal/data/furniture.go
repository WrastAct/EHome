package data

import "github.com/WrastAct/EHome/internal/validator"

type Shape int

const (
	Rectangle Shape = iota
	Circle
)

type Furniture struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description,omitempty"`
	X           int64   `json:"x"` // X coordinate relative to Room's left position.
	Y           int64   `json:"y"` // Y coordinate relative to Room's top position.
	Width       int64   `json:"width"`
	Height      int64   `json:"height"`
	Image       string  `json:"image,omitempty"` // Path to the image
	Shape       Shape   `json:"shape"`           // To improve collision detection
}

func ValidateFurniture(v *validator.Validator, furniture *Furniture) {
	v.Check(furniture.ID != 0, "furniture_list_id", "must be provided")
	v.Check(furniture.ID > 0, "furniture_list_id", "must be positive number")

	v.Check(len(furniture.Name) <= 50, "furniture_list_name", "must not be more than 50 bytes long")

	v.Check(len(furniture.Description) <= 400, "furniture_list_description", "must not be more than 400 bytes long")

	v.Check(len(furniture.Image) <= 200, "furniture_list_image", "must not be more than 200 bytes long")

	v.Check(furniture.Width != 0, "furniture_list_width", "must be provided")
	v.Check(furniture.Height != 0, "furniture_list_height", "must be provided")
}
