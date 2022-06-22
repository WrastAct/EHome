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
		FurnitureID int64 `json:"furniture_id"`
		RoomID      int64 `json:"-"`
		X           int64 `json:"x"`
		Y           int64 `json:"y"`
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

	furnitureIDs, err := app.models.Furniture.GetAllID()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var furnitureList []data.FurnitureList
	for _, val := range input.FurnitureList {
		if !validator.InInts(val.FurnitureID, furnitureIDs...) {
			err = errors.New("no furniture with this id")
			app.badRequestResponse(w, r, err)
			return
		}

		furnitureList = append(furnitureList, data.FurnitureList{
			FurnitureID: val.FurnitureID,
			X:           val.X,
			Y:           val.Y,
		})
	}

	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	room := &data.Room{
		OwnerID:       user.ID,
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

	for key := range furnitureList {
		furnitureList[key].RoomID = room.ID
	}

	err = app.models.FurnitureList.InsertTransaction(furnitureList)
	if err != nil {
		app.serverErrorResponse(w, r, err)
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

	furnitureList, err := app.models.FurnitureList.GetAll(room.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	room.FurnitureList = furnitureList

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

	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	if room.OwnerID != user.ID {
		app.foreignRoomResponse(w, r)
		return
	}

	furnitureList, err := app.models.FurnitureList.GetAll(room.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	room.FurnitureList = furnitureList

	type furnitureInput struct {
		FurnitureID int64 `json:"furniture_id"`
		RoomID      int64 `json:"room_id"`
		X           int64 `json:"x"`
		Y           int64 `json:"y"`
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

	v := validator.New()

	if data.ValidateRoom(v, room); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if input.FurnitureList != nil {
		furnitureIDs, err := app.models.Furniture.GetAllID()
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		var furnitureList []data.FurnitureList
		for _, val := range input.FurnitureList {
			if !validator.InInts(val.FurnitureID, furnitureIDs...) {
				err = errors.New("no furniture with this id")
				app.badRequestResponse(w, r, err)
				return
			}

			furnitureList = append(furnitureList, data.FurnitureList{
				FurnitureID: val.FurnitureID,
				RoomID:      room.ID,
				X:           val.X,
				Y:           val.Y,
			})
		}

		for _, val := range furnitureList {
			if data.ValidateFurnitureList(v, &val); !v.Valid() {
				app.failedValidationResponse(w, r, v.Errors)
				return
			}
		}

		room.FurnitureList = furnitureList
		app.models.FurnitureList.Delete(room.ID)
		err = app.models.FurnitureList.InsertTransaction(furnitureList)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
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

	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	if room.OwnerID != user.ID {
		app.foreignRoomResponse(w, r)
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

func (app *application) listRoomHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string `json:"title"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
		data.Filters
	}
	v := validator.New()
	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Width = app.readInt(qs, "width", 0, v)
	input.Height = app.readInt(qs, "height", 0, v)

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafelist = []string{"id", "title", "room_width", "room_height", "-id", "-title", "-room_width", "-room_height"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	rooms, metadata, err := app.models.Room.GetAll(input.Title, input.Width, input.Height, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	for _, val := range rooms {
		furnitureList, err := app.models.FurnitureList.GetAll(val.ID)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		val.FurnitureList = furnitureList
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"rooms": rooms, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
