package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/WrastAct/EHome/internal/validator"
)

type FurnitureList struct {
	FurnitureID int64 `json:"furniture_id"`
	RoomID      int64 `json:"room_id"`
	X           int64 `json:"x"`
	Y           int64 `json:"y"`
}

func ValidateFurnitureList(v *validator.Validator, flist *FurnitureList) {
	v.Check(flist.X >= 0, "x", "must be positive")
	v.Check(flist.Y >= 0, "y", "must be positive")
}

type FurnitureListModel struct {
	DB *sql.DB
}

func (fl FurnitureListModel) GetAll(id int64) ([]FurnitureList, error) {
	query := `
		SELECT furniture_id, x, y
		FROM room_furniture 
		WHERE room_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := fl.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	furnitureLists := []FurnitureList{}

	for rows.Next() {

		var furnitureList FurnitureList

		err := rows.Scan(
			&furnitureList.FurnitureID,
			&furnitureList.X,
			&furnitureList.Y,
		)
		if err != nil {
			return nil, err
		}

		furnitureLists = append(furnitureLists, furnitureList)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return furnitureLists, nil
}

func (fl FurnitureListModel) Insert(flist *FurnitureList) error {
	query := `
		INSERT INTO room_furniture (furniture_id, room_id, x, y)
		VALUES ($1, $2, $3, $4)`

	args := []interface{}{
		flist.FurnitureID,
		flist.RoomID,
		flist.X,
		flist.Y,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := fl.DB.ExecContext(ctx, query, args...)

	return err
}

func (fl FurnitureListModel) InsertTransaction(flist []FurnitureList) error {
	query := `
		INSERT INTO room_furniture (furniture_id, room_id, x, y)
		VALUES ($1, $2, $3, $4)`

	tx, err := fl.DB.Begin()
	if err != nil {
		return err
	}

	for _, val := range flist {
		args := []interface{}{
			val.FurnitureID,
			val.RoomID,
			val.X,
			val.Y,
		}

		_, err = tx.Exec(query, args...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	return err
}

func (fl FurnitureListModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM room_furniture
		WHERE room_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := fl.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
