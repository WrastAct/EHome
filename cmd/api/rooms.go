package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/WrastAct/EHome/internal/data"
	"github.com/WrastAct/EHome/internal/validator"
)

func (app *application) createRoomHandler(w http.ResponseWriter, r *http.Request) {
	type furnitureInput struct {
		ID int64
		X  int64
		Y  int64
	}

	var input struct {
		Description   string           `json:"description"`
		Title         string           `json:"title"`
		Width         int64            `json:"width"`
		Height        int64            `json:"height"`
		FurnitureList []furnitureInput `json:"furniture_list"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var furnitureList []data.Furniture
	for _, val := range input.FurnitureList {
		furnitureList = append(furnitureList, data.Furniture{ID: val.ID})
	}

	room := &data.Room{
		Description:   input.Description,
		Title:         input.Title,
		Width:         input.Width,
		Height:        input.Height,
		FurnitureList: furnitureList,
	}

	v := validator.New()

	if data.ValidateRoom(v, room); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Room.Insert(room)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/rooms/%d", room.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"room": room}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showRoomHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	room, err := app.models.Room.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"room": room}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateRoomHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	room, err := app.models.Room.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	type furnitureInput struct {
		ID int64
		X  int64
		Y  int64
	}

	var input struct {
		Description   *string          `json:"description"`
		Title         *string          `json:"title"`
		Width         *int64           `json:"width"`
		Height        *int64           `json:"height"`
		FurnitureList []furnitureInput `json:"furniture_list"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Description != nil {
		room.Description = *input.Description
	}

	if input.Title != nil {
		room.Title = *input.Title
	}

	if input.Width != nil {
		room.Width = *input.Width
	}

	if input.Height != nil {
		room.Height = *input.Height
	}

	if input.FurnitureList != nil {
		var furnitureList []data.Furniture
		for _, val := range input.FurnitureList {
			furnitureList = append(furnitureList, data.Furniture{ID: val.ID})
		}
		room.FurnitureList = furnitureList
	}

	v := validator.New()

	if data.ValidateRoom(v, room); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Room.Update(room)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"room": room}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteRoomHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Room.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "room successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
