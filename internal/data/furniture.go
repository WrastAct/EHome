package data

import (
	"database/sql"

	"github.com/WrastAct/EHome/internal/validator"
)

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

type FurnitureModel struct {
	DB *sql.DB
}

func (f FurnitureModel) Insert(furniture *Furniture) error {
	return nil
}

func (f FurnitureModel) Get(id int64) (*Furniture, error) {
	return nil, nil
}

func (f FurnitureModel) Update(furniture *Furniture) error {
	return nil
}

func (f FurnitureModel) Delete(id int64) error {
	return nil
}
