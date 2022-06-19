package data

import (
	"database/sql"
	"errors"
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

type RoomModel struct {
	DB *sql.DB
}

func (r RoomModel) Insert(room *Room) error {
	query := `
		INSERT INTO room (room_description, title, room_width, room_height)
		VALUES ($1, $2, $3, $4)
		RETURNING room_id, date`

	args := []interface{}{room.Description, room.Title, room.Width, room.Height}

	return r.DB.QueryRow(query, args...).Scan(&room.ID, &room.Date)
}

func (r RoomModel) Get(id int64) (*Room, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT room_id, date, room_description, title, room_width, room_height
		FROM room
		WHERE room_id = $1`

	var room Room

	err := r.DB.QueryRow(query, id).Scan(
		&room.ID,
		&room.Date,
		&room.Description,
		&room.Title,
		&room.Width,
		&room.Height,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &room, nil
}

func (r RoomModel) Update(room *Room) error {
	query := `
		UPDATE room
		SET room_description = $1, title = $2, room_width = $3, room_height = $4
		WHERE room_id = $5`

	args := []interface{}{
		room.Description,
		room.Title,
		room.Width,
		room.Height,
		room.ID,
	}

	_, err := r.DB.Exec(query, args...)

	return err
}

func (r RoomModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM room
		WHERE room_id = $1`

	result, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
