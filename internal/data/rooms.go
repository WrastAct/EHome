package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return r.DB.QueryRowContext(ctx, query, args...).Scan(&room.ID, &room.Date)
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, id).Scan(
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

func (r RoomModel) GetAll(title string, width int, height int, filters Filters) ([]*Room, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), room_id, date, room_description, title, room_width, room_height
		FROM room
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (room_width <= $2 OR $2 = 0)
		AND (room_height <= $3 OR $3 = 0)
		ORDER BY %s %s, room_id ASC
		LIMIT $4 OFFSET $5`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{title, width, height, filters.limit(), filters.offset()}

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	rooms := []*Room{}

	for rows.Next() {

		var room Room

		err := rows.Scan(
			&totalRecords,
			&room.ID,
			&room.Date,
			&room.Description,
			&room.Title,
			&room.Width,
			&room.Height,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		rooms = append(rooms, &room)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return rooms, metadata, nil
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.DB.ExecContext(ctx, query, args...)

	return err
}

func (r RoomModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM room
		WHERE room_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := r.DB.ExecContext(ctx, query, id)
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
