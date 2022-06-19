package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/WrastAct/EHome/internal/data"
	"github.com/WrastAct/EHome/internal/validator"
)

func (app *application) editRoomHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "open room editor")
}

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

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showRoomHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	furniture1 := data.Furniture{
		ID:          1,
		Name:        "Chair",
		Price:       423.75,
		Description: "heh",
		Width:       25,
		Height:      25,
		Image:       "../img",
		Shape:       data.Circle,
	}

	furniture2 := data.Furniture{
		ID:     2,
		Name:   "Table",
		Price:  103.24,
		Width:  25,
		Height: 25,
		Image:  "../img",
		Shape:  data.Circle,
	}

	room := data.Room{
		ID:            id,
		Date:          time.Now(),
		Description:   "Hehw",
		Title:         "Custom room",
		Width:         500,
		Height:        300,
		FurnitureList: []data.Furniture{furniture1, furniture2},
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"room": room}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
