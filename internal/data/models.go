package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Furniture FurnitureModel
	Room      RoomModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Furniture: FurnitureModel{DB: db},
		Room:      RoomModel{DB: db},
	}
}
