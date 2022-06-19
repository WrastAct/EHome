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
	Width       int64   `json:"width"`
	Height      int64   `json:"height"`
	Image       string  `json:"image,omitempty"` // Path to the image
	Shape       Shape   `json:"shape"`           // To improve collision detection
}

func ValidateFurniture(v *validator.Validator, furniture *Furniture) {
	v.Check(furniture.ID != 0, "furniture_list_id", "must be provided")
	v.Check(furniture.ID > 0, "furniture_list_id", "must be positive number")
}
