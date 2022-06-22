package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

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
	v.Check(furniture.ID != 0, "furniture_id", "must be provided")
	v.Check(furniture.ID > 0, "furniture_id", "must be positive number")

	v.Check(furniture.Name != "", "furniture_name", "must be provided")
	v.Check(furniture.Price > 0, "furniture_price", "must be positive number")

	v.Check(furniture.Width != 0, "furniture_width", "must be provided")
	v.Check(furniture.Height != 0, "furniture_height", "must be provided")

	v.Check(furniture.Image != "", "furniture_image", "must be provided")
	v.Check(furniture.Shape == Rectangle ||
		furniture.Shape == Circle, "furniture_shape", "must be a correct value")
}

type FurnitureModel struct {
	DB *sql.DB
}

func (f FurnitureModel) Insert(furniture *Furniture) error {
	query := `
		INSERT INTO furniture (name, price, furniture_description, 
			furniture_width, furniture_height, image, shape)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING furniture_id`

	args := []interface{}{
		furniture.Name,
		furniture.Price,
		furniture.Description,
		furniture.Width,
		furniture.Height,
		furniture.Image,
		furniture.Shape,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return f.DB.QueryRowContext(ctx, query, args...).Scan(&furniture.ID)
}

func (f FurnitureModel) Get(id int64) (*Furniture, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT furniture_id, name, price, furniture_description,
			furniture_width, furniture_height, image, shape
		FROM furniture
		WHERE furniture_id = $1`

	var furniture Furniture

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := f.DB.QueryRowContext(ctx, query, id).Scan(
		&furniture.ID,
		&furniture.Name,
		&furniture.Price,
		&furniture.Description,
		&furniture.Width,
		&furniture.Height,
		&furniture.Image,
		&furniture.Shape,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &furniture, nil
}

func (f FurnitureModel) GetAll() ([]*Furniture, error) {
	query := `
		SELECT furniture_id, name, price, furniture_description,
			furniture_width, furniture_height, image, shape
		FROM furniture`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := f.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	furnitures := []*Furniture{} // I know it is uncountable

	for rows.Next() {

		var furniture Furniture

		err := rows.Scan(
			&furniture.ID,
			&furniture.Name,
			&furniture.Price,
			&furniture.Description,
			&furniture.Width,
			&furniture.Height,
			&furniture.Image,
			&furniture.Shape,
		)
		if err != nil {
			return nil, err
		}

		furnitures = append(furnitures, &furniture)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return furnitures, nil
}

func (f FurnitureModel) Update(furniture *Furniture) error {
	query := `
		UPDATE furniture
		SET name = $1, price = $2, furniture_description = $3, 
			furniture_width = $4, furniture_height = $5, image = $6,
			shape = $7
		WHERE furniture_id = $8`

	args := []interface{}{
		furniture.Name,
		furniture.Price,
		furniture.Description,
		furniture.Width,
		furniture.Height,
		furniture.Image,
		furniture.Shape,
		furniture.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := f.DB.ExecContext(ctx, query, args...)

	return err
}

func (f FurnitureModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM furniture
		WHERE furniture_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := f.DB.ExecContext(ctx, query, id)
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
